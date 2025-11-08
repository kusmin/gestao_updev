package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
)

var publicPrefixes = []string{
	"/v1/healthz",
	"/v1/auth/signup",
	"/v1/auth/login",
	"/v1/auth/refresh",
	"/swagger",
}

// TenantEnforcer garante que requisições autenticadas contenham o cabeçalho do tenant.
func TenantEnforcer(headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, prefix := range publicPrefixes {
			if strings.HasPrefix(path, prefix) {
				c.Next()
				return
			}
		}

		if headerName == "" {
			headerName = "X-Tenant-ID"
		}

		tenantID := c.GetHeader(headerName)
		if tenantID == "" {
			response.Error(c, http.StatusBadRequest, "TENANT_ID_REQUIRED", "Cabeçalho de tenant ausente", nil)
			c.Abort()
			return
		}

		c.Set(ContextTenantIDKey, tenantID)
		c.Next()
	}
}
