package redisrepo

import (
	"github.com/Edbeer/restapi/config"
	"github.com/go-redis/redis/v9"
)

type Storage struct {
	Auth    *AuthStorage
	News    *NewsStorage
	Session *SessionStorage
}

func NewStorage(redis *redis.Client, config *config.Config) *Storage {
	return &Storage{
		Auth:    NewAuthStorage(redis),
		News:    NewNewsStorage(redis),
		Session: NewSessionStorage(redis, config),
	}
}
