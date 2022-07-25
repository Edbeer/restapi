package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/internal/service"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Auth sessions middleware using redis
func (mw *MiddlewareManager) AuthSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(mw.config.Server.CookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(err))
			}
			mw.logger.Errorf("AuthSessionMiddleware RequestID: %s, Error: %v",
				utils.GetRequestID(c),
				err.Error(),
			)
			return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(httpe.Unauthorized))
		}

		sessionId := cookie.Value
		session, err := mw.sessionService.GetSessionByID(c.Request().Context(), sessionId)
		if err != nil {
			mw.logger.Errorf("GetSessionByID RequestID: %s, Cookie value: %s, Error: %v",
				utils.GetRequestID(c),
				sessionId,
				err.Error(),
			)
			return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(httpe.Unauthorized))
		}

		user, err := mw.authService.GetUserByID(c.Request().Context(), session.UserID)
		if err != nil {
			mw.logger.Errorf("GetUserByID RequestID: %s, Error: %v",
				utils.GetRequestID(c),
				err.Error(),
			)
			return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(httpe.Unauthorized))
		}

		c.Set("sid", sessionId)
		c.Set("uid", session.UserID)
		c.Set("user", user)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))

		mw.logger.Info(
			"SessionMiddleware, RequestID: %s, UserID: %s, CookieSessionID: %s",
			utils.GetRequestID(c),
			user.ID.String(),
			sessionId,
		)

		return next(c)
	}
}

// Check auth middleware
func (mw *MiddlewareManager) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			mw.logger.Errorf("CheckAuth.c.Cookie: %s, Cookie: %#v, Error: %v",
				utils.GetRequestID(c),
				cookie,
				err,
			)
			return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(err))
		}
		sessionId := cookie.Value

		session, err := mw.sessionService.GetSessionByID(c.Request().Context(), sessionId)
		if err != nil {
			// Cookie is invalid, delete it from browser
			newCookie := http.Cookie{
				Name: "session_id", 
				Value: sessionId, 
				Expires: time.Now().AddDate(-1, 0, 0),
			}
			c.SetCookie(&newCookie)

			mw.logger.Errorf("CheckAuth.sessUC.GetSessionByID: %s, Cookie: %#v, Error: %s",
				utils.GetRequestID(c),
				cookie,
				err,
			)
			return c.JSON(http.StatusUnauthorized, httpe.NoCookie)
		}

		c.Set("uid", session.UserID.String())
		c.Set("sid", sessionId)
		return next(c)
	}
}

func (mw *MiddlewareManager) OwnerOrAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*entity.User)
			if !ok {
				return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(httpe.Unauthorized))
			}

			if *user.Role == "admin" {
				return next(c)
			}

			if user.ID.String() != c.Param("user_id") {
				mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, Error: %s",
					utils.GetRequestID(c),
					user.ID.String(),
					"invalid user context",
				)
				return c.JSON(http.StatusForbidden, httpe.NewForbiddenError(httpe.Forbidden))
			}

			return next(c)
		}
	}
}

// Role based auth middleware, using ctx user
func (mw *MiddlewareManager) RoleBasedAuthMiddleware(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*entity.User)
			if !ok {
				mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, Error: %s,",
					utils.GetRequestID(c),
					user.ID.String(),
					"invalid user context",
				)
				return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(httpe.Unauthorized))
			}

			for _, role := range roles {
				if role == *user.Role {
					return next(c)
				}
			}

			mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, Error: %s,",
				utils.GetRequestID(c),
				user.ID.String(),
				"invalid user context",
			)

			return c.JSON(http.StatusForbidden, httpe.NewForbiddenError(httpe.PermissionDenied))
		}
	}
}

// JWT way of auth using cookie or Authorization header
func AuthJWTMiddleware(authService service.AuthService, config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")
			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					return c.JSON(httpe.ErrorResponse(httpe.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := validateJWTToken(tokenString, authService, c, config); err != nil {
					return c.JSON(httpe.ErrorResponse(err))
				}
				return next(c)
			} else {
				cookie, err := c.Cookie("jwt-token")
				if err != nil {
					return c.JSON(httpe.ErrorResponse(err))
				}

				if err := validateJWTToken(cookie.Value, authService, c, config); err != nil {
					return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(httpe.Unauthorized))
				}
				return next(c)
			}
		}
	}
}

// Admin role
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*entity.User)
		if !ok || *user.Role != "admin" {
			return c.JSON(http.StatusForbidden, httpe.NewUnauthorizedError(httpe.PermissionDenied))
		}
		return next(c)
	}
}

func validateJWTToken(tokenString string, authService service.AuthService, c echo.Context, config *config.Config) error {
	if tokenString == "" {
		return httpe.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", t.Header["alg"])
		}
		secret := []byte(config.Server.JwtSecretKey)
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return httpe.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return httpe.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		u, err := authService.GetUserByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), "user", u)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
