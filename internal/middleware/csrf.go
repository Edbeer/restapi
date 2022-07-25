package middleware

import (
	"net/http"

	"github.com/Edbeer/restapi/pkg/csrf"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/labstack/echo/v4"
)

// CSRF middleware
func (mw *MiddlewareManager) CSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !mw.config.Server.CSRF {
			return next(c)
		}

		token := c.Request().Header.Get(csrf.CSRFHeader)
		if token == "" {
			mw.logger.Errorf("CSRF Middleware get CSRF header, Token: %s, Error: %s, RequestId: %s",
				token,
				"empty CSRF token",
				utils.GetRequestID(c),
			)
			return c.JSON(http.StatusForbidden, httpe.NewRestError(http.StatusForbidden, "Invalid CSRF Token", "no CSRF Token"))
		}

		sessionId, ok := c.Get("sid").(string)
		if !csrf.ValidateToken(token, sessionId, mw.logger) || !ok {
			mw.logger.Errorf("CSRF Middleware csrf.ValidateToken Token: %s, Error: %s, RequestId: %s",
				token,
				"empty token",
				utils.GetRequestID(c),
			)
			return c.JSON(http.StatusForbidden, httpe.NewRestError(http.StatusForbidden, "Invalid CSRF Token", "no CSRF Token"))
		}
		return next(c)
	}
}
