package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

// Auth Storage
type AuthStorage struct {
	redis *redis.Client
}

// Auth Storage constructor
func NewAuthStorage(redis *redis.Client) *AuthStorage {
	return &AuthStorage{redis: redis}
}

// Get user by id
func (s *AuthStorage) GetByIDCtx(ctx context.Context, key string) (*entity.User, error) {
	userBytes, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "AuthStorageRedis.GetByIDCtx.Get")
	}
	user := &entity.User{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, errors.Wrap(err, "AuthStorageRedis.GetByIDCtx.Unmarshal")
	}
	return user, nil
}

// Cache user with duration in seconds
func (s *AuthStorage) SetUserCtx(ctx context.Context, key string, seconds int, user *entity.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "AuthStorageRedis.SetUserCtx.Marshal")
	}

	if err := s.redis.Set(ctx, key, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "AuthStorageRedis.SetUserCtx.Set")
	}

	return nil
}

// Delete user by key
func (s *AuthStorage) DeleteUserCtx(ctx context.Context, key string) error {
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "AuthStorageRedis.DeleteUserCtx.Del")
	}
	return nil
}
