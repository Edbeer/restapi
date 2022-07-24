package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/internal/service"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

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

func OwnerOrAdminMiddleware() echo.MiddlewareFunc {
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
				return c.JSON(http.StatusForbidden, httpe.NewForbiddenError(httpe.Forbidden))
			}

			return next(c)
		}
	}
}

// Role based auth middleware, using ctx user
func RoleBasedAuthMiddleware(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*entity.User)
			if !ok {
				return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(httpe.Unauthorized))
			}

			for _, role := range roles {
				if role == *user.Role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, httpe.NewForbiddenError(httpe.PermissionDenied))
		}
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
		c.Request().WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
