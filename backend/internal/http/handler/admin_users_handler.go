package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/service"
	"net/http"
	"strconv"
)

func (h *API) RegisterAdminUserRoutes(router *gin.RouterGroup) {
	router.GET("/users", h.ListAllUsers)
	router.POST("/users", h.AdminCreateUser)
	router.PUT("/users/:id", h.AdminUpdateUser)
	router.DELETE("/users/:id", h.AdminDeleteUser)
}

func (h *API) ListAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	role := c.Query("role")

	filter := service.UsersFilter{
		Role:    role,
		Page:    page,
		PerPage: perPage,
	}

	users, _, err := h.svc.ListAllUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

type AdminCreateUserInput struct {
	service.CreateUserInput
	TenantID string `json:"tenant_id" binding:"required"`
}

func (h *API) AdminCreateUser(c *gin.Context) {
	var input AdminCreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	user, err := h.svc.AdminCreateUser(c.Request.Context(), input.CreateUserInput, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *API) AdminUpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input service.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.AdminUpdateUser(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *API) AdminDeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.AdminDeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
