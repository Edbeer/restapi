package service

import (
	"context"
	"testing"

	"github.com/Edbeer/restapi/internal/entity"
	mockredis "github.com/Edbeer/restapi/internal/storage/redis/mock"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateSession(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionredis(ctrl)
	sessionService := NewSessionService(nil, mockSessionRedis, nil)

	ctx := context.Background()
	session := &entity.Session{}
	sid := "session id"

	mockSessionRedis.EXPECT().CreateSession(gomock.Any(), gomock.Eq(session), 10).Return(sid, nil)

	createdSession, err := sessionService.CreateSession(ctx, session, 10)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotEqual(t, createdSession, "")
}

func TestService_GetSessionByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionredis(ctrl)
	sessionService := NewSessionService(nil, mockSessionRedis, nil)

	ctx := context.Background()
	session := &entity.Session{}
	sid := "session id"

	mockSessionRedis.EXPECT().GetSessionByID(gomock.Any(), gomock.Eq(sid)).Return(session, nil)

	session, err := sessionService.GetSessionByID(ctx, sid)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, session)
}

func TestService_DeleteSessionByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionredis(ctrl)
	sessionService := NewSessionService(nil, mockSessionRedis, nil)

	ctx := context.Background()
	sid := "session id"

	mockSessionRedis.EXPECT().DeleteSessionByID(gomock.Any(), gomock.Eq(sid)).Return(nil)

	err := sessionService.DeleteSessionByID(ctx, sid)
	require.NoError(t, err)
	require.Nil(t, err)
}