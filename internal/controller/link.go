package controller

import (
    "context"
    "errors"
    "fmt"

    "link-shortener/internal/storage"

    "crypto/rand"
    "math/big"
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

    shortID := generateShortID()
    lc.storage.AddShortURL(shortID, originalURL)

    if fileStore, ok := lc.storage.(interface{ SaveToFile() error }); ok {
        if err := fileStore.SaveToFile(); err != nil {
            return "", fmt.Errorf("ошибка сохранения в файл: %w", err)
        }
    }

    return fmt.Sprintf("http://localhost:8080/%s", shortID), nil
}

func (lc *LinkController) GetOriginalLink(ctx context.Context, shortID string) (string, error) {
    url, exists := lc.storage.GetOriginalURL(shortID)
    if !exists {
        return "", errors.New("ссылка не найдена")
    }
    return url, nil
}

func generateShortID() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    const length = 6
    id := make([]byte, length)
    for i := range id {
        num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        id[i] = charset[num.Int64()]
    }
    return string(id)
}

func (lc *LinkController) Ping(ctx context.Context) error {
    return lc.storage.Ping()
}
