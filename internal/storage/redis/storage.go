//go:generate mockgen -source storage.go -destination mock/redis_storage_mock.go -package mock
package redisrepo

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/go-redis/redis/v9"
)

// News StorageRedis interface
type NewsRedis interface {
	GetNewsByIDCtx(ctx context.Context, key string) (*entity.NewsBase, error)
	SetNewsCtx(ctx context.Context, key string, seconds int, news *entity.NewsBase) error
	DeleteNewsCtx(ctx context.Context, key string) error
}

// Auth StorageRedis interface
type AuthRedis interface {
	GetByIDCtx(ctx context.Context, key string) (*entity.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *entity.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}

// Session redis storage interface
type Sessionredis interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
}

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
