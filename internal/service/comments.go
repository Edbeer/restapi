package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/internal/storage/psql"
	"github.com/google/uuid"
)

// Comments storage interface
type CommentsStorage interface {
	Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	Delete(ctx context.Context, commentID uuid.UUID) error
}

// Comments service
type CommentsService struct {
	config          *config.Config
	commentsStorage *psql.CommentsStorage
}

// Comments service constructor
func NewCommentsService(config *config.Config, commentsStorage *psql.CommentsStorage) *CommentsService {
	return &CommentsService{config: config, commentsStorage: commentsStorage}
}

// Create comments
func (c *CommentsService) Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error) {
	comments, err := c.commentsStorage.Create(ctx, comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Delete comments
func (c *CommentsService) Delete(ctx context.Context, commentID uuid.UUID) error {
	if err := c.commentsStorage.Delete(ctx, commentID); err != nil {
		return err
	}
	return nil
}