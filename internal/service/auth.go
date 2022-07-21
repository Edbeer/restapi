package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"

	"github.com/google/uuid"
)

const (
	baseAuthPrefix    = "api-auth:"
	cacheAuthDuration = 3600
)

// Auth StoragePsql interface
type AuthPsql interface {
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindUsersByName(ctx context.Context, name string, pq *utils.PaginationQuery) (*entity.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error)
	FindUserByEmail(ctx context.Context, user *entity.User) (*entity.User, error)
}

// Auth StorageRedis interface
type AuthRedis interface {
	GetByIDCtx(ctx context.Context, key string) (*entity.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *entity.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}

// Auth service
type AuthService struct {
	logger       logger.Logger
	config       *config.Config
	storagePsql  AuthPsql
	storageRedis AuthRedis
}

// Auth service constructor
func NewAuthService(config *config.Config, storagePsql AuthPsql, storageRedis AuthRedis, logger logger.Logger) *AuthService {
	return &AuthService{
		config:       config,
		storagePsql:  storagePsql,
		storageRedis: storageRedis,
		logger:       logger,
	}
}

// Register new user
func (a *AuthService) Register(ctx context.Context, user *entity.User) (*entity.UserWithToken, error) {
	existsUser, err := a.storagePsql.FindUserByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpe.NewRestErrorWithMessage(http.StatusBadRequest, httpe.ErrEmailAlreadyExists, nil)
	}

	if err := user.PrepareCreate(); err != nil {
		return nil, httpe.NewBadRequestError(errors.Wrap(err, "AuthService.Register.PrepareCreate"))
	}

	if err := utils.ValidateStruct(ctx, user); err != nil {
		return nil, err
	}

	createdUser, err := a.storagePsql.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, a.config)
	if err != nil {
		return nil, httpe.NewInternalServerError(errors.Wrap(err, "AuthService.Register.GenerateJWTToken"))
	}

	return &entity.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Update user
func (a *AuthService) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := utils.ValidateStruct(ctx, user); err != nil {
		return nil, err
	}

	if err := user.PrepareUpdate(); err != nil {
		return nil, httpe.NewBadRequestError(errors.Wrap(err, "AuthService.Update.PrepareCreate"))
	}

	updatedUser, err := a.storagePsql.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	updatedUser.SanitizePassword()

	if err = a.storageRedis.DeleteUserCtx(ctx, a.generateUserKey(user.ID.String())); err != nil {
		a.logger.Errorf("AuthService.Update.DeleteUserCtx: %v", err)
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

// Delete user
func (a *AuthService) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := a.storagePsql.Delete(ctx, userID); err != nil {
		return err
	}
	if err := a.storageRedis.DeleteUserCtx(ctx, a.generateUserKey(userID.String())); err != nil {
		a.logger.Errorf("AuthService.Delete.DeleteUserCtx: %v", err)
	}
	return nil
}

// Get user by id
func (a *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	cachedUser, err := a.storageRedis.GetByIDCtx(ctx, a.generateUserKey(userID.String()))
	if err != nil {
		a.logger.Errorf("AuthService.GetUserByID.GetByIDCtx")
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := a.storagePsql.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := a.storageRedis.SetUserCtx(ctx, a.generateUserKey(userID.String()), cacheAuthDuration, user); err != nil {
		a.logger.Errorf("AuthService.GetByID.SetUserCtx: %v", err)
	}
	user.SanitizePassword()
	return user, nil
}

// Find users by name
func (a *AuthService) FindUsersByName(ctx context.Context, name string,
	pq *utils.PaginationQuery) (*entity.UsersList, error) {
	users, err := a.storagePsql.FindUsersByName(ctx, name, pq)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Get users
func (a *AuthService) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error) {
	users, err := a.storagePsql.GetUsers(ctx, pq)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Login user, returns user model with jwt token
func (a *AuthService) Login(ctx context.Context, user *entity.User) (*entity.UserWithToken, error) {
	foundUser, err := a.storagePsql.FindUserByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := foundUser.ComparePassword(user.Password); err != nil {
		return nil, httpe.NewUnauthorizedError(errors.Wrap(err, "AuthService.Login.ComparePassword"))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, a.config)
	if err != nil {
		return nil, httpe.NewInternalServerError(errors.Wrap(err, "AuthService.Login.GenerateJWTToken"))
	}

	return &entity.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (a *AuthService) generateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", baseAuthPrefix, userID)
}
