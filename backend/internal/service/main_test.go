package service

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"github.com/kusmin/gestao_updev/backend/internal/repository"
)

var (
	testDB       *gorm.DB
	testSvc      *Service
	schemaModels = []interface{}{
		&domain.Company{},
		&domain.User{},
		&domain.Client{},
		&domain.Service{},
		&domain.Product{},
		&domain.Booking{},
		&domain.SalesOrder{},
		&domain.SalesItem{},
		&domain.Payment{},
		&domain.InventoryMovement{},
	}
)

func TestMain(m *testing.M) {
	// Setup
	db, err := setupTestDatabase()
	if err != nil {
		log.Fatalf("could not set up test database: %v", err)
	}
	testDB = db

	if err := ensureCompanyDocumentConstraint(db); err != nil {
		log.Fatalf("could not adjust company constraint: %v", err)
	}
	if err := autoMigrateIfNeeded(db); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}
	clearAllData()

	// Create a mock service
	// In a real scenario, you might want to mock dependencies like JWT manager
	repo := repository.New(testDB)
	testSvc = New(&config.Config{}, repo, nil, nil) // Adjust as needed

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
func createTestTenant() (*domain.Company, error) {
	tenant := &domain.Company{
		Name:     "Test Tenant",
		Document: uuid.New().String(),
	}
	if err := testDB.Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

// Helper function to clear all data from tables
func clearAllData() {
	tables := []string{
		"availability_rules",
		"payments",
		"sales_items",
		"sales_orders",
		"inventory_movements",
		"bookings",
		"services",
		"products",
		"clients",
		"users",
		"professionals",
		"audit_logs",
		"companies",
	}
	stmt := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE"
	if err := testDB.Exec(stmt).Error; err != nil {
		log.Printf("failed to truncate test tables: %v", err)
	}
}

func autoMigrateIfNeeded(db *gorm.DB) error {
	if os.Getenv("SKIP_AUTO_MIGRATE") == "1" && schemaAlreadyPresent(db) {
		return nil
	}
	return db.AutoMigrate(schemaModels...)
}

func schemaAlreadyPresent(db *gorm.DB) bool {
	migrator := db.Migrator()
	for _, model := range schemaModels {
		if !migrator.HasTable(model) {
			return false
		}
	}
	return true
}

func ensureCompanyDocumentConstraint(db *gorm.DB) error {
	const stmt = `
DO $$
BEGIN
	IF to_regclass('public.companies') IS NULL THEN
		RETURN;
	END IF;

	IF EXISTS (
		SELECT 1 FROM pg_constraint
		WHERE conname = 'companies_document_key'
		  AND conrelid = 'public.companies'::regclass
	) THEN
		ALTER TABLE companies RENAME CONSTRAINT companies_document_key TO uni_companies_document;
	END IF;

	IF NOT EXISTS (
		SELECT 1 FROM pg_constraint
		WHERE conname = 'uni_companies_document'
		  AND conrelid = 'public.companies'::regclass
	) THEN
		ALTER TABLE companies ADD CONSTRAINT uni_companies_document UNIQUE (document);
	END IF;
END $$;
`
	return db.Exec(stmt).Error
}
