package service

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/auth"
	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/repository"
)

// ErrInvalidCredentials sinaliza falha de autenticação.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Service agrega casos de uso de negócio.
type Service struct {
	repo   *repository.Repository
	jwt    *auth.JWTManager
	cfg    *config.Config
	logger *zap.Logger
}

// New instancia o service layer.
func New(cfg *config.Config, repo *repository.Repository, jwt *auth.JWTManager, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		jwt:    jwt,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Service) dbWithContext(ctx context.Context) *gorm.DB {
	return s.repo.DB().WithContext(ctx)
}

func (s *Service) sanitizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func (s *Service) clampPagination(page, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}
