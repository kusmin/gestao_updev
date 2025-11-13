package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

func TestGetProduct(t *testing.T) {
	// Setup: Ensure clean state and create a tenant
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)

	// Create a product to be fetched
	product := &domain.Product{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		Name:        "Test Product",
		SKU:         "TST-001",
		Price:       99.99,
	}
	err = testDB.Create(product).Error
	require.NoError(t, err)

	t.Run("should get product successfully", func(t *testing.T) {
		// Act
		foundProduct, err := testSvc.GetProduct(context.Background(), tenant.ID, product.ID)

		// Assert
		assert.NoError(t, err)
		require.NotNil(t, foundProduct)
		assert.Equal(t, product.ID, foundProduct.ID)
		assert.Equal(t, "Test Product", foundProduct.Name)
	})

	t.Run("should return error for non-existent product", func(t *testing.T) {
		// Act
		nonExistentID := uuid.New()
		foundProduct, err := testSvc.GetProduct(context.Background(), tenant.ID, nonExistentID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundProduct)
	})

	t.Run("should return error for product in another tenant", func(t *testing.T) {
		// Setup: create another tenant and product
		otherTenant, err := createTestTenant()
		require.NoError(t, err)

		// Act
		foundProduct, err := testSvc.GetProduct(context.Background(), otherTenant.ID, product.ID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundProduct)
	})
}

func TestGetService(t *testing.T) {
	// Setup: Ensure clean state and create a tenant
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)

	// Create a service to be fetched
	service := &domain.Service{
		TenantModel:     domain.TenantModel{TenantID: tenant.ID},
		Name:            "Test Service",
		DurationMinutes: 60,
		Price:           150.00,
	}
	err = testDB.Create(service).Error
	require.NoError(t, err)

	t.Run("should get service successfully", func(t *testing.T) {
		// Act
		foundService, err := testSvc.GetService(context.Background(), tenant.ID, service.ID)

		// Assert
		assert.NoError(t, err)
		require.NotNil(t, foundService)
		assert.Equal(t, service.ID, foundService.ID)
		assert.Equal(t, "Test Service", foundService.Name)
	})

	t.Run("should return error for non-existent service", func(t *testing.T) {
		// Act
		nonExistentID := uuid.New()
		foundService, err := testSvc.GetService(context.Background(), tenant.ID, nonExistentID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundService)
	})

	t.Run("should return error for service in another tenant", func(t *testing.T) {
		// Setup: create another tenant
		otherTenant, err := createTestTenant()
		require.NoError(t, err)

		// Act
		foundService, err := testSvc.GetService(context.Background(), otherTenant.ID, service.ID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundService)
	})
}
