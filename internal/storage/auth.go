package storage

import (
	"github.com/Edbeer/restapi/pkg/logger"
)

// Auth Storage
type AuthStorage struct {
	logger *logger.Logger
}

// Auth Storage constructor
func NewAuthStorage(logger *logger.Logger) *AuthStorage {
	return &AuthStorage{logger: logger}
}

// Create user
func (s *AuthStorage) Create() error {
	return nil
}
