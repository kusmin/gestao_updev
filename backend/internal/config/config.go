package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

const (
	defaultDatabaseURL      = "postgres://postgres:postgres@localhost:5432/gestao_updev?sslmode=disable"
	defaultJWTAccessSecret  = "dev-access-secret"
	defaultJWTRefreshSecret = "dev-refresh-secret"
)

// Config centraliza parâmetros de execução do backend.
type Config struct {
	AppEnv             string        `env:"APP_ENV" envDefault:"development"`
	ServiceName        string        `env:"SERVICE_NAME" envDefault:"gestao-api"`
	HTTPPort           int           `env:"HTTP_PORT" envDefault:"8080"`
	TenantHeader       string        `env:"TENANT_HEADER" envDefault:"X-Tenant-ID"`
	LogLevel           string        `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL        string        `env:"DATABASE_URL"`
	DBMaxIdleConns     int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	DBMaxOpenConns     int           `env:"DB_MAX_OPEN_CONNS" envDefault:"20"`
	DBConnMaxLifetime  time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1h"`
	JWTAccessSecret    string        `env:"JWT_ACCESS_SECRET"`
	JWTRefreshSecret   string        `env:"JWT_REFRESH_SECRET"`
	JWTAccessTTL       time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	JWTRefreshTTL      time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h"`
	BcryptCost         int           `env:"BCRYPT_COST" envDefault:"12"`
	RefreshTokenLength int           `env:"REFRESH_TOKEN_LENGTH" envDefault:"64"`
	TelemetryEnabled   bool          `env:"OTEL_ENABLED" envDefault:"false"`
	OTLPEndpoint       string        `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTLPHeaders        string        `env:"OTEL_EXPORTER_OTLP_HEADERS"`
	OTLPInsecure       bool          `env:"OTEL_EXPORTER_OTLP_INSECURE" envDefault:"false"`
	MetricsRoute       string        `env:"METRICS_ROUTE" envDefault:"/metrics"`
}

// Load lê variáveis de ambiente e monta a configuração.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}
	if err := cfg.ensureProtectedValues(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Address retorna host:port para o servidor HTTP.
func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.HTTPPort)
}

func (c *Config) ensureProtectedValues() error {
	isProd := strings.EqualFold(c.AppEnv, "production")

	switch {
	case c.DatabaseURL == "":
		if isProd {
			return fmt.Errorf("DATABASE_URL must be defined in production")
		}
		c.DatabaseURL = defaultDatabaseURL
	case isProd && c.DatabaseURL == defaultDatabaseURL:
		return fmt.Errorf("DATABASE_URL must not use the development default in production")
	}

	switch {
	case c.JWTAccessSecret == "":
		if isProd {
			return fmt.Errorf("JWT_ACCESS_SECRET must be defined in production")
		}
		c.JWTAccessSecret = defaultJWTAccessSecret
	case isProd && c.JWTAccessSecret == defaultJWTAccessSecret:
		return fmt.Errorf("JWT_ACCESS_SECRET must not use the development default in production")
	}

	switch {
	case c.JWTRefreshSecret == "":
		if isProd {
			return fmt.Errorf("JWT_REFRESH_SECRET must be defined in production")
		}
		c.JWTRefreshSecret = defaultJWTRefreshSecret
	case isProd && c.JWTRefreshSecret == defaultJWTRefreshSecret:
		return fmt.Errorf("JWT_REFRESH_SECRET must not use the development default in production")
	}

	return nil
}
