package contextutil

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/middleware"
)

// TenantID extrai o tenant_id do contexto.
func TenantID(c *gin.Context) (uuid.UUID, error) {
	value, exists := c.Get(middleware.ContextTenantIDKey)
	if !exists {
		return uuid.Nil, fmt.Errorf("tenant_id not found in context")
	}
	tenantStr, ok := value.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("tenant_id inválido")
	}
	return uuid.Parse(tenantStr)
}

// UserID extrai o identificador do usuário autenticado.
func UserID(c *gin.Context) (uuid.UUID, error) {
	value, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}
	userStr, ok := value.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id inválido")
	}
	return uuid.Parse(userStr)
}
