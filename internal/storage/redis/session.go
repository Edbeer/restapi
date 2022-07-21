package redisrepo

import (
	"github.com/Edbeer/restapi/config"
	"github.com/go-redis/redis/v9"
)

// Session storage
type SessionStorage struct {
	config *config.Config
	redis  *redis.Client
}

// Session storage constructor
func NewSessionStorage(config *config.Config, redis *redis.Client) *SessionStorage {
	return &SessionStorage{config: config, redis: redis}
}

