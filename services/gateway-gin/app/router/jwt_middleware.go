package router

import (
	"errors"
	"net/http"
	"strings"

	"github.com/adrputra/face-recognition-svc/gateway-gin/app/model"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(signingKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			_ = utils.LogError(c, model.ThrowError(http.StatusUnauthorized, errors.New("missing bearer token")), nil)
			return
		}

		claims := new(model.JwtCustomClaims)
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (any, error) {
				return []byte(signingKey), nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)
		if err != nil || !token.Valid {
			_ = utils.LogError(c, model.ThrowError(http.StatusUnauthorized, errors.New("invalid token")), nil)
			return
		}

		c.Set("user", token)
		c.Next()
	}
}

func extractBearerToken(headerValue string) string {
	parts := strings.SplitN(strings.TrimSpace(headerValue), " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
