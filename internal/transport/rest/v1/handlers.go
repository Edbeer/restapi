package v1

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/service"
	"github.com/Edbeer/restapi/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	service *service.Services
	config *config.Config
	logger logger.Logger
}

func (h *Handlers) InitHandlers(e *echo.Group) {
	v1 := e.Group("/v1")
	{
		h.initAuthHandlers(v1)
		h.initNewsHandlers(v1)
	}
}

func NewHandlers(service *service.Services, config *config.Config, logger logger.Logger) *Handlers {
	return &Handlers{service: service, config: config, logger: logger}
}

