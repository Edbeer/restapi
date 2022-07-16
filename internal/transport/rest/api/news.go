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
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*entity.News, error)
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

// Create news
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

// Update news
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

// Delete news by id
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

// Get news
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

// Get single news by id
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

// Find news by title
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