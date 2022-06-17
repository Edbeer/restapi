package redis

import (
	"github.com/go-redis/redis/v9"
)

// Auth Storage
type AuthStorage struct {
	redis  *redis.Client
}

// Auth Storage constructor
func NewAuthStorage(redis *redis.Client) *AuthStorage {
	return &AuthStorage{redis: redis}
}

// Create user
func (s *AuthStorage) Create() error {
	return nil
}