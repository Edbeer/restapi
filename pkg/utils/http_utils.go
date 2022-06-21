package utils

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Get ctx with timeout and request id from echo context
func GetCtxWithReqID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5 * time.Second)
	ctx = context.WithValue(ctx, "ReqID", GetRequestID(c))
	return ctx, cancel
}