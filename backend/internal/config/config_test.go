package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var configEnvKeys = []string{
	"APP_ENV",
	"SERVICE_NAME",
	"HTTP_PORT",
	"TENANT_HEADER",
	"LOG_LEVEL",
	"DATABASE_URL",
	"DB_MAX_IDLE_CONNS",
	"DB_MAX_OPEN_CONNS",
	"DB_CONN_MAX_LIFETIME",
	"JWT_ACCESS_SECRET",
	"JWT_REFRESH_SECRET",
	"JWT_ACCESS_TTL",
	"JWT_REFRESH_TTL",
	"BCRYPT_COST",
	"REFRESH_TOKEN_LENGTH",
	"OTEL_ENABLED",
	"OTEL_EXPORTER_OTLP_ENDPOINT",
	"OTEL_EXPORTER_OTLP_HEADERS",
	"OTEL_EXPORTER_OTLP_INSECURE",
	"METRICS_ROUTE",
}

func unsetConfigEnv() {
	for _, key := range configEnvKeys {
		os.Unsetenv(key)
	}
}

func TestLoad(t *testing.T) {
	unsetConfigEnv()
	// Test with environment variables set
	os.Setenv("APP_ENV", "production")
	os.Setenv("SERVICE_NAME", "gestao-prod")
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
	os.Setenv("OTEL_ENABLED", "true")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://telemetry.updev.local:4318")
	os.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "api-key=secret,tenant=main")
	os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "false")
	os.Setenv("METRICS_ROUTE", "/v1/metrics")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, "gestao-prod", cfg.ServiceName)
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
	assert.True(t, cfg.TelemetryEnabled)
	assert.Equal(t, "https://telemetry.updev.local:4318", cfg.OTLPEndpoint)
	assert.Equal(t, "api-key=secret,tenant=main", cfg.OTLPHeaders)
	assert.False(t, cfg.OTLPInsecure)
	assert.Equal(t, "/v1/metrics", cfg.MetricsRoute)

	// Test with default values
	unsetConfigEnv()

	cfg, err = Load()
	assert.NoError(t, err)
	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, "gestao-api", cfg.ServiceName)
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
	assert.False(t, cfg.TelemetryEnabled)
	assert.Equal(t, "", cfg.OTLPEndpoint)
	assert.Equal(t, "", cfg.OTLPHeaders)
	assert.False(t, cfg.OTLPInsecure)
	assert.Equal(t, "/metrics", cfg.MetricsRoute)
}

func TestAddress(t *testing.T) {
	cfg := &Config{HTTPPort: 8888}
	assert.Equal(t, ":8888", cfg.Address())
}

func TestLoadProductionValidations(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr string
	}{
		{
			name: "missing database url",
			setup: func() {
				unsetConfigEnv()
				os.Setenv("APP_ENV", "production")
				os.Setenv("JWT_ACCESS_SECRET", "prod-access")
				os.Setenv("JWT_REFRESH_SECRET", "prod-refresh")
			},
			wantErr: "DATABASE_URL must be defined in production",
		},
		{
			name: "default database url",
			setup: func() {
				unsetConfigEnv()
				os.Setenv("APP_ENV", "production")
				os.Setenv("DATABASE_URL", defaultDatabaseURL)
				os.Setenv("JWT_ACCESS_SECRET", "prod-access")
				os.Setenv("JWT_REFRESH_SECRET", "prod-refresh")
			},
			wantErr: "DATABASE_URL must not use the development default in production",
		},
		{
			name: "default access secret",
			setup: func() {
				unsetConfigEnv()
				os.Setenv("APP_ENV", "production")
				os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/prod?sslmode=disable")
				os.Setenv("JWT_ACCESS_SECRET", defaultJWTAccessSecret)
				os.Setenv("JWT_REFRESH_SECRET", "prod-refresh")
			},
			wantErr: "JWT_ACCESS_SECRET must not use the development default in production",
		},
		{
			name: "missing refresh secret",
			setup: func() {
				unsetConfigEnv()
				os.Setenv("APP_ENV", "production")
				os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/prod?sslmode=disable")
				os.Setenv("JWT_ACCESS_SECRET", "prod-access")
			},
			wantErr: "JWT_REFRESH_SECRET must be defined in production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := Load()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
