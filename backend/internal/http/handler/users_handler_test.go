package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/auth"
	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"github.com/kusmin/gestao_updev/backend/internal/middleware"
	"github.com/kusmin/gestao_updev/backend/internal/repository"
	"github.com/kusmin/gestao_updev/backend/internal/service"
)

var testDB *gorm.DB

// setupTestAPI cria um ambiente de teste completo com um tenant e um usuário admin.
// Retorna o motor do Gin, o tenant criado e um token de autenticação válido.
func setupTestAPI(t *testing.T) (*gin.Engine, *domain.Company, string) {
	// Reutiliza a lógica de setup do banco de dados dos testes de serviço
	var err error
	testDB, err = setupTestDatabase()
	require.NoError(t, err)

	// Garante que o banco esteja limpo
	clearAllData()

	// Roda as migrações
	err = testDB.AutoMigrate(
		&domain.Company{}, &domain.User{}, &domain.Client{}, &domain.Service{},
		&domain.Product{}, &domain.Booking{}, &domain.SalesOrder{}, &domain.SalesItem{},
		&domain.Payment{}, &domain.InventoryMovement{},
	)
	require.NoError(t, err)

	// Cria o tenant
	tenant, err := createTestTenant()
	require.NoError(t, err)

	// Configuração
	cfg := &config.Config{
		JWTAccessSecret: "test-secret",
		TenantHeader:    "X-Tenant-Id",
	}
	logger := zap.NewNop()
	repo := repository.New(testDB)
	jwtManager := auth.NewJWTManager(cfg.JWTAccessSecret, "", 0, 0)
	svc := service.New(cfg, repo, jwtManager, logger)

	// Cria um usuário admin para autenticação
	adminUser, err := svc.CreateUser(context.Background(), tenant.ID, service.CreateUserInput{
		Name:     "Admin",
		Email:    "admin@test.com",
		Password: "password",
		Role:     "admin",
	})
	require.NoError(t, err)

	// Gera o token
	token, err := jwtManager.GenerateAccessToken(adminUser.ID.String(), tenant.ID.String(), adminUser.Role)
	require.NoError(t, err)

	// Cria o handler e o motor do Gin
	apiHandler := New(svc, logger)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/v1")
	api.Use(middleware.TenantEnforcer(cfg.TenantHeader))
	api.Use(middleware.Auth(jwtManager, cfg.TenantHeader))
	registerUserRoutes(api, apiHandler) // Função helper para registrar apenas rotas de usuário

	return router, tenant, token
}

// registerUserRoutes é uma função helper para registrar apenas as rotas de usuário para os testes.
func registerUserRoutes(api *gin.RouterGroup, h *API) {
	api.POST("/users", h.CreateUser)
	api.GET("/users/:id", h.GetUser)
	api.PATCH("/users/:id", h.UpdateUser)
	api.DELETE("/users/:id", h.DeleteUser)
}

func TestCreateUserHandler(t *testing.T) {
	router, tenant, token := setupTestAPI(t)

	t.Run("should create user successfully", func(t *testing.T) {
		// Payload
		userPayload := CreateUserRequest{
			Name:     "New User",
			Email:    "newuser@test.com",
			Password: "password123",
			Role:     "professional",
		}
		payloadBytes, _ := json.Marshal(userPayload)

		// Request
		req, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Tenant-Id", tenant.ID.String())

		// Executa o request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Asserts
		assert.Equal(t, http.StatusCreated, w.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err)

		data, ok := responseBody["data"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "New User", data["name"])
		assert.Equal(t, "newuser@test.com", data["email"])
		assert.Equal(t, "professional", data["role"])
		assert.NotContains(t, data, "password_hash", "A senha não deve ser retornada")
	})

	t.Run("should return validation error for invalid payload", func(t *testing.T) {
		// Payload inválido (sem nome)
		userPayload := CreateUserRequest{
			Email:    "invalid@test.com",
			Password: "password123",
			Role:     "professional",
		}
		payloadBytes, _ := json.Marshal(userPayload)

		// Request
		req, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Tenant-Id", tenant.ID.String())

		// Executa o request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Asserts
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err)
		assert.Equal(t, "VALIDATION_ERROR", responseBody["code"])
	})

	t.Run("should return error for duplicate email", func(t *testing.T) {
		// Cria um usuário uma vez
		userPayload := CreateUserRequest{
			Name:     "Duplicate User",
			Email:    "duplicate@test.com",
			Password: "password123",
			Role:     "professional",
		}
		payloadBytes, _ := json.Marshal(userPayload)
		req, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Tenant-Id", tenant.ID.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		// Tenta criar o mesmo usuário de novo
		req2, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(payloadBytes))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer "+token)
		req2.Header.Set("X-Tenant-Id", tenant.ID.String())
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		// Asserts
		assert.Equal(t, http.StatusInternalServerError, w2.Code) // GORM retorna erro de constraint
	})
}

func TestGetUserHandler(t *testing.T) {
	router, tenant, token := setupTestAPI(t)

	// Cria um usuário para buscar
	user := &domain.User{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		Name:        "Find Me",
		Email:       "findme@test.com",
		Role:        "viewer",
	}
	require.NoError(t, testDB.Create(user).Error)

	t.Run("should get user successfully", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/users/"+user.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Tenant-Id", tenant.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err)
		data, _ := responseBody["data"].(map[string]interface{})
		assert.Equal(t, "Find Me", data["name"])
	})

	t.Run("should return error for non-existent user", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		req, _ := http.NewRequest(http.MethodGet, "/v1/users/"+nonExistentID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Tenant-Id", tenant.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code) // GORM ErrRecordNotFound
	})
}

func TestUpdateUserHandler(t *testing.T) {
	router, tenant, token := setupTestAPI(t)

	user := &domain.User{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		Name:        "Before Update",
		Email:       "update@test.com",
		Role:        "viewer",
	}
	require.NoError(t, testDB.Create(user).Error)

	t.Run("should update user successfully", func(t *testing.T) {
		newName := "After Update"
		payload := UpdateUserRequest{Name: &newName}
		payloadBytes, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPatch, "/v1/users/"+user.ID.String(), bytes.NewBuffer(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Tenant-Id", tenant.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err)
		data, _ := responseBody["data"].(map[string]interface{})
		assert.Equal(t, "After Update", data["name"])
	})
}

func TestDeleteUserHandler(t *testing.T) {
	router, tenant, token := setupTestAPI(t)

	user := &domain.User{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		Name:        "To Be Deleted",
		Email:       "delete@test.com",
		Role:        "viewer",
	}
	require.NoError(t, testDB.Create(user).Error)

	t.Run("should delete user successfully", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/v1/users/"+user.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Tenant-Id", tenant.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verifica se o usuário foi realmente soft-deleted
		var count int64
		testDB.Model(&domain.User{}).Where("id = ?", user.ID).Count(&count)
		assert.Equal(t, int64(0), count, "User should not be found by default scope")

		var deletedUser domain.User
		testDB.Unscoped().First(&deletedUser, user.ID)
		assert.NotNil(t, deletedUser.DeletedAt, "DeletedAt should be set")
	})
}

// --- Funções Helper de Teste ---

func setupTestDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://testuser:testpassword@localhost:5433/testdb?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTestTenant() (*domain.Company, error) {
	tenant := &domain.Company{
		Name:     "Test Tenant",
		Document: uuid.New().String(),
	}
	if err := testDB.Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

func clearAllData() {
	tables := []string{
		"payments", "sales_items", "sales_orders", "inventory_movements",
		"bookings", "products", "services", "clients", "users", "companies",
	}
	for _, table := range tables {
		testDB.Exec("DELETE FROM " + table)
	}
}
