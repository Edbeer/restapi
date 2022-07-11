package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/google/uuid"
)

// News StoragePsql interface
type NewsPsql interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error)
	Update(ctx context.Context, news *entity.News) (*entity.News, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
}

//  News service
type NewsService struct {
	config       *config.Config
	storagePsql  NewsPsql
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