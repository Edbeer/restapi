package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/internal/storage/psql"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
)

// Comments storage interface
type CommentsStorage interface {
	Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	Update(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentBase, error)
	GetAllByNewsID(ctx context.Context, newsID uuid.UUID, pq *utils.PaginationQuery) (*entity.CommentsList, error)
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

// Update comments
func (c *CommentsService) Update(ctx context.Context, comments *entity.Comment) (*entity.Comment, error) {
	comments, err := c.commentsStorage.Update(ctx, comments)
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

// Get comment by id
func (c *CommentsService) GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentBase, error) {
	comment, err := c.commentsStorage.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// Get comments by news id
func (c *CommentsService) GetAllByNewsID(ctx context.Context, 
	newsID uuid.UUID, pq *utils.PaginationQuery) (*entity.CommentsList, error) {

	comments, err := c.commentsStorage.GetAllByNewsID(ctx, newsID, pq)
	if err != nil {
		return nil, err
	}
	return comments, nil
}