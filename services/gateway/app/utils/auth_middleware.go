package utils

import (
	"errors"
	"net/http"
	"strings"

	"face-recognition-svc/gateway/app/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type InterfaceAuthMiddleware interface {
	IsAuthorized() echo.MiddlewareFunc
}

type AuthMiddleware struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAuthMiddleware(db *gorm.DB, redis *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		db:    db,
		redis: redis,
	}
}

func (m *AuthMiddleware) IsAuthorized() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*model.JwtCustomClaims)

			ctx := c.Request().Context()
			requiredPermission := strings.TrimSpace(c.Request().Header.Get("app-permission"))
			if requiredPermission != "" {
				if len(claims.RoleIDs) == 0 {
					return LogError(c, model.ThrowError(http.StatusForbidden, errors.New("missing role assignment")), nil)
				}
				var permissions []string
				query := `
					SELECT p.name
					FROM permission p
					JOIN role_permission rp ON rp.permission_id = p.id
					WHERE rp.role_id IN ?
					AND p.is_active = TRUE`
				err := m.db.WithContext(ctx).Raw(query, claims.RoleIDs).Scan(&permissions).Error
				if err != nil {
					return LogError(c, model.ThrowError(http.StatusInternalServerError, err), nil)
				}
				if !Contains(permissions, requiredPermission) {
					return LogError(c, model.ThrowError(http.StatusForbidden, errors.New("permission denied")), nil)
				}
			}

			md := metadata.New(map[string]string{
				"user_id":        claims.UserID,
				"username":       claims.Username,
				"role_ids":       strings.Join(claims.RoleIDs, ","),
				"institution_id": claims.InstitutionID,
			})

			c.SetRequest(c.Request().WithContext(metadata.NewIncomingContext(c.Request().Context(), md)))

			return next(c)
		}
	}
}
