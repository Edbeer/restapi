package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/logger"
)

// Session redis storage interface
type SessionRedis interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
}

// Session service
type SessionService struct {
	config         *config.Config
	logger         logger.Logger
	sessionStorage SessionRedis
}

// SessionService constructor
func NewSessionService(config *config.Config, sessionStorage SessionRedis, logger logger.Logger) *SessionService {
	return &SessionService{
		config:         config,
		logger:         logger,
		sessionStorage: sessionStorage,
	}
}

// Create session
func (s *SessionService) CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error) {
	return s.sessionStorage.CreateSession(ctx, session, expire)
}

// Get session by id
func (s *SessionService) GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error) {
	return s.sessionStorage.GetSessionByID(ctx, sessionID)
}

// Delete session by id
func (s *SessionService) DeleteSessionByID(ctx context.Context, sessionID string) error {
	return s.sessionStorage.DeleteSessionByID(ctx, sessionID)
}

