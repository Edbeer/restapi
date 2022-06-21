package psql

import (
	"context"

	"github.com/Edbeer/restapi/internal/entity"
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
func (a *AuthStorage) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	var u entity.User
	if err := a.psql.QueryRow(ctx, createUserQuery, 
		&user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Role, &user.Avatar, 
		&user.PhoneNumber, &user.Address, &user.City, 
		&user.Country,	&user.Postcode, &user.CreatedAt, 
		&user.UpdatedAt,
	).Scan(&u); err != nil {
		return nil, err
	}

	return &u, nil
}
