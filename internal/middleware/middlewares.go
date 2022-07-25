package middleware

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
)

// Session service interface
type SessionService interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
}

// Auth Service interface
type AuthService interface {
	Register(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindUsersByName(ctx context.Context, name string, pq *utils.PaginationQuery) (*entity.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error)
	Login(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
}

// Middleware manager
type MiddlewareManager struct {
	sessionService SessionService
	authService    AuthService
	config         *config.Config
	origins        []string
	logger         logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(sessionService SessionService, authService AuthService, config *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		sessionService: sessionService,
		authService:    authService,
		config:         config,
		origins:        origins,
		logger:         logger,
	}
}
