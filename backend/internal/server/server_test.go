package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/kusmin/gestao_updev/backend/internal/config"
)

func TestHealthz(t *testing.T) {
	// setup
	cfg := &config.Config{
		AppEnv: "test",
	}
	logger, _ := zap.NewDevelopment()
	s := New(cfg, logger)

	// execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/healthz", nil)
	s.engine.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"data":{"env":"test","status":"ok"},"error":null}`, w.Body.String())
}

func TestPlaceholderRoutes(t *testing.T) {
	// setup
	cfg := &config.Config{
		AppEnv:       "test",
		TenantHeader: "X-Tenant-ID",
	}
	logger, _ := zap.NewDevelopment()
	s := New(cfg, logger)
	router := s.engine

	// routes to test
	placeholderRoutes := []struct {
		Method string
		Path   string
	}{
		{"POST", "/v1/auth/signup"},
		{"POST", "/v1/auth/login"},
		{"POST", "/v1/auth/refresh"},
		{"GET", "/v1/companies/me"},
		{"PUT", "/v1/companies/me"},
		{"GET", "/v1/users"},
		{"POST", "/v1/users"},
		{"PATCH", "/v1/users/:id"},
		{"DELETE", "/v1/users/:id"},
		{"GET", "/v1/clients"},
		{"POST", "/v1/clients"},
		{"GET", "/v1/clients/:id"},
		{"PUT", "/v1/clients/:id"},
		{"DELETE", "/v1/clients/:id"},
		{"GET", "/v1/professionals"},
		{"GET", "/v1/services"},
		{"POST", "/v1/services"},
		{"PUT", "/v1/services/:id"},
		{"DELETE", "/v1/services/:id"},
		{"GET", "/v1/products"},
		{"POST", "/v1/products"},
		{"PUT", "/v1/products/:id"},
		{"DELETE", "/v1/products/:id"},
		{"GET", "/v1/inventory/movements"},
		{"POST", "/v1/inventory/movements"},
		{"GET", "/v1/bookings"},
		{"POST", "/v1/bookings"},
		{"PATCH", "/v1/bookings/:id"},
		{"POST", "/v1/bookings/:id/cancel"},
		{"GET", "/v1/sales/orders"},
		{"POST", "/v1/sales/orders"},
		{"PATCH", "/v1/sales/orders/:id"},
		{"POST", "/v1/sales/orders/:id/payments"},
		{"GET", "/v1/payments"},
		{"GET", "/v1/dashboard/daily"},
	}

	for _, route := range placeholderRoutes {
		t.Run(fmt.Sprintf("%s_%s", route.Method, route.Path), func(t *testing.T) {
			// execute
			w := httptest.NewRecorder()
			path := strings.Replace(route.Path, ":id", "123", -1)
			req, _ := http.NewRequest(route.Method, path, nil)
			req.Header.Set(cfg.TenantHeader, "test-tenant")
			router.ServeHTTP(w, req)

			// assert
			assert.Equal(t, http.StatusNotImplemented, w.Code)
			assert.Contains(t, w.Body.String(), "NOT_IMPLEMENTED")
		})
	}
}