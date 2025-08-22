package storage

import (
    "encoding/json"
    "os"
    "sync"
)

type LinkRepository interface {
    GetOriginalURL(shortID string) (string, bool)
    AddShortURL(shortID, originalURL string) error
    Ping() error
}

type FileSaver interface {
    SaveToFile() error
}

type Storage struct {
    mutex    sync.RWMutex
    links    map[string]string
    filename string
}
func NewStorage(filename string) (*Storage, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        if os.IsNotExist(err) {
            return &Storage{
                links:    make(map[string]string),
                filename: filename,
            }, nil
        }
        return nil, err
    }

    var links map[string]string
    if err := json.Unmarshal(data, &links); err != nil {
        return nil, err
    }

    return &Storage{
        mutex:    sync.RWMutex{},
        links:    links,
        filename: filename,
    }, nil
}

func (s *Storage) GetOriginalURL(shortID string) (string, bool) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()

    url, exists := s.links[shortID]
    return url, exists
}

func (s *Storage) AddShortURL(shortID, originalURL string) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    s.links[shortID] = originalURL
    return nil
}

func (s *Storage) SaveToFile() error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    data, err := json.Marshal(s.links)
    if err != nil {
        return err
    }

    return os.WriteFile(s.filename, data, 0644)
}

func (s *Storage) Ping() error {
    return nil
}

type InMemoryStorage struct {
    mutex sync.RWMutex
    links map[string]string
}

func NewInMemory() *InMemoryStorage {
    return &InMemoryStorage{
        links: make(map[string]string),
    }
}

func (s *InMemoryStorage) GetOriginalURL(shortID string) (string, bool) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    url, exists := s.links[shortID]
    return url, exists
}

func (s *InMemoryStorage) AddShortURL(shortID, originalURL string) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    s.links[shortID] = originalURL
    return nil
}

func (s *InMemoryStorage) Ping() error {
    return nil
}