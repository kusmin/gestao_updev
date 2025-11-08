package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone"`
	Role     string `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role"`
	Active   *bool   `json:"active"`
	Password *string `json:"password"`
}

// ListUsers
// @Summary Lista usuários do tenant
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param role query string false "Filtro por role"
// @Param page query int false "Página" default(1)
// @Param per_page query int false "Itens por página" default(20)
// @Success 200 {object} response.APIResponse
// @Router /users [get]
func (api *API) ListUsers(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	role := c.Query("role")

	users, total, err := api.svc.ListUsers(c.Request.Context(), tenantID, service.UsersFilter{
		Role:    role,
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, users, metaPagination(page, perPage, total))
}

// CreateUser
// @Summary Cria um usuário
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param request body CreateUserRequest true "Novo usuário"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /users [post]
func (api *API) CreateUser(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	user, err := api.svc.CreateUser(c.Request.Context(), tenantID, service.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Role:     req.Role,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, user, nil)
}

// UpdateUser
// @Summary Atualiza parcialmente um usuário
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "Campos editáveis"
// @Success 200 {object} response.APIResponse
// @Router /users/{id} [patch]
func (api *API) UpdateUser(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	user, err := api.svc.UpdateUser(c.Request.Context(), tenantID, userID, service.UpdateUserInput{
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     req.Role,
		Active:   req.Active,
		Password: req.Password,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, user, nil)
}

// DeleteUser
// @Summary Remove (soft delete) um usuário
// @Tags Users
// @Security BearerAuth
// @Security TenantHeader
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Router /users/{id} [delete]
func (api *API) DeleteUser(c *gin.Context) {
	tenantID, ok := api.tenantID(c)
	if !ok {
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID inválido", nil)
		return
	}

	if err := api.svc.DeleteUser(c.Request.Context(), tenantID, userID); err != nil {
		api.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
