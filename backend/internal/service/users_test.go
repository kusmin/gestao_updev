package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

func TestGetUser(t *testing.T) {
	t.Run("should get user successfully", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)
		user := &domain.User{
			TenantModel:  domain.TenantModel{TenantID: tenant.ID},
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: "some-hash",
			Role:         "admin",
			Active:       true,
		}
		err = testDB.Create(user).Error
		require.NoError(t, err)

		// Act
		foundUser, err := testSvc.GetUser(context.Background(), tenant.ID, user.ID)

		// Assert
		assert.NoError(t, err)
		require.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, "Test User", foundUser.Name)
		assert.Equal(t, "test@example.com", foundUser.Email)
	})

	t.Run("should return error for non-existent user", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)

		// Act
		nonExistentID := uuid.New()
		foundUser, err := testSvc.GetUser(context.Background(), tenant.ID, nonExistentID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("should return error for user in another tenant", func(t *testing.T) {
		clearAllData()
		tenant, err := createTestTenant()
		require.NoError(t, err)
		user := &domain.User{
			TenantModel:  domain.TenantModel{TenantID: tenant.ID},
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: "some-hash",
			Role:         "admin",
			Active:       true,
		}
		err = testDB.Create(user).Error
		require.NoError(t, err)

		otherTenant, err := createTestTenant()
		require.NoError(t, err)

		// Act
		foundUser, err := testSvc.GetUser(context.Background(), otherTenant.ID, user.ID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})
}
