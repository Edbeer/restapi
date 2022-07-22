package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Auth Service interface
type AuthService interface {
	Register(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindUsersByName(ctx context.Context, name string, pq *utils.PaginationQuery) (*entity.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error)
	Login(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
}

// AuthSessionService interface
type AuthSessionService interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
}

// AuthHandler
type AuthHandler struct {
	authService    AuthService
	sessionService AuthSessionService
	config         *config.Config
	logger         logger.Logger
}

// AuthHandler constructor
func NewAuthHandler(authService AuthService, config *config.Config, logger logger.Logger, sessionService AuthSessionService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		config:         config,
		logger:         logger,
		sessionService: sessionService,
	}
}

// Register new user
func (h *AuthHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		user := &entity.User{}
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		createdUser, err := h.authService.Register(ctx, user)
		if err != nil {
			return c.JSON(httpe.ParseErrors(err).Status(), httpe.ParseErrors(err))
		}

		session, err := h.sessionService.CreateSession(ctx, &entity.Session{
			UserID: createdUser.User.ID,
		}, h.config.Session.Expire)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		c.SetCookie(utils.ConfigureSessionCookie(h.config, session))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

// Update User
func (h *AuthHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpe.NewBadRequestError(err.Error()))
		}
		u := &entity.User{}
		u.ID = uID

		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusBadRequest, httpe.BadRequest)
		}

		updatedUser, err := h.authService.Update(ctx, u)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}

}

// Delete User
func (h *AuthHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpe.NewBadRequestError(err.Error()))
		}

		if err := h.authService.Delete(ctx, uID); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// Get user by id
func (h *AuthHandler) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpe.NewBadRequestError(err.Error()))
		}

		user, err := h.authService.GetUserByID(ctx, uID)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Find users by name
func (h *AuthHandler) FindUsersByName() echo.HandlerFunc {
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

		users, err := h.authService.FindUsersByName(ctx, c.QueryParam("name"), pq)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, users)
	}
}

func (h *AuthHandler) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		users, err := h.authService.GetUsers(ctx, pq)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, users)
	}
}

// Login
func (h *AuthHandler) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		login := &Login{}
		userWithToken, err := h.authService.Login(ctx, &entity.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		session, err := h.sessionService.CreateSession(ctx, &entity.Session{
			UserID: userWithToken.User.ID,
		}, h.config.Session.Expire)

		c.SetCookie(utils.ConfigureSessionCookie(h.config, session))

		return c.JSON(http.StatusOK, userWithToken)
	}
}

// Logout
func (h *AuthHandler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		cookie, err := c.Cookie("session-id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusInternalServerError, httpe.NewInternalServerError(err))
		}
		if err = h.sessionService.DeleteSessionByID(ctx, cookie.Value); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}
		utils.DeleteSessionCookie(c, h.config.Session.Name)

		return c.NoContent(http.StatusOK)
	}
}

// Load current user from ctx with auth middleware
func (h *AuthHandler) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*entity.User)
		if !ok {
			httpe.NewUnauthorizedError(httpe.Unauthorized)
		}

		return c.JSON(http.StatusOK, user)
	}
}
