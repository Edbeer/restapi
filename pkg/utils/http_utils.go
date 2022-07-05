package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/Edbeer/restapi/config"
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
