package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/internal/service"
	mockservice "github.com/Edbeer/restapi/internal/service/mock"
	"github.com/Edbeer/restapi/pkg/converter"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHandler_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommentsService := mockservice.NewMockComments(ctrl)
	commentsService := service.NewCommentsService(nil, mockCommentsService, apiLogger)

	commHandlers := NewCommentsHandler(commentsService, nil, apiLogger)
	handlerFunc := 	commHandlers.Create()

	userID := uuid.New()
	newsID := uuid.New()
	comment := &entity.Comment{
		AuthorID: userID,
		Message:  "message Key: 'Comment.Message' Error:Field validation for 'Message' failed on the 'gte' tag",
		NewsID:   newsID,
	}

	buffer, err := converter.AnyToBytesBuffer(comment)
	require.NoError(t, err)
	require.NotNil(t, buffer)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/comments", strings.NewReader(buffer.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	u := &entity.User{
		ID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	req = req.WithContext(ctxWithValue)

	e := echo.New()
	ctx := e.NewContext(req, res)

	mockComm := &entity.Comment{
		AuthorID: userID,
		NewsID:   comment.NewsID,
		Message:  "message",
	}

	fmt.Printf("COMMENT: %#v\n", comment)
	fmt.Printf("MOCK COMMENT: %#v\n", mockComm)

	mockCommentsService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockComm, nil)

	err = handlerFunc(ctx)
	require.NoError(t, err)
}


func TestHandler_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommentsService := mockservice.NewMockComments(ctrl)
	commentsService := service.NewCommentsService(nil, mockCommentsService, apiLogger)

	commHandlers := NewCommentsHandler(commentsService, nil, apiLogger)
	handlerFunc := 	commHandlers.GetByID()

	r := httptest.NewRequest(http.MethodGet, "/api/comments/5c9a9d67-ad38-499c-9858-086bfdeaf7d2", nil)
	w := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(r, w)
	c.SetParamNames("comment_id")
	c.SetParamValues("5c9a9d67-ad38-499c-9858-086bfdeaf7d2")
	
	comm := &entity.CommentBase{}

	mockCommentsService.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(comm, nil)

	err := handlerFunc(c)
	require.NoError(t, err)
}

func TestHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommentsService := mockservice.NewMockComments(ctrl)
	commentsService := service.NewCommentsService(nil, mockCommentsService, apiLogger)

	commHandlers := NewCommentsHandler(commentsService, nil, apiLogger)
	handlerFunc := commHandlers.Delete()

	userID := uuid.New()
	commID := uuid.New()
	comm := &entity.CommentBase{
		CommentID: commID,
		AuthorID:  userID,
	}

	r := httptest.NewRequest(http.MethodDelete, "/api/comments/5c9a9d67-ad38-499c-9858-086bfdeaf7d2", nil)
	w := httptest.NewRecorder()
	u := &entity.User{
		ID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	r = r.WithContext(ctxWithValue)
	e := echo.New()
	c := e.NewContext(r, w)
	c.SetParamNames("comment_id")
	c.SetParamValues(commID.String())

	mockCommentsService.EXPECT().GetByID(gomock.Any(), commID).Return(comm, nil)
	mockCommentsService.EXPECT().Delete(gomock.Any(), commID).Return(nil)

	err := handlerFunc(c)
	require.NoError(t, err)
}
