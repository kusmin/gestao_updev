package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// InventoryFilter para listagem.
type InventoryFilter struct {
	ProductID *uuid.UUID
	Type      string
	StartDate *time.Time
	EndDate   *time.Time
}

// InventoryInput para criação manual.
type InventoryInput struct {
	ProductID uuid.UUID
	OrderID   *uuid.UUID
	Type      string
	Quantity  int
	Reason    string
}

var ErrInvalidInventoryType = errors.New("tipo de movimentação inválido")

func (s *Service) ListInventoryMovements(ctx context.Context, tenantID uuid.UUID, filter InventoryFilter) ([]domain.InventoryMovement, error) {
	query := s.dbWithContext(ctx).
		Where("tenant_id = ?", tenantID)

	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	var movements []domain.InventoryMovement
	if err := query.Order("created_at DESC").Find(&movements).Error; err != nil {
		return nil, err
	}
	return movements, nil
}

func (s *Service) CreateInventoryMovement(ctx context.Context, tenantID uuid.UUID, input InventoryInput) (*domain.InventoryMovement, error) {
	if input.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	switch input.Type {
	case domain.InventoryMovementIn, domain.InventoryMovementOut, domain.InventoryMovementAdjustment:
	default:
		return nil, ErrInvalidInventoryType
	}

	movement := &domain.InventoryMovement{
		TenantModel: domain.TenantModel{
			TenantID: tenantID,
		},
		ProductID: input.ProductID,
		Type:      input.Type,
		Quantity:  input.Quantity,
		Reason:    input.Reason,
	}
	if input.OrderID != nil {
		movement.OrderID = input.OrderID
	}

	if err := s.dbWithContext(ctx).Create(movement).Error; err != nil {
		return nil, err
	}
	return movement, nil
}
