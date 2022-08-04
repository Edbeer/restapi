package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestHandlers_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsService := mockservice.NewMockNews(ctrl)
	newsHandlers := NewNewsHandler(mockNewsService, nil, apiLogger)

	handlerFunc := newsHandlers.Create()

	userID := uuid.New()

	news := &entity.News{
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content some text content",
	}

	buffer, err := converter.AnyToBytesBuffer(news)
	require.NoError(t, err)
	require.NotNil(t, buffer)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/news/create", strings.NewReader(buffer.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	u := &entity.User{
		ID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	req = req.WithContext(ctxWithValue)
	e := echo.New()
	ctx := e.NewContext(req, res)
	ctxWithReqID := utils.GetRequestCtx(ctx)

	mockNews := &entity.News{
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content asdasdsadsadadsad",
	}

	mockNewsService.EXPECT().Create(ctxWithReqID, gomock.Any()).Return(mockNews, nil)

	err = handlerFunc(ctx)
	require.NoError(t, err)
}

func TestNewsHandlers_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsService := mockservice.NewMockNews(ctrl)
	newsHandlers := NewNewsHandler(mockNewsService, nil, apiLogger)

	handlerFunc := newsHandlers.Update()

	userID := uuid.New()

	news := &entity.News{
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content asdasdsadsadadsad",
	}

	buffer, err := converter.AnyToBytesBuffer(news)
	require.NoError(t, err)
	require.NotNil(t, buffer)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/news/f8a3cc26-fbe1-4713-98be-a2927201356e", strings.NewReader(buffer.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	u := &entity.User{
		ID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	req = req.WithContext(ctxWithValue)
	e := echo.New()
	ctx := e.NewContext(req, res)
	ctx.SetParamNames("news_id")
	ctx.SetParamValues("f8a3cc26-fbe1-4713-98be-a2927201356e")
	ctxWithReqID := utils.GetRequestCtx(ctx)

	mockNews := &entity.News{
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content asdasdsadsadadsad",
	}

	mockNewsService.EXPECT().Update(ctxWithReqID, gomock.Any()).Return(mockNews, nil)

	err = handlerFunc(ctx)
	require.NoError(t, err)
}

func TestHandlers_GetNewsByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsService := mockservice.NewMockNews(ctrl)
	newsHandlers := NewNewsHandler(mockNewsService, nil, apiLogger)

	handlerFunc := newsHandlers.GetNewsByID()

	userID := uuid.New()
	newsID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/news/f8a3cc26-fbe1-4713-98be-a2927201356e", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	u := &entity.User{
		ID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	req = req.WithContext(ctxWithValue)
	e := echo.New()
	ctx := e.NewContext(req, res)
	ctx.SetParamNames("news_id")
	ctx.SetParamValues(newsID.String())
	ctxWithReqID := utils.GetRequestCtx(ctx)

	mockNews := &entity.NewsBase{
		NewsID:   newsID,
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content asdasdsadsadadsad",
	}

	mockNewsService.EXPECT().GetNewsByID(ctxWithReqID, newsID).Return(mockNews, nil)

	err := handlerFunc(ctx)
	require.NoError(t, err)
}

