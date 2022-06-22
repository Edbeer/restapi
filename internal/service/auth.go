package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/utils"
)

// Auth StoragePsql interface
type AuthPsql interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
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

	return createdUser, nil
}

// Update user
func (a *AuthService) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := utils.ValidateStruct(ctx, user); err != nil {
		return nil, err
	}

	if err := user.PrepareUpdate(); err != nil {
		return nil, err
	}

	updatedUser, err := a.storagePsql.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}