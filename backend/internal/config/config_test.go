package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Test with environment variables set
	os.Setenv("APP_ENV", "production")
	os.Setenv("HTTP_PORT", "9090")
	os.Setenv("TENANT_HEADER", "X-My-Tenant-ID")
	os.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, 9090, cfg.HTTPPort)
	assert.Equal(t, "X-My-Tenant-ID", cfg.TenantHeader)
	assert.Equal(t, "debug", cfg.LogLevel)

	// Test with default values
	os.Unsetenv("APP_ENV")
	os.Unsetenv("HTTP_PORT")
	os.Unsetenv("TENANT_HEADER")
	os.Unsetenv("LOG_LEVEL")

	cfg, err = Load()
	assert.NoError(t, err)
	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, 8080, cfg.HTTPPort)
	assert.Equal(t, "X-Tenant-ID", cfg.TenantHeader)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestAddress(t *testing.T) {
	cfg := &Config{HTTPPort: 8888}
	assert.Equal(t, ":8888", cfg.Address())
}
