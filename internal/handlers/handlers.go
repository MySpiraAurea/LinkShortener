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

    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    contentType := r.Header.Get("Content-Type")
    if contentType != "application/json" {
        http.Error(w, "Требуется Content-Type: application/json", http.StatusUnsupportedMediaType)
        return
    }

    var req struct {
        URL string `json:"url"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Некорректный JSON", http.StatusBadRequest)
        logger.WarnContext(ctx, "Ошибка декодирования JSON", "error", err)
        return
    }

    if req.URL == "" {
        http.Error(w, "URL не может быть пустым", http.StatusBadRequest)
        logger.WarnContext(ctx, "URL отсутствует")
        return
    }

    if _, err := url.ParseRequestURI(req.URL); err != nil {
        http.Error(w, "Некорректный URL", http.StatusBadRequest)
        logger.WarnContext(ctx, "Некорректный URL", "url", req.URL)
        return
    }

    shortURL, err := h.controller.CreateShortLink(ctx, req.URL)
    if err != nil {
        http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
        logger.ErrorContext(ctx, "Ошибка создания ссылки", "error", err)
        return
    }

    response := map[string]string{"result": shortURL}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    logger.InfoContext(ctx, "Ссылка сокращена", "short_url", shortURL)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    logger := slog.With("method", "GET", "path", r.URL.Path)
    logger.InfoContext(ctx, "Начало обработки запроса")

    shortID := r.URL.Path[1:]
    originalURL, err := h.controller.GetOriginalLink(ctx, shortID)

    if err != nil {
        http.NotFound(w, r)
        slog.Warn("Ссылка не найдена", "short_id", shortID)
        return
    }

    http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
    slog.Info("Редирект выполнен", "short_id", shortID, "url", originalURL)
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