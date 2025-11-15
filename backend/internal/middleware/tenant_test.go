package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTenantEnforcer_PublicRoute(t *testing.T) {
	// setup
	router := gin.New()
	router.Use(TenantEnforcer("X-Tenant-ID"))
	router.GET("/v1/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/healthz", nil)
	router.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTenantEnforcer_PrivateRoute_NoHeader(t *testing.T) {
	// setup
	router := gin.New()
	router.Use(TenantEnforcer("X-Tenant-ID"))
	router.GET("/private", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/private", nil)
	router.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"data":{},"meta":{},"error":{"code":"TENANT_ID_REQUIRED","message":"Cabeçalho de tenant ausente"}}`, w.Body.String())
}

func TestTenantEnforcer_PrivateRoute_WithHeader(t *testing.T) {
	// setup
	router := gin.New()
	router.Use(TenantEnforcer("X-Tenant-ID"))
	router.GET("/private", func(c *gin.Context) {
		tenantID, exists := c.Get("tenant_id")
		assert.True(t, exists)
		assert.Equal(t, "test-tenant", tenantID)
		c.Status(http.StatusOK)
	})

	// execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/private", nil)
	req.Header.Set("X-Tenant-ID", "test-tenant")
	router.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}
