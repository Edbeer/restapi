package redisrepo

import (
	"context"
	"log"
	"testing"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func SetupSessionRedis() *SessionStorage {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	sessionRedisStorage := NewSessionStorage(client, nil)
	return sessionRedisStorage
}

func Test_CreateSession(t *testing.T) {
	t.Parallel()

	sessionRedisStorage := SetupSessionRedis()

	t.Run("CreateSession", func(t *testing.T) {
		sessionId := uuid.New()
		session := &entity.Session{
			SessionID: sessionId.String(),
			UserID: sessionId,
		}

		s, err := sessionRedisStorage.CreateSession(context.Background(), session, 10)
		require.NoError(t, err)
		require.NotEqual(t, s, "")
	})
}

func Test_GetSessionByID(t *testing.T) {
	t.Parallel()

	sessionRedisStorage := SetupSessionRedis()

	t.Run("GetSessionByID", func(t *testing.T) {
		sessionId := uuid.New()
		session := &entity.Session{
			SessionID: sessionId.String(),
			UserID: sessionId,
		}

		createdSession, err := sessionRedisStorage.CreateSession(context.Background(), session, 10)
		require.NoError(t, err)
		require.NotEqual(t, createdSession, "")

		s, err := sessionRedisStorage.GetSessionByID(context.Background(), createdSession)
		require.NoError(t, err)
		require.NotNil(t, s)
	})
}

func Test_DeleteSessionByID(t *testing.T) {
	t.Parallel()

	sessionRedisStorage := SetupSessionRedis()

	t.Run("DeleteSessionByID", func(t *testing.T) {
		sessionId := uuid.New()
		session := &entity.Session{
			SessionID: sessionId.String(),
			UserID: sessionId,
		}

		createdSession, err := sessionRedisStorage.CreateSession(context.Background(), session, 10)
		require.NoError(t, err)
		require.NotEqual(t, createdSession, "")

		err = sessionRedisStorage.DeleteSessionByID(context.Background(), createdSession)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}