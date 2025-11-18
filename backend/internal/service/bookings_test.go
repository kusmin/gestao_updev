package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

func TestCreateBookingCalculatesEndAndPreventsOverlap(t *testing.T) {
	setupTest(t)
	clearAllData()
	tenant, _ := createTestTenant()
	client := seedClientRecord(t, tenant.ID, "Booking Client", "booking@example.com", nil)
	pro := seedProfessionalRecord(t, tenant.ID, "Pro 1")
	service := seedServiceRecord(t, tenant.ID, "Haircut", 30)

	start := time.Now().UTC().Truncate(time.Minute).Add(time.Hour)
	created, err := testSvc.CreateBooking(context.Background(), tenant.ID, BookingInput{
		ClientID:       client.ID,
		ProfessionalID: pro.ID,
		ServiceID:      service.ID,
		StartAt:        start,
		Notes:          "Initial booking",
	})
	require.NoError(t, err)
	require.Equal(t, start.Add(30*time.Minute), created.EndAt)
	assert.Equal(t, domain.BookingStatusPending, created.Status)

	// attempt overlapping booking
	_, err = testSvc.CreateBooking(context.Background(), tenant.ID, BookingInput{
		ClientID:       client.ID,
		ProfessionalID: pro.ID,
		ServiceID:      service.ID,
		StartAt:        start.Add(10 * time.Minute),
	})
	require.ErrorIs(t, err, ErrBookingConflict)
}

func TestUpdateBookingWithConflictDetection(t *testing.T) {
	setupTest(t)
	clearAllData()
	tenant, _ := createTestTenant()
	client := seedClientRecord(t, tenant.ID, "Update Booking", "update@example.com", nil)
	pro := seedProfessionalRecord(t, tenant.ID, "Pro 2")
	service := seedServiceRecord(t, tenant.ID, "Massage", 60)

	slot1 := time.Now().UTC().Truncate(time.Minute).Add(2 * time.Hour)
	slot2 := slot1.Add(2 * time.Hour)

	first, _ := testSvc.CreateBooking(context.Background(), tenant.ID, BookingInput{
		ClientID:       client.ID,
		ProfessionalID: pro.ID,
		ServiceID:      service.ID,
		StartAt:        slot1,
		Status:         domain.BookingStatusConfirmed,
	})
	second, _ := testSvc.CreateBooking(context.Background(), tenant.ID, BookingInput{
		ClientID:       client.ID,
		ProfessionalID: pro.ID,
		ServiceID:      service.ID,
		StartAt:        slot2,
	})

	moveStart := slot1.Add(15 * time.Minute)
	_, err := testSvc.UpdateBooking(context.Background(), tenant.ID, second.ID, BookingUpdateInput{
		StartAt: &moveStart,
	})
	require.ErrorIs(t, err, ErrBookingConflict)

	newStatus := domain.BookingStatusDone
	updated, err := testSvc.UpdateBooking(context.Background(), tenant.ID, first.ID, BookingUpdateInput{
		Status: &newStatus,
	})
	require.NoError(t, err)
	assert.Equal(t, newStatus, updated.Status)
}

func TestCancelBookingPersistsReason(t *testing.T) {
	setupTest(t)
	clearAllData()
	tenant, _ := createTestTenant()
	client := seedClientRecord(t, tenant.ID, "Cancel Booking", "cancel@example.com", nil)
	pro := seedProfessionalRecord(t, tenant.ID, "Pro Cancel")
	service := seedServiceRecord(t, tenant.ID, "Color", 45)

	start := time.Now().UTC().Truncate(time.Minute).Add(3 * time.Hour)
	booking, _ := testSvc.CreateBooking(context.Background(), tenant.ID, BookingInput{
		ClientID:       client.ID,
		ProfessionalID: pro.ID,
		ServiceID:      service.ID,
		StartAt:        start,
		Status:         domain.BookingStatusConfirmed,
	})

	reason := "No show"
	canceled, err := testSvc.CancelBooking(context.Background(), tenant.ID, booking.ID, reason)
	require.NoError(t, err)
	assert.Equal(t, domain.BookingStatusCanceled, canceled.Status)
	assert.Equal(t, reason, canceled.Metadata["cancel_reason"])
}

func TestListBookingsFilters(t *testing.T) {
	setupTest(t)
	clearAllData()
	tenant, _ := createTestTenant()
	client := seedClientRecord(t, tenant.ID, "List Booking", "list@example.com", nil)
	service := seedServiceRecord(t, tenant.ID, "Spa", 30)
	pro1 := seedProfessionalRecord(t, tenant.ID, "Pro Filter 1")
	pro2 := seedProfessionalRecord(t, tenant.ID, "Pro Filter 2")

	date := time.Now().UTC().Truncate(24 * time.Hour)
	_, _ = testSvc.CreateBooking(context.Background(), tenant.ID, BookingInput{
		ClientID:       client.ID,
		ProfessionalID: pro1.ID,
		ServiceID:      service.ID,
		StartAt:        date.Add(10 * time.Hour),
		Status:         domain.BookingStatusDone,
	})
	_, _ = testSvc.CreateBooking(context.Background(), tenant.ID, BookingInput{
		ClientID:       client.ID,
		ProfessionalID: pro2.ID,
		ServiceID:      service.ID,
		StartAt:        date.Add(26 * time.Hour),
	})

	pro1ID := pro1.ID
	filter := BookingFilter{Date: &date, ProfessionalID: &pro1ID, Status: domain.BookingStatusDone}
	bookings, err := testSvc.ListBookings(context.Background(), tenant.ID, filter)
	require.NoError(t, err)
	require.Len(t, bookings, 1)
	assert.Equal(t, pro1ID, bookings[0].ProfessionalID)
}

func TestAdminBookingLifecycle(t *testing.T) {
	setupTest(t)
	clearAllData()
	tenant, _ := createTestTenant()
	client := seedClientRecord(t, tenant.ID, "Admin Lifecycle", "adminlifecyle@example.com", nil)
	service := seedServiceRecord(t, tenant.ID, "Admin Service", 30)
	pro := seedProfessionalRecord(t, tenant.ID, "Pro Admin")

	start := time.Now().UTC().Truncate(time.Minute)
	adminBooking, err := testSvc.AdminCreateBooking(context.Background(), AdminBookingInput{
		BookingInput: BookingInput{
			ClientID:       client.ID,
			ProfessionalID: pro.ID,
			ServiceID:      service.ID,
			StartAt:        start,
		},
		TenantID: tenant.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, tenant.ID, adminBooking.TenantID)

	newStatus := domain.BookingStatusConfirmed
	updated, err := testSvc.AdminUpdateBooking(context.Background(), adminBooking.ID, BookingUpdateInput{
		Status: &newStatus,
	})
	require.NoError(t, err)
	assert.Equal(t, newStatus, updated.Status)

	require.NoError(t, testSvc.AdminDeleteBooking(context.Background(), adminBooking.ID))
	var count int64
	testDB.Model(&domain.Booking{}).Where("id = ?", adminBooking.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestListAllBookings(t *testing.T) {
	setupTest(t)
	clearAllData()
	tenant1, _ := createTestTenant()
	tenant2, _ := createTestTenant()

	client1 := seedClientRecord(t, tenant1.ID, "Client 1", "client1@example.com", nil)
	pro1 := seedProfessionalRecord(t, tenant1.ID, "Pro 1")
	service1 := seedServiceRecord(t, tenant1.ID, "Service 1", 30)

	client2 := seedClientRecord(t, tenant2.ID, "Client 2", "client2@example.com", nil)
	pro2 := seedProfessionalRecord(t, tenant2.ID, "Pro 2")
	service2 := seedServiceRecord(t, tenant2.ID, "Service 2", 60)

	date1 := time.Now().UTC().Truncate(24 * time.Hour)
	date2 := date1.Add(24 * time.Hour)

	// Bookings for tenant 1
	_, _ = testSvc.CreateBooking(context.Background(), tenant1.ID, BookingInput{
		ClientID:       client1.ID,
		ProfessionalID: pro1.ID,
		ServiceID:      service1.ID,
		StartAt:        date1.Add(9 * time.Hour),
		Status:         domain.BookingStatusConfirmed,
	})
	_, _ = testSvc.CreateBooking(context.Background(), tenant1.ID, BookingInput{
		ClientID:       client1.ID,
		ProfessionalID: pro1.ID,
		ServiceID:      service1.ID,
		StartAt:        date1.Add(10 * time.Hour),
		Status:         domain.BookingStatusPending,
	})

	// Booking for tenant 2
	_, _ = testSvc.CreateBooking(context.Background(), tenant2.ID, BookingInput{
		ClientID:       client2.ID,
		ProfessionalID: pro2.ID,
		ServiceID:      service2.ID,
		StartAt:        date2.Add(11 * time.Hour),
		Status:         domain.BookingStatusConfirmed,
	})

	// Test ListAllBookings without filters
	bookings, err := testSvc.ListAllBookings(context.Background(), BookingFilter{})
	require.NoError(t, err)
	assert.Len(t, bookings, 3)

	// Test ListAllBookings with status filter
	filter := BookingFilter{Status: domain.BookingStatusConfirmed}
	bookings, err = testSvc.ListAllBookings(context.Background(), filter)
	require.NoError(t, err)
	assert.Len(t, bookings, 2)

	// Test ListAllBookings with professional ID filter
	pro1ID := pro1.ID
	filter = BookingFilter{ProfessionalID: &pro1ID}
	bookings, err = testSvc.ListAllBookings(context.Background(), filter)
	require.NoError(t, err)
	assert.Len(t, bookings, 2)

	// Test ListAllBookings with date filter
	filter = BookingFilter{Date: &date1}
	bookings, err = testSvc.ListAllBookings(context.Background(), filter)
	require.NoError(t, err)
	assert.Len(t, bookings, 2)
}
