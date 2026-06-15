package middleware

import (
	"fmt"
	"strings"

	"github.com/felipe/dev-test-api/internal/auth"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/felipe/dev-test-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Problem(c, apierr.ErrUnauthorized("missing authorization header", c.Request.URL.Path))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			response.Problem(c, apierr.ErrUnauthorized("invalid authorization header format", c.Request.URL.Path))
			return
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			response.Problem(c, apierr.ErrUnauthorized("invalid or expired token", c.Request.URL.Path))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Problem(c, apierr.ErrUnauthorized("invalid token claims", c.Request.URL.Path))
			return
		}

		c.Set("user_claims", &claims)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("user_claims")
		if !exists {
			response.Problem(c, apierr.ErrForbidden("admin access required", c.Request.URL.Path))
			return
		}

		mapClaims, ok := claims.(*jwt.MapClaims)
		if !ok || !auth.IsAdmin(mapClaims) {
			response.Problem(c, apierr.ErrForbidden("admin access required", c.Request.URL.Path))
			return
		}

		c.Next()
	}
}
