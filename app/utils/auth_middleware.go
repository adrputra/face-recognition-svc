package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"face-recognition-svc/app/model"

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

			var menuMapping map[string]string
			var access []string

			// Try to get menu mapping from Redis
			ctx := c.Request().Context()
			roleID := claims.Role
			redisVal := m.redis.Get(ctx, roleID).Val()

			if redisVal != "" {
				// Unmarshal the menu mapping from Redis
				if err := json.Unmarshal([]byte(redisVal), &menuMapping); err == nil {
					// Use menu mapping from Redis
					menuID := c.Request().Header.Get("app-menu-id")
					if accessMethod, exists := menuMapping[menuID]; exists {
						access = strings.Split(accessMethod, ",")
					}
				}
			} else {
				var role []*model.MenuRoleMapping
				err := m.db.Debug().Raw("SELECT map.id, map.menu_id, menu.menu_name, map.role_id, menu.menu_route, map.access_method FROM menu_mapping AS map LEFT JOIN menu ON map.menu_id = menu.id LEFT JOIN role ON map.role_id = role.id WHERE role_id = ? ORDER BY map.id ASC", roleID).Scan(&role).Error
				if err != nil {
					return LogError(c, model.ThrowError(http.StatusInternalServerError, err), nil)
				}

				if len(role) < 1 {
					return LogError(c, model.ThrowError(http.StatusBadRequest, errors.New("menu role mapping not found")), nil)
				}

				menuMapping = make(map[string]string)
				for _, v := range role {
					menuMapping[v.MenuID] = v.AccessMethod
				}

				access = strings.Split(menuMapping[c.Request().Header.Get("app-menu-id")], ",")

				menuMappingJSON, _ := json.Marshal(menuMapping)
				m.redis.Set(ctx, roleID, menuMappingJSON, 0).Err()
			}

			if len(access) == 0 || !Contains(access, c.Request().Method) {
				return LogError(c, model.ThrowError(http.StatusForbidden, errors.New("anda tidak memiliki akses")), nil)
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
