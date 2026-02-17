package utils

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/adrputra/face-recognition-svc/gateway-gin/app/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type InterfaceAuthMiddleware interface {
	IsAuthorized() gin.HandlerFunc
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

func (m *AuthMiddleware) IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenVal, exists := c.Get("user")
		if !exists {
			_ = LogError(c, model.ThrowError(http.StatusUnauthorized, errors.New("missing auth token")), nil)
			return
		}

		token, ok := tokenVal.(*jwt.Token)
		if !ok || token == nil {
			_ = LogError(c, model.ThrowError(http.StatusUnauthorized, errors.New("invalid auth token")), nil)
			return
		}

		claims, ok := token.Claims.(*model.JwtCustomClaims)
		if !ok || claims == nil {
			_ = LogError(c, model.ThrowError(http.StatusUnauthorized, errors.New("invalid token claims")), nil)
			return
		}

		ctx := c.Request.Context()
		requiredPermission := strings.TrimSpace(c.GetHeader("app-permission"))
		if requiredPermission != "" {
			if len(claims.RoleIDs) == 0 {
				_ = LogError(c, model.ThrowError(http.StatusForbidden, errors.New("missing role assignment")), nil)
				return
			}

			permissions, err := m.getPermissions(ctx, claims.RoleIDs)
			if err != nil {
				_ = LogError(c, model.ThrowError(http.StatusInternalServerError, err), nil)
				return
			}
			if !Contains(permissions, requiredPermission) {
				_ = LogError(c, model.ThrowError(http.StatusForbidden, errors.New("permission denied")), nil)
				return
			}
		}

		md := metadata.New(map[string]string{
			"user_id":        claims.UserID,
			"username":       claims.Username,
			"role_ids":       strings.Join(claims.RoleIDs, ","),
			"institution_id": claims.InstitutionID,
		})
		c.Request = c.Request.WithContext(metadata.NewIncomingContext(c.Request.Context(), md))
		c.Next()
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
