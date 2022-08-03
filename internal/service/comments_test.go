package service

import (
	"context"
	"testing"

	"github.com/Edbeer/restapi/internal/entity"
	mockstorage "github.com/Edbeer/restapi/internal/storage/psql/mock"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommStorage := mockstorage.NewMockCommentsPsql(ctrl)
	commentsService := NewCommentsService(nil, mockCommStorage, apiLogger)

	comment := &entity.Comment{}

	ctx := context.Background()

	mockCommStorage.EXPECT().Create(ctx, gomock.Eq(comment)).Return(comment, nil)

	createdComment, err := commentsService.Create(ctx, comment)
	require.NoError(t, err)
	require.NotNil(t, createdComment)
}

func TestService_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommStorage := mockstorage.NewMockCommentsPsql(ctrl)
	commentsService := NewCommentsService(nil, mockCommStorage, apiLogger)

	comment := &entity.Comment{
		CommentID: uuid.New(),
	}

	commentBase := &entity.CommentBase{}

	ctx := context.Background()

	mockCommStorage.EXPECT().GetByID(ctx, gomock.Eq(comment.CommentID)).Return(commentBase, nil)

	cb, err := commentsService.GetByID(ctx, comment.CommentID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, cb)
}

func TestService_UpdateComment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommStorage := mockstorage.NewMockCommentsPsql(ctrl)
	commentsService := NewCommentsService(nil, mockCommStorage, apiLogger)

	authorID := uuid.New()

	comment := &entity.Comment{
		CommentID: uuid.New(),
		AuthorID:  authorID,
	}

	commentBase := &entity.CommentBase{
		AuthorID: authorID,
	}

	user := &entity.User{
		ID: authorID,
	}

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)

	mockCommStorage.EXPECT().GetByID(ctx, gomock.Eq(comment.CommentID)).Return(commentBase, nil)
	mockCommStorage.EXPECT().Update(ctx, gomock.Eq(comment)).Return(comment, nil)

	updatedComment, err := commentsService.Update(ctx, comment)
	require.NoError(t, err)
	require.NotNil(t, updatedComment)
}

func TestService_DeleteComment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommStorage := mockstorage.NewMockCommentsPsql(ctrl)
	commentsService := NewCommentsService(nil, mockCommStorage, apiLogger)

	authorID := uuid.New()

	comment := &entity.Comment{
		CommentID: uuid.New(),
		AuthorID:  authorID,
	}

	commentBase := &entity.CommentBase{
		AuthorID: authorID,
	}

	user := &entity.User{
		ID: authorID,
	}

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)

	mockCommStorage.EXPECT().GetByID(ctx, gomock.Eq(comment.CommentID)).Return(commentBase, nil)
	mockCommStorage.EXPECT().Delete(ctx, gomock.Eq(comment.CommentID)).Return(nil)

	err := commentsService.Delete(ctx, comment.CommentID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestService_GetAllByNewsID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommStorage := mockstorage.NewMockCommentsPsql(ctrl)
	commentsService := NewCommentsService(nil, mockCommStorage, apiLogger)

	newsID := uuid.New()

	comment := &entity.Comment{
		CommentID: uuid.New(),
		NewsID:    newsID,
	}

	commentsList := &entity.CommentsList{}

	ctx := context.Background()

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	mockCommStorage.EXPECT().GetAllByNewsID(ctx, gomock.Eq(comment.NewsID), query).Return(commentsList, nil)

	comments, err := commentsService.GetAllByNewsID(ctx, comment.NewsID, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, comments)
}
