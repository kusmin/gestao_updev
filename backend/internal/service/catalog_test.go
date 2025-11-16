package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

func TestGetProduct(t *testing.T) {
	t.Run("should get product successfully", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)
		product := &domain.Product{
			TenantModel: domain.TenantModel{TenantID: tenant.ID},
			Name:        "Test Product",
			SKU:         "TST-001",
			Price:       99.99,
		}
		err = testDB.Create(product).Error
		require.NoError(t, err)

		// Act
		foundProduct, err := testSvc.GetProduct(context.Background(), tenant.ID, product.ID)

		// Assert
		assert.NoError(t, err)
		require.NotNil(t, foundProduct)
		assert.Equal(t, product.ID, foundProduct.ID)
		assert.Equal(t, "Test Product", foundProduct.Name)
	})

	t.Run("should return error for non-existent product", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)

		// Act
		nonExistentID := uuid.New()
		foundProduct, err := testSvc.GetProduct(context.Background(), tenant.ID, nonExistentID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundProduct)
	})

	t.Run("should return error for product in another tenant", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)
		product := &domain.Product{
			TenantModel: domain.TenantModel{TenantID: tenant.ID},
			Name:        "Test Product",
			SKU:         "TST-001",
			Price:       99.99,
		}
		err = testDB.Create(product).Error
		require.NoError(t, err)

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
	t.Run("should get service successfully", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)
		service := &domain.Service{
			TenantModel:     domain.TenantModel{TenantID: tenant.ID},
			Name:            "Test Service",
			DurationMinutes: 60,
			Price:           150.00,
		}
		err = testDB.Create(service).Error
		require.NoError(t, err)

		// Act
		foundService, err := testSvc.GetService(context.Background(), tenant.ID, service.ID)

		// Assert
		assert.NoError(t, err)
		require.NotNil(t, foundService)
		assert.Equal(t, service.ID, foundService.ID)
		assert.Equal(t, "Test Service", foundService.Name)
	})

	t.Run("should return error for non-existent service", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)

		// Act
		nonExistentID := uuid.New()
		foundService, err := testSvc.GetService(context.Background(), tenant.ID, nonExistentID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundService)
	})

	t.Run("should return error for service in another tenant", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)
		service := &domain.Service{
			TenantModel:     domain.TenantModel{TenantID: tenant.ID},
			Name:            "Test Service",
			DurationMinutes: 60,
			Price:           150.00,
		}
		err = testDB.Create(service).Error
		require.NoError(t, err)

		otherTenant, err := createTestTenant()
		require.NoError(t, err)

		// Act
		foundService, err := testSvc.GetService(context.Background(), otherTenant.ID, service.ID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundService)
	})
}

func TestListProducts(t *testing.T) {
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)
	otherTenant, err := createTestTenant()
	require.NoError(t, err)

	_, err = testSvc.CreateProduct(context.Background(), tenant.ID, ProductInput{
		Name:     "Beta Product",
		SKU:      "SKU-002",
		Price:    200,
		Cost:     120,
		StockQty: 5,
		Metadata: map[string]interface{}{"tier": "beta"},
	})
	require.NoError(t, err)
	_, err = testSvc.CreateProduct(context.Background(), tenant.ID, ProductInput{
		Name:     "Alpha Product",
		SKU:      "SKU-001",
		Price:    150,
		Cost:     80,
		StockQty: 2,
		Metadata: map[string]interface{}{"tier": "alpha"},
	})
	require.NoError(t, err)
	_, err = testSvc.CreateProduct(context.Background(), otherTenant.ID, ProductInput{
		Name:  "Other Tenant Product",
		SKU:   "SKU-999",
		Price: 500,
	})
	require.NoError(t, err)

	products, err := testSvc.ListProducts(context.Background(), tenant.ID)

	assert.NoError(t, err)
	require.Len(t, products, 2)
	assert.Equal(t, "Alpha Product", products[0].Name)
	assert.Equal(t, "Beta Product", products[1].Name)
}

func TestCreateProduct(t *testing.T) {
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)

	input := ProductInput{
		Name:        "Created Product",
		SKU:         "CP-123",
		Price:       75.5,
		Cost:        30.2,
		StockQty:    7,
		MinStock:    2,
		Description: "A brand new product",
		Metadata:    map[string]interface{}{"origin": "local"},
	}

	product, err := testSvc.CreateProduct(context.Background(), tenant.ID, input)

	require.NoError(t, err)
	require.NotNil(t, product)
	assert.Equal(t, tenant.ID, product.TenantID)
	assert.Equal(t, input.Name, product.Name)
	assert.Equal(t, input.SKU, product.SKU)
	assert.Equal(t, input.Price, product.Price)
	assert.Equal(t, input.Cost, product.Cost)
	assert.Equal(t, input.StockQty, product.StockQty)
	assert.Equal(t, input.MinStock, product.MinStock)
	assert.Equal(t, input.Description, product.Description)
	assert.Equal(t, datatypes.JSONMap{"origin": "local"}, product.Metadata)
}

func TestUpdateProduct(t *testing.T) {
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)
	created, err := testSvc.CreateProduct(context.Background(), tenant.ID, ProductInput{
		Name:        "Initial Product",
		SKU:         "IP-001",
		Price:       40,
		Cost:        10,
		StockQty:    4,
		MinStock:    1,
		Description: "Initial description",
		Metadata:    map[string]interface{}{"stage": "initial"},
	})
	require.NoError(t, err)

	updated, err := testSvc.UpdateProduct(context.Background(), tenant.ID, created.ID, ProductInput{
		Name:        "Updated Product",
		SKU:         "IP-002",
		Price:       55.5,
		Cost:        15.5,
		StockQty:    10,
		MinStock:    3,
		Description: "Updated description",
		Metadata:    map[string]interface{}{"stage": "updated"},
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated Product", updated.Name)
	assert.Equal(t, "IP-002", updated.SKU)
	assert.Equal(t, 55.5, updated.Price)
	assert.Equal(t, 15.5, updated.Cost)
	assert.Equal(t, 10, updated.StockQty)
	assert.Equal(t, 3, updated.MinStock)
	assert.Equal(t, "Updated description", updated.Description)
	assert.Equal(t, datatypes.JSONMap{"stage": "updated"}, updated.Metadata)
}

func TestDeleteProduct(t *testing.T) {
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)
	created, err := testSvc.CreateProduct(context.Background(), tenant.ID, ProductInput{
		Name:  "Disposable Product",
		SKU:   "DP-001",
		Price: 20,
	})
	require.NoError(t, err)

	err = testSvc.DeleteProduct(context.Background(), tenant.ID, created.ID)

	assert.NoError(t, err)
	deleted, getErr := testSvc.GetProduct(context.Background(), tenant.ID, created.ID)
	assert.Nil(t, deleted)
	assert.ErrorIs(t, getErr, gorm.ErrRecordNotFound)
}

func TestCreateService(t *testing.T) {
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)

	input := Input{Name: "Created Service",
		Category:        "Wellness",
		Description:     "A relaxing service",
		DurationMinutes: 90,
		Price:           250.75,
		Color:           "#FFFFFF",
		Metadata:        map[string]interface{}{"level": "premium"},
	}

	service, err := testSvc.CreateService(context.Background(), tenant.ID, input)

	require.NoError(t, err)
	require.NotNil(t, service)
	assert.Equal(t, tenant.ID, service.TenantID)
	assert.Equal(t, input.Name, service.Name)
	assert.Equal(t, input.Category, service.Category)
	assert.Equal(t, input.Description, service.Description)
	assert.Equal(t, input.DurationMinutes, service.DurationMinutes)
	assert.Equal(t, input.Price, service.Price)
	assert.Equal(t, input.Color, service.Color)
	assert.Equal(t, datatypes.JSONMap{"level": "premium"}, service.Metadata)
}

func TestUpdateService(t *testing.T) {
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)
	created, err := testSvc.CreateService(context.Background(), tenant.ID, Input{Name: "Initial Service",
		Category:        "Initial",
		Description:     "Initial desc",
		DurationMinutes: 45,
		Price:           120,
		Color:           "#000000",
		Metadata:        map[string]interface{}{"phase": "initial"},
	})
	require.NoError(t, err)

	updated, err := testSvc.UpdateService(context.Background(), tenant.ID, created.ID, Input{Name: "Updated Service",
		Category:        "Updated",
		Description:     "Updated desc",
		DurationMinutes: 60,
		Price:           180.5,
		Color:           "#FF00FF",
		Metadata:        map[string]interface{}{"phase": "updated"},
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated Service", updated.Name)
	assert.Equal(t, "Updated", updated.Category)
	assert.Equal(t, "Updated desc", updated.Description)
	assert.Equal(t, 60, updated.DurationMinutes)
	assert.Equal(t, 180.5, updated.Price)
	assert.Equal(t, "#FF00FF", updated.Color)
	assert.Equal(t, datatypes.JSONMap{"phase": "updated"}, updated.Metadata)
}

func TestDeleteService(t *testing.T) {
	clearAllData()
	tenant, err := createTestTenant()
	require.NoError(t, err)
	created, err := testSvc.CreateService(context.Background(), tenant.ID, Input{Name: "Disposable Service",
		DurationMinutes: 30,
		Price:           50,
	})
	require.NoError(t, err)

	err = testSvc.DeleteService(context.Background(), tenant.ID, created.ID)

	assert.NoError(t, err)
	deleted, getErr := testSvc.GetService(context.Background(), tenant.ID, created.ID)
	assert.Nil(t, deleted)
	assert.ErrorIs(t, getErr, gorm.ErrRecordNotFound)
}

func TestListAllServicesAndProducts(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()
	_, _ = testSvc.CreateService(context.Background(), tenant.ID, Input{Name: "Service B"})
	_, _ = testSvc.CreateService(context.Background(), tenant.ID, Input{Name: "Service A"})
	_, _ = testSvc.CreateProduct(context.Background(), tenant.ID, ProductInput{Name: "Product B", SKU: "B"})
	_, _ = testSvc.CreateProduct(context.Background(), tenant.ID, ProductInput{Name: "Product A", SKU: "A"})

	services, err := testSvc.ListAllServices(context.Background())
	require.NoError(t, err)
	require.Len(t, services, 2)
	assert.Equal(t, "Service A", services[0].Name)

	products, err := testSvc.ListAllProducts(context.Background())
	require.NoError(t, err)
	require.Len(t, products, 2)
	assert.Equal(t, "Product A", products[0].Name)
}

func TestAdminProductLifecycle(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()

	product, err := testSvc.AdminCreateProduct(context.Background(), AdminProductInput{
		ProductInput: ProductInput{
			Name:     "Admin Prod",
			SKU:      "ADMIN-1",
			Price:    99,
			Metadata: map[string]interface{}{"source": "admin"},
		},
		TenantID: tenant.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, tenant.ID, product.TenantID)

	updated, err := testSvc.AdminUpdateProduct(context.Background(), product.ID, ProductInput{
		Name:  "Admin Prod Updated",
		SKU:   "ADMIN-2",
		Price: 120,
	})
	require.NoError(t, err)
	assert.Equal(t, "Admin Prod Updated", updated.Name)
	assert.Equal(t, "ADMIN-2", updated.SKU)

	require.NoError(t, testSvc.AdminDeleteProduct(context.Background(), product.ID))
	var count int64
	testDB.Model(&domain.Product{}).Where("id = ?", product.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestAdminServiceLifecycle(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()

	service, err := testSvc.AdminCreateService(context.Background(), AdminServiceInput{
		Input: Input{Name: "Admin Service",
			DurationMinutes: 30,
			Price:           70,
		},
		TenantID: tenant.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, tenant.ID, service.TenantID)

	updated, err := testSvc.AdminUpdateService(context.Background(), service.ID, Input{Name: "Admin Service Updated",
		DurationMinutes: 45,
		Price:           90,
	})
	require.NoError(t, err)
	assert.Equal(t, "Admin Service Updated", updated.Name)
	assert.Equal(t, 45, updated.DurationMinutes)

	require.NoError(t, testSvc.AdminDeleteService(context.Background(), service.ID))
	var count int64
	testDB.Model(&domain.Service{}).Where("id = ?", service.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
