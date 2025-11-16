package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	_ "github.com/kusmin/gestao_updev/backend/docs"
	"github.com/kusmin/gestao_updev/backend/internal/auth"
	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/http/handler"
	"github.com/kusmin/gestao_updev/backend/internal/http/response"
	"github.com/kusmin/gestao_updev/backend/internal/middleware"
	"github.com/kusmin/gestao_updev/backend/internal/repository"
	"github.com/kusmin/gestao_updev/backend/internal/service"
	"github.com/kusmin/gestao_updev/backend/pkg/telemetry"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server representa a API HTTP.
type Server struct {
	cfg       *config.Config
	logger    *zap.Logger
	engine    *gin.Engine
	db        *gorm.DB
	telemetry *telemetry.Telemetry
}

// New cria uma instância do servidor HTTP.
func New(cfg *config.Config, logger *zap.Logger, db *gorm.DB, telem *telemetry.Telemetry) *Server {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestID())
	if telem != nil && telem.TracerProvider() != nil {
		engine.Use(otelgin.Middleware(cfg.ServiceName, otelgin.WithTracerProvider(telem.TracerProvider())))
	}
	engine.Use(middleware.Logger(logger))

	// Repositories
	repo := repository.New(db)
	companyRepo := repository.NewCompanyRepository(db)

	// Services
	jwtManager := auth.NewJWTManager(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	svc := service.New(cfg, repo, jwtManager, logger)
	companySvc := service.NewCompanyService(companyRepo)

	// Handlers
	apiHandler := handler.New(svc, logger)
	companyHandler := handler.NewCompanyHandler(companySvc)

	api := engine.Group("/v1")
	api.Use(middleware.TenantEnforcer(cfg.TenantHeader))
	registerRoutes(api, cfg, apiHandler, companyHandler, jwtManager)

	engine.GET("/v1/healthz", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{
			"status": "ok",
			"env":    cfg.AppEnv,
		}, nil)
	})

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if telem != nil && telem.MetricsHandler() != nil && cfg.MetricsRoute != "" {
		engine.GET(cfg.MetricsRoute, gin.WrapH(telem.MetricsHandler()))
	}

	return &Server{
		cfg:       cfg,
		logger:    logger,
		engine:    engine,
		db:        db,
		telemetry: telem,
	}
}

// Run inicia o servidor HTTP e bloqueia até receber sinal de desligamento.
func (s *Server) Run(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:    s.cfg.Address(),
		Handler: s.engine,
		ReadHeaderTimeout: 5 * time.Second, // Mitigação contra ataques Slowloris
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

// Router expõe a instância do gin.Engine para middlewares externos.
func (s *Server) Router() *gin.Engine {
	return s.engine
}

func registerRoutes(api *gin.RouterGroup, cfg *config.Config, h *handler.API, companyHandler *handler.CompanyHandler, jwtManager *auth.JWTManager) {
	authGroup := api.Group("/auth")
	authGroup.POST("/signup", h.Signup)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/refresh", h.RefreshToken)

	protected := api.Group("/")
	protected.Use(middleware.Auth(jwtManager, cfg.TenantHeader))

	protected.GET("/companies/me", h.GetCompany)
	protected.PUT("/companies/me", h.UpdateCompany)

	protected.GET("/users", h.ListUsers)
	protected.POST("/users", h.CreateUser)
	protected.GET("/users/:id", h.GetUser)
	protected.PATCH("/users/:id", h.UpdateUser)
	protected.DELETE("/users/:id", h.DeleteUser)

	protected.GET("/clients", h.ListClients)
	protected.POST("/clients", h.CreateClient)
	protected.GET("/clients/:id", h.GetClient)
	protected.PUT("/clients/:id", h.UpdateClient)
	protected.DELETE("/clients/:id", h.DeleteClient)

	protected.GET("/professionals", h.ListProfessionals)

	protected.GET("/services", h.ListServices)
	protected.POST("/services", h.CreateService)
	protected.GET("/services/:id", h.GetService)
	protected.PUT("/services/:id", h.UpdateService)
	protected.DELETE("/services/:id", h.DeleteService)

	protected.GET("/products", h.ListProducts)
	protected.POST("/products", h.CreateProduct)
	protected.GET("/products/:id", h.GetProduct)
	protected.PUT("/products/:id", h.UpdateProduct)
	protected.DELETE("/products/:id", h.DeleteProduct)

	protected.GET("/inventory/movements", h.ListInventoryMovements)
	protected.POST("/inventory/movements", h.CreateInventoryMovement)

	protected.GET("/bookings", h.ListBookings)
	protected.POST("/bookings", h.CreateBooking)
	protected.PATCH("/bookings/:id", h.UpdateBooking)
	protected.POST("/bookings/:id/cancel", h.CancelBooking)

	protected.GET("/sales/orders", h.ListSalesOrders)
	protected.POST("/sales/orders", h.CreateSalesOrder)
	protected.PATCH("/sales/orders/:id", h.UpdateSalesOrder)
	protected.POST("/sales/orders/:id/payments", h.CreatePayment)
	protected.GET("/payments", h.ListPayments)

	protected.GET("/dashboard/daily", h.DashboardDaily)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.Auth(jwtManager, cfg.TenantHeader), middleware.Admin())
	companyHandler.RegisterRoutes(admin)
	h.RegisterAdminUserRoutes(admin)
	h.RegisterAdminProductRoutes(admin)
	h.RegisterAdminServiceRoutes(admin)
	h.RegisterAdminClientRoutes(admin)
	h.RegisterAdminBookingRoutes(admin)
	h.RegisterAdminSalesRoutes(admin)
	h.RegisterAdminDashboardRoutes(admin)
}
