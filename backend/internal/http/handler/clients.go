package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

type ClientRequest struct {
	Name    string                 `json:"name" binding:"required"`
	Email   string                 `json:"email"`
	Phone   string                 `json:"phone"`
	Notes   string                 `json:"notes"`
	Tags    []string               `json:"tags"`
	Contact map[string]interface{} `json:"contact"`
}

// ListClients
// @Summary Lista clientes
// @Tags Clients
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param search query string false "Filtro por nome/email/telefone"
// @Param tags query string false "Lista de tags separadas por vírgula"
// @Param page query int false "Página" default(1)
// @Param per_page query int false "Itens por página" default(20)
// @Success 200 {object} response.APIResponse
// @Router /clients [get]
func (api *API) ListClients(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	var tags []string
	if raw := c.Query("tags"); raw != "" {
		for _, tag := range strings.Split(raw, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	clients, total, err := api.svc.ListClients(c.Request.Context(), tenantID, service.ClientsFilter{
		Search:  search,
		Tags:    tags,
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, clients, metaPagination(page, perPage, total))
}

// CreateClient
// @Summary Cria cliente
// @Tags Clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param request body ClientRequest true "Cliente"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /clients [post]
func (api *API) CreateClient(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var req ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	client, err := api.svc.CreateClient(c.Request.Context(), tenantID, service.ClientInput{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Notes:   req.Notes,
		Tags:    req.Tags,
		Contact: req.Contact,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, client, nil)
}

// GetClient
// @Summary Busca cliente
// @Tags Clients
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Client ID"
// @Success 200 {object} response.APIResponse
// @Router /clients/{id} [get]
func (api *API) GetClient(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	client, stats, err := api.svc.GetClient(c.Request.Context(), tenantID, clientID)
	if err != nil {
		api.handleError(c, err)
		return
	}

	data := gin.H{"client": client}
	if stats != nil {
		data["stats"] = stats
	}
	response.Success(c, http.StatusOK, data, nil)
}

// UpdateClient
// @Summary Atualiza cliente
// @Tags Clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Client ID"
// @Param request body ClientRequest true "Cliente"
// @Success 200 {object} response.APIResponse
// @Router /clients/{id} [put]
func (api *API) UpdateClient(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	var req ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	client, err := api.svc.UpdateClient(c.Request.Context(), tenantID, clientID, service.ClientInput{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Notes:   req.Notes,
		Tags:    req.Tags,
		Contact: req.Contact,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, client, nil)
}

// DeleteClient
// @Summary Remove cliente
// @Tags Clients
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "Client ID"
// @Success 204 "No Content"
// @Router /clients/{id} [delete]
func (api *API) DeleteClient(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	if err := api.svc.DeleteClient(c.Request.Context(), tenantID, clientID); err != nil {
		api.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
