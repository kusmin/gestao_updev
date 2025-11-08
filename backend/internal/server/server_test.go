package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/config"
)

func setupTestServer(t *testing.T) *Server {
	cfg := &config.Config{
		AppEnv:           "test",
		TenantHeader:     "X-Test-Tenant",
		JWTAccessSecret:  "test-access",
		JWTRefreshSecret: "test-refresh",
		JWTAccessTTL:     15 * time.Minute,
		JWTRefreshTTL:    24 * time.Hour,
	}
	logger := zap.NewNop()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	return New(cfg, logger, db)
}

func TestHealthz(t *testing.T) {
	s := setupTestServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/healthz", nil)
	s.engine.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"data":{"env":"test","status":"ok"},"error":null}`, w.Body.String())
}

func TestProtectedRouteRequiresAuth(t *testing.T) {
	s := setupTestServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/clients", nil)
	req.Header.Set("X-Test-Tenant", "tenant")
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
