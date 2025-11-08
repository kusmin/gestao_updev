package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

type SignupRequest struct {
	Company struct {
		Name     string `json:"name" binding:"required"`
		Document string `json:"document"`
		Phone    string `json:"phone"`
	} `json:"company" binding:"required"`
	User struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Phone    string `json:"phone"`
	} `json:"user" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Signup
// @Summary Cria empresa e usuário administrador
// @Description Fluxo inicial da plataforma: cria empresa, usuário admin e retorna tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SignupRequest true "Dados do signup"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /auth/signup [post]
func (api *API) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	result, err := api.svc.Signup(c.Request.Context(), service.SignupInput{
		CompanyName:     req.Company.Name,
		CompanyDocument: req.Company.Document,
		CompanyPhone:    req.Company.Phone,
		UserName:        req.User.Name,
		UserEmail:       req.User.Email,
		UserPassword:    req.User.Password,
		UserPhone:       req.User.Phone,
	})
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"tenant_id":     result.TenantID,
		"user_id":       result.UserID,
		"access_token":  result.Tokens.AccessToken,
		"refresh_token": result.Tokens.RefreshToken,
		"expires_in":    result.Tokens.ExpiresIn,
	}, nil)
}

// Login
// @Summary Autentica usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credenciais"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/login [post]
func (api *API) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	_, tokens, err := api.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, tokens, nil)
}

// RefreshToken
// @Summary Atualiza tokens de acesso
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/refresh [post]
func (api *API) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	tokens, err := api.svc.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		api.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, tokens, nil)
}
