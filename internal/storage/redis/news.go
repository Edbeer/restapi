package redisrepo

import "github.com/go-redis/redis/v9"

// News storage
type NewsStorage struct {
	redis *redis.Client
}

// News storage constructor
func NewNewsStorage(redis *redis.Client) *NewsStorage {
	return &NewsStorage{redis: redis}
}

func (n *NewsStorage) Create() error {
	return nil
}