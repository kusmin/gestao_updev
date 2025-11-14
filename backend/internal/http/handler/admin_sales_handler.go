package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/service"
	"net/http"
)

func (h *API) RegisterAdminSalesRoutes(router *gin.RouterGroup) {
	router.GET("/sales/orders", h.ListAllSalesOrders)
	router.POST("/sales/orders", h.AdminCreateSalesOrder)
	router.PUT("/sales/orders/:id", h.AdminUpdateSalesOrder)
	router.DELETE("/sales/orders/:id", h.AdminDeleteSalesOrder)
}

func (h *API) ListAllSalesOrders(c *gin.Context) {
	salesOrders, err := h.svc.ListAllSalesOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": salesOrders})
}

type AdminCreateSalesOrderInput struct {
	service.SalesOrderInput
	TenantID string `json:"tenant_id" binding:"required"`
}

func (h *API) AdminCreateSalesOrder(c *gin.Context) {
	var input AdminCreateSalesOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	adminInput := service.AdminSalesOrderInput{
		SalesOrderInput: input.SalesOrderInput,
		TenantID:        tenantID,
	}

	salesOrder, err := h.svc.AdminCreateSalesOrder(c.Request.Context(), adminInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": salesOrder})
}

func (h *API) AdminUpdateSalesOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input service.SalesOrderUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	salesOrder, err := h.svc.AdminUpdateSalesOrder(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": salesOrder})
}

func (h *API) AdminDeleteSalesOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.AdminDeleteSalesOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
