package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/kusmin/gestao_updev/backend/internal/http/contextutil"
	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

// API agrega os handlers HTTP.
type API struct {
	svc    *service.Service
	logger *zap.Logger
}

// New cria um handler básico.
func New(svc *service.Service, logger *zap.Logger) *API {
	return &API{
		svc:    svc,
		logger: logger,
	}
}

func (api *API) tenantID(c *gin.Context) (uuid.UUID, bool) {
	tenantID, err := contextutil.TenantID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "TENANT_ID_REQUIRED", "Tenant inválido", nil)
		return uuid.Nil, false
	}
	return tenantID, true
}

func (api *API) handleError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrInvalidCredentials) {
		response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Credenciais inválidas", nil)
		return
	}
	if errors.Is(err, service.ErrBookingConflict) {
		response.Error(c, http.StatusConflict, "BOOKING_CONFLICT", err.Error(), nil)
		return
	}
	response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
}

func metaPagination(page, perPage int, total int64) gin.H {
	return gin.H{
		"pagination": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	}
}
