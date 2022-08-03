package service

import (
	"context"
	"net/http"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Comments storage interface
type CommentsPsql interface {
	Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	Update(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentBase, error)
	GetAllByNewsID(ctx context.Context, newsID uuid.UUID, pq *utils.PaginationQuery) (*entity.CommentsList, error)
	Delete(ctx context.Context, commentID uuid.UUID) error
}

// Comments service
type CommentsService struct {
	logger          logger.Logger
	config          *config.Config
	commentsStorage CommentsPsql
}

// Comments service constructor
func NewCommentsService(config *config.Config, commentsStorage CommentsPsql, logger logger.Logger) *CommentsService {
	return &CommentsService{config: config, commentsStorage: commentsStorage, logger: logger}
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
func (c *CommentsService) Update(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	commByID, err := c.commentsStorage.GetByID(ctx, comment.CommentID)
	if err != nil {
		return nil, err
	}

	if err = utils.ValidateIsOwner(ctx, commByID.AuthorID.String(), c.logger); err != nil {
		return nil, httpe.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "CommentService.Update.ValidateIsOwner"))
	}

	comments, err := c.commentsStorage.Update(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Delete comments
func (c *CommentsService) Delete(ctx context.Context, commentID uuid.UUID) error {
	commentByID, err := c.commentsStorage.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if err = utils.ValidateIsOwner(ctx, commentByID.AuthorID.String(), c.logger); err != nil {
		return httpe.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "CommentService.Delete.ValidateIsOwner"))
	}

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
