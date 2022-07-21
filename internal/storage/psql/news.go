package psql

import (
	"context"
	"github.com/pkg/errors"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type NewsStorage struct {
	psql *pgxpool.Pool
}

// News storage constructor
func NewNewsStorage(psql *pgxpool.Pool) *NewsStorage {
	return &NewsStorage{psql: psql}
}

// Create news
func (s *NewsStorage) Create(ctx context.Context, news *entity.News) (*entity.News, error) {
	n := &entity.News{}
	if err := s.psql.QueryRow(ctx,
		createNews,
		&news.AuthorID,
		&news.Title,
		&news.Content,
		&news.Category,
	).Scan(n); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.Create.QueryRow")
	}

	return n, nil
}

// Update news item
func (s *NewsStorage) Update(ctx context.Context, news *entity.News) (*entity.News, error) {
	n := &entity.News{}
	if err := s.psql.QueryRow(ctx,
		updateNews,
		&news.Title,
		&news.Content,
		&news.ImageURL,
		&news.Category,
		&news.NewsID,
	).Scan(n); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.Update.QueryRow")
	}

	return n, nil
}

// Get news
func (s *NewsStorage) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*entity.NewsList, error) {
	// Start a transaction to ensure user count
	tx, err := s.psql.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.GetNews.Begin")
	}
	defer tx.Rollback(ctx)

	var totalCount int
	err = tx.QueryRow(ctx, getTotalNewsCount).Scan(&totalCount)
	if err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.GetNews.QueryRow")
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
	rows, err := tx.Query(ctx, getNews, pq.GetDifference(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.GetNews.Query")
	}
	defer rows.Close()

	for rows.Next() {
		var news entity.News
		if err := rows.Scan(&news); err != nil {
			return nil, errors.Wrap(err, "NewsStoragePsql.GetNews.Scan")
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
	if _, err := s.psql.Exec(ctx, deleteNews, newsID); err != nil {
		return errors.Wrap(err, "NewsStoragePsql.Delete.Exec")
	}
	return nil
}

// Get single news by id
func (s *NewsStorage) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.NewsBase, error) {
	news := &entity.NewsBase{}
	if err := s.psql.QueryRow(ctx, getNewsByID, newsID).Scan(news); err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.GetNewsByID.QueryRow")
	}
	return news, nil
}

// Find news by title
func (s *NewsStorage) SearchNews(ctx context.Context, pq *utils.PaginationQuery, title string) (*entity.NewsList, error) {
	tx, err := s.psql.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.SearchNews.Begin")
	}
	defer tx.Rollback(ctx)

	var totalCount int
	err = tx.QueryRow(ctx, getTitleCount).Scan(&totalCount)
	if err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.SearchNews.QueryRow")
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
	rows, err := tx.Query(ctx, findByTitle, title, pq.GetDifference(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "NewsStoragePsql.SearchNews.Query")
	}
	defer rows.Close()

	for rows.Next() {
		var news *entity.News
		if err := rows.Scan(&news); err != nil {
			return nil, errors.Wrap(err, "NewsStoragePsql.SearchNews.Scan")
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
