package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/go-redis/redis/v9"
)

// News storage
type NewsStorage struct {
	redis *redis.Client
}

// News storage constructor
func NewNewsStorage(redis *redis.Client) *NewsStorage {
	return &NewsStorage{redis: redis}
}

// Get news by id
func (n *NewsStorage) GetNewsByIDCtx(ctx context.Context, key string) (*entity.NewsBase, error) {
	newsBytes, err := n.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	news := &entity.NewsBase{}
	if err := json.Unmarshal(newsBytes, news); err != nil {
		return nil, err
	}

	return news, nil
}

// Cache news item
func (n *NewsStorage) SetNewsCtx(ctx context.Context, key string, seconds int, news *entity.NewsBase) error {
	newsBytes, err := json.Marshal(news)
	if err != nil {
		return err
	}

	if err := n.redis.Set(ctx, key, newsBytes, time.Second * time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}