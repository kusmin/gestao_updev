package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config centraliza parâmetros de execução do backend.
type Config struct {
	AppEnv       string `env:"APP_ENV" envDefault:"development"`
	HTTPPort     int    `env:"HTTP_PORT" envDefault:"8080"`
	TenantHeader string `env:"TENANT_HEADER" envDefault:"X-Tenant-ID"`
	LogLevel     string `env:"LOG_LEVEL" envDefault:"info"`
	DBURL        string `env:"DATABASE_URL,required"`
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
