package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/kusmin/gestao_updev/backend/docs"
)

// Server representa a API HTTP.
type Server struct {
	cfg    *config.Config
	logger *zap.Logger
	engine *gin.Engine
}

// New cria uma instância do servidor HTTP.
func New(cfg *config.Config, logger *zap.Logger) *Server {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestID())
	engine.Use(gin.Logger())

	api := engine.Group("/v1")
	api.Use(middleware.TenantEnforcer(cfg.TenantHeader))
	registerRoutes(api)

	engine.GET("/v1/healthz", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{
			"status": "ok",
			"env":    cfg.AppEnv,
		}, nil)
	})

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{
		cfg:    cfg,
		logger: logger,
		engine: engine,
	}
}

// Run inicia o servidor HTTP e bloqueia até receber sinal de desligamento.
func (s *Server) Run(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:    s.cfg.Address(),
		Handler: s.engine,
	}

	errCh := make(chan error, 1)

	go func() {
		s.logger.Info("HTTP server starting", zap.String("addr", s.cfg.Address()))
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func registerRoutes(api *gin.RouterGroup) {
	placeholderRoutes := []struct {
		Method string
		Path   string
		Name   string
	}{
		{"POST", "/auth/signup", "auth.signup"},
		{"POST", "/auth/login", "auth.login"},
		{"POST", "/auth/refresh", "auth.refresh"},
		{"GET", "/companies/me", "companies.current.get"},
		{"PUT", "/companies/me", "companies.current.update"},
		{"GET", "/users", "users.list"},
		{"POST", "/users", "users.create"},
		{"PATCH", "/users/:id", "users.update"},
		{"DELETE", "/users/:id", "users.delete"},
		{"GET", "/clients", "clients.list"},
		{"POST", "/clients", "clients.create"},
		{"GET", "/clients/:id", "clients.get"},
		{"PUT", "/clients/:id", "clients.update"},
		{"DELETE", "/clients/:id", "clients.delete"},
		{"GET", "/professionals", "professionals.list"},
		{"GET", "/services", "services.list"},
		{"POST", "/services", "services.create"},
		{"PUT", "/services/:id", "services.update"},
		{"DELETE", "/services/:id", "services.delete"},
		{"GET", "/products", "products.list"},
		{"POST", "/products", "products.create"},
		{"PUT", "/products/:id", "products.update"},
		{"DELETE", "/products/:id", "products.delete"},
		{"GET", "/inventory/movements", "inventory.list"},
		{"POST", "/inventory/movements", "inventory.create"},
		{"GET", "/bookings", "bookings.list"},
		{"POST", "/bookings", "bookings.create"},
		{"PATCH", "/bookings/:id", "bookings.update"},
		{"POST", "/bookings/:id/cancel", "bookings.cancel"},
		{"GET", "/sales/orders", "sales.list"},
		{"POST", "/sales/orders", "sales.create"},
		{"PATCH", "/sales/orders/:id", "sales.update"},
		{"POST", "/sales/orders/:id/payments", "sales.payments.create"},
		{"GET", "/payments", "payments.list"},
		{"GET", "/dashboard/daily", "dashboard.daily"},
	}

	for _, route := range placeholderRoutes {
		handler := placeholder(route.Name)
		switch route.Method {
		case http.MethodGet:
			api.GET(route.Path, handler)
		case http.MethodPost:
			api.POST(route.Path, handler)
		case http.MethodPatch:
			api.PATCH(route.Path, handler)
		case http.MethodPut:
			api.PUT(route.Path, handler)
		case http.MethodDelete:
			api.DELETE(route.Path, handler)
		default:
			panic(fmt.Sprintf("unsupported method %s for %s", route.Method, route.Path))
		}
	}
}

func placeholder(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Error(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", fmt.Sprintf("Endpoint %s em desenvolvimento", name), nil)
	}
}
