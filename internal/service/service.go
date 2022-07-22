package service

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/storage/psql"
	"github.com/Edbeer/restapi/internal/storage/redis"
	"github.com/Edbeer/restapi/pkg/logger"
)

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
