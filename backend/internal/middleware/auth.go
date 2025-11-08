package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kusmin/gestao_updev/backend/internal/auth"
	"github.com/kusmin/gestao_updev/backend/internal/http/response"
)

const (
	ContextUserIDKey   = "user_id"
	ContextUserRoleKey = "user_role"
	ContextTenantIDKey = "tenant_id"
)

// Auth valida o JWT e sincroniza tenant/token.
func Auth(jwtManager *auth.JWTManager, tenantHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Token ausente", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Cabeçalho Authorization inválido", nil)
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Token inválido ou expirado", nil)
			c.Abort()
			return
		}

		headerTenant := c.GetHeader(tenantHeader)
		if headerTenant != "" && !strings.EqualFold(headerTenant, claims.TenantID) {
			response.Error(c, http.StatusForbidden, "TENANT_MISMATCH", "Tenant informado não pertence ao token", nil)
			c.Abort()
			return
		}

		c.Set(ContextTenantIDKey, claims.TenantID)
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Next()
	}
}
