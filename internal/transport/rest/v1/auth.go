package v1

import (
	"context"
	"net/http"

	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/internal/middleware"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Auth Service interface
type Auth interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindUsersByName(ctx context.Context, name string, pq *utils.PaginationQuery) (*entity.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error)
}

func (h *Handlers) initAuthHandlers(g *echo.Group) {
	authGroup := g.Group("/auth")
	{
		authGroup.POST("/create", h.Create())
		authGroup.GET("/:user_id", h.GetUserByID())
		authGroup.GET("/find", h.FindUsersByName())
		authGroup.GET("/all", h.GetUsers())
		authGroup.Use(middleware.AuthJWTMiddleware(h.service.Auth, h.config))
		authGroup.PUT("/:user_id", h.Update())
		authGroup.DELETE("/:user_id", h.Delete())
	}
}

// Create new user
func (h *Handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		var user entity.User
		// Method POST
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
		
		// Method PUT
		if err := c.Bind(&u); err != nil {
			return c.JSON(http.StatusBadRequest, httpe.BadRequest)
		}

		err = h.service.Auth.Update(ctx, &u)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, "user updated")
	}

}

// Delete User
func (h *Handlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpe.NewBadRequestError(err.Error()))
		}

		if err := h.service.Auth.Delete(ctx, uID); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// Get user by id
func (h *Handlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpe.NewBadRequestError(err.Error()))
		}

		user, err := h.service.Auth.GetUserByID(ctx, uID)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Find users by name
func (h *Handlers) FindUsersByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		if c.QueryParam("name") == "" {
			return c.JSON(http.StatusBadRequest, httpe.NewBadRequestError("name query param is required"))
		}

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		users, err := h.service.Auth.FindUsersByName(ctx, c.QueryParam("name"), pq);
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, users)
	}
}

func (h *Handlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		users, err := h.service.Auth.GetUsers(ctx, pq)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, users)
	}
}