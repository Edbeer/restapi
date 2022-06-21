package v1

import (
	"github.com/Edbeer/restapi/internal/service"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	logger *logger.Logger
	service *service.Services
}

func (h *Handlers) InitHandlers(e *echo.Group) {
	v1 := e.Group("/v1")
	{
		h.initAuthHandlers(v1)
	}
}

func NewHandlers(logger *logger.Logger, service *service.Services) *Handlers {
	return &Handlers{logger: logger, service: service}
}

