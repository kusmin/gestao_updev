package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims representa os dados adicionais embutidos no token.
type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair agrupa access + refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// JWTManager encapsula geração e validação de tokens.
type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// NewJWTManager cria o gerador/validador padrão.
func NewJWTManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

// GenerateTokens devolve o par access/refresh.
func (m *JWTManager) GenerateTokens(userID, tenantID, role string) (*TokenPair, error) {
	access, err := m.GenerateAccessToken(userID, tenantID, role)
	if err != nil {
		return nil, err
	}
	refresh, err := m.GenerateRefreshToken(userID, tenantID, role)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(m.accessTTL.Seconds()),
	}, nil
}

// GenerateAccessToken cria um token de curta duração.
func (m *JWTManager) GenerateAccessToken(userID, tenantID, role string) (string, error) {
	return m.generateToken(userID, tenantID, role, m.accessTTL, m.accessSecret)
}

// GenerateRefreshToken cria um token de renovação.
func (m *JWTManager) GenerateRefreshToken(userID, tenantID, role string) (string, error) {
	return m.generateToken(userID, tenantID, role, m.refreshTTL, m.refreshSecret)
}

func (m *JWTManager) generateToken(userID, tenantID, role string, ttl time.Duration, secret []byte) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ValidateAccessToken valida e extrai as claims do access token.
func (m *JWTManager) ValidateAccessToken(token string) (*Claims, error) {
	return m.parseClaims(token, m.accessSecret)
}

// ValidateRefreshToken valida e extrai as claims do refresh token.
func (m *JWTManager) ValidateRefreshToken(token string) (*Claims, error) {
	return m.parseClaims(token, m.refreshSecret)
}

func (m *JWTManager) parseClaims(token string, secret []byte) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
