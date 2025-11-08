package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// CompanyUpdateInput representa campos editáveis da empresa.
type CompanyUpdateInput struct {
	Name     *string
	Timezone *string
	Phone    *string
	Email    *string
	Settings map[string]interface{}
}

// GetCompany retorna dados da empresa corrente (tenant).
func (s *Service) GetCompany(ctx context.Context, tenantID uuid.UUID) (*domain.Company, error) {
	var company domain.Company
	if err := s.dbWithContext(ctx).First(&company, "id = ?", tenantID).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

// UpdateCompany aplica mudanças parciais nos dados da empresa.
func (s *Service) UpdateCompany(ctx context.Context, tenantID uuid.UUID, input CompanyUpdateInput) (*domain.Company, error) {
	var company domain.Company
	if err := s.dbWithContext(ctx).First(&company, "id = ?", tenantID).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Timezone != nil {
		updates["timezone"] = *input.Timezone
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Settings != nil {
		updates["settings"] = input.Settings
	}

	if len(updates) == 0 {
		return &company, nil
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Company{}).
		Where("id = ?", tenantID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).First(&company, "id = ?", tenantID).Error; err != nil {
		return nil, err
	}

	return &company, nil
}
