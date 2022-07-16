package psql

import (
	"context"

	"github.com/Edbeer/restapi/internal/entity"
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

// Delete comments
func (s *CommentsStorage) Delete(ctx context.Context, commentID uuid.UUID) error {
	if _, err := s.psql.Exec(ctx, deleteComment, commentID); err != nil {
		return err
	}
	return nil
}