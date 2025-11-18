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

func TestListClientsFiltersBySearchAndTags(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()
	otherTenant, _ := createTestTenant()

	target := seedClientRecord(t, tenant.ID, "Ana Paula", "ana@example.com", []string{"vip"})
	seedClientRecord(t, tenant.ID, "Carlos Souza", "carlos@example.com", []string{"basic"})
	seedClientRecord(t, otherTenant.ID, "Outro Tenant", "other@example.com", []string{"vip"})

	result, total, err := testSvc.ListClients(context.Background(), tenant.ID, ClientsFilter{
		Search: "ana",
		Tags:   []string{"vip"},
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, result, 1)
	assert.Equal(t, target.ID, result[0].ID)
}

func TestCreateAndUpdateClientPersistsStructuredData(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()
	input := ClientInput{
		Name:    "Primeiro Cliente",
		Email:   "cliente@example.com",
		Phone:   "1234-5678",
		Notes:   "Anotações iniciais",
		Tags:    []string{"vip", "recorrente"},
		Contact: map[string]interface{}{"instagram": "@cliente"},
	}

	created, err := testSvc.CreateClient(context.Background(), tenant.ID, input)
	require.NoError(t, err)
	require.Equal(t, "Primeiro Cliente", created.Name)
	assert.JSONEq(t, `["vip","recorrente"]`, string(created.Tags))
	assert.Equal(t, "cliente@example.com", created.Email)
	assert.Equal(t, "1234-5678", created.Phone)

	updateInput := ClientInput{
		Name:    "Cliente Atualizado",
		Email:   "atualizado@example.com",
		Phone:   "9999-0000",
		Notes:   "Observação alterada",
		Tags:    []string{"vip"},
		Contact: map[string]interface{}{"instagram": "@novo"},
	}

	updated, err := testSvc.UpdateClient(context.Background(), tenant.ID, created.ID, updateInput)
	require.NoError(t, err)
	assert.Equal(t, "Cliente Atualizado", updated.Name)
	assert.Equal(t, "atualizado@example.com", updated.Email)
	assert.JSONEq(t, `["vip"]`, string(updated.Tags))
	assert.Equal(t, "9999-0000", updated.Phone)
	assert.Equal(t, "Observação alterada", updated.Notes)
	assert.Equal(t, "@novo", updated.Contact["instagram"])
}

func TestGetClientReturnsStats(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()
	client := seedClientRecord(t, tenant.ID, "Cliente Stats", "stats@example.com", nil)
	pro := seedProfessionalRecord(t, tenant.ID, "Pro Stats")
	service := seedServiceRecord(t, tenant.ID, "Serv Stats", 60)

	// Seed bookings
	book1 := &domain.Booking{
		TenantModel:    domain.TenantModel{TenantID: tenant.ID},
		ClientID:       client.ID,
		ProfessionalID: pro.ID,
		ServiceID:      service.ID,
		Status:         domain.BookingStatusConfirmed,
		StartAt:        time.Now().Add(-48 * time.Hour),
		EndAt:          time.Now().Add(-47 * time.Hour),
	}
	require.NoError(t, testDB.Create(book1).Error)

	book2 := &domain.Booking{
		TenantModel:    domain.TenantModel{TenantID: tenant.ID},
		ClientID:       client.ID,
		ProfessionalID: pro.ID,
		ServiceID:      service.ID,
		Status:         domain.BookingStatusDone,
		StartAt:        time.Now().Add(-24 * time.Hour),
		EndAt:          time.Now().Add(-23 * time.Hour),
	}
	require.NoError(t, testDB.Create(book2).Error)

	// Seed payment linked to sales order
	order := &domain.SalesOrder{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		ClientID:    client.ID,
		Status:      domain.SalesOrderStatusPaid,
		Total:       150,
	}
	require.NoError(t, testDB.Create(order).Error)

	payment := &domain.Payment{
		TenantModel: domain.TenantModel{TenantID: tenant.ID},
		OrderID:     order.ID,
		Method:      "pix",
		Amount:      150,
		PaidAt:      time.Now(),
	}
	require.NoError(t, testDB.Create(payment).Error)

	result, stats, err := testSvc.GetClient(context.Background(), tenant.ID, client.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, stats)
	assert.Equal(t, int64(2), stats.TotalBookings)
	assert.Equal(t, 150.0, stats.TotalSpent)
	require.NotNil(t, stats.LastBooking)
	assert.True(t, stats.LastBooking.After(book1.StartAt))
}

func TestDeleteClientSoftDelete(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()
	client := seedClientRecord(t, tenant.ID, "Excluir", "del@example.com", nil)

	err := testSvc.DeleteClient(context.Background(), tenant.ID, client.ID)
	require.NoError(t, err)

	var count int64
	testDB.Model(&domain.Client{}).Where("id = ?", client.ID).Count(&count)
	assert.Equal(t, int64(0), count)

	var deleted domain.Client
	testDB.Unscoped().First(&deleted, client.ID)
	assert.NotNil(t, deleted.DeletedAt)
}

func TestListAllClients(t *testing.T) {
	clearAllData()
	tenant1, _ := createTestTenant()
	tenant2, _ := createTestTenant()

	client1Record := seedClientRecord(t, tenant1.ID, "Client 1", "client1@example.com", []string{"vip"})
	_ = seedClientRecord(t, tenant2.ID, "Client 2", "client2@example.com", []string{"normal"})
	_ = seedClientRecord(t, tenant1.ID, "Client 3", "client3@example.com", []string{"vip", "new"})

	// Test ListAllClients without filters
	clients, total, err := testSvc.ListAllClients(context.Background(), ClientsFilter{})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.Len(t, clients, 3)

	// Test ListAllClients with search filter
	filter := ClientsFilter{Search: "client 1"}
	clients, total, err = testSvc.ListAllClients(context.Background(), filter)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, clients, 1)
	assert.Equal(t, client1Record.ID, clients[0].ID)

	// Test ListAllClients with tags filter
	filter = ClientsFilter{Tags: []string{"vip"}}
	clients, total, err = testSvc.ListAllClients(context.Background(), filter)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, clients, 2)

	// Test ListAllClients with pagination
	filter = ClientsFilter{Page: 1, PerPage: 2}
	clients, total, err = testSvc.ListAllClients(context.Background(), filter)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.Len(t, clients, 2)

	filter = ClientsFilter{Page: 2, PerPage: 2}
	clients, total, err = testSvc.ListAllClients(context.Background(), filter)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.Len(t, clients, 1)
}

func TestAdminCreateClient(t *testing.T) {
	clearAllData()
	tenant, _ := createTestTenant()

	input := AdminClientInput{
		ClientInput: ClientInput{
			Name:    "Admin Client",
			Email:   "admin@example.com",
			Phone:   "1111-2222",
			Notes:   "Admin created client",
			Tags:    []string{"admin"},
			Contact: map[string]interface{}{"source": "internal"},
		},
		TenantID: tenant.ID,
	}

	created, err := testSvc.AdminCreateClient(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, tenant.ID, created.TenantID)
	assert.Equal(t, input.Name, created.Name)
	assert.Equal(t, input.Email, created.Email)
	assert.Equal(t, input.Phone, created.Phone)
	assert.Equal(t, input.Notes, created.Notes)
	assert.JSONEq(t, `["admin"]`, string(created.Tags))
	assert.Equal(t, "internal", created.Contact["source"])
}

func seedClientRecord(t *testing.T, tenantID uuid.UUID, name, email string, tags []string) *domain.Client {
	t.Helper()
	client := &domain.Client{
		TenantModel: domain.TenantModel{TenantID: tenantID},
		Name:        name,
		Email:       email,
		Phone:       "0000-0000",
		Tags:        marshalTags(tags),
		Contact:     marshalContact(map[string]interface{}{"channel": "email"}),
	}
	require.NoError(t, testDB.Create(client).Error)
	return client
}
