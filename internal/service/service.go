package service

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/storage/psql"
	"github.com/Edbeer/restapi/internal/storage/redis"
)

type Services struct {
	Auth *AuthService
}

type Deps struct {
	Config       *config.Config
	PsqlStorage  *psql.Storage
	RedisStorage *redisrepo.Storage
}

func NewService(deps Deps) *Services {
	authService := NewAuthService(deps.Config, deps.PsqlStorage.Auth, deps.RedisStorage.Auth)
	return &Services{
		Auth: authService,
	}
}
