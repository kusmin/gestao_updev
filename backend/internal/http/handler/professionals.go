package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
)

// ListProfessionals
// @Summary Lista profissionais ativos
// @Tags Professionals
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Success 200 {object} response.APIResponse
// @Router /professionals [get]
func (api *API) ListProfessionals(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	professionals, err := api.svc.ListProfessionals(c.Request.Context(), tenantID)
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, professionals, nil)
}
