package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetOriginalURL(shortID string) (string, bool) {
	args := m.Called(shortID)
	return args.String(0), args.Bool(1)
}

func (m *MockStorage) AddShortURL(shortID, originalURL string) error {
	args := m.Called(shortID, originalURL)
	return args.Error(0)
}

func (m *MockStorage) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateShortLink_Success(t *testing.T) {
	mockStore := new(MockStorage)
	mockStore.On("GetOriginalURL", mock.Anything).Return("", false)
	mockStore.On("AddShortURL", mock.Anything, "https://example.com").Return(nil)

	ctrl := NewLinkController(mockStore)
	shortID, err := ctrl.CreateShortLink(context.Background(), "https://example.com")

	assert.NoError(t, err)
	assert.Len(t, shortID, 6)
	mockStore.AssertExpectations(t)
}