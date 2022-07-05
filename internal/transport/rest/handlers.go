package rest

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/service"
	"github.com/Edbeer/restapi/internal/transport/rest/v1"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handlers struct {
	service *service.Services
	config  *config.Config
	logger logger.Logger
}

func NewHandlers(service *service.Services, config *config.Config, logger logger.Logger) *Handlers {
	return &Handlers{service: service, config: config, logger: logger}
}

func (h *Handlers) Init(e *echo.Echo) error {
	e.Pre(middleware.HTTPSRedirect())
	e.Use(middleware.RequestID())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	// Request ID middleware generates a unique id for a request.
	// echo.Use(middleware.CSRF())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	h.initApi(e)

	return nil
}

func (h *Handlers) initApi(e *echo.Echo) {
	handlerV1 := v1.NewHandlers(h.service, h.config, h.logger)
	api := e.Group("/api")
	{
		handlerV1.InitHandlers(api)
	}
}
