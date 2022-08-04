package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/csrf"
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

// Session service interface
type SessionService interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
}

// AuthHandler
type AuthHandler struct {
	config         *config.Config
	authService    AuthService
	sessionService SessionService
	logger         logger.Logger
}

// AuthHandler constructor
func NewAuthHandler(config *config.Config, authService AuthService, sessionService SessionService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		config:         config,
		authService:    authService,
		sessionService: sessionService,
		logger:         logger,
	}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} entity.User
// @Router /auth/register [post]
func (h *AuthHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		user := &entity.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
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

// Update godoc
// @Summary Update user
// @Description update existing user
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {object} entity.User
// @Router /auth/{id} [put]
func (h *AuthHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

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

// Delete godoc
// @Summary Delete user account
// @Description some description
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpe.RestError
// @Router /auth/{id} [delete]
func (h *AuthHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

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

// GetUserByID godoc
// @Summary get user by id
// @Description get string by ID
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param id path int true "user_id"
// @Success 200 {object} entity.User
// @Failure 500 {object} httpe.RestError
// @Router /auth/{id} [get]
func (h *AuthHandler) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

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

// FindUsersByName godoc
// @Summary Find by name
// @Description Find user by name
// @Tags Auth
// @Accept json
// @Param name query string false "username" Format(username)
// @Produce json
// @Success 200 {object} entity.UsersList
// @Failure 500 {object} httpe.RestError
// @Router /auth/find [get]
func (h *AuthHandler) FindUsersByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

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

// GetUsers godoc
// @Summary Get users
// @Description Get the list of all users
// @Tags Auth
// @Accept json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Produce json
// @Success 200 {object} entity.UsersList
// @Failure 500 {object} httpe.RestError
// @Router /auth/all [get]
func (h *AuthHandler) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

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

// Login godoc
// @Summary Login new user
// @Description login user, returns user and set session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
// @Router /auth/login [post]
func (h *AuthHandler) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}
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

// Logout godoc
// @Summary Logout user
// @Description logout user removing session
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

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

// GetMe godoc
// @Summary Get user by id
// @Description Get current user by id
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
// @Failure 500 {object} httpe.RestError
// @Router /auth/me [get]
func (h *AuthHandler) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*entity.User)
		if !ok {
			httpe.NewUnauthorizedError(httpe.Unauthorized)
		}

		return c.JSON(http.StatusOK, user)
	}
}

// GetCSRFToken godoc
// @Summary Get CSRF token
// @Description Get CSRF token, required auth session cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "Ok"
// @Failure 500 {object} httpe.RestError
// @Router /auth/token [get]
func (h *AuthHandler) GetCSRFToken() echo.HandlerFunc {
	return func(c echo.Context) error {

		sessionId, ok := c.Get("sid").(string)
		if !ok {
			httpe.NewUnauthorizedError(httpe.Unauthorized)
		}

		token := csrf.MakeToken(sessionId, h.logger)
		c.Response().Header().Set(csrf.CSRFHeader, token)
		c.Response().Header().Set("Access-Control-Expose-Headres", csrf.CSRFHeader)

		return c.NoContent(http.StatusOK)
	}
}