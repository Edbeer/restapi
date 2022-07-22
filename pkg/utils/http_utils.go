package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/labstack/echo/v4"
)

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Get ctx with timeout and request id from echo context
func GetCtxWithReqID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	ctx = context.WithValue(ctx, "ReqID", GetRequestID(c))
	return ctx, cancel
}

// Configure JWT cookie
func ConfigureJWTCookie(cfg *config.Config, jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:       cfg.Cookie.Name,
		Value:      jwtToken,
		Path:       "/",
		// Domain:     "/",
		RawExpires: "",
		MaxAge:     cfg.Cookie.MaxAge,
		Secure:     cfg.Cookie.Secure,
		// it is better to keep jwtokens on httpOnly from XSS attacks
		HttpOnly:   cfg.Cookie.HTTPOnly, // true
		SameSite:   0,
	}
}

// Configure Session Cookie
func ConfigureSessionCookie(cfg *config.Config, session string) *http.Cookie {
	return &http.Cookie{
		Name: cfg.Session.Name,
		Value: session,
		Path: "/",
		RawExpires: "",
		MaxAge:     cfg.Session.Expire,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HTTPOnly,
		SameSite:   0,
	}
}

func GetUserFromCtx(ctx context.Context) (*entity.User, error) {
	user, ok := ctx.Value("user").(*entity.User)
	if !ok {
		return nil, httpe.Unauthorized
	}
	return user, nil
}