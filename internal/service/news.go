package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
)

// News StoragePsql interface
type NewsPsql interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error)
	Update(ctx context.Context, news *entity.News) (*entity.News, error)
	GetNews(ctx context.Context, pq *utils.PaginationQuery) (*entity.NewsList, error)
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.News, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
}

//  News service
type NewsService struct {
	config      *config.Config
	storagePsql NewsPsql
}

// News service constructor
func NewNewsService(config *config.Config, storagePsql NewsPsql) *NewsService {
	return &NewsService{config: config, storagePsql: storagePsql}
}

// Create news
func (n *NewsService) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	news, err := n.storagePsql.Create(ctx, news)
	if err != nil {
		return nil, err
	}
	return news, nil
}

// Update news items
func (n *NewsService) Update(ctx context.Context, news *entity.News) (*entity.News, error) {
	updatedNews, err := n.storagePsql.Update(ctx, news)
	if err != nil {
		return nil, err
	}
	return updatedNews, err
}

// Delete news by id
func (n *NewsService) Delete(ctx context.Context, newsID uuid.UUID) error {
	if err := n.storagePsql.Delete(ctx, newsID); err != nil {
		return err
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
func (n *NewsService) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.News, error) {
	news, err := n.storagePsql.GetNewsByID(ctx, newsID)
	if err != nil {
		return nil, err
	}
	return news, nil
}