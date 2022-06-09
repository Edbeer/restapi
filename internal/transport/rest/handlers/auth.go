package handlers

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/labstack/echo/v4"
)

// Auth Service interface
type AuthService interface {
	Create() error
}

// Auth HTTP Transport interface
type Handlers interface {
	Create() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
}

// Auth handler 
type AuthHandler struct {
	config *config.Config
	service AuthService
	logger logger.Logger
}

// Map auth routes
func MapAuthRoutes(ag *echo.Group, h Handlers, cfg *config.Config, l logger.Logger) {
	ag.GET("/:user_id", h.GetUserByID())
	ag.POST("", h.Create())
}

// Auth handlers constructor
func NewAuthHandler(config *config.Config, service AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{config: config, service: service, logger: logger}
}

// Create new user
func (h *AuthHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(201, nil)
	}
}

// Fet user by id
func (h *AuthHandler) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, nil)
	}
}