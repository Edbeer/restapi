package redisrepo

import "github.com/go-redis/redis/v9"

type Storage struct {
	Auth *AuthStorage
}

func NewStorage(redis *redis.Client) *Storage {
	return &Storage{
		Auth: NewAuthStorage(redis),
	}
}