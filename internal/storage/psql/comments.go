package psql

import (
	"context"
	"database/sql"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Comments storage
type CommentsStorage struct {
	psql PgxClient
}

// Comments storage constructor
func NewCommentsStorage(psql PgxClient) *CommentsStorage {
	return &CommentsStorage{psql: psql}
}

// Create comments
func (s *CommentsStorage) Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error) {
	c := &entity.Comment{}
	if err := s.psql.QueryRowxContext(ctx,
		createComments,
		&comments.AuthorID,
		&comments.NewsID,
		&comments.Message,
	).StructScan(c); err != nil {
		return nil, errors.Wrap(err, "CommentsStoragePsql.Create.StructScan")
	}
	return c, nil
}

// Update comments
func (s *CommentsStorage) Update(ctx context.Context, comments *entity.Comment) (*entity.Comment, error) {
	c := &entity.Comment{}
	if err := s.psql.QueryRowxContext(
		ctx,
		updateComment,
		&comments.CommentID,
		&comments.Message,
	).StructScan(c); err != nil {
		return nil, errors.Wrap(err, "CommentsStoragePsql.Update.StructScan")
	}
	return c, nil
}

// Delete comments
func (s *CommentsStorage) Delete(ctx context.Context, commentID uuid.UUID) error {
	result, err := s.psql.ExecContext(ctx, deleteComment, commentID)
	if err != nil {
		return errors.Wrap(err, "CommentsStoragePsql.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "CommentsStoragePsql.Delete.RowsAffected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "CommentsStoragePsql.Delete.rowsAffected")
	}

	return nil
}

// Get by id comment
func (s *CommentsStorage) GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentBase, error) {
	comment := &entity.CommentBase{}
	if err := s.psql.GetContext(ctx, comment, getCommentByID, commentID); err != nil {
		return nil, errors.Wrap(err, "CommentsStoragePsql.GetByID.GetContext")
	}
	return comment, nil
}

// Get all comments by news id
func (s *CommentsStorage) GetAllByNewsID(ctx context.Context,
	newsID uuid.UUID, pq *utils.PaginationQuery) (*entity.CommentsList, error) {

	var totalCount int
	if err := s.psql.QueryRowxContext(ctx, getCommentsCount, newsID).Scan(&totalCount); err != nil {
		return nil, errors.Wrap(err, "CommentsStoragePsql.GetAllByNewsID.QueryRowxContext")
	}

	if totalCount == 0 {
		return &entity.CommentsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Comments:   make([]*entity.CommentBase, 0),
		}, nil
	}

	rows, err := s.psql.QueryxContext(ctx, getCommentsByNewsID, newsID, pq.GetDifference(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "CommentsStoragePsql.GetAllByNewsID.QueryxContext")
	}
	defer rows.Close()

	var commentsList = make([]*entity.CommentBase, 0, pq.GetSize())
	for rows.Next() {
		comment := &entity.CommentBase{}
		if err := rows.StructScan(comment); err != nil {
			return nil, errors.Wrap(err, "CommentsStoragePsql.GetAllByNewsID.StructScan")
		}
		commentsList = append(commentsList, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "CommentsStoragePsql.GetAllByNewsID.rows.Err")
	}

	return &entity.CommentsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Comments:   commentsList,
	}, nil
}
