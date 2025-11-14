package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/service"
	"net/http"
	"time"
)

func (h *API) RegisterAdminBookingRoutes(router *gin.RouterGroup) {
	router.GET("/bookings", h.ListAllBookings)
	router.POST("/bookings", h.AdminCreateBooking)
	router.PUT("/bookings/:id", h.AdminUpdateBooking)
	router.DELETE("/bookings/:id", h.AdminDeleteBooking)
}

func (h *API) ListAllBookings(c *gin.Context) {
	var filter service.BookingFilter

	if dateStr := c.Query("date"); dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			filter.Date = &date
		}
	}
	if profIDStr := c.Query("professional_id"); profIDStr != "" {
		profID, err := uuid.Parse(profIDStr)
		if err == nil {
			filter.ProfessionalID = &profID
		}
	}
	filter.Status = c.Query("status")

	bookings, err := h.svc.ListAllBookings(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

type AdminCreateBookingInput struct {
	service.BookingInput
	TenantID string `json:"tenant_id" binding:"required"`
}

func (h *API) AdminCreateBooking(c *gin.Context) {
	var input AdminCreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	adminInput := service.AdminBookingInput{
		BookingInput: input.BookingInput,
		TenantID:     tenantID,
	}

	booking, err := h.svc.AdminCreateBooking(c.Request.Context(), adminInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": booking})
}

func (h *API) AdminUpdateBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input service.BookingUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.svc.AdminUpdateBooking(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func (h *API) AdminDeleteBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.AdminDeleteBooking(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
