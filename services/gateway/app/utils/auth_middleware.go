package utils

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	model "github.com/adrputra/face-recognition-svc/gateway/app/domain"

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
				permissions, err := m.getPermissions(ctx, claims.RoleIDs)
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

func (m *AuthMiddleware) getPermissions(ctx context.Context, roleIDs []string) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	var permissions []string
	cacheKey := ""
	if m.redis != nil {
		sortedRoleIDs := append([]string{}, roleIDs...)
		sort.Strings(sortedRoleIDs)
		cacheKey = "permissions:roles:" + strings.Join(sortedRoleIDs, ",")
		cached, err := m.redis.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			return strings.Split(cached, ","), nil
		}
	}

	query := `
		SELECT p.name
		FROM permission p
		JOIN role_permission rp ON rp.permission_id = p.id
		WHERE rp.role_id IN ?
		AND p.is_active = TRUE`
	if err := m.db.WithContext(ctx).Raw(query, roleIDs).Scan(&permissions).Error; err != nil {
		return nil, err
	}

	if m.redis != nil && cacheKey != "" && len(permissions) > 0 {
		_ = m.redis.Set(ctx, cacheKey, strings.Join(permissions, ","), 5*time.Minute).Err()
	}

	return permissions, nil
}

type permissionClaims struct {
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// RequirePermission enforces RBAC based on JWT permissions.
func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if permission == "" {
				return next(c)
			}

			tokenString := extractBearerToken(c.Request().Header.Get("Authorization"))
			if tokenString == "" {
				return forbiddenError(c, permission)
			}

			claims := &permissionClaims{}
			_, _, err := jwt.NewParser().ParseUnverified(tokenString, claims)
			if err != nil {
				return forbiddenError(c, permission)
			}

			if !hasPermission(claims.Permissions, permission) {
				return forbiddenError(c, permission)
			}

			// Optional: forward permission for downstream services.
			c.Request().Header.Set("app-permission", permission)
			return next(c)
		}
	}
}

func extractBearerToken(headerValue string) string {
	if headerValue == "" {
		return ""
	}
	parts := strings.SplitN(headerValue, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func hasPermission(permissions []string, required string) bool {
	for _, p := range permissions {
		if p == required {
			return true
		}
	}
	return false
}

func forbiddenError(c echo.Context, permission string) error {
	return c.JSON(http.StatusForbidden, map[string]string{
		"error":   "forbidden",
		"message": "missing permission: " + permission,
	})
}
