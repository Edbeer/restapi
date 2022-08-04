package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	mockservice "github.com/Edbeer/restapi/internal/service/mock"
	"github.com/Edbeer/restapi/pkg/converter"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mockservice.NewMockAuth(ctrl)
	mockSessionService := mockservice.NewMockSession(ctrl)

	config := &config.Config{
		Session: config.SessionConfig {
			Expire: 10,
		},
		Logger: config.Logger {
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(config)
	authHandler := NewAuthHandler(config, mockAuthService, mockSessionService, apiLogger)

	user := &entity.User{
		FirstName: "Pavel",
		LastName: "Volkov",
		Email: "edbeermtn@gmail.com",
		Password: "12345678",
	}

	buffer, err := converter.AnyToBytesBuffer(user)
	require.NoError(t, err)
	require.NotNil(t, buffer)
	require.Nil(t, err)

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(buffer.String()))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	c := e.NewContext(request, recorder)
	ctx := utils.GetRequestCtx(c)

	handlerFunc := authHandler.Register()

	userID := uuid.New()
	userWithToken := &entity.UserWithToken{
		User: &entity.User{
			ID: userID,
		},
	}
	sess := &entity.Session{
		UserID: userID,
	}
	session := "session"

	mockAuthService.EXPECT().Register(ctx, gomock.Eq(user)).Return(userWithToken, nil)
	mockSessionService.EXPECT().CreateSession(ctx, gomock.Eq(sess), 10).Return(session, nil)

	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mockservice.NewMockAuth(ctrl)
	mockSessionService := mockservice.NewMockSession(ctrl)

	config := &config.Config{
		Session: config.SessionConfig {
			Expire: 10,
		},
		Logger: config.Logger {
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(config)
	authHandler := NewAuthHandler(config, mockAuthService, mockSessionService, apiLogger)

	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	login := &Login{
		Email:    "edbeermtn@gmail.com",
		Password: "12345678",
	}

	user := &entity.User{
		Email:    login.Email,
		Password: login.Password,
	}

	buffer, err := converter.AnyToBytesBuffer(user)
	require.NoError(t, err)
	require.NotNil(t, buffer)
	require.Nil(t, err)

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(buffer.String()))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	c := e.NewContext(request, recorder)
	ctx := utils.GetRequestCtx(c)

	handlerFunc := authHandler.Login()

	userID := uuid.New()
	userWithToken := &entity.UserWithToken{
		User: &entity.User{
			ID: userID,
		},
	}
	sess := &entity.Session{
		UserID: userID,
	}
	session := "session"

	mockAuthService.EXPECT().Login(ctx, gomock.Eq(user)).Return(userWithToken, nil)
	mockSessionService.EXPECT().CreateSession(ctx, gomock.Eq(sess), 10).Return(session, nil)

	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestHandler_Logout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mockservice.NewMockAuth(ctrl)
	mockSessionService := mockservice.NewMockSession(ctrl)

	config := &config.Config{
		Session: config.SessionConfig {
			Expire: 10,
		},
		Logger: config.Logger {
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(config)
	authHandler := NewAuthHandler(config, mockAuthService, mockSessionService, apiLogger)
	sessionKey := "session-id"
	cookieValue := "cookieValue"

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.AddCookie(&http.Cookie{Name: sessionKey, Value: cookieValue})
	recorder := httptest.NewRecorder()

	c := e.NewContext(request, recorder)
	ctx := utils.GetRequestCtx(c)

	logout := authHandler.Logout()

	cookie, err := request.Cookie(sessionKey)
	require.NoError(t, err)
	require.NotNil(t, cookie)
	require.NotEqual(t, cookie.Value, "")
	require.Equal(t, cookie.Value, cookieValue)

	mockSessionService.EXPECT().DeleteSessionByID(ctx, gomock.Eq(cookie.Value)).Return(nil)

	err = logout(c)
	require.NoError(t, err)
	require.Nil(t, err)
}