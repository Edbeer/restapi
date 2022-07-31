package psql

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
)

type NewsStorage struct {
	psql PgxClient
}

// News storage constructor
func NewNewsStorage(psql PgxClient) *NewsStorage {
	return &NewsStorage{psql: psql}
}

// Create news
func (s *NewsStorage) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	n := &entity.News{}
	if err := s.psql.QueryRowxContext(ctx,
		createNews,
		&news.AuthorID,
		&news.Title,
		&news.Content,
		&news.Category,
	).StructScan(n); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.Create.StructScan")
	}

	return n, nil
}

// Update news item
func (s *NewsStorage) Update(ctx context.Context, news *entity.News) (*entity.News, error) {
	n := &entity.News{}
	if err := s.psql.QueryRowxContext(
		ctx,
		updateNews,
		&news.Title,
		&news.Content,
		&news.ImageURL,
		&news.Category,
		&news.NewsID,
	).StructScan(n); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.Update.StructScan")
	}

	return n, nil
}

// Get news
func (s *NewsStorage) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*entity.NewsList, error) {

	var totalCount int
	if err := s.psql.GetContext(ctx, &totalCount, getTotalNewsCount); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.GetNews.GetContext")
	}

	if totalCount == 0 {
		return &entity.NewsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			News:       make([]*entity.News, 0),
		}, nil
	}

	var newsList = make([]*entity.News, 0, pq.GetSize())
	rows, err := s.psql.QueryxContext(ctx, getNews, pq.GetDifference(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.GetNews.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		var news entity.News
		if err := rows.StructScan(&news); err != nil {
			return nil, errors.Wrap(err, "NewsStoragePsql.GetNews.StructScan")
		}
		newsList = append(newsList, &news)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.GetNews.rows.Err")
	}

	return &entity.NewsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		News:       newsList,
	}, nil
}

// Delete news
func (s *NewsStorage) Delete(ctx context.Context, newsID uuid.UUID) error {
	result, err := s.psql.ExecContext(ctx, deleteNews, newsID)
	if err != nil {
		return errors.Wrap(err, "NewsStoragePsql.Delete.Exec")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "NewsStoragePsql.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "NewsStoragePsql.Delete.rowsAffected")
	}

	return nil
}

// Get single news by id
func (s *NewsStorage) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.NewsBase, error) {
	news := &entity.NewsBase{}
	if err := s.psql.GetContext(ctx, news, getNewsByID, newsID); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.GetNewsByID.GetContext")
	}
	return news, nil
}

// Find news by title
func (s *NewsStorage) SearchNews(ctx context.Context, title string, pq *utils.PaginationQuery) (*entity.NewsList, error) {

	var totalCount int
	if err := s.psql.GetContext(ctx, &totalCount, getTitleCount, title); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.SearchNews.GetContext")
	}

	if totalCount == 0 {
		return &entity.NewsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			News:       make([]*entity.News, 0),
		}, nil
	}

	var newsList = make([]*entity.News, 0, pq.GetSize())
	rows, err := s.psql.QueryxContext(ctx, findByTitle, title, pq.GetDifference(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.SearchNews.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		var news *entity.News
		if err := rows.StructScan(&news); err != nil {
			return nil, errors.Wrap(err, "NewsStoragePsql.SearchNews.StructScan")
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.SearchNews.rows.Err")
	}

	return &entity.NewsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		News:       newsList,
	}, nil
}
