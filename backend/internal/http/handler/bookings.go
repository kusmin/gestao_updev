package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

type BookingRequest struct {
	ClientID       uuid.UUID  `json:"client_id" binding:"required"`
	ProfessionalID uuid.UUID  `json:"professional_id" binding:"required"`
	ServiceID      uuid.UUID  `json:"service_id" binding:"required"`
	Status         string     `json:"status"`
	StartAt        time.Time  `json:"start_at" binding:"required"`
	EndAt          *time.Time `json:"end_at"`
	Notes          string     `json:"notes"`
}

type BookingUpdateRequest struct {
	Status  *string    `json:"status"`
	StartAt *time.Time `json:"start_at"`
	EndAt   *time.Time `json:"end_at"`
	Notes   *string    `json:"notes"`
}

type BookingCancelRequest struct {
	Reason string `json:"reason"`
}

// ListBookings
// @Summary Lista agendamentos
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param date query string false "Data (YYYY-MM-DD)"
// @Param professional_id query string false "Profissional"
// @Param status query string false "Status"
// @Success 200 {object} response.APIResponse
// @Router /bookings [get]
func (api *API) ListBookings(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var (
		datePtr *time.Time
		profID  *uuid.UUID
	)
	if raw := c.Query("date"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			datePtr = &t
		}
	}
	if raw := c.Query("professional_id"); raw != "" {
		if id, err := uuid.Parse(raw); err == nil {
			profID = &id
		}
	}

	bookings, err := api.svc.ListBookings(c.Request.Context(), tenantID, service.BookingFilter{
		Date:           datePtr,
		ProfessionalID: profID,
		Status:         c.Query("status"),
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, bookings, nil)
}

// CreateBooking
// @Summary Cria agendamento
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param request body BookingRequest true "Agendamento"
// @Success 201 {object} response.APIResponse
// @Router /bookings [post]
func (api *API) CreateBooking(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	booking, err := api.svc.CreateBooking(c.Request.Context(), tenantID, service.BookingInput{
		ClientID:       req.ClientID,
		ProfessionalID: req.ProfessionalID,
		ServiceID:      req.ServiceID,
		Status:         req.Status,
		StartAt:        req.StartAt,
		EndAt:          req.EndAt,
		Notes:          req.Notes,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, booking, nil)
}

// UpdateBooking
// @Summary Atualiza agendamento
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Booking ID"
// @Param request body BookingUpdateRequest true "Campos editáveis"
// @Success 200 {object} response.APIResponse
// @Router /bookings/{id} [patch]
func (api *API) UpdateBooking(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	var req BookingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	booking, err := api.svc.UpdateBooking(c.Request.Context(), tenantID, bookingID, service.BookingUpdateInput{
		Status:  req.Status,
		StartAt: req.StartAt,
		EndAt:   req.EndAt,
		Notes:   req.Notes,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, booking, nil)
}

// CancelBooking
// @Summary Cancela agendamento
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Booking ID"
// @Param request body BookingCancelRequest true "Motivo"
// @Success 200 {object} response.APIResponse
// @Router /bookings/{id}/cancel [post]
func (api *API) CancelBooking(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	var req BookingCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	booking, err := api.svc.CancelBooking(c.Request.Context(), tenantID, bookingID, req.Reason)
	if err != nil {
		api.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, booking, nil)
}
