package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

func (h *API) RegisterAdminClientRoutes(router *gin.RouterGroup) {
	router.GET("/clients", h.ListAllClients)
	router.POST("/clients", h.AdminCreateClient)
	router.PUT("/clients/:id", h.AdminUpdateClient)
	router.DELETE("/clients/:id", h.AdminDeleteClient)
}

func (h *API) ListAllClients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	tags := c.Query("tags")

	var tagsFilter []string
	if tags != "" {
		tagsFilter = strings.Split(tags, ",")
	}

	filter := service.ClientsFilter{
		Search:  search,
		Tags:    tagsFilter,
		Page:    page,
		PerPage: perPage,
	}

	clients, _, err := h.svc.ListAllClients(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": clients})
}

type AdminCreateClientInput struct {
	service.ClientInput
	TenantID string `json:"tenant_id" binding:"required"`
}

func (h *API) AdminCreateClient(c *gin.Context) {
	var input AdminCreateClientInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	adminInput := service.AdminClientInput{
		ClientInput: input.ClientInput,
		TenantID:    tenantID,
	}

	client, err := h.svc.AdminCreateClient(c.Request.Context(), adminInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": client})
}

func (h *API) AdminUpdateClient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input service.ClientInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.svc.AdminUpdateClient(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": client})
}

func (h *API) AdminDeleteClient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.AdminDeleteClient(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
