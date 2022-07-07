package v1

import (
	"context"
	"net/http"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/labstack/echo/v4"
)

// News service interface
type NewsService interface {
	Create(ctx context.Context, news *entity.News) (*entity.News, error) 
}

func (h *Handlers) initNewsHandlers(g *echo.Group) {
	newsGroup := g.Group("/news")
	{
		newsGroup.POST("/Create", h.Create())
	}
}

func (h *Handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, 200)
	}
}