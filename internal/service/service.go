//go:generate mockgen -source service.go -destination mock/service_mock.go -package mock
package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/internal/storage/psql"
	"github.com/Edbeer/restapi/internal/storage/redis"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
)

// Auth Service interface
type Auth interface {
	Register(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindUsersByName(ctx context.Context, name string, pq *utils.PaginationQuery) (*entity.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error)
	Login(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
}

// News service interface
type News interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error)
	Update(ctx context.Context, news *entity.News) (*entity.News, error)
	GetNews(ctx context.Context, pq *utils.PaginationQuery) (*entity.NewsList, error)
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.NewsBase, error)
	SearchNews(ctx context.Context, pq *utils.PaginationQuery, title string) (*entity.NewsList, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
}

// Comments Service interface
type Comments interface {
	Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	Update(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	GetAllByNewsID(ctx context.Context, newsID uuid.UUID, pq *utils.PaginationQuery) (*entity.CommentsList, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentBase, error)
	Delete(ctx context.Context, commentID uuid.UUID) error
}

// Session service interface
type Session interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
}

type Services struct {
	Auth     *AuthService
	News     *NewsService
	Comments *CommentsService
	Session  *SessionService
}

type Deps struct {
	Logger       logger.Logger
	Config       *config.Config
	PsqlStorage  *psql.Storage
	RedisStorage *redisrepo.Storage
}

func NewService(deps Deps) *Services {
	authService := NewAuthService(deps.Config, deps.PsqlStorage.Auth, deps.RedisStorage.Auth, deps.Logger)
	newsService := NewNewsService(deps.Config, deps.PsqlStorage.News, deps.RedisStorage.News, deps.Logger)
	commentsService := NewCommentsService(deps.Config, deps.PsqlStorage.Comments, deps.Logger)
	sessionService := NewSessionService(deps.Config, deps.RedisStorage.Session, deps.Logger)
	return &Services{
		Auth:     authService,
		News:     newsService,
		Comments: commentsService,
		Session:  sessionService,
	}
}
