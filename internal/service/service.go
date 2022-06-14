package service

import (
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

}

func NewService(deps Deps) *Services {
	authStorage := storage.NewAuthStorage()
	return &Services{
		Auth: authStorage,
	}
}