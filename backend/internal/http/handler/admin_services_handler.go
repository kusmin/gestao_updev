package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/service"
	"net/http"
)

func (h *API) RegisterAdminServiceRoutes(router *gin.RouterGroup) {
	router.GET("/services", h.ListAllServices)
	router.POST("/services", h.AdminCreateService)
	router.PUT("/services/:id", h.AdminUpdateService)
	router.DELETE("/services/:id", h.AdminDeleteService)
}

func (h *API) ListAllServices(c *gin.Context) {
	services, err := h.svc.ListAllServices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": services})
}

type AdminCreateServiceInput struct {
	service.ServiceInput
	TenantID string `json:"tenant_id" binding:"required"`
}

func (h *API) AdminCreateService(c *gin.Context) {
	var input AdminCreateServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	adminInput := service.AdminServiceInput{
		ServiceInput: input.ServiceInput,
		TenantID:     tenantID,
	}

	service, err := h.svc.AdminCreateService(c.Request.Context(), adminInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": service})
}

func (h *API) AdminUpdateService(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input service.ServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.svc.AdminUpdateService(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": service})
}

func (h *API) AdminDeleteService(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.AdminDeleteService(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
