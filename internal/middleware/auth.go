package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/service"
	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func AuthJWTMiddleware(authService service.AuthPsql, config *config.Config) echo.MiddlewareFunc {
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
					return c.JSON(httpe.ErrorResponse(err))
				}
				return next(c)
			}
		}
	}
}

func validateJWTToken(tokenString string, authService service.AuthPsql, c echo.Context, config *config.Config) error {
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

		ctx := context.WithValue(c.Request().Context(), "user", u)
		c.Request().WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}