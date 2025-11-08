package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Test with environment variables set
	os.Setenv("APP_ENV", "production")
	os.Setenv("HTTP_PORT", "9090")
	os.Setenv("TENANT_HEADER", "X-My-Tenant-ID")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db?sslmode=disable")
	os.Setenv("DB_MAX_IDLE_CONNS", "11")
	os.Setenv("DB_MAX_OPEN_CONNS", "33")
	os.Setenv("DB_CONN_MAX_LIFETIME", "30m")
	os.Setenv("JWT_ACCESS_SECRET", "access-secret")
	os.Setenv("JWT_REFRESH_SECRET", "refresh-secret")
	os.Setenv("JWT_ACCESS_TTL", "30m")
	os.Setenv("JWT_REFRESH_TTL", "168h")
	os.Setenv("BCRYPT_COST", "14")
	os.Setenv("REFRESH_TOKEN_LENGTH", "96")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, 9090, cfg.HTTPPort)
	assert.Equal(t, "X-My-Tenant-ID", cfg.TenantHeader)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db?sslmode=disable", cfg.DatabaseURL)
	assert.Equal(t, 11, cfg.DBMaxIdleConns)
	assert.Equal(t, 33, cfg.DBMaxOpenConns)
	assert.Equal(t, 30*time.Minute, cfg.DBConnMaxLifetime)
	assert.Equal(t, "access-secret", cfg.JWTAccessSecret)
	assert.Equal(t, "refresh-secret", cfg.JWTRefreshSecret)
	assert.Equal(t, 30*time.Minute, cfg.JWTAccessTTL)
	assert.Equal(t, 168*time.Hour, cfg.JWTRefreshTTL)
	assert.Equal(t, 14, cfg.BcryptCost)
	assert.Equal(t, 96, cfg.RefreshTokenLength)

	// Test with default values
	os.Unsetenv("APP_ENV")
	os.Unsetenv("HTTP_PORT")
	os.Unsetenv("TENANT_HEADER")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("DB_MAX_IDLE_CONNS")
	os.Unsetenv("DB_MAX_OPEN_CONNS")
	os.Unsetenv("DB_CONN_MAX_LIFETIME")
	os.Unsetenv("JWT_ACCESS_SECRET")
	os.Unsetenv("JWT_REFRESH_SECRET")
	os.Unsetenv("JWT_ACCESS_TTL")
	os.Unsetenv("JWT_REFRESH_TTL")
	os.Unsetenv("BCRYPT_COST")
	os.Unsetenv("REFRESH_TOKEN_LENGTH")

	cfg, err = Load()
	assert.NoError(t, err)
	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, 8080, cfg.HTTPPort)
	assert.Equal(t, "X-Tenant-ID", cfg.TenantHeader)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/gestao_updev?sslmode=disable", cfg.DatabaseURL)
	assert.Equal(t, 5, cfg.DBMaxIdleConns)
	assert.Equal(t, 20, cfg.DBMaxOpenConns)
	assert.Equal(t, time.Hour, cfg.DBConnMaxLifetime)
	assert.Equal(t, "dev-access-secret", cfg.JWTAccessSecret)
	assert.Equal(t, "dev-refresh-secret", cfg.JWTRefreshSecret)
	assert.Equal(t, 15*time.Minute, cfg.JWTAccessTTL)
	assert.Equal(t, 720*time.Hour, cfg.JWTRefreshTTL)
	assert.Equal(t, 12, cfg.BcryptCost)
	assert.Equal(t, 64, cfg.RefreshTokenLength)
}

func TestAddress(t *testing.T) {
	cfg := &Config{HTTPPort: 8888}
	assert.Equal(t, ":8888", cfg.Address())
}
