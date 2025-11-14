package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"

	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/server"
	"github.com/kusmin/gestao_updev/backend/pkg/database"
	"github.com/kusmin/gestao_updev/backend/pkg/logger"
	"github.com/kusmin/gestao_updev/backend/pkg/telemetry"
)

// @title Gestão UpDev API
// @version 1.0
// @description API para a plataforma de gestão de comércios locais.
// @contact.name API Support
// @contact.email support@updev.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Token JWT no formato `Bearer <token>`
// @securityDefinitions.apikey TenantHeader
// @in header
// @name X-Tenant-ID
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	zapLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer zapLogger.Sync()

	if sentryDSN := os.Getenv("SENTRY_DSN_BACKEND"); sentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			TracesSampleRate: 1.0,
			EnableTracing:    true,
		}); err != nil {
			zapLogger.Fatal("sentry.Init failed", zap.Error(err))
		}
		defer sentry.Flush(2 * time.Second)
		zapLogger.Info("Sentry initialized for backend")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tel, err := telemetry.Init(ctx, telemetry.Config{
		ServiceName:   cfg.ServiceName,
		Environment:   cfg.AppEnv,
		EnableTracing: cfg.TelemetryEnabled,
		OTLPEndpoint:  cfg.OTLPEndpoint,
		OTLPHeaders:   cfg.OTLPHeaders,
		OTLPInsecure:  cfg.OTLPInsecure,
	})
	if err != nil {
		zapLogger.Fatal("failed to initialize telemetry", zap.Error(err))
	}
	defer func() {
		if err := tel.Shutdown(context.Background()); err != nil {
			zapLogger.Warn("error shutting down telemetry", zap.Error(err))
		}
	}()

	db, err := database.New(database.Config{
		URL:             cfg.DatabaseURL,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		MaxOpenConns:    cfg.DBMaxOpenConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
		LogMode:         mapGormLogLevel(cfg.LogLevel),
	})
	if err != nil {
		zapLogger.Fatal("failed to initialize database", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		zapLogger.Fatal("failed to expose database handle", zap.Error(err))
	}
	defer sqlDB.Close()

	srv := server.New(cfg, zapLogger, db, tel)

	if sentry.CurrentHub().Client() != nil {
		srv.Router.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
		defer func() {
			if r := recover(); r != nil {
				sentry.CurrentHub().Recover(r)
				sentry.Flush(2 * time.Second)
				panic(r)
			}
		}()
	}

	if err := srv.Run(ctx); err != nil {
		zapLogger.Fatal("server stopped with error", zap.Any("error", err))
	}
}

func mapGormLogLevel(level string) gormLogger.LogLevel {
	switch level {
	case "debug":
		return gormLogger.Info
	case "warn":
		return gormLogger.Warn
	case "error":
		return gormLogger.Error
	default:
		return gormLogger.Silent
	}
}
