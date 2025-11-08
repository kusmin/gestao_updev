package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/auth"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// SignupInput representa os dados necessários para criar uma empresa.
type SignupInput struct {
	CompanyName     string
	CompanyDocument string
	CompanyPhone    string
	UserName        string
	UserEmail       string
	UserPassword    string
	UserPhone       string
}

// SignupResult é o retorno padrão do fluxo de cadastro.
type SignupResult struct {
	TenantID uuid.UUID
	UserID   uuid.UUID
	Tokens   *auth.TokenPair
}

// Signup cria a empresa, usuário admin e retorna os tokens iniciais.
func (s *Service) Signup(ctx context.Context, input SignupInput) (*SignupResult, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), s.cfg.BcryptCost)
	if err != nil {
		return nil, err
	}

	company := &domain.Company{
		Name:     input.CompanyName,
		Document: input.CompanyDocument,
		Phone:    input.CompanyPhone,
	}

	user := &domain.User{
		Name:         input.UserName,
		Email:        s.sanitizeEmail(input.UserEmail),
		Phone:        input.UserPhone,
		Role:         "admin",
		PasswordHash: string(passwordHash),
		Active:       true,
	}

	err = s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(company).Error; err != nil {
			return err
		}
		user.TenantID = company.ID
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	tokenPair, err := s.jwt.GenerateTokens(user.ID.String(), company.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	return &SignupResult{
		TenantID: company.ID,
		UserID:   user.ID,
		Tokens:   tokenPair,
	}, nil
}

// Login realiza autenticação via email/senha e retorna os tokens.
func (s *Service) Login(ctx context.Context, email, password string) (*domain.User, *auth.TokenPair, error) {
	var user domain.User
	if err := s.dbWithContext(ctx).
		Where("lower(email) = ?", s.sanitizeEmail(email)).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	if !user.Active {
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	tokenPair, err := s.jwt.GenerateTokens(user.ID.String(), user.TenantID.String(), user.Role)
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	_ = s.dbWithContext(ctx).Model(&user).Update("last_login_at", now).Error

	return &user, tokenPair, nil
}

// RefreshTokens valida o refresh token e devolve novos tokens.
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	claims, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	var user domain.User
	if err := s.dbWithContext(ctx).
		Where("id = ?", claims.UserID).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !user.Active {
		return nil, ErrInvalidCredentials
	}

	tokenPair, err := s.jwt.GenerateTokens(user.ID.String(), user.TenantID.String(), user.Role)
	if err != nil {
		return nil, err
	}
	return tokenPair, nil
}
