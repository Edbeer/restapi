package psql

import (
	"context"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/google/uuid"
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

// Update user
func (a *AuthStorage) Update(ctx context.Context, user *entity.User) error {
	if _, err := a.psql.Exec(ctx, updateUserQuery, 
		&user.FirstName, &user.LastName, &user.Email,
		&user.Role, &user.Avatar, &user.PhoneNumber,
		&user.Address, &user.City, &user.Country,
		&user.Postcode, &user.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

// Delete user
func (a *AuthStorage) Delete(ctx context.Context, userID uuid.UUID) error {
	if _, err := a.psql.Exec(ctx, deleteUserQuery, userID); err != nil {
		return httpe.NotFound
	}

	return nil
}

// Get user by id
func (a *AuthStorage) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	var u entity.User
	if err := a.psql.QueryRow(ctx, getUserByID, userID).Scan(&u); err != nil {
		return nil, err
	}

	return &u, nil
}