package v1

import (
	"github.com/labstack/echo/v4"
)

// Auth Service interface
type Auth interface {
	Create() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
}

func (h *Handlers) initAuthHandlers(g *echo.Group) {
	authGroup := g.Group("/auth")
	{
		authGroup.GET("/:user_id", h.GetUserByID())
		authGroup.POST("", h.Create())
	}
}

// Create new user
func (h *Handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(201, nil)
	}
}

// Fet user by id
func (h *Handlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, nil)
	}
}