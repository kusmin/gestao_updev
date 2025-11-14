package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/service"
	"net/http"
)

func (h *API) RegisterAdminProductRoutes(router *gin.RouterGroup) {
	router.GET("/products", h.ListAllProducts)
	router.POST("/products", h.AdminCreateProduct)
	router.PUT("/products/:id", h.AdminUpdateProduct)
	router.DELETE("/products/:id", h.AdminDeleteProduct)
}

func (h *API) ListAllProducts(c *gin.Context) {
	products, err := h.svc.ListAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

type AdminCreateProductInput struct {
	service.ProductInput
	TenantID string `json:"tenant_id" binding:"required"`
}

func (h *API) AdminCreateProduct(c *gin.Context) {
	var input AdminCreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	adminInput := service.AdminProductInput{
		ProductInput: input.ProductInput,
		TenantID:     tenantID,
	}

	product, err := h.svc.AdminCreateProduct(c.Request.Context(), adminInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

func (h *API) AdminUpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input service.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.svc.AdminUpdateProduct(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *API) AdminDeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.AdminDeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
