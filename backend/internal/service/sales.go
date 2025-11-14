package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// SalesOrderFilter filtros de listagem.
type SalesOrderFilter struct {
	Status   string
	ClientID *uuid.UUID
	Date     *time.Time
}

// SalesItemInput representa itens da venda.
type SalesItemInput struct {
	Type      string
	RefID     uuid.UUID
	Quantity  int
	UnitPrice float64
}

// SalesOrderInput payload de criação.
type SalesOrderInput struct {
	ClientID  uuid.UUID
	BookingID *uuid.UUID
	Items     []SalesItemInput
	Discount  float64
	Notes     string
}

// SalesOrderUpdateInput campos editáveis.
type SalesOrderUpdateInput struct {
	Status *string
	Notes  *string
}

// PaymentInput dados para registrar pagamentos.
type PaymentInput struct {
	Method  string
	Amount  float64
	PaidAt  time.Time
	Details map[string]interface{}
}

// PaymentFilter filtros de listagem de pagamentos.
type PaymentFilter struct {
	Method    string
	StartDate *time.Time
	EndDate   *time.Time
}

func (s *Service) ListSalesOrders(ctx context.Context, tenantID uuid.UUID, filter SalesOrderFilter) ([]domain.SalesOrder, error) {
	query := s.dbWithContext(ctx).
		Preload("Items").
		Where("tenant_id = ?", tenantID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ClientID != nil {
		query = query.Where("client_id = ?", *filter.ClientID)
	}
	if filter.Date != nil {
		start := filter.Date.Truncate(24 * time.Hour)
		end := start.Add(24 * time.Hour)
		query = query.Where("created_at >= ? AND created_at < ?", start, end)
	}

	var orders []domain.SalesOrder
	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Service) CreateSalesOrder(ctx context.Context, tenantID uuid.UUID, input SalesOrderInput) (*domain.SalesOrder, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("ao menos um item é obrigatório")
	}
	if err := s.ensureTenantRecord(ctx, &domain.Client{}, tenantID, input.ClientID); err != nil {
		return nil, err
	}
	if input.BookingID != nil {
		if err := s.ensureTenantRecord(ctx, &domain.Booking{}, tenantID, *input.BookingID); err != nil {
			return nil, err
		}
	}
	if err := s.ensureSalesItems(ctx, tenantID, input.Items); err != nil {
		return nil, err
	}

	order := &domain.SalesOrder{
		TenantModel: domain.TenantModel{
			TenantID: tenantID,
		},
		ClientID: input.ClientID,
		Status:   domain.SalesOrderStatusDraft,
		Discount: input.Discount,
		Notes:    input.Notes,
	}
	if input.BookingID != nil {
		order.BookingID = input.BookingID
	}

	err := s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		var total float64
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for _, item := range input.Items {
			if item.Quantity <= 0 {
				return errors.New("quantidade inválida em item")
			}
			salesItem := domain.SalesItem{
				TenantModel: domain.TenantModel{
					TenantID: tenantID,
				},
				OrderID:   order.ID,
				ItemType:  item.Type,
				ItemRefID: item.RefID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			}
			total += float64(item.Quantity) * item.UnitPrice
			if err := tx.Create(&salesItem).Error; err != nil {
				return err
			}
		}

		order.Total = total - input.Discount
		if order.Total < 0 {
			order.Total = 0
		}
		return tx.Model(order).Updates(map[string]interface{}{"total": order.Total}).Error
	})
	if err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		Preload("Items").
		First(order, "id = ?", order.ID).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) UpdateSalesOrder(ctx context.Context, tenantID, orderID uuid.UUID, input SalesOrderUpdateInput) (*domain.SalesOrder, error) {
	var order domain.SalesOrder
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, orderID).
		Preload("Items").
		First(&order).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}

	if len(updates) == 0 {
		return &order, nil
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.SalesOrder{}).
		Where("tenant_id = ? AND id = ?", tenantID, orderID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, orderID).
		Preload("Items").
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *Service) AddPayment(ctx context.Context, tenantID, orderID uuid.UUID, input PaymentInput) (*domain.Payment, error) {
	if err := s.ensureTenantRecord(ctx, &domain.SalesOrder{}, tenantID, orderID); err != nil {
		return nil, err
	}
	payment := &domain.Payment{
		TenantModel: domain.TenantModel{
			TenantID: tenantID,
		},
		OrderID: orderID,
		Method:  input.Method,
		Amount:  input.Amount,
		PaidAt:  input.PaidAt,
		Details: datatypes.JSONMap(input.Details),
	}
	if payment.Details == nil {
		payment.Details = datatypes.JSONMap{}
	}

	if err := s.dbWithContext(ctx).Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service) ListPayments(ctx context.Context, tenantID uuid.UUID, filter PaymentFilter) ([]domain.Payment, error) {
	query := s.dbWithContext(ctx).
		Where("tenant_id = ?", tenantID)

	if filter.Method != "" {
		query = query.Where("method = ?", filter.Method)
	}
	if filter.StartDate != nil {
		query = query.Where("paid_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("paid_at <= ?", *filter.EndDate)
	}

	var payments []domain.Payment
	if err := query.Order("paid_at DESC").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (s *Service) ensureSalesItems(ctx context.Context, tenantID uuid.UUID, items []SalesItemInput) error {
	for _, item := range items {
		switch item.Type {
		case "service":
			if err := s.ensureTenantRecord(ctx, &domain.Service{}, tenantID, item.RefID); err != nil {
				return err
			}
		case "product":
			if err := s.ensureTenantRecord(ctx, &domain.Product{}, tenantID, item.RefID); err != nil {
				return err
			}
		default:
			return fmt.Errorf("tipo de item %q não suportado", item.Type)
		}
	}
	return nil
}

func (s *Service) ListAllSalesOrders(ctx context.Context) ([]domain.SalesOrder, error) {
	var orders []domain.SalesOrder
	if err := s.dbWithContext(ctx).
		Preload("Items").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

type AdminSalesOrderInput struct {
	SalesOrderInput
	TenantID uuid.UUID
}

func (s *Service) AdminCreateSalesOrder(ctx context.Context, input AdminSalesOrderInput) (*domain.SalesOrder, error) {
	return s.CreateSalesOrder(ctx, input.TenantID, input.SalesOrderInput)
}

func (s *Service) AdminUpdateSalesOrder(ctx context.Context, orderID uuid.UUID, input SalesOrderUpdateInput) (*domain.SalesOrder, error) {
	var order domain.SalesOrder
	if err := s.dbWithContext(ctx).First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return s.UpdateSalesOrder(ctx, order.TenantID, orderID, input)
}

func (s *Service) AdminDeleteSalesOrder(ctx context.Context, orderID uuid.UUID) error {
	return s.dbWithContext(ctx).Delete(&domain.SalesOrder{}, "id = ?", orderID).Error
}
