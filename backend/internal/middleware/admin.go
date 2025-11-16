package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kusmin/gestao_updev/backend/internal/auth"
	"github.com/kusmin/gestao_updev/backend/internal/http/response"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized", nil)
			c.Abort()
			return
		}

		jwtClaims, ok := claims.(*auth.Claims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized", nil)
			c.Abort()
			return
		}

		if jwtClaims.Role != "admin" {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "forbidden", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
