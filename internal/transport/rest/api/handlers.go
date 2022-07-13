package api

import (
	"github.com/Edbeer/restapi/config"
	middle "github.com/Edbeer/restapi/internal/middleware"
	"github.com/Edbeer/restapi/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Deps struct {
	AuthService AuthService
	NewsService NewsService
	Config      *config.Config
	Logger      logger.Logger
}

type Handlers struct {
	Auth *AuthHandler
	News *NewsHandler
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
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
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
	api := e.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register())
			auth.POST("/login", h.Auth.Login())
			auth.POST("/logout", h.Auth.Logout())
			auth.GET("/:user_id", h.Auth.GetUserByID())
			auth.GET("/find", h.Auth.FindUsersByName())
			auth.GET("/all", h.Auth.GetUsers())
			// auth.Use(middleware.AuthJWTMiddleware(*h.service.Auth, h.config))
			auth.PUT("/:user_id", h.Auth.UpdateUser(), middle.OwnerOrAdminMiddleware())
			auth.DELETE("/:user_id", h.Auth.DeleteUser(), middle.RoleBasedAuthMiddleware([]string{"admin"}))
			auth.GET("/me", h.Auth.GetMe())
		}
		news := api.Group("/news")
		{
			news.POST("/create", h.News.CreateNews())
			news.GET("/all", h.News.GetNews())
			news.PUT("/:news_id", h.News.UpdateNews())
			news.DELETE("/:news_id", h.News.DeleteNews())
		}
	}
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(deps.AuthService, deps.Config, deps.Logger),
		News: NewNewsHandler(deps.NewsService, deps.Config, deps.Logger),
	}
}
