package utils

import (
	"errors"
	"net/http"
	"strings"

	"face-recognition-svc/app/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
)

func IsAuthorized() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*model.JwtCustomClaims)
			access := strings.Split(claims.MenuMapping[c.Request().Header.Get("app-menu-id")], ",")
			if c.Request().Header.Get("app-menu-id") != claims.Role {
				return LogError(c, model.ThrowError(http.StatusForbidden, errors.New("Anda Tidak Memiliki Akses")), nil)
			}
			if !Contains(access, c.Request().Method) {
				return LogError(c, model.ThrowError(http.StatusForbidden, errors.New("Method Not Allowed")), nil)
			}

			md := metadata.New(map[string]string{
				"username": claims.Name,
				"role_id":  claims.Role,
			})

			c.SetRequest(c.Request().WithContext(metadata.NewIncomingContext(c.Request().Context(), md)))

			return next(c)
		}
	}
}
