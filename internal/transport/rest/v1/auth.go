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

// Auth Service interface
type Auth interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
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
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		var user entity.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		createdUser, err := h.service.Auth.Create(ctx, &user)
		if err != nil {
			return c.JSON(httpe.ParseErrors(err).Status(), httpe.ParseErrors(err))
		}

		return c.JSON(http.StatusCreated, createdUser)
	}
}

// Update User
func (h *Handlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		var u entity.User
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpe.NewBadRequestError(err.Error()))
		}
		u.ID = uID
		
		if err := c.Bind(&u); err != nil {
			return c.JSON(http.StatusBadRequest, httpe.BadRequest)
		}

		updatedUser, err := h.service.Auth.Update(ctx, &u)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, updatedUser)
	}

}

// Fet user by id
func (h *Handlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, nil)
	}
}