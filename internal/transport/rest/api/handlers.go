package api

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/docs"
	middle "github.com/Edbeer/restapi/internal/middleware"
	"github.com/Edbeer/restapi/pkg/csrf"
	"github.com/Edbeer/restapi/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
)

type Deps struct {
	AuthService     AuthService
	NewsService     NewsService
	CommentsService CommentsService
	SessionService 	SessionService
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
		auth:     NewAuthHandler(deps.Config, deps.AuthService, deps.SessionService, deps.Logger),
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
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	// Request ID middleware generates a unique id for a request.
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	// Middleware Manager
	mw := middle.NewMiddlewareManager(h.auth.sessionService, 
		h.auth.authService, 
		h.auth.config, 
		[]string{"*"}, 
		h.auth.logger,
	)
	docs.SwaggerInfo.Title = "Go example restapi"
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	
	h.initApi(e, mw)

	return nil
}

func (h *Handlers) initApi(e *echo.Echo, mw *middle.MiddlewareManager) {
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
			auth.Use(mw.AuthSessionMiddleware)
			auth.GET("/token", h.auth.GetCSRFToken())
			auth.PUT("/:user_id", h.auth.Update(), mw.OwnerOrAdminMiddleware(), mw.CSRF)
			auth.DELETE("/:user_id", h.auth.Delete(), mw.RoleBasedAuthMiddleware([]string{"admin"}))
			auth.GET("/me", h.auth.GetMe())
		}

		news := api.Group("/news")
		{
			news.POST("/create", h.news.Create(), mw.AuthSessionMiddleware, mw.CSRF)
			news.PUT("/:news_id", h.news.Update(), mw.AuthSessionMiddleware, mw.CSRF)
			news.DELETE("/:news_id", h.news.Delete(), mw.AuthSessionMiddleware, mw.CSRF)
			news.GET("/all", h.news.GetNews())
			news.GET("/:news_id", h.news.GetNewsByID())
			news.GET("/search", h.news.SearchNews())
		}

		comments := api.Group("/comments")
		{
			comments.POST("", h.comments.Create(), mw.AuthSessionMiddleware, mw.CSRF)
			comments.PUT("/:comments_id", h.comments.Update(), mw.AuthSessionMiddleware, mw.CSRF)
			comments.DELETE("/delete", h.comments.Delete(), mw.AuthSessionMiddleware, mw.CSRF)
			comments.GET("/:comments_id", h.comments.GetByID())
			comments.GET("/byNewsID/:news_id", h.comments.GetAllByNewsID())
		}
	}
}
