package psql

import (
	"github.com/jackc/pgx/v4/pgxpool"
)


type Storage struct {
	Auth *AuthStorage
	News *NewsStorage
}

func NewStorage(psql *pgxpool.Pool) *Storage {
	return &Storage{
		Auth: NewAuthStorage(psql),
		News: NewNewsStorage(psql),
	}
}