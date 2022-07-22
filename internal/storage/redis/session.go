package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Session storage
type SessionStorage struct {
	redis  *redis.Client
	config *config.Config
}

// Session storage constructor
func NewSessionStorage(redis *redis.Client, config *config.Config) *SessionStorage {
	return &SessionStorage{redis: redis, config: config}
}

// Create session in redis
func (s *SessionStorage) CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error) {

	session.SessionID = uuid.New().String()
	sessionKey := s.createSessionKey(session.SessionID)

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return "" , errors.Wrap(err, "SessionStorage.CreateSession.Marshal")
	}
	if err = s.redis.Set(ctx, sessionKey, sessionBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "" , errors.Wrap(err, "SessionStorage.CreateSession.Set")
	}

	return sessionKey, nil
}

func (s *SessionStorage) createSessionKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.config.Session.Prefix, sessionID)
} 