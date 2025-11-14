package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) ListAll(ctx context.Context) ([]domain.Booking, error) {
	var bookings []domain.Booking
	if err := r.db.WithContext(ctx).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	var booking domain.Booking
	if err := r.db.WithContext(ctx).First(&booking, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *BookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}

func (r *BookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Booking{}, id).Error
}
