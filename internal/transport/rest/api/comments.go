package api

import (
	"context"
	"net/http"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Comments Service interface
type CommentsService interface {
	Create(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	Update(ctx context.Context, comments *entity.Comment) (*entity.Comment, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentResp, error)
	Delete(ctx context.Context, commentID uuid.UUID) error
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
func (h *CommentsHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		comment := &entity.Comment{}
		comment.AuthorID = user.ID

		comments, err := h.commentsService.Create(ctx, comment)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusCreated, comments)
	}
}

// Update comments
func (h *CommentsHandler) Update() echo.HandlerFunc {
	type UpdatedComment struct {
		Message string `json:"message" db:"message" validate:"required,gte=0"`
		Likes   int64  `json:"likes" db:"likes" validate:"omitempty"`
	}
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		commentUUID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		comm := &UpdatedComment{}

		updatedComment, err := h.commentsService.Update(ctx, &entity.Comment{
			CommentID: commentUUID,
			Message: comm.Message,
			Likes: comm.Likes,
		})
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedComment)
	}
}

// Delete comments
func (h *CommentsHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		commentUUID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		if err := h.commentsService.Delete(ctx, commentUUID); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// Get comment by id
func (h *CommentsHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		commentID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		comment, err := h.commentsService.GetByID(ctx, commentID)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, comment)
	}
}