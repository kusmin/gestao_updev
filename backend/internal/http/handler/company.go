package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

type CompanyUpdateRequest struct {
	Name     *string                `json:"name"`
	Timezone *string                `json:"timezone"`
	Phone    *string                `json:"phone"`
	Email    *string                `json:"email"`
	Settings map[string]interface{} `json:"settings"`
}

// GetCompany
// @Summary Retorna dados da empresa atual
// @Tags Company
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /companies/me [get]
func (api *API) GetCompany(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	company, err := api.svc.GetCompany(c.Request.Context(), tenantID)
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, company, nil)
}

// UpdateCompany
// @Summary Atualiza dados da empresa atual
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param request body CompanyUpdateRequest true "Campos editáveis"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /companies/me [put]
func (api *API) UpdateCompany(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var req CompanyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	company, err := api.svc.UpdateCompany(c.Request.Context(), tenantID, service.CompanyUpdateInput{
		Name:     req.Name,
		Timezone: req.Timezone,
		Phone:    req.Phone,
		Email:    req.Email,
		Settings: req.Settings,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, company, nil)
}
