package handlers

import (
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"

    "link-shortener/internal/storage"
    "link-shortener/internal/utils"
)

type Handler struct {
    storage *storage.Storage
}

func NewHandler(storage *storage.Storage) *Handler {
    return &Handler{storage: storage}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    logger := slog.With("method", "POST", "path", "/api/shorten")
    logger.InfoContext(ctx, "Начало обработки запроса")

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

    shortID := utils.GenerateShortID()
    h.storage.AddShortURL(shortID, req.URL)
    h.storage.SaveToFile("links.json") // Сохраняем сразу

    response := map[string]string{"result": fmt.Sprintf("http://localhost:8080/%s", shortID)}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    logger.InfoContext(ctx, "Ссылка сокращена", "short_id", shortID)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
    shortID := r.URL.Path[1:] // Убираем первый символ "/"
    originalURL, exists := h.storage.GetOriginalURL(shortID)

    if !exists {
        http.NotFound(w, r)
        slog.Warn("Ссылка не найдена", "short_id", shortID)
        return
    }

    http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
    slog.Info("Редирект выполнен", "short_id", shortID, "url", originalURL)
}