package service

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/storage"
)

// Auth Service interface
type Auth interface {
	Create() error
}

type Services struct {
	Auth Auth
}

type Deps struct {
	Storage storage.Storage
	Config  *config.Config
}

func NewService(deps Deps) *Services {
	authService := NewAuthService(deps.Config, deps.Storage.Auth)
	return &Services{
		Auth: authService,
	}
}
