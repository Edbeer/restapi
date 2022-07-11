package psql

import (
	"context"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/jackc/pgx/v4/pgxpool"
)

type NewsStorage struct {
	psql *pgxpool.Pool
}

// News storage constructor
func NewNewsStorage(psql *pgxpool.Pool) *NewsStorage {
	return &NewsStorage{psql: psql}
}

// Create news
func (s *NewsStorage) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	var n entity.News
	if err := s.psql.QueryRow(ctx,
		createNews,
		&news.AuthorID,
		&news.Title,
		&news.Content,
		&news.Category,
	).Scan(&n); err != nil {
		return nil, err
	}

	return &n, nil
}

// Update news item 
func (s *NewsStorage) Update(ctx context.Context, news *entity.News) (*entity.News, error) {
	var n entity.News
	if err := s.psql.QueryRow(ctx, 
		updateNews, 
		&news.Title, 
		&news.Content, 
		&news.ImageURL, 
		&news.Category,
		&news.NewsID,
	).Scan(&n); err != nil {
		return nil, err
	}

	return &n, nil
}