package psql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func Test_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commentsStorage := NewCommentsStorage(sqlxDB)

	t.Run("Create comment", func(t *testing.T) {
		authorId := uuid.New()
		newsId := uuid.New()

		columns := []string{
			"author_id",
			"news_id",
			"message",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			authorId,
			newsId,
			"hello",
		)

		comment := &entity.Comment{
			AuthorID: authorId,
			NewsID:   newsId,
			Message:  "hello",
		}

		mock.ExpectQuery(createComments).WithArgs(
			&comment.AuthorID, &comment.NewsID, &comment.Message,
		).WillReturnRows(rows)

		createdComment, err := commentsStorage.Create(context.Background(), comment)
		require.NoError(t, err)
		require.NotNil(t, createdComment)
		require.Equal(t, createdComment, comment)
	})

	t.Run("Create err", func(t *testing.T) {
		newsUID := uuid.New()
		message := "message"
		createErr := errors.New("Create comment error")

		comment := &entity.Comment{
			NewsID:  newsUID,
			Message: message,
		}

		mock.ExpectQuery(createComments).WithArgs(&comment.AuthorID, &comment.NewsID, &comment.Message).WillReturnError(createErr)

		createdComment, err := commentsStorage.Create(context.Background(), comment)
		require.NotNil(t, err)
		require.Nil(t, createdComment)
	})
}

func Test_UpdateComment(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commentsStorage := NewCommentsStorage(sqlxDB)

	t.Run("Update comment", func(t *testing.T) {
		commentId := uuid.New()
		newsID := uuid.New()

		columns := []string{
			"news_id",
			"comment_id",
			"message",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			newsID,
			commentId,
			"hello",
		)

		comment := &entity.Comment{
			CommentID: commentId,
			Message:   "hello",
		}

		mock.ExpectQuery(updateComment).WithArgs(&comment.CommentID, &comment.Message).WillReturnRows(rows)

		updatedComment, err := commentsStorage.Update(context.Background(), comment)
		require.NoError(t, err)
		require.NotNil(t, updatedComment)
		require.Equal(t, updatedComment.Message, comment.Message)
	})

	t.Run("Update err", func(t *testing.T) {
		commentId := uuid.New()
		errUpdate := errors.New("Update comment err")

		comment := &entity.Comment{
			CommentID: commentId,
			Message:   "hello",
		}

		updatedComment, err := commentsStorage.Update(context.Background(), comment)

		mock.ExpectQuery(updateComment).WithArgs(&comment.CommentID, &comment.Message).WillReturnError(errUpdate)
		require.Error(t, err)
		require.Nil(t, updatedComment)
	})
}

func Test_DeleteComment(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commentsStorage := NewCommentsStorage(sqlxDB)

	t.Run("Delete comment", func(t *testing.T) {
		commID := uuid.New()
		mock.ExpectExec(deleteComment).WithArgs(commID).WillReturnResult(sqlmock.NewResult(1, 1))
		err := commentsStorage.Delete(context.Background(), commID)

		require.NoError(t, err)
	})

	t.Run("Delete err", func(t *testing.T) {
		commID := uuid.New()

		mock.ExpectExec(deleteComment).WithArgs(commID).WillReturnResult(sqlmock.NewResult(1, 0))

		err := commentsStorage.Delete(context.Background(), commID)
		require.NotNil(t, err)
	})
}

func Test_GetCommentByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commentsStorage := NewCommentsStorage(sqlxDB)

	t.Run("GetByID", func(t *testing.T) {
		commentId := uuid.New()

		columns := []string{
			"comment_id",
			"message",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			commentId,
			"hello",
		)

		comment := &entity.CommentBase{
			CommentID: commentId,
			Message:   "hello",
		}

		mock.ExpectQuery(getCommentByID).WithArgs(&comment.CommentID).WillReturnRows(rows)

		commentByID, err := commentsStorage.GetByID(context.Background(), commentId)
		require.NoError(t, err)
		require.NotNil(t, commentByID)
		require.Equal(t, commentByID, comment)
	})
}

func Test_GetAllByNewsID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commentsStorage := NewCommentsStorage(sqlxDB)

	t.Run("GetAllByNewsID", func(t *testing.T) {
		newsID := uuid.New()
		commentId := uuid.New()

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		columns := []string{
			"comment_id",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			commentId,
		)

		mock.ExpectQuery(getCommentsCount).WithArgs(newsID).WillReturnRows(totalCountRows)
		mock.ExpectQuery(getCommentsByNewsID).WithArgs(newsID, 0, 10).WillReturnRows(rows)

		commentsList, err := commentsStorage.GetAllByNewsID(context.Background(), newsID, &utils.PaginationQuery{
			Size:    10,
			Page:    0,
			OrderBy: "",
		})

		require.NoError(t, err)
		require.NotNil(t, commentsList)
	})
}
