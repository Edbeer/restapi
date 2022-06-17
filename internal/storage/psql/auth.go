package psql

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

// Auth Storage
type AuthStorage struct {
	psql  *pgxpool.Pool
}

// Auth Storage constructor
func NewAuthStorage(psql *pgxpool.Pool,) *AuthStorage {
	return &AuthStorage{psql: psql}
}

// Create user
func (s *AuthStorage) Create() error {
	return nil
}
