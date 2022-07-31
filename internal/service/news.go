package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	baseNewsPrefix    = "api-news:"
	cacheNewsDuration = 3600
)

// News StoragePsql interface
type NewsPsql interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error)
	Update(ctx context.Context, news *entity.News) (*entity.News, error)
	GetNews(ctx context.Context, pq *utils.PaginationQuery) (*entity.NewsList, error)
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.NewsBase, error)
	SearchNews(ctx context.Context, title string, pq *utils.PaginationQuery) (*entity.NewsList, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
}

// News StorageRedis interface
type NewsRedis interface {
	GetNewsByIDCtx(ctx context.Context, key string) (*entity.NewsBase, error)
	SetNewsCtx(ctx context.Context, key string, seconds int, news *entity.NewsBase) error
	DeleteNewsCtx(ctx context.Context, key string) error
}

//  News service
type NewsService struct {
	logger       logger.Logger
	config       *config.Config
	storagePsql  NewsPsql
	storageRedis NewsRedis
}

// News service constructor
func NewNewsService(config *config.Config, storagePsql NewsPsql, redis NewsRedis, logger logger.Logger) *NewsService {
	return &NewsService{
		config:       config,
		storagePsql:  storagePsql,
		storageRedis: redis,
		logger:       logger,
	}
}

// Create news
func (n *NewsService) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if user == nil {
		return nil, httpe.NewUnauthorizedError(errors.WithMessage(err, "NewsService.Create.GetUserFromCtx"))
	}
	news.AuthorID = user.ID

	if err = utils.ValidateStruct(ctx, news); err != nil {
		return nil, httpe.NewBadRequestError(errors.WithMessage(err, "NewsService.Create.ValidateStruct"))
	}

	news, err = n.storagePsql.Create(ctx, news)
	if err != nil {
		return nil, err
	}
	return news, nil
}

// Update news items
func (n *NewsService) Update(ctx context.Context, news *entity.News) (*entity.News, error) {
	newsByID, err := n.storagePsql.GetNewsByID(ctx, news.NewsID)
	if err != nil {
		return nil, err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String(), n.logger); err != nil {
		return nil, httpe.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "NewsService.Update.ValidateIsOwner"))
	}

	updatedNews, err := n.storagePsql.Update(ctx, news)
	if err != nil {
		return nil, err
	}
	if err := n.storageRedis.DeleteNewsCtx(ctx, news.NewsID.String()); err != nil {
		n.logger.Errorf("NewsService.Update.DeleteNewsCtx: %v", err)
	}
	return updatedNews, err
}

// Delete news by id
func (n *NewsService) Delete(ctx context.Context, newsID uuid.UUID) error {
	newsByID, err := n.storagePsql.GetNewsByID(ctx, newsID)
	if err != nil {
		return err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String(), n.logger); err != nil {
		return httpe.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "NewsService.Delete.ValidateIsOwner"))
	}

	if err := n.storagePsql.Delete(ctx, newsID); err != nil {
		return err
	}

	if err := n.storageRedis.DeleteNewsCtx(ctx, newsID.String()); err != nil {
		n.logger.Errorf("NewsService.Delete.DeleteNewsCtx: %v", err)
	}
	return nil
}

// Get news
func (n *NewsService) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*entity.NewsList, error) {
	newsList, err := n.storagePsql.GetNews(ctx, pq)
	if err != nil {
		return nil, err
	}
	return newsList, err
}

// Get single news by id
func (n *NewsService) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.NewsBase, error) {
	cachedNews, err := n.storageRedis.GetNewsByIDCtx(ctx, n.generateNewsKey(newsID.String()))
	if err != nil {
		return nil, err
	}
	if cachedNews != nil {
		return cachedNews, nil
	}

	news, err := n.storagePsql.GetNewsByID(ctx, newsID)
	if err != nil {
		return nil, err
	}

	if err := n.storageRedis.SetNewsCtx(ctx, newsID.String(), cacheNewsDuration, news); err != nil {
		n.logger.Errorf("NewsService.GetNewsByID.SetNewsCtx: %v", err)
	}

	return news, nil
}

// Find news by title
func (n *NewsService) SearchNews(ctx context.Context, pq *utils.PaginationQuery, title string) (*entity.NewsList, error) {
	news, err := n.storagePsql.SearchNews(ctx, title, pq)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (n *NewsService) generateNewsKey(newsID string) string {
	return fmt.Sprintf("%s: %s", baseNewsPrefix, newsID)
}
