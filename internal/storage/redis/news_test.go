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

func SetupNewsRedis() *NewsStorage {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	newsRedisStorage := NewNewsStorage(client)
	return newsRedisStorage 
}

func Test_SetNewsCtx(t *testing.T) {
	t.Parallel()

	newsRedisStorage := SetupNewsRedis()

	t.Run("SetNewsCtx", func(t *testing.T) {
		key := uuid.New().String()
		newsId := uuid.New()
		n := &entity.NewsBase{
			NewsID: newsId,
			Title: "title",
		}
		err := newsRedisStorage.SetNewsCtx(context.Background(), key, 10, n)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}

func Test_GetNewsByIDCtx(t *testing.T) {
	t.Parallel()

	newsRedisStorage := SetupNewsRedis()

	t.Run("GetNewsByIDCtx", func(t *testing.T) {
		key := uuid.New().String()
		newsId := uuid.New()
		n := &entity.NewsBase{
			NewsID: newsId,
			Title: "title",
		}

		news, err := newsRedisStorage.GetNewsByIDCtx(context.Background(), key)
		require.Nil(t, news)
		require.NotNil(t, err)

		err = newsRedisStorage.SetNewsCtx(context.Background(), key, 10, n)
		require.NoError(t, err)
		require.Nil(t, err)

		news, err = newsRedisStorage.GetNewsByIDCtx(context.Background(), key)
		require.NoError(t, err)
		require.NotNil(t, news)
		require.Nil(t, err)
	})
}

func Test_DeleteNewsCtx(t *testing.T) {
	t.Parallel()

	newsRedisStorage := SetupNewsRedis()

	t.Run("DeleteNewsCtx", func(t *testing.T) {
		key := uuid.New().String()
		err := newsRedisStorage.DeleteNewsCtx(context.Background(), key)
		require.NoError(t, err)
		require.Nil(t, err)	
	})
}