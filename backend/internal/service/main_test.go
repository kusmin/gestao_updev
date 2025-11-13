package service

import (
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

var (
	testDB  *gorm.DB
	testSvc *Service
)

func TestMain(m *testing.M) {
	// Setup
	db, err := setupTestDatabase()
	if err != nil {
		log.Fatalf("could not set up test database: %v", err)
	}
	testDB = db

	// Run migrations
	err = db.AutoMigrate(
		&domain.Tenant{},
		&domain.User{},
		&domain.Client{},
		&domain.Service{},
		&domain.Product{},
		&domain.Booking{},
		&domain.SalesOrder{},
		&domain.SalesOrderItem{},
		&domain.Payment{},
		&domain.InventoryMovement{},
	)
	if err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	// Create a mock service
	// In a real scenario, you might want to mock dependencies like JWT manager
	testSvc = New(&config.Config{}, nil, nil, nil) // Adjust as needed
	testSvc.db = testDB

	// Run tests
	code := m.Run()

	// Teardown
	// No teardown needed for the test container, it will be destroyed.

	os.Exit(code)
}

func setupTestDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for local testing if needed
		dsn = "postgres://testuser:testpassword@localhost:5433/testdb?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Helper function to create a tenant for tests
func createTestTenant() (*domain.Tenant, error) {
	tenant := &domain.Tenant{
		Name: "Test Tenant",
	}
	if err := testDB.Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

// Helper function to clear all data from tables
func clearAllData() {
	tables := []string{
		"payments", "sales_order_items", "sales_orders", "inventory_movements",
		"bookings", "products", "services", "clients", "users", "tenants",
	}
	for _, table := range tables {
		testDB.Exec("DELETE FROM " + table)
	}
}
