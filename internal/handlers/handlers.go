package handlers

import (
	"encoding/json"
	"link-shortener/internal/controller"
	"log/slog"
	"net/http"
	"net/url"
)

type Handler struct {
    controller *controller.LinkController
}

func NewHandler(ctrl *controller.LinkController) *Handler {
    return &Handler{controller: ctrl}
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
    shortURL, err := h.controller.CreateShortLink(ctx, req.URL)
    if err != nil {
        http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
        logger.ErrorContext(ctx, "Ошибка создания ссылки", "error", err)
        return
    }

    // Формирование ответа
    response := map[string]string{"result": shortURL}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    logger.InfoContext(ctx, "Ссылка сокращена", "short_url", shortURL)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    logger := slog.With("method", "GET", "path", r.URL.Path)
    logger.InfoContext(ctx, "Начало обработки запроса")

    shortID := r.URL.Path[1:] // Убираем первый символ "/"
    originalURL, err := h.controller.GetOriginalLink(ctx, shortID)

    if err != nil {
        http.NotFound(w, r)
        slog.Warn("Ссылка не найдена", "short_id", shortID)
        return
    }

    http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
    slog.Info("Редирект выполнен", "short_id", shortID, "url", originalURL)
}