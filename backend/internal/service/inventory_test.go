package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

func TestCreateInventoryMovementValidatesInput(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()
	product := seedProductRecord(t, tenant.ID, "Product INV", "SKU-INV")

	_, err := testSvc.CreateInventoryMovement(context.Background(), tenant.ID, InventoryInput{
		ProductID: product.ID,
		Type:      domain.InventoryMovementIn,
		Quantity:  0,
	})
	require.Error(t, err)

	_, err = testSvc.CreateInventoryMovement(context.Background(), tenant.ID, InventoryInput{
		ProductID: product.ID,
		Type:      "unknown",
		Quantity:  5,
	})
	require.ErrorIs(t, err, ErrInvalidInventoryType)
}

func TestCreateInventoryMovementReferencesOrderAndProduct(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()
	product := seedProductRecord(t, tenant.ID, "Product", "INV-001")
	service := seedServiceRecord(t, tenant.ID, "Service", 30)
	client := seedClientRecord(t, tenant.ID, "Order Client", "order@example.com", nil)

	order := &domain.SalesOrder{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		ClientID:    client.ID,
		Status:      domain.SalesOrderStatusConfirmed,
	}
	require.NoError(t, testDB.Create(order).Error)

	movement, err := testSvc.CreateInventoryMovement(context.Background(), tenant.ID, InventoryInput{
		ProductID: product.ID,
		OrderID:   &order.ID,
		Type:      domain.InventoryMovementOut,
		Quantity:  2,
		Reason:    "Order shipment",
	})
	require.NoError(t, err)
	assert.Equal(t, order.ID, *movement.OrderID)
	assert.Equal(t, product.ID, movement.ProductID)
	assert.Equal(t, "Order shipment", movement.Reason)

	// ensure service is used to calculate if needed
	_ = service
}

func TestListInventoryMovementsFilters(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()
	product := seedProductRecord(t, tenant.ID, "Product F", "INV-002")
	otherProduct := seedProductRecord(t, tenant.ID, "Product G", "INV-003")

	moveIn := &domain.InventoryMovement{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		ProductID:   product.ID,
		Type:        domain.InventoryMovementIn,
		Quantity:    3,
		Reason:      "Restock",
	}
	require.NoError(t, testDB.Create(moveIn).Error)
	require.NoError(t, testDB.First(moveIn, moveIn.ID).Error)
	time.Sleep(10 * time.Millisecond)
	moveOut := &domain.InventoryMovement{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		ProductID:   otherProduct.ID,
		Type:        domain.InventoryMovementOut,
		Quantity:    1,
	}
	require.NoError(t, testDB.Create(moveOut).Error)
	require.NoError(t, testDB.First(moveOut, moveOut.ID).Error)

	productFilter := InventoryFilter{ProductID: &product.ID}
	productMovements, err := testSvc.ListInventoryMovements(context.Background(), tenant.ID, productFilter)
	require.NoError(t, err)
	require.Len(t, productMovements, 1)
	assert.Equal(t, product.ID, productMovements[0].ProductID)

	startRange := moveOut.CreatedAt.Add(-time.Second)
	endRange := moveOut.CreatedAt.Add(time.Second)
	filter := InventoryFilter{Type: domain.InventoryMovementOut, StartDate: &startRange, EndDate: &endRange}
	outMovements, err := testSvc.ListInventoryMovements(context.Background(), tenant.ID, filter)
	require.NoError(t, err)
	require.Len(t, outMovements, 1)
	assert.Equal(t, domain.InventoryMovementOut, outMovements[0].Type)
}

func seedProductRecord(t *testing.T, tenantID uuid.UUID, name, sku string) *domain.Product {
	t.Helper()
	product := &domain.Product{
		TenantModel: domain.TenantModel{TenantID: tenantID},
		Name:        name,
		SKU:         sku,
		Price:       100,
	}
	require.NoError(t, testDB.Create(product).Error)
	return product
}
