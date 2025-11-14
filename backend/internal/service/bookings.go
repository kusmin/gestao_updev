package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// BookingFilter define filtros básicos.
type BookingFilter struct {
	Date           *time.Time
	ProfessionalID *uuid.UUID
	Status         string
}

// BookingInput dados obrigatórios para criação.
type BookingInput struct {
	ClientID       uuid.UUID
	ProfessionalID uuid.UUID
	ServiceID      uuid.UUID
	Status         string
	StartAt        time.Time
	EndAt          *time.Time
	Notes          string
}

// BookingUpdateInput campos permitidos na edição.
type BookingUpdateInput struct {
	Status  *string
	StartAt *time.Time
	EndAt   *time.Time
	Notes   *string
}

var ErrBookingConflict = errors.New("já existe agendamento no horário selecionado")

func (s *Service) ListBookings(ctx context.Context, tenantID uuid.UUID, filter BookingFilter) ([]domain.Booking, error) {
	query := s.dbWithContext(ctx).
		Where("tenant_id = ?", tenantID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ProfessionalID != nil {
		query = query.Where("professional_id = ?", *filter.ProfessionalID)
	}
	if filter.Date != nil {
		start := filter.Date.Truncate(24 * time.Hour)
		end := start.Add(24 * time.Hour)
		query = query.Where("start_at >= ? AND start_at < ?", start, end)
	}

	var bookings []domain.Booking
	if err := query.Order("start_at ASC").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *Service) CreateBooking(ctx context.Context, tenantID uuid.UUID, input BookingInput) (*domain.Booking, error) {
	if err := s.ensureTenantRecord(ctx, &domain.Client{}, tenantID, input.ClientID); err != nil {
		return nil, err
	}
	if err := s.ensureTenantRecord(ctx, &domain.Professional{}, tenantID, input.ProfessionalID); err != nil {
		return nil, err
	}
	if input.Status == "" {
		input.Status = domain.BookingStatusPending
	}

	start := input.StartAt
	end := input.EndAt
	if end == nil || end.Before(start) {
		var service domain.Service
		if err := s.dbWithContext(ctx).First(&service, "tenant_id = ? AND id = ?", tenantID, input.ServiceID).Error; err != nil {
			return nil, err
		}
		calculated := start.Add(time.Duration(service.DurationMinutes) * time.Minute)
		end = &calculated
	}

	if err := s.checkBookingConflict(ctx, tenantID, input.ProfessionalID, start, *end, nil); err != nil {
		return nil, err
	}

	booking := &domain.Booking{
		TenantModel: domain.TenantModel{
			TenantID: tenantID,
		},
		ClientID:       input.ClientID,
		ProfessionalID: input.ProfessionalID,
		ServiceID:      input.ServiceID,
		Status:         input.Status,
		StartAt:        start,
		EndAt:          *end,
		Notes:          input.Notes,
		Metadata:       datatypes.JSONMap{},
	}

	if err := s.dbWithContext(ctx).Create(booking).Error; err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *Service) UpdateBooking(ctx context.Context, tenantID, bookingID uuid.UUID, input BookingUpdateInput) (*domain.Booking, error) {
	var booking domain.Booking
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, bookingID).
		First(&booking).Error; err != nil {
		return nil, err
	}

	if input.StartAt != nil || input.EndAt != nil {
		start := booking.StartAt
		end := booking.EndAt
		if input.StartAt != nil {
			start = *input.StartAt
		}
		if input.EndAt != nil {
			end = *input.EndAt
		}
		if err := s.checkBookingConflict(ctx, tenantID, booking.ProfessionalID, start, end, &booking.ID); err != nil {
			return nil, err
		}
	}

	updates := map[string]interface{}{}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.StartAt != nil {
		updates["start_at"] = *input.StartAt
	}
	if input.EndAt != nil {
		updates["end_at"] = *input.EndAt
	}
	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}

	if len(updates) == 0 {
		return &booking, nil
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Booking{}).
		Where("tenant_id = ? AND id = ?", tenantID, bookingID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, bookingID).
		First(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *Service) CancelBooking(ctx context.Context, tenantID, bookingID uuid.UUID, reason string) (*domain.Booking, error) {
	var booking domain.Booking
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, bookingID).
		First(&booking).Error; err != nil {
		return nil, err
	}

	metadata := booking.Metadata
	if metadata == nil {
		metadata = datatypes.JSONMap{}
	}
	if reason != "" {
		metadata["cancel_reason"] = reason
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Booking{}).
		Where("tenant_id = ? AND id = ?", tenantID, bookingID).
		Updates(map[string]interface{}{
			"status":   domain.BookingStatusCanceled,
			"metadata": metadata,
		}).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, bookingID).
		First(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *Service) checkBookingConflict(ctx context.Context, tenantID, professionalID uuid.UUID, start, end time.Time, ignoreID *uuid.UUID) error {
	query := s.dbWithContext(ctx).
		Model(&domain.Booking{}).
		Where("tenant_id = ? AND professional_id = ? AND status != ?", tenantID, professionalID, domain.BookingStatusCanceled).
		Where("start_at < ? AND end_at > ?", end, start)

	if ignoreID != nil {
		query = query.Where("id <> ?", *ignoreID)
	}

	var conflict domain.Booking
	if err := query.First(&conflict).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return ErrBookingConflict
}

func (s *Service) ListAllBookings(ctx context.Context, filter BookingFilter) ([]domain.Booking, error) {
	query := s.dbWithContext(ctx).Model(&domain.Booking{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ProfessionalID != nil {
		query = query.Where("professional_id = ?", *filter.ProfessionalID)
	}
	if filter.Date != nil {
		start := filter.Date.Truncate(24 * time.Hour)
		end := start.Add(24 * time.Hour)
		query = query.Where("start_at >= ? AND start_at < ?", start, end)
	}

	var bookings []domain.Booking
	if err := query.Order("start_at ASC").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

type AdminBookingInput struct {
	BookingInput
	TenantID uuid.UUID
}

func (s *Service) AdminCreateBooking(ctx context.Context, input AdminBookingInput) (*domain.Booking, error) {
	return s.CreateBooking(ctx, input.TenantID, input.BookingInput)
}

func (s *Service) AdminUpdateBooking(ctx context.Context, bookingID uuid.UUID, input BookingUpdateInput) (*domain.Booking, error) {
	var booking domain.Booking
	if err := s.dbWithContext(ctx).First(&booking, "id = ?", bookingID).Error; err != nil {
		return nil, err
	}
	return s.UpdateBooking(ctx, booking.TenantID, bookingID, input)
}

func (s *Service) AdminDeleteBooking(ctx context.Context, bookingID uuid.UUID) error {
	return s.dbWithContext(ctx).Delete(&domain.Booking{}, "id = ?", bookingID).Error
}
