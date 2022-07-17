package psql

import (
	"context"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Comments storage
type CommentsStorage struct {
	psql *pgxpool.Pool
}

// Comments storage constructor
func NewCommentsStorage(psql *pgxpool.Pool) *CommentsStorage {
	return &CommentsStorage{psql: psql}
}

// Create comments
func (s *CommentsStorage) Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error) {
	var c entity.Comment
	if err := s.psql.QueryRow(ctx,
		createComments,
		&comments.AuthorID,
		&comments.NewsID,
		&comments.Message,
	).Scan(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Update comments
func (s *CommentsStorage) Update(ctx context.Context, comments *entity.Comment) (*entity.Comment, error) {
	var c entity.Comment
	if err := s.psql.QueryRow(ctx,
		updateComment,
		&comments.Message,
		&comments.CommentID,
	).Scan(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Delete comments
func (s *CommentsStorage) Delete(ctx context.Context, commentID uuid.UUID) error {
	if _, err := s.psql.Exec(ctx, deleteComment, commentID); err != nil {
		return err
	}
	return nil
}

// Get by id comment
func (s *CommentsStorage) GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentResp, error) {
	var comment entity.CommentResp
	if err := s.psql.QueryRow(ctx, getCommentByID, commentID).Scan(&comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// Get all comments by news id
func (s *CommentsStorage) GetAllByNewsID(ctx context.Context,
	newsID uuid.UUID, pq *utils.PaginationQuery) (*entity.CommentsList, error) {

	tx, err := s.psql.Begin(ctx)
	if err != nil {
		return nil, err
	}

	var totalCount int
	if err := tx.QueryRow(ctx, getCommentsCount, newsID).Scan(&totalCount); err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &entity.CommentsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetPage()),
			Comments:   make([]*entity.CommentResp, 0),
		}, nil
	}

	rows, err := tx.Query(ctx, getCommentsByNewsID, newsID, pq.GetDifference(), pq.GetLimit())
	if err != nil {
		return nil, err
	}
	var commentsList = make([]*entity.CommentResp, 0, pq.GetSize())
	for rows.Next() {
		comment := &entity.CommentResp{}
		if err := rows.Scan(&comment); err != nil {
			return nil, err
		}
		commentsList = append(commentsList, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &entity.CommentsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetPage()),
		Comments:   commentsList,
	}, nil
}
