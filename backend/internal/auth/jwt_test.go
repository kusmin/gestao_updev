package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateTokensReturnsPair(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("access-secret", "refresh-secret", 2*time.Minute, time.Hour)

	pair, err := manager.GenerateTokens("user-123", "tenant-abc", "admin")

	require.NoError(t, err)
	require.NotEmpty(t, pair.AccessToken)
	require.NotEmpty(t, pair.RefreshToken)
	require.Equal(t, int64((2 * time.Minute).Seconds()), pair.ExpiresIn)
}

func TestValidateAccessToken(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("access-secret", "refresh-secret", time.Minute, time.Hour)
	token, err := manager.GenerateAccessToken("user-123", "tenant-abc", "admin")
	require.NoError(t, err)

	claims, err := manager.ValidateAccessToken(token)

	require.NoError(t, err)
	require.Equal(t, "user-123", claims.UserID)
	require.Equal(t, "tenant-abc", claims.TenantID)
	require.Equal(t, "admin", claims.Role)
	require.WithinDuration(t, time.Now().Add(time.Minute), claims.RegisteredClaims.ExpiresAt.Time, 5*time.Second)
}

func TestValidateRefreshTokenRequiresMatchingSecret(t *testing.T) {
	t.Parallel()

	issuer := NewJWTManager("access-secret", "refresh-secret", time.Minute, time.Hour)
	token, err := issuer.GenerateRefreshToken("user-123", "tenant-abc", "admin")
	require.NoError(t, err)

	otherManager := NewJWTManager("access-secret", "other-refresh", time.Minute, time.Hour)
	claims, err := otherManager.ValidateRefreshToken(token)

	require.Error(t, err)
	require.Nil(t, claims)
}

func TestValidateAccessTokenFailsWhenExpired(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("access-secret", "refresh-secret", -time.Minute, time.Hour)
	token, err := manager.GenerateAccessToken("user-123", "tenant-abc", "admin")
	require.NoError(t, err)

	claims, err := manager.ValidateAccessToken(token)

	require.Error(t, err)
	require.Nil(t, claims)
}
