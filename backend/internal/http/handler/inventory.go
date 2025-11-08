package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

type InventoryRequest struct {
	ProductID uuid.UUID  `json:"product_id" binding:"required"`
	OrderID   *uuid.UUID `json:"order_id"`
	Type      string     `json:"type" binding:"required"`
	Quantity  int        `json:"quantity" binding:"required"`
	Reason    string     `json:"reason"`
}

// ListInventoryMovements
// @Summary Lista movimentações de estoque
// @Tags Inventory
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param product_id query string false "Filtro por produto"
// @Param type query string false "Tipo (in|out|adjustment)"
// @Param start_date query string false "Data inicial (RFC3339)"
// @Param end_date query string false "Data final (RFC3339)"
// @Success 200 {object} response.APIResponse
// @Router /inventory/movements [get]
func (api *API) ListInventoryMovements(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var (
		productID *uuid.UUID
		startDate *time.Time
		endDate   *time.Time
	)

	if raw := c.Query("product_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err == nil {
			productID = &id
		}
	}
	if raw := c.Query("start_date"); raw != "" {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			startDate = &t
		}
	}
	if raw := c.Query("end_date"); raw != "" {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			endDate = &t
		}
	}
	movements, err := api.svc.ListInventoryMovements(c.Request.Context(), tenantID, service.InventoryFilter{
		ProductID: productID,
		Type:      c.Query("type"),
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, movements, nil)
}

// CreateInventoryMovement
// @Summary Registra movimento de estoque
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param request body InventoryRequest true "Movimentação"
// @Success 201 {object} response.APIResponse
// @Router /inventory/movements [post]
func (api *API) CreateInventoryMovement(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var req InventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	movement, err := api.svc.CreateInventoryMovement(c.Request.Context(), tenantID, service.InventoryInput{
		ProductID: req.ProductID,
		OrderID:   req.OrderID,
		Type:      req.Type,
		Quantity:  req.Quantity,
		Reason:    req.Reason,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, movement, nil)
}
