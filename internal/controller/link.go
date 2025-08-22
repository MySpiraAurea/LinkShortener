package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"link-shortener/internal/storage"
	"link-shortener/internal/utils"
	"log/slog"
)

type LinkController struct {
	storage storage.LinkRepository
}

func NewLinkController(store storage.LinkRepository) *LinkController {
	return &LinkController{storage: store}
}

func (lc *LinkController) CreateShortLink(ctx context.Context, originalURL string) (string, error) {
	if originalURL == "" {
		return "", errors.New("original URL не может быть пустым")
	}

	var shortID string
	var err error

	// Ретрай до 5 раз
	for i := 0; i < 5; i++ {
		shortID = utils.GenerateShortID()
		slog.Debug("Генерация короткого ID", "попытка", i+1, "id", shortID)

		// Проверяем, существует ли уже такой ID
		_, exists := lc.storage.GetOriginalURL(shortID)
		if exists {
			slog.Debug("ID уже существует, пробуем снова", "id", shortID)
			continue
		}

		// Пытаемся сохранить
		if err = lc.storage.AddShortURL(shortID, originalURL); err == nil {
			slog.Info("Ссылка успешно создана", "short_id", shortID, "url", originalURL)
			break
		}

		// Если ошибка не связана с коллизией — выходим
		if !isConflictError(err) {
			return "", fmt.Errorf("ошибка сохранения ссылки: %w", err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	if err != nil {
		return "", fmt.Errorf("не удалось создать уникальный ID за 5 попыток: %w", err)
	}

	// Сохраняем в файл, если хранилище поддерживает
	if saver, ok := lc.storage.(storage.FileSaver); ok {
		if saveErr := saver.SaveToFile(); saveErr != nil {
			return "", fmt.Errorf("ошибка сохранения в файл: %w", saveErr)
		}
	}

	return shortID, nil
}

func (lc *LinkController) Ping(ctx context.Context) error {
	return lc.storage.Ping()
}

// isConflictError проверяет, является ли ошибка конфликтом (например, дубликат)
func isConflictError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key")
}

func (lc *LinkController) GetOriginalLink(ctx context.Context, shortID string) (string, error) {
    url, exists := lc.storage.GetOriginalURL(shortID)
    if !exists {
        return "", errors.New("ссылка не найдена")
    }
    return url, nil
}