package service

import (
	"context"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/internal/storage/psql"
)

// Comments storage interface
type CommentsStorage interface {
	CreateComments(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
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

func (c *CommentsService) CreateComments(ctx context.Context, comments *entity.Comment) (*entity.Comment, error) {
	comments, err := c.commentsStorage.CreateComments(ctx, comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
