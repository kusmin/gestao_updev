package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"gorm.io/gorm"
)

type SalesRepository struct {
	db *gorm.DB
}

func NewSalesRepository(db *gorm.DB) *SalesRepository {
	return &SalesRepository{db: db}
}

func (r *SalesRepository) ListAll(ctx context.Context) ([]domain.SalesOrder, error) {
	var salesOrders []domain.SalesOrder
	if err := r.db.WithContext(ctx).Find(&salesOrders).Error; err != nil {
		return nil, err
	}
	return salesOrders, nil
}

func (r *SalesRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	var salesOrder domain.SalesOrder
	if err := r.db.WithContext(ctx).First(&salesOrder, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &salesOrder, nil
}

func (r *SalesRepository) Create(ctx context.Context, salesOrder *domain.SalesOrder) error {
	return r.db.WithContext(ctx).Create(salesOrder).Error
}

func (r *SalesRepository) Update(ctx context.Context, salesOrder *domain.SalesOrder) error {
	return r.db.WithContext(ctx).Save(salesOrder).Error
}

func (r *SalesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.SalesOrder{}, id).Error
}
