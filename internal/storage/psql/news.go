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

// create news
func (n *NewsStorage) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	return &entity.News{}, nil
}