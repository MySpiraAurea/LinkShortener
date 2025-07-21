package storage

import (
    "encoding/json"
    "os"
    "sync"
)

type Storage struct {
    mutex sync.RWMutex
    links map[string]string
}

func NewStorage(filename string) (*Storage, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        if os.IsNotExist(err) {
            return &Storage{
                links: make(map[string]string),
            }, nil
        }
        return nil, err
    }

    var links map[string]string
    if err := json.Unmarshal(data, &links); err != nil {
        return nil, err
    }

    return &Storage{
        mutex: sync.RWMutex{},
        links: links,
    }, nil
}

func (s *Storage) GetOriginalURL(shortID string) (string, bool) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()

    url, exists := s.links[shortID]
    return url, exists
}

func (s *Storage) AddShortURL(shortID, originalURL string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    s.links[shortID] = originalURL
}

func (s *Storage) SaveToFile(filename string) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    data, err := json.Marshal(s.links)
    if err != nil {
        return err
    }

    return os.WriteFile(filename, data, 0644)
}