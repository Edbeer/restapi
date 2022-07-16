package api

import (
	"context"
	"net/http"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/labstack/echo/v4"
)

// Comments Service interface
type CommentsService interface {
	CreateComments(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
}

// Comments Handler
type CommentsHandler struct {
	commentsService CommentsService
	config          *config.Config
	logger          logger.Logger
}

// Comments Handler constructor
func NewCommentsHandler(commentsService CommentsService, config *config.Config, logger logger.Logger) *CommentsHandler {
	return &CommentsHandler{
		commentsService: commentsService,
		config:          config,
		logger:          logger,
	}
}

// Create Comments
func (h *CommentsHandler) CreateComments() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		comment := &entity.Comment{}
		comment.AuthorID = user.ID

		comments, err := h.commentsService.CreateComments(ctx, comment)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusCreated, comments)
	}
}
