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
func (n *NewsStorage) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	var cn entity.News
	if err := n.psql.QueryRow(ctx,
		createNews,
		&news.AuthorID,
		&news.Title,
		&news.Content,
		&news.Category,
	).Scan(&cn); err != nil {
		return nil, err
	}

	return &cn, nil
}