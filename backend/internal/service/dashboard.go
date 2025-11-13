package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// DashboardDailyDTO estrutura retorno do endpoint.
type DashboardDailyDTO struct {
	Date        string       `json:"date"`
	Bookings    int64        `json:"bookings"`
	Completed   int64        `json:"completed"`
	Revenue     float64      `json:"revenue"`
	TopServices []TopService `json:"top_services"`
}

// TopService agrega ranking diário.
type TopService struct {
	ServiceID uuid.UUID `json:"service_id"`
	Name      string    `json:"name"`
	Quantity  int64     `json:"quantity"`
}

func (s *Service) DashboardDaily(ctx context.Context, tenantID uuid.UUID, date time.Time, professionalID *uuid.UUID) (*DashboardDailyDTO, error) {
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	bookingQuery := s.dbWithContext(ctx).Model(&domain.Booking{}).
		Where("tenant_id = ? AND start_at >= ? AND start_at < ?", tenantID, start, end)
	if professionalID != nil {
		bookingQuery = bookingQuery.Where("professional_id = ?", *professionalID)
	}

	var total, completed int64
	if err := bookingQuery.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := bookingQuery.Where("status = ?", domain.BookingStatusDone).Count(&completed).Error; err != nil {
		return nil, err
	}

	paymentQuery := s.dbWithContext(ctx).Model(&domain.Payment{}).
		Where("tenant_id = ? AND paid_at >= ? AND paid_at < ?", tenantID, start, end)
	var rev struct {
		Total float64
	}
	if err := paymentQuery.Select("COALESCE(SUM(amount),0) as total").Scan(&rev).Error; err != nil {
		return nil, err
	}
	revenue := rev.Total

	type result struct {
		ServiceID uuid.UUID
		Name      string
		Quantity  int64
	}
	var top []result
	err := s.dbWithContext(ctx).
		Table("sales_items").
		Select("sales_items.item_ref_id as service_id, services.name, SUM(sales_items.quantity) as quantity").
		Joins("JOIN sales_orders ON sales_orders.id = sales_items.order_id").
		Joins("JOIN services ON services.id = sales_items.item_ref_id").
		Where("sales_items.tenant_id = ? AND sales_items.item_type = ? AND sales_orders.created_at >= ? AND sales_orders.created_at < ?", tenantID, "service", start, end).
		Group("sales_items.item_ref_id, services.name").
		Order("quantity DESC").
		Limit(5).
		Scan(&top).Error
	if err != nil {
		return nil, err
	}

	topServices := make([]TopService, len(top))
	for i, item := range top {
		topServices[i] = TopService{
			ServiceID: item.ServiceID,
			Name:      item.Name,
			Quantity:  item.Quantity,
		}
	}

	return &DashboardDailyDTO{
		Date:        start.Format("2006-01-02"),
		Bookings:    total,
		Completed:   completed,
		Revenue:     revenue,
		TopServices: topServices,
	}, nil
}
