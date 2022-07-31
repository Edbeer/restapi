package psql

import (
	"context"
	"database/sql"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Auth Storage
type AuthStorage struct {
	psql PgxClient
}

// Auth Storage constructor
func NewAuthStorage(psql PgxClient) *AuthStorage {
	return &AuthStorage{psql: psql}
}

// Register user
func (a *AuthStorage) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	u := &entity.User{}
	if err := a.psql.QueryRowxContext(ctx, createUserQuery,
		&user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Role, &user.Avatar,
		&user.PhoneNumber, &user.Address, &user.City,
		&user.Country, &user.Postcode,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.Register.StructScan")
	}
	return u, nil
}

// Update user
func (a *AuthStorage) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	u := &entity.User{}
	if err := a.psql.GetContext(ctx, u, updateUserQuery,
		&user.FirstName, &user.LastName, &user.Email,
		&user.Role, &user.Avatar, &user.PhoneNumber,
		&user.Address, &user.City, &user.Country,
		&user.Postcode, &user.ID,
	); err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.Update.GetContext")
	}

	return u, nil
}

// Delete user
func (a *AuthStorage) Delete(ctx context.Context, userID uuid.UUID) error {

	result, err := a.psql.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return errors.Wrap(err, "AuthStoragePsql.Delete.Context")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "AuthStoragePsql.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "AuthStoragePsql.Delete.rowsAffected")
	}

	return nil
}

// Get user by id
func (a *AuthStorage) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	u := &entity.User{}
	if err := a.psql.QueryRowxContext(ctx, getUserByID, userID).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.GetUserByID.StructScan")
	}

	return u, nil
}

// Find users by name
func (a *AuthStorage) FindUsersByName(ctx context.Context,
	name string, pq *utils.PaginationQuery) (*entity.UsersList, error) {

	var totalCount int
	if err := a.psql.GetContext(ctx, &totalCount, getTotalCount, name); err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.FindUsersByName.GetContext")
	}

	if totalCount == 0 {
		return &entity.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Users:      make([]*entity.User, 0),
		}, nil
	}

	rows, err := a.psql.QueryxContext(ctx, findUsersByName, name, pq.GetDifference(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.FindUsersByName.QueryxContext")
	}
	defer rows.Close()

	users := make([]*entity.User, 0, pq.GetSize())
	for rows.Next() {
		var user entity.User
		if err := rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, "AuthStoragePsql.FindUsersByName.StructScan")
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.FindUsersByName.rows.Err")
	}

	return &entity.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      users,
	}, nil
}

// Get users
func (a *AuthStorage) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error) {

	var totalCount int
	if err := a.psql.GetContext(ctx, &totalCount, getTotal); err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.GetUsers.GetContext")
	}

	if totalCount == 0 {
		return &entity.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Users:      make([]*entity.User, 0),
		}, nil
	}

	users := make([]*entity.User, 0, pq.GetSize())
	if err := a.psql.SelectContext(
		ctx,
		&users,
		getUsers,
		pq.GetOrderBy(),
		pq.GetDifference(),
		pq.GetLimit(),
	); err != nil {
		return nil, errors.Wrap(err, "AuthStorage.GetUsers.SelectContext")
	}

	return &entity.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      users,
	}, nil
}

// Find user by email
func (a *AuthStorage) FindUserByEmail(ctx context.Context, user *entity.User) (*entity.User, error) {
	foundUser := &entity.User{}
	if err := a.psql.QueryRowxContext(ctx, findUserByEmail, user.Email).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.FindUserByEmail.StructScan")
	}
	return foundUser, nil
}
