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
	StoragePsql  psql.Storage
	StorageRedis redis.Storage
	Config       *config.Config
}

func NewService(deps Deps) *Services {
	authService := NewAuthService(deps.Config, deps.StoragePsql.Auth, deps.StorageRedis.Auth)
	return &Services{
		Auth: authService,
	}
}
