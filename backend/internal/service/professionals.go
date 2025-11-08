package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// ListProfessionals retorna profissionais ativos.
func (s *Service) ListProfessionals(ctx context.Context, tenantID uuid.UUID) ([]domain.Professional, error) {
	var professionals []domain.Professional
	if err := s.dbWithContext(ctx).
		Preload("Availability").
		Where("tenant_id = ? AND active = true", tenantID).
		Order("name ASC").
		Find(&professionals).Error; err != nil {
		return nil, err
	}
	return professionals, nil
}
