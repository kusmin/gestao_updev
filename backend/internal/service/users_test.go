package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"github.com/kusmin/gestao_updev/backend/internal/testutil"
)

func TestGetUser(t *testing.T) {
	t.Run("should get user successfully", func(t *testing.T) {
		setupTest(t)
		tenant, err := createTestTenant()
		require.NoError(t, err)
		user := &domain.User{
			TenantModel:  domain.TenantModel{TenantID: tenant.ID},
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: randomPasswordHash(t),
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
		setupTest(t)
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
		setupTest(t)
		tenant, err := createTestTenant()
		require.NoError(t, err)
		user := &domain.User{
			TenantModel:  domain.TenantModel{TenantID: tenant.ID},
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: randomPasswordHash(t),
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

func TestListUsersWithPaginationAndRoleFilter(t *testing.T) {
	setupTest(t)
	tenant, err := createTestTenant()
	require.NoError(t, err)

	// Seed 3 users for the tenant
	createTestUser(t, tenant.ID, "Admin A", "adminA@example.com", "admin")
	time.Sleep(10 * time.Millisecond)
	createTestUser(t, tenant.ID, "Member B", "member@example.com", "member")
	time.Sleep(10 * time.Millisecond)
	lastUser := createTestUser(t, tenant.ID, "Admin C", "adminC@example.com", "admin")

	// Seed another tenant user (should be ignored)
	otherTenant, _ := createTestTenant()
	createTestUser(t, otherTenant.ID, "Other Tenant", "other@example.com", "admin")

	users, total, err := testSvc.ListUsers(
		context.Background(),
		tenant.ID,
		UsersFilter{Page: 1, PerPage: 2},
	)

	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, users, 2, "should respect pagination limit")
	assert.Equal(t, lastUser.ID, users[0].ID, "expects newest first")

	admins, adminTotal, err := testSvc.ListUsers(
		context.Background(),
		tenant.ID,
		UsersFilter{Role: "admin"},
	)
	require.NoError(t, err)
	require.Equal(t, int64(2), adminTotal)
	require.Len(t, admins, 2)
}

func TestCreateUserSanitizesEmailAndHashesPassword(t *testing.T) {
	setupTest(t)
	tenant, err := createTestTenant()
	require.NoError(t, err)

	password := testutil.RandomPassword()
	user, err := testSvc.CreateUser(context.Background(), tenant.ID, CreateUserInput{
		Name:     "New User",
		Email:    "  USER@Example.com ",
		Password: password,
		Phone:    "123",
		Role:     "admin",
	})

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "user@example.com", user.Email)
	assert.True(t, user.Active)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)))
}

func TestUpdateUserAllowsPartialChanges(t *testing.T) {
	setupTest(t)
	tenant, err := createTestTenant()
	require.NoError(t, err)
	user := createTestUser(t, tenant.ID, "Original", "original@example.com", "member")

	newName := "Updated"
	newPhone := "555-1234"
	newRole := "admin"
	active := false
	newPassword := testutil.RandomPassword()

	updated, err := testSvc.UpdateUser(context.Background(), tenant.ID, user.ID, UpdateUserInput{
		Name:     &newName,
		Phone:    &newPhone,
		Role:     &newRole,
		Active:   &active,
		Password: &newPassword,
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, newPhone, updated.Phone)
	assert.Equal(t, newRole, updated.Role)
	assert.Equal(t, active, updated.Active)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte(newPassword)))
}

func TestDeleteUserPerformsSoftDelete(t *testing.T) {
	setupTest(t)
	tenant, err := createTestTenant()
	require.NoError(t, err)
	user := createTestUser(t, tenant.ID, "To Delete", "delete@example.com", "member")

	err = testSvc.DeleteUser(context.Background(), tenant.ID, user.ID)
	require.NoError(t, err)

	var count int64
	testDB.Model(&domain.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)

	var deleted domain.User
	testDB.Unscoped().First(&deleted, user.ID)
	assert.NotNil(t, deleted.DeletedAt)
}

func createTestUser(t *testing.T, tenantID uuid.UUID, name, email, role string) *domain.User {
	t.Helper()
	hash := randomPasswordHash(t)
	user := &domain.User{
		TenantModel:  domain.TenantModel{TenantID: tenantID},
		Name:         name,
		Email:        email,
		Role:         role,
		PasswordHash: hash,
		Active:       true,
	}
	require.NoError(t, testDB.Create(user).Error)
	return user
}

func randomPasswordHash(t *testing.T) string {
	t.Helper()
	password := testutil.RandomPassword()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}
