package service

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/internal/storage"
)

// Auth service
type AuthService struct {
	logger *logger.Logger
	config *config.Config
	storage storage.AuthStorage
}

// Auth service constructor
func NewAuthService(logger *logger.Logger, config *config.Config, storage storage.AuthStorage) *AuthService {
	return &AuthService{logger: logger, config: config, storage: storage}
}

// Create new user
func (a *AuthService) Create() error {
	return a.storage.Create()
}