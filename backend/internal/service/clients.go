package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// ClientsFilter define filtros disponíveis.
type ClientsFilter struct {
	Search  string
	Tags    []string
	Page    int
	PerPage int
}

// ClientInput concentra dados editáveis.
type ClientInput struct {
	Name    string
	Email   string
	Phone   string
	Notes   string
	Tags    []string
	Contact map[string]interface{}
}

// ClientStats agrega dados de histórico.
type ClientStats struct {
	TotalBookings int64      `json:"total_bookings"`
	TotalSpent    float64    `json:"total_spent"`
	LastBooking   *time.Time `json:"last_booking"`
}

// ListClients retorna clientes com paginação/filtros básicos.
func (s *Service) ListClients(ctx context.Context, tenantID uuid.UUID, filter ClientsFilter) ([]domain.Client, int64, error) {
	var clients []domain.Client
	var total int64

	page, perPage := s.clampPagination(filter.Page, filter.PerPage)

	query := s.dbWithContext(ctx).
		Model(&domain.Client{}).
		Where("tenant_id = ?", tenantID)

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("(name ILIKE ? OR email ILIKE ? OR phone ILIKE ?)", like, like, like)
	}

	for _, tag := range filter.Tags {
		tagJSON, _ := json.Marshal([]string{tag})
		query = query.Where("tags @> ?", datatypes.JSON(tagJSON))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Client{}, 0, nil
	}

	if err := query.
		Order("created_at DESC").
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

// CreateClient adiciona um novo cliente.
func (s *Service) CreateClient(ctx context.Context, tenantID uuid.UUID, input ClientInput) (*domain.Client, error) {
	client := &domain.Client{
		TenantModel: domain.TenantModel{
			TenantID: tenantID,
		},
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Notes:   input.Notes,
		Tags:    marshalTags(input.Tags),
		Contact: marshalContact(input.Contact),
	}
	if err := s.dbWithContext(ctx).Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

// GetClient retorna cliente + estatísticas básicas.
func (s *Service) GetClient(ctx context.Context, tenantID, clientID uuid.UUID) (*domain.Client, *ClientStats, error) {
	var client domain.Client
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, clientID).
		First(&client).Error; err != nil {
		return nil, nil, err
	}

	var stats ClientStats

	err := s.dbWithContext(ctx).
		Model(&domain.Booking{}).
		Where("tenant_id = ? AND client_id = ?", tenantID, clientID).
		Count(&stats.TotalBookings).Error
	if err != nil {
		return &client, nil, err
	}

	var lastBooking domain.Booking
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND client_id = ?", tenantID, clientID).
		Order("start_at DESC").
		First(&lastBooking).Error; err == nil {
		stats.LastBooking = &lastBooking.StartAt
	} else if err != gorm.ErrRecordNotFound {
		return &client, nil, err
	}

	type totalSpent struct {
		Sum float64
	}
	var ts totalSpent
	if err := s.dbWithContext(ctx).
		Model(&domain.Payment{}).
		Select("COALESCE(SUM(amount),0) AS sum").
		Joins("JOIN sales_orders so ON so.id = payments.order_id").
		Where("payments.tenant_id = ? AND so.client_id = ?", tenantID, clientID).
		Scan(&ts).Error; err != nil {
		return &client, nil, err
	}
	stats.TotalSpent = ts.Sum

	return &client, &stats, nil
}

// UpdateClient realiza alterações completas.
func (s *Service) UpdateClient(ctx context.Context, tenantID, clientID uuid.UUID, input ClientInput) (*domain.Client, error) {
	var client domain.Client
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, clientID).
		First(&client).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":    input.Name,
		"email":   input.Email,
		"phone":   input.Phone,
		"notes":   input.Notes,
		"contact": marshalContact(input.Contact),
		"tags":    marshalTags(input.Tags),
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Client{}).
		Where("tenant_id = ? AND id = ?", tenantID, clientID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, clientID).
		First(&client).Error; err != nil {
		return nil, err
	}

	return &client, nil
}

// DeleteClient faz soft delete.
func (s *Service) DeleteClient(ctx context.Context, tenantID, clientID uuid.UUID) error {
	return s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, clientID).
		Delete(&domain.Client{}).Error
}

func (s *Service) ListAllClients(ctx context.Context, filter ClientsFilter) ([]domain.Client, int64, error) {
	var clients []domain.Client
	var total int64

	page, perPage := s.clampPagination(filter.Page, filter.PerPage)

	query := s.dbWithContext(ctx).Model(&domain.Client{})

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("(name ILIKE ? OR email ILIKE ? OR phone ILIKE ?)", like, like, like)
	}

	for _, tag := range filter.Tags {
		tagJSON, _ := json.Marshal([]string{tag})
		query = query.Where("tags @> ?", datatypes.JSON(tagJSON))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Client{}, 0, nil
	}

	if err := query.
		Order("created_at DESC").
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

type AdminClientInput struct {
	ClientInput
	TenantID uuid.UUID
}

func (s *Service) AdminCreateClient(ctx context.Context, input AdminClientInput) (*domain.Client, error) {
	client := &domain.Client{
		TenantModel: domain.TenantModel{
			TenantID: input.TenantID,
		},
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Notes:   input.Notes,
		Tags:    marshalTags(input.Tags),
		Contact: marshalContact(input.Contact),
	}
	if err := s.dbWithContext(ctx).Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Service) AdminUpdateClient(ctx context.Context, clientID uuid.UUID, input ClientInput) (*domain.Client, error) {
	var client domain.Client
	if err := s.dbWithContext(ctx).
		First(&client, "id = ?", clientID).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":    input.Name,
		"email":   input.Email,
		"phone":   input.Phone,
		"notes":   input.Notes,
		"contact": marshalContact(input.Contact),
		"tags":    marshalTags(input.Tags),
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Client{}).
		Where("id = ?", clientID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		First(&client, "id = ?", clientID).Error; err != nil {
		return nil, err
	}

	return &client, nil
}

func (s *Service) AdminDeleteClient(ctx context.Context, clientID uuid.UUID) error {
	return s.dbWithContext(ctx).
		Delete(&domain.Client{}, "id = ?", clientID).Error
}


func marshalTags(tags []string) datatypes.JSON {
	if tags == nil {
		return datatypes.JSON([]byte("[]"))
	}
	b, _ := json.Marshal(tags)
	return datatypes.JSON(b)
}

func marshalContact(contact map[string]interface{}) datatypes.JSONMap {
	if contact == nil {
		return datatypes.JSONMap{}
	}
	return datatypes.JSONMap(contact)
}
