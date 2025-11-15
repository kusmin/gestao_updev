package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/kusmin/gestao_updev/backend/internal/auth"
	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"github.com/kusmin/gestao_updev/backend/internal/repository"
	"github.com/kusmin/gestao_updev/backend/internal/testutil"
)

func newAuthTestService(t *testing.T) *Service {
	t.Helper()
	cfg := &config.Config{
		BcryptCost:       4,
		JWTAccessSecret:  "test-access",
		JWTRefreshSecret: "test-refresh",
		JWTAccessTTL:     time.Minute,
		JWTRefreshTTL:    time.Hour,
	}
	repo := repository.New(testDB)
	jwtMgr := auth.NewJWTManager(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	return New(cfg, repo, jwtMgr, nil)
}

func TestSignupCreatesCompanyAndAdmin(t *testing.T) {
	clearAllData()
	svc := newAuthTestService(t)
	adminPassword := testutil.RandomPassword()

	result, err := svc.Signup(context.Background(), SignupInput{
		CompanyName:     "Tech Corp",
		CompanyDocument: "1234567890",
		CompanyPhone:    "+5511999999999",
		UserName:        "Owner",
		UserEmail:       "OWNER@example.com ",
		UserPassword:    adminPassword,
		UserPhone:       "+5511988887777",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.TenantID)
	assert.NotEqual(t, uuid.Nil, result.UserID)
	require.NotNil(t, result.Tokens)
	assert.NotEmpty(t, result.Tokens.AccessToken)
	assert.NotEmpty(t, result.Tokens.RefreshToken)

	var company domain.Company
	require.NoError(t, testDB.First(&company, "id = ?", result.TenantID).Error)
	assert.Equal(t, "Tech Corp", company.Name)
	assert.Equal(t, "1234567890", company.Document)
	assert.Equal(t, "+5511999999999", company.Phone)

	var user domain.User
	require.NoError(t, testDB.First(&user, "id = ?", result.UserID).Error)
	assert.Equal(t, "Owner", user.Name)
	assert.Equal(t, "owner@example.com", user.Email)
	assert.Equal(t, "admin", user.Role)
	assert.Equal(t, company.ID, user.TenantID)
	assert.True(t, user.Active)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(adminPassword)))
}

func TestLoginHandlesCredentialScenarios(t *testing.T) {
	clearAllData()
	svc := newAuthTestService(t)
	tenant, _ := createTestTenant()

	loginPassword := testutil.RandomPassword()
	hash, err := bcrypt.GenerateFromPassword([]byte(loginPassword), svc.cfg.BcryptCost)
	require.NoError(t, err)
	user := &domain.User{
		TenantModel:  domain.TenantModel{TenantID: tenant.ID},
		Name:         "Login User",
		Email:        "login@example.com",
		PasswordHash: string(hash),
		Role:         "admin",
		Active:       true,
	}
	require.NoError(t, testDB.Create(user).Error)

	t.Run("successfully authenticates and updates last login", func(t *testing.T) {
		found, tokens, err := svc.Login(context.Background(), " LOGIN@example.com ", loginPassword)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, user.ID, found.ID)
		require.NotNil(t, tokens)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)

		var updated domain.User
		require.NoError(t, testDB.First(&updated, user.ID).Error)
		assert.NotNil(t, updated.LastLoginAt)
	})

	t.Run("fails with invalid password", func(t *testing.T) {
		_, _, err := svc.Login(context.Background(), user.Email, testutil.RandomPassword())
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("fails when user inactive", func(t *testing.T) {
		require.NoError(t, testDB.Model(user).Update("active", false).Error)
		_, _, err := svc.Login(context.Background(), user.Email, loginPassword)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
}

func TestRefreshTokensGeneratesNewPair(t *testing.T) {
	clearAllData()
	svc := newAuthTestService(t)
	tenant, _ := createTestTenant()

	user := &domain.User{
		TenantModel:  domain.TenantModel{TenantID: tenant.ID},
		Name:         "Refresh User",
		Email:        "refresh@example.com",
		PasswordHash: "hash",
		Role:         "member",
		Active:       true,
	}
	require.NoError(t, testDB.Create(user).Error)

	refreshToken, err := svc.jwt.GenerateRefreshToken(user.ID.String(), tenant.ID.String(), user.Role)
	require.NoError(t, err)

	tokens, err := svc.RefreshTokens(context.Background(), refreshToken)
	require.NoError(t, err)
	require.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}
