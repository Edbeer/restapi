package psql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type PgxClient interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Storage struct {
	Auth     *AuthStorage
	News     *NewsStorage
	Comments *CommentsStorage
}

func NewStorage(psql PgxClient) *Storage {
	return &Storage{
		Auth:     NewAuthStorage(psql),
		News:     NewNewsStorage(psql),
		Comments: NewCommentsStorage(psql),
	}
}
