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
	GetAllByNewsID(ctx context.Context, newsID uuid.UUID, pq *utils.PaginationQuery) (*entity.CommentsList, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (*entity.CommentBase, error)
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

// Create godoc
// @Summary Create new comment
// @Description create new comment
// @Tags Comments
// @Accept  json
// @Produce  json
// @Success 201 {object} entity.Comment
// @Failure 500 {object} httpe.RestErr
// @Router /comments [post]
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

// Update godoc
// @Summary Update comment
// @Description update new comment
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {object} entity.Comment
// @Failure 500 {object} httpe.RestErr
// @Router /comments/{id} [put]
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

// Delete godoc
// @Summary Delete comment
// @Description delete comment
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpe.RestErr
// @Router /comments/{id} [delete]
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

// GetByID godoc
// @Summary Get comment
// @Description Get comment by id
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {object} entity.Comment
// @Failure 500 {object} httpe.RestErr
// @Router /comments/{id} [get]
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

// GetAllByNewsID godoc
// @Summary Get comments by news
// @Description Get all comment by news id
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param id path int true "news_id"
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} entity.CommentsList
// @Failure 500 {object} httpe.RestErr
// @Router /comments/byNewsId/{id} [get]
func (h *CommentsHandler) GetAllByNewsID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		newsID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		comentsList, err := h.commentsService.GetAllByNewsID(ctx, newsID, pq)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, comentsList)
	}
}