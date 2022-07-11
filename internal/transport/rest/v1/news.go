package v1

import (
	"context"
	"net/http"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/labstack/echo/v4"
)

// News service interface
type NewsService interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error) 
}

func (h *Handlers) initNewsHandlers(g *echo.Group) {
	newsGroup := g.Group("/news")
	{
		newsGroup.POST("/create", h.Create())
	}
}

func (h *Handlers) Create() echo.HandlerFunc {
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