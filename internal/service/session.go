package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/logger"
)

// Session redis storage interface
type SessionStorage interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
}

// Session service
type SessionService struct {
	config         *config.Config
	logger         logger.Logger
	sessionStorage SessionStorage
}

// SessionService constructor
func NewSessionService(config *config.Config, sessionStorage SessionStorage, logger logger.Logger) *SessionService {
	return &SessionService{
		config:         config,
		logger:         logger,
		sessionStorage: sessionStorage,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error) {
	return s.sessionStorage.CreateSession(ctx, session, expire)
}

