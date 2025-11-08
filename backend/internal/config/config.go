package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config centraliza parâmetros de execução do backend.
type Config struct {
	AppEnv             string        `env:"APP_ENV" envDefault:"development"`
	HTTPPort           int           `env:"HTTP_PORT" envDefault:"8080"`
	TenantHeader       string        `env:"TENANT_HEADER" envDefault:"X-Tenant-ID"`
	LogLevel           string        `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL        string        `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/gestao_updev?sslmode=disable"`
	DBMaxIdleConns     int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	DBMaxOpenConns     int           `env:"DB_MAX_OPEN_CONNS" envDefault:"20"`
	DBConnMaxLifetime  time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1h"`
	JWTAccessSecret    string        `env:"JWT_ACCESS_SECRET" envDefault:"dev-access-secret"`
	JWTRefreshSecret   string        `env:"JWT_REFRESH_SECRET" envDefault:"dev-refresh-secret"`
	JWTAccessTTL       time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	JWTRefreshTTL      time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h"`
	BcryptCost         int           `env:"BCRYPT_COST" envDefault:"12"`
	RefreshTokenLength int           `env:"REFRESH_TOKEN_LENGTH" envDefault:"64"`
}

// Load lê variáveis de ambiente e monta a configuração.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}
	return cfg, nil
}

// Address retorna host:port para o servidor HTTP.
func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.HTTPPort)
}
