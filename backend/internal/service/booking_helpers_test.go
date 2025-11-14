package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

func seedServiceRecord(t *testing.T, tenantID uuid.UUID, name string, duration int) *domain.Service {
	t.Helper()
	ensureTable(t, &domain.Service{})
	service := &domain.Service{
		TenantModel:     domain.TenantModel{TenantID: tenantID},
		Name:            name,
		DurationMinutes: duration,
		Price:           100,
	}
	require.NoError(t, testDB.Create(service).Error)
	return service
}

func seedProfessionalRecord(t *testing.T, tenantID uuid.UUID, name string) *domain.Professional {
	t.Helper()
	ensureTable(t, &domain.Professional{})
	pro := &domain.Professional{
		TenantModel: domain.TenantModel{TenantID: tenantID},
		Name:        name,
		MaxParallel: 1,
	}
	require.NoError(t, testDB.Create(pro).Error)
	return pro
}

func ensureTable(t *testing.T, model interface{}) {
	t.Helper()
	require.NoError(t, testDB.AutoMigrate(model))
}
