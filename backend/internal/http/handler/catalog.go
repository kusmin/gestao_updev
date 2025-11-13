package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

type ServiceRequest struct {
	Name            string                 `json:"name" binding:"required"`
	Category        string                 `json:"category"`
	Description     string                 `json:"description"`
	DurationMinutes int                    `json:"duration_minutes" binding:"required"`
	Price           float64                `json:"price" binding:"required"`
	Color           string                 `json:"color"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type ProductRequest struct {
	Name        string                 `json:"name" binding:"required"`
	SKU         string                 `json:"sku" binding:"required"`
	Price       float64                `json:"price" binding:"required"`
	Cost        float64                `json:"cost"`
	StockQty    int                    `json:"stock_qty"`
	MinStock    int                    `json:"min_stock"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ListServices
// @Summary Lista serviços
// @Tags Services
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Success 200 {object} response.APIResponse
// @Router /services [get]
func (api *API) ListServices(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	servicesList, err := api.svc.ListServices(c.Request.Context(), tenantID)
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, servicesList, nil)
}

// CreateService
// @Summary Cria serviço
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param request body ServiceRequest true "Serviço"
// @Success 201 {object} response.APIResponse
// @Router /services [post]
func (api *API) CreateService(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	serviceObj, err := api.svc.CreateService(c.Request.Context(), tenantID, service.ServiceInput{
		Name:            req.Name,
		Category:        req.Category,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Price:           req.Price,
		Color:           req.Color,
		Metadata:        req.Metadata,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, serviceObj, nil)
}

// UpdateService
// @Summary Atualiza serviço
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Service ID"
// @Param request body ServiceRequest true "Serviço"
// @Success 200 {object} response.APIResponse
// @Router /services/{id} [put]
func (api *API) UpdateService(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	serviceObj, err := api.svc.UpdateService(c.Request.Context(), tenantID, serviceID, service.ServiceInput{
		Name:            req.Name,
		Category:        req.Category,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Price:           req.Price,
		Color:           req.Color,
		Metadata:        req.Metadata,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, serviceObj, nil)
}

// DeleteService
// @Summary Remove serviço
// @Tags Services
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Service ID"
// @Success 204 "No Content"
// @Router /services/{id} [delete]
func (api *API) DeleteService(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	if err := api.svc.DeleteService(c.Request.Context(), tenantID, serviceID); err != nil {
		api.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetService
// @Summary Busca serviço por ID
// @Tags Services
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Service ID"
// @Success 200 {object} response.APIResponse
// @Router /services/{id} [get]
func (api *API) GetService(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID de serviço inválido", nil)
		return
	}

	service, err := api.svc.GetService(c.Request.Context(), tenantID, serviceID)
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, service, nil)
}

// ListProducts
// @Summary Lista produtos
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Success 200 {object} response.APIResponse
// @Router /products [get]
func (api *API) ListProducts(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	products, err := api.svc.ListProducts(c.Request.Context(), tenantID)
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, products, nil)
}

// CreateProduct
// @Summary Cria produto
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param request body ProductRequest true "Produto"
// @Success 201 {object} response.APIResponse
// @Router /products [post]
func (api *API) CreateProduct(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	product, err := api.svc.CreateProduct(c.Request.Context(), tenantID, service.ProductInput{
		Name:        req.Name,
		SKU:         req.SKU,
		Price:       req.Price,
		Cost:        req.Cost,
		StockQty:    req.StockQty,
		MinStock:    req.MinStock,
		Description: req.Description,
		Metadata:    req.Metadata,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, product, nil)
}

// UpdateProduct
// @Summary Atualiza produto
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Product ID"
// @Param request body ProductRequest true "Produto"
// @Success 200 {object} response.APIResponse
// @Router /products/{id} [put]
func (api *API) UpdateProduct(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	product, err := api.svc.UpdateProduct(c.Request.Context(), tenantID, productID, service.ProductInput{
		Name:        req.Name,
		SKU:         req.SKU,
		Price:       req.Price,
		Cost:        req.Cost,
		StockQty:    req.StockQty,
		MinStock:    req.MinStock,
		Description: req.Description,
		Metadata:    req.Metadata,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, product, nil)
}

// DeleteProduct
// @Summary Remove produto
// @Tags Products
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Router /products/{id} [delete]
func (api *API) DeleteProduct(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	if err := api.svc.DeleteProduct(c.Request.Context(), tenantID, productID); err != nil {
		api.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetProduct
// @Summary Busca produto por ID
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Product ID"
// @Success 200 {object} response.APIResponse
// @Router /products/{id} [get]
func (api *API) GetProduct(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID de produto inválido", nil)
		return
	}

	product, err := api.svc.GetProduct(c.Request.Context(), tenantID, productID)
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, product, nil)
}
