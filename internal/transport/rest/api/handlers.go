package api

import (
	"github.com/Edbeer/restapi/config"
	middle "github.com/Edbeer/restapi/internal/middleware"
	"github.com/Edbeer/restapi/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Deps struct {
	AuthService     AuthService
	NewsService     NewsService
	CommentsService CommentsService
	Config          *config.Config
	Logger          logger.Logger
}

type Handlers struct {
	auth     *AuthHandler
	news     *NewsHandler
	comments *CommentsHandler
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		auth:     NewAuthHandler(deps.AuthService, deps.Config, deps.Logger),
		news:     NewNewsHandler(deps.NewsService, deps.Config, deps.Logger),
		comments: NewCommentsHandler(deps.CommentsService, deps.Config, deps.Logger),
	}
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
			auth.POST("/register", h.auth.Register())
			auth.POST("/login", h.auth.Login())
			auth.POST("/logout", h.auth.Logout())
			auth.GET("/:user_id", h.auth.GetUserByID())
			auth.GET("/find", h.auth.FindUsersByName())
			auth.GET("/all", h.auth.GetUsers())
			// auth.Use(middleware.AuthJWTMiddleware(*h.service.Auth, h.config))
			auth.PUT("/:user_id", h.auth.Update(), middle.OwnerOrAdminMiddleware())
			auth.DELETE("/:user_id", h.auth.Delete(), middle.RoleBasedAuthMiddleware([]string{"admin"}))
			auth.GET("/me", h.auth.GetMe())
		}

		news := api.Group("/news")
		{
			news.POST("/create", h.news.Create())
			news.GET("/all", h.news.GetNews())
			news.GET("/:news_id", h.news.GetNewsByID())
			news.GET("/search", h.news.SearchNews())
			news.PUT("/:news_id", h.news.Update())
			news.DELETE("/:news_id", h.news.Delete())
		}

		comments := api.Group("/comments")
		{
			comments.POST("", h.comments.Create())
			comments.GET("/:comments_id", h.comments.GetByID())
			comments.PUT("/:comments_id", h.comments.Update())
			comments.DELETE("/delete", h.comments.Delete())
		}

	}
}
