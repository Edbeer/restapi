package psql

import (
	"context"
	"database/sql"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Auth Storage
type AuthStorage struct {
	psql *pgxpool.Pool
}

// Auth Storage constructor
func NewAuthStorage(psql *pgxpool.Pool) *AuthStorage {
	return &AuthStorage{psql: psql}
}

// Create user
func (a *AuthStorage) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	var u entity.User
	if err := a.psql.QueryRow(ctx, createUserQuery,
		&user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Role, &user.Avatar,
		&user.PhoneNumber, &user.Address, &user.City,
		&user.Country, &user.Postcode, &user.CreatedAt,
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
	if err := a.psql.QueryRow(ctx, getUserByID, userID).Scan(&u); err == sql.ErrNoRows {
		return nil, httpe.NotFound
	}

	return &u, nil
}

// Find users by name
func (a *AuthStorage) FindUsersByName(ctx context.Context, name string) ([]*entity.User, error) {
	rows, err := a.psql.Query(ctx, findUsersByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Get users
func (a *AuthStorage) GetUsers(ctx context.Context, pg *utils.PaginationQuery) ([]*entity.User, error) {
	rows, err := a.psql.Query(ctx, getUsers, pg.GetDifference(), pg.GetOrderBy(), pg.GetLimit())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
