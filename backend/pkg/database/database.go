package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config encapsula parâmetros do pool de conexão.
type Config struct {
	URL             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogMode         logger.LogLevel
}

// New cria uma instância do GORM conectada ao PostgreSQL.
func New(cfg Config) (*gorm.DB, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("database url is required")
	}

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogMode),
	}

	db, err := gorm.Open(postgres.Open(cfg.URL), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("extract sql.DB: %w", err)
	}

	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	return db, nil
}
