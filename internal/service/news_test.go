package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/Edbeer/restapi/internal/entity"
	mockstorage "github.com/Edbeer/restapi/internal/storage/psql/mock"
	mockredis "github.com/Edbeer/restapi/internal/storage/redis/mock"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestService_CreateNews(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsStorage := mockstorage.NewMockNewsPsql(ctrl)
	newsService := NewNewsService(nil, mockNewsStorage, nil, apiLogger)

	userID := uuid.New()

	news := &entity.News{
		AuthorID: userID,
		Title:    "TitleTitleTitleTitleTitleTitleTitle",
		Content:  "ContentContentContentContentContent",
	}

	user := &entity.User{
		ID: userID,
	}

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)

	mockNewsStorage.EXPECT().Create(ctx, news).Return(news, nil)

	createdNews, err := newsService.Create(ctx, news)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, createdNews)
}

func TestService_UpdateNews(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsStorage := mockstorage.NewMockNewsPsql(ctrl)
	mockNewsRedis := mockredis.NewMockNewsRedis(ctrl)
	newsService := NewNewsService(nil, mockNewsStorage, mockNewsRedis, apiLogger)

	userID := uuid.New()
	newsID := uuid.New()
	news := &entity.News{
		NewsID:   newsID,
		AuthorID: userID,
		Title:    "TitleTitleTitleTitleTitleTitleTitle",
		Content:  "ContentContentContentContentContent",
	}

	newsBase := &entity.NewsBase{
		NewsID:   newsID,
		AuthorID: userID,
		Title:    "TitleTitleTitleTitleTitleTitleTitleTitle",
		Content:  "ContentContentContentContentContent",
	}

	user := &entity.User{
		ID: userID,
	}

	cacheKey := fmt.Sprintf("%s: %s", baseNewsPrefix, news.NewsID)

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)

	mockNewsStorage.EXPECT().GetNewsByID(ctx, gomock.Eq(news.NewsID)).Return(newsBase, nil)
	mockNewsStorage.EXPECT().Update(ctx, gomock.Eq(news)).Return(news, nil)
	mockNewsRedis.EXPECT().DeleteNewsCtx(ctx, gomock.Eq(cacheKey)).Return(nil)

	updatedNews, err := newsService.Update(ctx, news)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, updatedNews)
}

func TestService_GetNewsByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsStorage := mockstorage.NewMockNewsPsql(ctrl)
	mockNewsRedis := mockredis.NewMockNewsRedis(ctrl)
	newsService := NewNewsService(nil, mockNewsStorage, mockNewsRedis, apiLogger)

	newsID := uuid.New()
	newsBase := &entity.NewsBase{
		NewsID: newsID,
	}
	ctx := context.Background()

	cacheKey := fmt.Sprintf("%s: %s", baseNewsPrefix, newsID)

	mockNewsRedis.EXPECT().GetNewsByIDCtx(ctx, gomock.Eq(cacheKey)).Return(nil, nil)
	mockNewsStorage.EXPECT().GetNewsByID(ctx, gomock.Eq(newsID)).Return(newsBase, nil)
	mockNewsRedis.EXPECT().SetNewsCtx(ctx, cacheKey, cacheNewsDuration, newsBase).Return(nil)

	newsById, err := newsService.GetNewsByID(ctx, newsBase.NewsID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, newsById)
}

func TestService_DeleteNews(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsStorage := mockstorage.NewMockNewsPsql(ctrl)
	mockNewsRedis := mockredis.NewMockNewsRedis(ctrl)
	newsService := NewNewsService(nil, mockNewsStorage, mockNewsRedis, apiLogger)

	newsID := uuid.New()
	userID := uuid.New()
	newsBase := &entity.NewsBase{
		NewsID:   newsID,
		AuthorID: userID,
	}
	cacheKey := fmt.Sprintf("%s: %s", baseNewsPrefix, newsID)

	user := &entity.User{
		ID: userID,
	}

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)

	mockNewsStorage.EXPECT().GetNewsByID(ctx, gomock.Eq(newsBase.NewsID)).Return(newsBase, nil)
	mockNewsStorage.EXPECT().Delete(ctx, gomock.Eq(newsID)).Return(nil)
	mockNewsRedis.EXPECT().DeleteNewsCtx(ctx, gomock.Eq(cacheKey)).Return(nil)

	err := newsService.Delete(ctx, newsBase.NewsID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestService_GetNews(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsStorage := mockstorage.NewMockNewsPsql(ctrl)
	mockNewsRedis := mockredis.NewMockNewsRedis(ctrl)
	newsService := NewNewsService(nil, mockNewsStorage, mockNewsRedis, apiLogger)

	ctx := context.Background()

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	newsList := &entity.NewsList{}

	mockNewsStorage.EXPECT().GetNews(ctx, query).Return(newsList, nil)

	news, err := newsService.GetNews(ctx, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, news)
}

func TestService_SearchNews(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsStorage := mockstorage.NewMockNewsPsql(ctrl)
	mockNewsRedis := mockredis.NewMockNewsRedis(ctrl)
	newsService := NewNewsService(nil, mockNewsStorage, mockNewsRedis, apiLogger)

	ctx := context.Background()

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	newsList := &entity.NewsList{}
	title := "title"

	mockNewsStorage.EXPECT().SearchNews(ctx, title, query).Return(newsList, nil)

	news, err := newsService.SearchNews(ctx, title, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, news)
}