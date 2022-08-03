package service

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	mockstorage "github.com/Edbeer/restapi/internal/storage/psql/mock"
	mockredis "github.com/Edbeer/restapi/internal/storage/redis/mock"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(config)
	mockAuthStorage := mockstorage.NewMockAuthPsql(ctrl)
	authService := NewAuthService(config, mockAuthStorage, nil, apiLogger)

	user := &entity.User{
		Password: "12345678",
		Email: "edbeermtn@gmail.com",
	}

	ctx := context.Background()

	mockAuthStorage.EXPECT().FindUserByEmail(ctx, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	mockAuthStorage.EXPECT().Register(ctx, gomock.Eq(user)).Return(user, nil)

	createdUser, err := authService.Register(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Nil(t, err)
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(config)
	mockAuthStorage := mockstorage.NewMockAuthPsql(ctrl)
	mockAuthRedis := mockredis.NewMockAuthRedis(ctrl)
	authService := NewAuthService(config, mockAuthStorage, mockAuthRedis, apiLogger)

	user := &entity.User{
		Password: "12345678",
		Email: "edbeermtn@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", baseAuthPrefix, user.ID)

	ctx := context.Background()

	mockAuthStorage.EXPECT().Update(ctx, gomock.Eq(user)).Return(user, nil)
	mockAuthRedis.EXPECT().DeleteUserCtx(ctx, key).Return(nil)

	updatedUser, err := authService.Update(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Nil(t, err)
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(config)
	mockAuthStorage := mockstorage.NewMockAuthPsql(ctrl)
	mockAuthRedis := mockredis.NewMockAuthRedis(ctrl)
	authService := NewAuthService(config, mockAuthStorage, mockAuthRedis, apiLogger)

	user := &entity.User{
		Password: "12345678",
		Email: "edbeermtn@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", baseAuthPrefix, user.ID)

	ctx := context.Background()

	mockAuthStorage.EXPECT().Delete(ctx, gomock.Eq(user.ID)).Return(nil)
	mockAuthRedis.EXPECT().DeleteUserCtx(ctx, key).Return(nil)

	err := authService.Delete(ctx, user.ID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestService_GetUserByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(config)
	mockAuthStorage := mockstorage.NewMockAuthPsql(ctrl)
	mockAuthRedis := mockredis.NewMockAuthRedis(ctrl)
	authService := NewAuthService(config, mockAuthStorage, mockAuthRedis, apiLogger)

	user := &entity.User{
		Password: "12345678",
		Email: "edbeermtn@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", baseAuthPrefix, user.ID)

	ctx := context.Background()

	mockAuthRedis.EXPECT().GetByIDCtx(ctx, key).Return(nil, nil)
	mockAuthStorage.EXPECT().GetUserByID(ctx, gomock.Eq(user.ID)).Return(user, nil)
	mockAuthRedis.EXPECT().SetUserCtx(ctx, key, cacheAuthDuration, user).Return(nil)

	u, err := authService.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, u)
}

func TestService_FindUserByName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(config)
	mockAuthStorage := mockstorage.NewMockAuthPsql(ctrl)
	mockAuthRedis := mockredis.NewMockAuthRedis(ctrl)
	authService := NewAuthService(config, mockAuthStorage, mockAuthRedis, apiLogger)

	userName := "name"
	query := &utils.PaginationQuery{
		Size: 10,
		Page: 1,
		OrderBy: "",
	}

	ctx := context.Background()

	usersList := &entity.UsersList{}

	mockAuthStorage.EXPECT().FindUsersByName(ctx, gomock.Eq(userName), query).Return(usersList, nil)

	users, err := authService.FindUsersByName(ctx, userName, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, users)
}

func TestService_GetUsers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(config)
	mockAuthStorage := mockstorage.NewMockAuthPsql(ctrl)
	mockAuthRedis := mockredis.NewMockAuthRedis(ctrl)
	authService := NewAuthService(config, mockAuthStorage, mockAuthRedis, apiLogger)

	query := &utils.PaginationQuery{
		Size: 10,
		Page: 1,
		OrderBy: "",
	}

	ctx := context.Background()

	usersList := &entity.UsersList{}

	mockAuthStorage.EXPECT().GetUsers(ctx, query).Return(usersList, nil)

	users, err := authService.GetUsers(ctx, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, users)
}

func TestService_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(config)
	mockAuthStorage := mockstorage.NewMockAuthPsql(ctrl)
	mockAuthRedis := mockredis.NewMockAuthRedis(ctrl)
	authService := NewAuthService(config, mockAuthStorage, mockAuthRedis, apiLogger)

	user := &entity.User{
		Password: "12345678",
		Email: "edbeermtn@gmail.com",
	}

	ctx := context.Background()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &entity.User{
		Password: string(hashPassword),
		Email:    "edbeermtn@gmail.com",
	}

	mockAuthStorage.EXPECT().FindUserByEmail(ctx, gomock.Eq(user)).Return(mockUser, nil)

	userWithToken, err := authService.Login(ctx, user)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userWithToken)
}

