package v1

import (
	"context"
	"net/http"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// News service interface
type NewsService interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error) 
	Update(ctx context.Context, news *entity.News) (*entity.News, error)
}

func (h *Handlers) initNewsHandlers(g *echo.Group) {
	newsGroup := g.Group("/news")
	{
		newsGroup.POST("/create", h.CreateNews())
		newsGroup.PUT("/:news_id", h.UpdateNews())
		newsGroup.DELETE("/:news_id", h.DeleteNews())
	}
}

// Create news
func (h *Handlers) CreateNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		n := &entity.News{}
		if err := c.Bind(n); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		createdNews, err := h.service.News.Create(ctx, n)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdNews) 
	}
}

// Update news 
func (h *Handlers) UpdateNews() echo.HandlerFunc {
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

		updatedNews, err := h.service.News.Update(ctx, n)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedNews)
	}
}

// Delete news by id
func (h *Handlers) DeleteNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		if err := h.service.News.Delete(ctx, newsUUID); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}