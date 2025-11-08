package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

type SalesOrderItemRequest struct {
	Type      string    `json:"type" binding:"required"`
	RefID     uuid.UUID `json:"ref_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
	UnitPrice float64   `json:"unit_price" binding:"required"`
}

type SalesOrderRequest struct {
	ClientID  uuid.UUID               `json:"client_id" binding:"required"`
	BookingID *uuid.UUID              `json:"booking_id"`
	Items     []SalesOrderItemRequest `json:"items" binding:"required,dive"`
	Discount  float64                 `json:"discount"`
	Notes     string                  `json:"notes"`
}

type SalesOrderUpdateRequest struct {
	Status *string `json:"status"`
	Notes  *string `json:"notes"`
}

type PaymentRequest struct {
	Method  string                 `json:"method" binding:"required"`
	Amount  float64                `json:"amount" binding:"required"`
	PaidAt  time.Time              `json:"paid_at" binding:"required"`
	Details map[string]interface{} `json:"details"`
}

// ListSalesOrders
// @Summary Lista pedidos/vendas
// @Tags Sales
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param status query string false "Status"
// @Param client_id query string false "Cliente"
// @Param date query string false "Data (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router /sales/orders [get]
func (api *API) ListSalesOrders(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}
	var (
		clientID *uuid.UUID
		datePtr  *time.Time
	)
	if raw := c.Query("client_id"); raw != "" {
		if id, err := uuid.Parse(raw); err == nil {
			clientID = &id
		}
	}
	if raw := c.Query("date"); raw != "" {
		if d, err := time.Parse("2006-01-02", raw); err == nil {
			datePtr = &d
		}
	}

	orders, err := api.svc.ListSalesOrders(c.Request.Context(), tenantID, service.SalesOrderFilter{
		Status:   c.Query("status"),
		ClientID: clientID,
		Date:     datePtr,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, orders, nil)
}

// CreateSalesOrder
// @Summary Cria pedido/venda
// @Tags Sales
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param request body SalesOrderRequest true "Pedido"
// @Success 201 {object} response.APIResponse
// @Router /sales/orders [post]
func (api *API) CreateSalesOrder(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var req SalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	items := make([]service.SalesItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = service.SalesItemInput{
			Type:      item.Type,
			RefID:     item.RefID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	order, err := api.svc.CreateSalesOrder(c.Request.Context(), tenantID, service.SalesOrderInput{
		ClientID:  req.ClientID,
		BookingID: req.BookingID,
		Items:     items,
		Discount:  req.Discount,
		Notes:     req.Notes,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, order, nil)
}

// UpdateSalesOrder
// @Summary Atualiza pedido
// @Tags Sales
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Order ID"
// @Param request body SalesOrderUpdateRequest true "Campos editáveis"
// @Success 200 {object} response.APIResponse
// @Router /sales/orders/{id} [patch]
func (api *API) UpdateSalesOrder(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	var req SalesOrderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	order, err := api.svc.UpdateSalesOrder(c.Request.Context(), tenantID, orderID, service.SalesOrderUpdateInput{
		Status: req.Status,
		Notes:  req.Notes,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, order, nil)
}

// CreatePayment
// @Summary Registra pagamento
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Order ID"
// @Param request body PaymentRequest true "Pagamento"
// @Success 201 {object} response.APIResponse
// @Router /sales/orders/{id}/payments [post]
func (api *API) CreatePayment(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	payment, err := api.svc.AddPayment(c.Request.Context(), tenantID, orderID, service.PaymentInput{
		Method:  req.Method,
		Amount:  req.Amount,
		PaidAt:  req.PaidAt,
		Details: req.Details,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, payment, nil)
}

// ListPayments
// @Summary Lista pagamentos
// @Tags Payments
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param method query string false "Método"
// @Param start_date query string false "Data inicial RFC3339"
// @Param end_date query string false "Data final RFC3339"
// @Success 200 {object} response.APIResponse
// @Router /payments [get]
func (api *API) ListPayments(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var (
		startDate *time.Time
		endDate   *time.Time
	)
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

	payments, err := api.svc.ListPayments(c.Request.Context(), tenantID, service.PaymentFilter{
		Method:    c.Query("method"),
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, payments, nil)
}
