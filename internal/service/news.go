package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
)

// News StoragePsql interface
type NewsPsql interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error)
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
