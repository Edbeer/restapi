package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/Edbeer/restapi/internal/service"
)

type Handlers struct {
	service *service.Services
}

func (h *Handlers) InitHandlers(e *echo.Group) {
	v1 := e.Group("/v1")
	{
		h.initAuthHandlers(v1)
	}
}

func NewHandlers(service *service.Services) *Handlers {
	return &Handlers{service: service}
}

