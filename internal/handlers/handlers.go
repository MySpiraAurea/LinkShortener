package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"link-shortener/internal/controller"
)

type Handler struct {
	controller *controller.LinkController
	baseURL    string
}

func NewHandler(ctrl *controller.LinkController, baseURL string) *Handler {
	return &Handler{
		controller: ctrl,
		baseURL:    baseURL,
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("method", r.Method, "path", r.URL.Path)
	logger.InfoContext(ctx, "Начало обработки запроса")

	// Проверка метода
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Требуется Content-Type: application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Декодирование JSON
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		logger.WarnContext(ctx, "Ошибка декодирования JSON", "error", err)
		return
	}

	// Проверка, что URL не пустой
	if req.URL == "" {
		http.Error(w, "URL не может быть пустым", http.StatusBadRequest)
		logger.WarnContext(ctx, "URL отсутствует")
		return
	}

	// Валидация URL
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "Некорректный URL", http.StatusBadRequest)
		logger.WarnContext(ctx, "Некорректный URL", "url", req.URL)
		return
	}

	// Вызов контроллера
	shortID, err := h.controller.CreateShortLink(ctx, req.URL)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		logger.ErrorContext(ctx, "Ошибка создания ссылки", "error", err)
		return
	}

	// Формирование полного URL
	shortURL := fmt.Sprintf("%s/%s", h.baseURL, shortID)

	// Формирование ответа
	response := map[string]string{"result": shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(response)
	logger.InfoContext(ctx, "Ссылка создана", "short_id", shortID, "original_url", req.URL)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("method", "GET", "path", r.URL.Path)
	logger.InfoContext(ctx, "Начало обработки запроса")

	shortID := r.URL.Path[1:]

	// Валидация shortID
	if !isValidShortID(shortID) {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		logger.Warn("Некорректный формат shortID", "short_id", shortID)
		return
	}

	originalURL, err := h.controller.GetOriginalLink(ctx, shortID)
	if err != nil {
		http.NotFound(w, r)
		logger.Warn("Ссылка не найдена", "short_id", shortID)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	logger.Info("Редирект выполнен", "short_id", shortID, "url", originalURL)
}

// isValidShortID проверяет, что ID состоит из 6 букв/цифр
func isValidShortID(id string) bool {
	if len(id) != 6 {
		return false
	}
	for _, c := range id {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("method", r.Method, "path", r.URL.Path)
	logger.InfoContext(ctx, "Проверка подключения к хранилищу")

	if err := h.controller.Ping(ctx); err != nil {
		logger.ErrorContext(ctx, "Проверка не пройдена", "error", err)
		http.Error(w, "DB not available", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}