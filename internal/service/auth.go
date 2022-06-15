package service

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/storage"
)

// Auth service
type AuthService struct {
	config  *config.Config
	storage storage.Auth
}

// Auth service constructor
func NewAuthService(config *config.Config, storage storage.Auth) *AuthService {
	return &AuthService{config: config, storage: storage}
}

// Create new user
func (a *AuthService) Create() error {
	return a.storage.Create()
}
