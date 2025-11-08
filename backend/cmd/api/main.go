package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/server"
	"github.com/kusmin/gestao_updev/backend/pkg/logger"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(cfg, zapLogger)

	if err := srv.Run(ctx); err != nil {
		zapLogger.Fatal("server stopped with error", zap.Any("error", err))
	}
}
