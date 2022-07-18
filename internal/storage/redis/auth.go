package redisrepo

import (
	"context"
	"encoding/json"

	"github.com/Edbeer/restapi/internal/entity"
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

// Get user by id
func (s *AuthStorage) GetByIDCtx(ctx context.Context, key string) (*entity.User, error) {
	userBytes, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	user := &entity.User{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, err
	}
	return user, nil
}