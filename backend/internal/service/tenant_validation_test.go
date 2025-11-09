package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/auth"
	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"github.com/kusmin/gestao_updev/backend/internal/repository"
)

func newTenantAwareService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	applyTestSchema(t, db)

	cfg := &config.Config{
		AppEnv:           "test",
		TenantHeader:     "X-Test-Tenant",
		JWTAccessSecret:  "test-access",
		JWTRefreshSecret: "test-refresh",
		JWTAccessTTL:     time.Minute,
		JWTRefreshTTL:    time.Hour,
	}
	jwtMgr := auth.NewJWTManager("test-access", "test-refresh", time.Minute, time.Hour)
	repo := repository.New(db)
	logger := zap.NewNop()

	return New(cfg, repo, jwtMgr, logger), db
}

func applyTestSchema(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS clients (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			tenant_id TEXT NOT NULL,
			name TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			notes TEXT,
			tags TEXT,
			contact TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS professionals (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			tenant_id TEXT NOT NULL,
			user_id TEXT,
			name TEXT NOT NULL,
			specialties TEXT,
			max_parallel INTEGER,
			active BOOLEAN
		);`,
		`CREATE TABLE IF NOT EXISTS services (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			tenant_id TEXT NOT NULL,
			name TEXT NOT NULL,
			category TEXT,
			description TEXT,
			duration_minutes INTEGER,
			price REAL,
			color TEXT,
			metadata TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			tenant_id TEXT NOT NULL,
			client_id TEXT NOT NULL,
			professional_id TEXT NOT NULL,
			service_id TEXT NOT NULL,
			status TEXT NOT NULL,
			start_at DATETIME NOT NULL,
			end_at DATETIME NOT NULL,
			notes TEXT,
			metadata TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS sales_orders (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			tenant_id TEXT NOT NULL,
			client_id TEXT NOT NULL,
			booking_id TEXT,
			status TEXT NOT NULL,
			payment_type TEXT,
			total REAL,
			discount REAL,
			notes TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS sales_items (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			tenant_id TEXT NOT NULL,
			order_id TEXT NOT NULL,
			item_type TEXT NOT NULL,
			item_ref_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price REAL NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			tenant_id TEXT NOT NULL,
			name TEXT NOT NULL,
			sku TEXT NOT NULL,
			price REAL,
			cost REAL,
			stock_qty INTEGER,
			min_stock INTEGER,
			description TEXT,
			metadata TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS inventory_movements (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			tenant_id TEXT NOT NULL,
			product_id TEXT NOT NULL,
			order_id TEXT,
			type TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			reason TEXT
		);`,
	}

	for _, stmt := range statements {
		require.NoError(t, db.Exec(stmt).Error)
	}
}

func seedClient(t *testing.T, db *gorm.DB, tenantID uuid.UUID, name string) domain.Client {
	t.Helper()
	client := domain.Client{
		TenantModel: domain.TenantModel{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			TenantID:  tenantID,
		},
		Name:    name,
		Contact: datatypes.JSONMap{},
		Tags:    datatypes.JSON([]byte("[]")),
	}
	require.NoError(t, db.Create(&client).Error)
	return client
}

func seedProfessional(t *testing.T, db *gorm.DB, tenantID uuid.UUID, name string) domain.Professional {
	t.Helper()
	pro := domain.Professional{
		TenantModel: domain.TenantModel{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			TenantID:  tenantID,
		},
		Name:        name,
		Specialties: datatypes.JSON([]byte("[]")),
		MaxParallel: 1,
		Active:      true,
	}
	require.NoError(t, db.Create(&pro).Error)
	return pro
}

func seedService(t *testing.T, db *gorm.DB, tenantID uuid.UUID, name string, duration int, price float64) domain.Service {
	t.Helper()
	service := domain.Service{
		TenantModel: domain.TenantModel{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			TenantID:  tenantID,
		},
		Name:            name,
		DurationMinutes: duration,
		Price:           price,
		Metadata:        datatypes.JSONMap{},
	}
	require.NoError(t, db.Create(&service).Error)
	return service
}

func seedBooking(t *testing.T, db *gorm.DB, tenantID, clientID, professionalID, serviceID uuid.UUID) domain.Booking {
	t.Helper()
	now := time.Now().UTC()
	booking := domain.Booking{
		TenantModel: domain.TenantModel{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			TenantID:  tenantID,
		},
		ClientID:       clientID,
		ProfessionalID: professionalID,
		ServiceID:      serviceID,
		Status:         domain.BookingStatusPending,
		StartAt:        now,
		EndAt:          now.Add(30 * time.Minute),
		Metadata:       datatypes.JSONMap{},
	}
	require.NoError(t, db.Create(&booking).Error)
	return booking
}

func seedProduct(t *testing.T, db *gorm.DB, tenantID uuid.UUID, name string) domain.Product {
	t.Helper()
	product := domain.Product{
		TenantModel: domain.TenantModel{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			TenantID:  tenantID,
		},
		Name:     name,
		SKU:      uuid.NewString(),
		Price:    100,
		Cost:     50,
		Metadata: datatypes.JSONMap{},
	}
	require.NoError(t, db.Create(&product).Error)
	return product
}

func seedSalesOrder(t *testing.T, db *gorm.DB, tenantID, clientID uuid.UUID) domain.SalesOrder {
	t.Helper()
	order := domain.SalesOrder{
		TenantModel: domain.TenantModel{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			TenantID:  tenantID,
		},
		ClientID: clientID,
		Status:   domain.SalesOrderStatusDraft,
		Total:    100,
		Discount: 0,
		Notes:    "seed order",
	}
	require.NoError(t, db.Create(&order).Error)
	return order
}

func TestCreateBookingRejectsCrossTenantClient(t *testing.T) {
	svc, db := newTenantAwareService(t)
	ctx := context.Background()

	tenantA := uuid.New()
	tenantB := uuid.New()

	clientB := seedClient(t, db, tenantB, "Client B")
	proA := seedProfessional(t, db, tenantA, "Pro A")
	serviceA := seedService(t, db, tenantA, "Service A", 30, 100)

	input := BookingInput{
		ClientID:       clientB.ID,
		ProfessionalID: proA.ID,
		ServiceID:      serviceA.ID,
		StartAt:        time.Now().UTC(),
		Status:         domain.BookingStatusPending,
	}

	_, err := svc.CreateBooking(ctx, tenantA, input)
	require.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCreateSalesOrderRejectsBookingFromAnotherTenant(t *testing.T) {
	svc, db := newTenantAwareService(t)
	ctx := context.Background()

	tenantA := uuid.New()
	tenantB := uuid.New()

	clientA := seedClient(t, db, tenantA, "Client A")
	serviceA := seedService(t, db, tenantA, "Service A", 60, 200)

	clientB := seedClient(t, db, tenantB, "Client B")
	proB := seedProfessional(t, db, tenantB, "Pro B")
	serviceB := seedService(t, db, tenantB, "Service B", 45, 150)
	bookingB := seedBooking(t, db, tenantB, clientB.ID, proB.ID, serviceB.ID)

	input := SalesOrderInput{
		ClientID:  clientA.ID,
		BookingID: &bookingB.ID,
		Items: []SalesItemInput{
			{
				Type:      "service",
				RefID:     serviceA.ID,
				Quantity:  1,
				UnitPrice: 200,
			},
		},
	}

	_, err := svc.CreateSalesOrder(ctx, tenantA, input)
	require.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCreateInventoryMovementRejectsOrderFromAnotherTenant(t *testing.T) {
	svc, db := newTenantAwareService(t)
	ctx := context.Background()

	tenantA := uuid.New()
	tenantB := uuid.New()

	productA := seedProduct(t, db, tenantA, "Product A")
	clientB := seedClient(t, db, tenantB, "Client B")
	orderB := seedSalesOrder(t, db, tenantB, clientB.ID)

	input := InventoryInput{
		ProductID: productA.ID,
		OrderID:   &orderB.ID,
		Type:      domain.InventoryMovementOut,
		Quantity:  1,
		Reason:    "cross-tenant link attempt",
	}

	_, err := svc.CreateInventoryMovement(ctx, tenantA, input)
	require.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
