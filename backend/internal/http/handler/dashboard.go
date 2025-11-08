package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
)

// DashboardDaily
// @Summary KPIs diários
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param date query string false "Data (YYYY-MM-DD)"
// @Param professional_id query string false "Profissional"
// @Success 200 {object} response.APIResponse
// @Router /dashboard/daily [get]
func (api *API) DashboardDaily(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	date := time.Now()
	if raw := c.Query("date"); raw != "" {
		if d, err := time.Parse("2006-01-02", raw); err == nil {
			date = d
		}
	}

	var professionalID *uuid.UUID
	if raw := c.Query("professional_id"); raw != "" {
		if id, err := uuid.Parse(raw); err == nil {
			professionalID = &id
		}
	}

	result, err := api.svc.DashboardDaily(c.Request.Context(), tenantID, date, professionalID)
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, result, nil)
}
