package psql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func Test_CreateNews(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsStorage := NewNewsStorage(sqlxDB)

	t.Run("Create news", func(t *testing.T) {
		authorId := uuid.New()
		category := "category"

		columns := []string{
			"author_id",
			"title",
			"content",
			"category",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			authorId,
			"title",
			"content",
			&category,
		)

		news := &entity.News{
			AuthorID: authorId,
			Title:    "title",
			Content:  "content",
			Category: &category,
		}

		mock.ExpectQuery(createNews).WithArgs(&news.AuthorID, &news.Title, &news.Content, &news.Category).WillReturnRows(rows)

		createdNews, err := newsStorage.Create(context.Background(), news)
		require.NoError(t, err)
		require.NotNil(t, createdNews)
		require.Equal(t, createdNews, news)
	})
}

func Test_UpdateNews(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsStorage := NewNewsStorage(sqlxDB)

	t.Run("Update news", func(t *testing.T) {
		newsId := uuid.New()
		imageUrl := "image"
		category := "category"

		columns := []string{
			"news_id",
			"title",
			"content",
			"image_url",
			"category",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			newsId,
			"title",
			"content",
			&imageUrl,
			&category,
		)

		news := &entity.News{
			NewsID:   newsId,
			Title:    "title",
			Content:  "content",
			ImageURL: &imageUrl,
			Category: &category,
		}

		mock.ExpectQuery(updateNews).WithArgs(
			&news.Title, &news.Content, &news.ImageURL, &news.Category, &news.NewsID,
		).WillReturnRows(rows)

		updatedNews, err := newsStorage.Update(context.Background(), news)
		require.NoError(t, err)
		require.NotNil(t, updatedNews)
		require.Equal(t, updatedNews, news)
	})
}

func Test_GetNews(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsStorage := NewNewsStorage(sqlxDB)

	t.Run("GetNews", func(t *testing.T) {
		newsId := uuid.New()

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		imageUrl := "image"
		category := "category"

		columns := []string{
			"news_id",
			"title",
			"content",
			"image_url",
			"category",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			newsId,
			"title",
			"content",
			&imageUrl,
			&category,
		)

		mock.ExpectQuery(getTotalNewsCount).WillReturnRows(totalCountRows)
		mock.ExpectQuery(getNews).WithArgs(0, 10).WillReturnRows(rows)

		newsList, err := newsStorage.GetNews(context.Background(), &utils.PaginationQuery{
			Size:    10,
			Page:    0,
			OrderBy: "",
		})
		require.NoError(t, err)
		require.NotNil(t, newsList)
	})
}

func Test_DeleteNews(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsStorage := NewNewsStorage(sqlxDB)

	t.Run("Delete news", func(t *testing.T) {
		newsId := uuid.New()
		mock.ExpectExec(deleteNews).WithArgs(newsId).WillReturnResult(sqlmock.NewResult(1, 1))
		err := newsStorage.Delete(context.Background(), newsId)
		require.NoError(t, err)
	})

	t.Run("Delete err", func(t *testing.T) {
		newsId := uuid.New()
		mock.ExpectExec(deleteNews).WithArgs(newsId).WillReturnResult(sqlmock.NewResult(1, 0))
		err := newsStorage.Delete(context.Background(), newsId)
		require.NotNil(t, err)
	})
}

func Test_GetNewsByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsStorage := NewNewsStorage(sqlxDB)

	t.Run("GetNewsByID", func(t *testing.T) {
		newsId := uuid.New()

		columns := []string{
			"news_id",
			"title",
			"content",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			newsId,
			"title",
			"content",
		)

		news := &entity.NewsBase{
			NewsID:  newsId,
			Title:   "title",
			Content: "content",
		}

		mock.ExpectQuery(getNewsByID).WithArgs(newsId).WillReturnRows(rows)

		newsById, err := newsStorage.GetNewsByID(context.Background(), newsId)
		require.NoError(t, err)
		require.NotNil(t, newsById)
		require.Equal(t, newsById, news)
	})
}

func Test_SearchNews(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsStorage := NewNewsStorage(sqlxDB)

	t.Run("SearchNews", func(t *testing.T) {
		newsId := uuid.New()

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		title := "title"
		columns := []string{
			"news_id",
			title,
			"content",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			newsId,
			title,
			"content",
		)

		mock.ExpectQuery(getTitleCount).WithArgs(title).WillReturnRows(totalCountRows)
		mock.ExpectQuery(findByTitle).WithArgs(title, 0, 10).WillReturnRows(rows)

		newsByTitle, err := newsStorage.SearchNews(context.Background(), title, &utils.PaginationQuery{
			Size:    10,
			Page:    0,
			OrderBy: "",
		})
		require.NoError(t, err)
		require.NotNil(t, newsByTitle)
	})
}
