package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) ListAll(ctx context.Context) ([]domain.Company, error) {
	var companies []domain.Company
	if err := r.db.WithContext(ctx).Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *CompanyRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	var company domain.Company
	if err := r.db.WithContext(ctx).First(&company, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}

func (r *CompanyRepository) Update(ctx context.Context, company *domain.Company) error {
	return r.db.WithContext(ctx).Save(company).Error
}

func (r *CompanyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Company{}, id).Error
}
