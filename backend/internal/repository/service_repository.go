package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) ListAll(ctx context.Context) ([]domain.Service, error) {
	var services []domain.Service
	if err := r.db.WithContext(ctx).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (r *ServiceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	var service domain.Service
	if err := r.db.WithContext(ctx).First(&service, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *ServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	return r.db.WithContext(ctx).Save(service).Error
}

func (r *ServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Service{}, id).Error
}
