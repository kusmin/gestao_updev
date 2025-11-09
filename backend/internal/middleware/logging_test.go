package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerMiddlewareIncludesCorrelationFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, observed := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	router.Use(RequestID())
	router.Use(func(c *gin.Context) {
		c.Set(ContextTenantIDKey, "tenant-123")
		c.Set(ContextUserIDKey, "user-456")
		c.Next()
	})
	router.Use(Logger(logger))
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/health?foo=bar", nil)
	req.Header.Set(requestIDHeader, "req-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	require.Equal(t, 1, observed.Len())
	entry := observed.All()[0]
	assert.Equal(t, "http_request", entry.Message)
	fields := entry.ContextMap()
	assert.Equal(t, "req-123", fields["request_id"])
	assert.Equal(t, "tenant-123", fields["tenant_id"])
	assert.Equal(t, "user-456", fields["user_id"])
	assert.Equal(t, "/health", fields["path"])
	assert.Equal(t, "foo=bar", fields["query"])
	assert.Equal(t, "GET", fields["method"])
	assert.EqualValues(t, http.StatusOK, fields["status"])
}
