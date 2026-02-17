package router

import "github.com/gin-gonic/gin"

// RequirePermission sets the permission header, then runs IsAuthorized.
func RequirePermission(permission string) gin.HandlerFunc {
	authorizer := GetFactory().Middleware.Auth.IsAuthorized()
	return func(c *gin.Context) {
		if permission != "" {
			c.Request.Header.Set("app-permission", permission)
		}
		authorizer(c)
	}
}
