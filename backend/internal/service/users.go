package service

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// UsersFilter parametriza listagem.
type UsersFilter struct {
	Role    string
	Page    int
	PerPage int
}

// CreateUserInput contém dados obrigatórios/ opcionais.
type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	Phone    string
	Role     string
}

// UpdateUserInput altera campos permitidos.
type UpdateUserInput struct {
	Name     *string
	Phone    *string
	Role     *string
	Active   *bool
	Password *string
}

// ListUsers retorna usuários do tenant com paginação.
func (s *Service) ListUsers(ctx context.Context, tenantID uuid.UUID, filter UsersFilter) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	page, perPage := s.clampPagination(filter.Page, filter.PerPage)
	query := s.dbWithContext(ctx).Model(&domain.User{}).
		Where("tenant_id = ?", tenantID)

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.User{}, 0, nil
	}

	if err := query.
		Order("created_at DESC").
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUser busca um usuário por ID.
func (s *Service) GetUser(ctx context.Context, tenantID, userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, userID).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser adiciona um novo colaborador ao tenant.
func (s *Service) CreateUser(ctx context.Context, tenantID uuid.UUID, input CreateUserInput) (*domain.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.cfg.BcryptCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		TenantModel: domain.TenantModel{
			TenantID: tenantID,
		},
		Name:         input.Name,
		Email:        s.sanitizeEmail(input.Email),
		Phone:        input.Phone,
		Role:         input.Role,
		PasswordHash: string(passwordHash),
		Active:       true,
	}

	if err := s.dbWithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser altera campos selecionados de um usuário existente.
func (s *Service) UpdateUser(ctx context.Context, tenantID, userID uuid.UUID, input UpdateUserInput) (*domain.User, error) {
	var user domain.User
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.Role != nil {
		updates["role"] = *input.Role
	}
	if input.Active != nil {
		updates["active"] = *input.Active
	}
	if input.Password != nil && *input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), s.cfg.BcryptCost)
		if err != nil {
			return nil, err
		}
		updates["password_hash"] = string(hash)
	}

	if len(updates) == 0 {
		return &user, nil
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.User{}).
		Where("tenant_id = ? AND id = ?", tenantID, userID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser realiza soft delete do usuário.
func (s *Service) DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error {
	return s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, userID).
		Delete(&domain.User{}).Error
}
