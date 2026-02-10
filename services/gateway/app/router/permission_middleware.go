package router

import "github.com/labstack/echo/v4"

// RequirePermission sets the permission header, then runs IsAuthorized.
func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if permission != "" {
				c.Request().Header.Set("app-permission", permission)
			}
			return GetFactory().Middleware.Auth.IsAuthorized()(next)(c)
		}
	}
}
