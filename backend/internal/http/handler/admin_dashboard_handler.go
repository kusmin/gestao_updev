package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *API) RegisterAdminDashboardRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard/metrics", h.GetOverallMetrics)
}

func (h *API) GetOverallMetrics(c *gin.Context) {
	metrics, err := h.svc.GetOverallMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": metrics})
}
