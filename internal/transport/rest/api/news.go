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

// News service interface
type NewsService interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error)
	Update(ctx context.Context, news *entity.News) (*entity.News, error)
	GetNews(ctx context.Context, pq *utils.PaginationQuery) (*entity.NewsList, error)
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.NewsBase, error)
	SearchNews(ctx context.Context, pq *utils.PaginationQuery, title string) (*entity.NewsList, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
}

// NewsHandler
type NewsHandler struct {
	newsService NewsService
	config      *config.Config
	logger      logger.Logger
}

// NewsHandler constructor
func NewNewsHandler(newsService NewsService, config *config.Config, logger logger.Logger) *NewsHandler {
	return &NewsHandler{
		newsService: newsService, 
		config: config,
		logger: logger,
	}
}

// Create godoc
// @Summary Create news
// @Description Create news handler
// @Tags News
// @Accept json
// @Produce json
// @Success 201 {object} entity.News
// @Router /news/create [post]
func (h *NewsHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		n := &entity.News{}
		if err := c.Bind(n); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		createdNews, err := h.newsService.Create(ctx, n)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdNews)
	}
}

// Update godoc
// @Summary Update news
// @Description Update news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {object} entity.News
// @Router /news/{id} [put]
func (h *NewsHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		n := &entity.News{}
		if err = c.Bind(n); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}
		n.NewsID = newsUUID

		updatedNews, err := h.newsService.Update(ctx, n)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedNews)
	}
}

// Delete godoc
// @Summary Delete news
// @Description Delete by id news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {string} string	"ok"
// @Router /news/{id} [delete]
func (h *NewsHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		if err := h.newsService.Delete(ctx, newsUUID); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// GetNews godoc
// @Summary Get all news
// @Description Get all news with pagination
// @Tags News
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} entity.NewsList
// @Router /news [get]
func (h *NewsHandler) GetNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		newsList, err := h.newsService.GetNews(ctx, pq)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newsList)
	}
}

// GetByID godoc
// @Summary Get by id news
// @Description Get by id news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {object} entity.News
// @Router /news/{id} [get]
func (h *NewsHandler) GetNewsByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		news, err := h.newsService.GetNewsByID(ctx, newsUUID)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, news)
	}
}

// SearchByTitle godoc
// @Summary Search by title
// @Description Search news by title
// @Tags News
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} entity.NewsList
// @Router /news/search [get]
func (h *NewsHandler) SearchNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		title := c.QueryParam("title")

		newsList, err := h.newsService.SearchNews(ctx, pq, title)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newsList)
	}
}