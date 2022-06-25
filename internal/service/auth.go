package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/utils"

	"github.com/google/uuid"
)

// Auth StoragePsql interface
type AuthPsql interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindUsersByName(ctx context.Context, name string) ([]*entity.User, error)
	GetUsers(ctx context.Context, pg *utils.PaginationQuery) (*entity.UsersList, error)
}

// Auth StorageRedis interface
type AuthRedis interface {
	Create() error
}

// Auth service
type AuthService struct {
	config  *config.Config
	storagePsql AuthPsql
	storageRedis AuthRedis
}

// Auth service constructor
func NewAuthService(config *config.Config, storagePsql AuthPsql, storageRedis AuthRedis) *AuthService {
	return &AuthService{config: config,storagePsql: storagePsql, storageRedis: storageRedis}
}

// Create new user
func (a *AuthService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := user.PrepareCreate(); err != nil {
		return nil, httpe.NewBadRequestError(err.Error())
	}

	if err := utils.ValidateStruct(ctx, user); err != nil {
		return nil, err
	}

	createdUser, err := a.storagePsql.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()
	return createdUser, nil
}

// Update user
func (a *AuthService) Update(ctx context.Context, user *entity.User) error {
	if err := utils.ValidateStruct(ctx, user); err != nil {
		return err
	}

	if err := user.PrepareUpdate(); err != nil {
		return err
	}

	err := a.storagePsql.Update(ctx, user)
	if err != nil {
		return err
	}
	user.SanitizePassword()
	return nil
}

// Delete user
func (a *AuthService) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := a.storagePsql.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}

// Get user by id
func (a *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := a.storagePsql.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.SanitizePassword()
	return user, nil
}

// Find users by name
func (a *AuthService) FindUsersByName(ctx context.Context, name string) ([]*entity.User, error) {
	users, err := a.storagePsql.FindUsersByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Get users
func (a *AuthService) GetUsers(ctx context.Context, pg *utils.PaginationQuery) (*entity.UsersList, error) {
	users, err := a.storagePsql.GetUsers(ctx, pg)
	if err != nil {
		return nil, err
	}

	return users, nil
}

