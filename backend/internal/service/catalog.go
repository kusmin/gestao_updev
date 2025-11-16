package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

// ServiceInput representa payload de serviços.
type Input struct {
	Name            string
	Category        string
	Description     string
	DurationMinutes int
	Price           float64
	Color           string
	Metadata        map[string]interface{}
}

// ProductInput representa payload de produtos.
type ProductInput struct {
	Name        string
	SKU         string
	Price       float64
	Cost        float64
	StockQty    int
	MinStock    int
	Description string
	Metadata    map[string]interface{}
}

func (s *Service) ListServices(ctx context.Context, tenantID uuid.UUID) ([]domain.Service, error) {
	var services []domain.Service
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("name ASC").
		Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (s *Service) GetService(ctx context.Context, tenantID, serviceID uuid.UUID) (*domain.Service, error) {
	var service domain.Service
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, serviceID).
		First(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *Service) CreateService(ctx context.Context, tenantID uuid.UUID, input Input) (*domain.Service, error) {
	service := &domain.Service{
		TenantModel: domain.TenantModel{
			TenantID: tenantID,
		},
		Name:            input.Name,
		Category:        input.Category,
		Description:     input.Description,
		DurationMinutes: input.DurationMinutes,
		Price:           input.Price,
		Color:           input.Color,
		Metadata:        datatypes.JSONMap(input.Metadata),
	}
	if err := s.dbWithContext(ctx).Create(service).Error; err != nil {
		return nil, err
	}
	return service, nil
}

func (s *Service) UpdateService(ctx context.Context, tenantID, serviceID uuid.UUID, input Input) (*domain.Service, error) {
	var service domain.Service
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, serviceID).
		First(&service).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":             input.Name,
		"category":         input.Category,
		"description":      input.Description,
		"duration_minutes": input.DurationMinutes,
		"price":            input.Price,
		"color":            input.Color,
		"metadata":         datatypes.JSONMap(input.Metadata),
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Service{}).
		Where("tenant_id = ? AND id = ?", tenantID, serviceID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, serviceID).
		First(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *Service) DeleteService(ctx context.Context, tenantID, serviceID uuid.UUID) error {
	return s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, serviceID).
		Delete(&domain.Service{}).Error
}

func (s *Service) ListProducts(ctx context.Context, tenantID uuid.UUID) ([]domain.Product, error) {
	var products []domain.Product
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("name ASC").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Service) GetProduct(ctx context.Context, tenantID, productID uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, productID).
		First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Service) CreateProduct(ctx context.Context, tenantID uuid.UUID, input ProductInput) (*domain.Product, error) {
	product := &domain.Product{
		TenantModel: domain.TenantModel{
			TenantID: tenantID,
		},
		Name:        input.Name,
		SKU:         input.SKU,
		Price:       input.Price,
		Cost:        input.Cost,
		StockQty:    input.StockQty,
		MinStock:    input.MinStock,
		Description: input.Description,
		Metadata:    datatypes.JSONMap(input.Metadata),
	}
	if err := s.dbWithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Service) UpdateProduct(ctx context.Context, tenantID, productID uuid.UUID, input ProductInput) (*domain.Product, error) {
	var product domain.Product
	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, productID).
		First(&product).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":        input.Name,
		"sku":         input.SKU,
		"price":       input.Price,
		"cost":        input.Cost,
		"stock_qty":   input.StockQty,
		"min_stock":   input.MinStock,
		"description": input.Description,
		"metadata":    datatypes.JSONMap(input.Metadata),
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Product{}).
		Where("tenant_id = ? AND id = ?", tenantID, productID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, productID).
		First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Service) DeleteProduct(ctx context.Context, tenantID, productID uuid.UUID) error {
	return s.dbWithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, productID).
		Delete(&domain.Product{}).Error
}

func (s *Service) ListAllProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := s.dbWithContext(ctx).
		Order("name ASC").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

type AdminProductInput struct {
	ProductInput
	TenantID uuid.UUID
}

func (s *Service) AdminCreateProduct(ctx context.Context, input AdminProductInput) (*domain.Product, error) {
	product := &domain.Product{
		TenantModel: domain.TenantModel{
			TenantID: input.TenantID,
		},
		Name:        input.Name,
		SKU:         input.SKU,
		Price:       input.Price,
		Cost:        input.Cost,
		StockQty:    input.StockQty,
		MinStock:    input.MinStock,
		Description: input.Description,
		Metadata:    datatypes.JSONMap(input.Metadata),
	}
	if err := s.dbWithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Service) AdminUpdateProduct(ctx context.Context, productID uuid.UUID, input ProductInput) (*domain.Product, error) {
	var product domain.Product
	if err := s.dbWithContext(ctx).
		First(&product, "id = ?", productID).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":        input.Name,
		"sku":         input.SKU,
		"price":       input.Price,
		"cost":        input.Cost,
		"stock_qty":   input.StockQty,
		"min_stock":   input.MinStock,
		"description": input.Description,
		"metadata":    datatypes.JSONMap(input.Metadata),
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", productID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		First(&product, "id = ?", productID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Service) AdminDeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return s.dbWithContext(ctx).
		Delete(&domain.Product{}, "id = ?", productID).Error
}

func (s *Service) ListAllServices(ctx context.Context) ([]domain.Service, error) {
	var services []domain.Service
	if err := s.dbWithContext(ctx).
		Order("name ASC").
		Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

type AdminServiceInput struct {
	Input
	TenantID uuid.UUID
}

func (s *Service) AdminCreateService(ctx context.Context, input AdminServiceInput) (*domain.Service, error) {
	service := &domain.Service{
		TenantModel: domain.TenantModel{
			TenantID: input.TenantID,
		},
		Name:            input.Name,
		Category:        input.Category,
		Description:     input.Description,
		DurationMinutes: input.DurationMinutes,
		Price:           input.Price,
		Color:           input.Color,
		Metadata:        datatypes.JSONMap(input.Metadata),
	}
	if err := s.dbWithContext(ctx).Create(service).Error; err != nil {
		return nil, err
	}
	return service, nil
}

func (s *Service) AdminUpdateService(ctx context.Context, serviceID uuid.UUID, input Input) (*domain.Service, error) {
	var service domain.Service
	if err := s.dbWithContext(ctx).
		First(&service, "id = ?", serviceID).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":             input.Name,
		"category":         input.Category,
		"description":      input.Description,
		"duration_minutes": input.DurationMinutes,
		"price":            input.Price,
		"color":            input.Color,
		"metadata":         datatypes.JSONMap(input.Metadata),
	}

	if err := s.dbWithContext(ctx).
		Model(&domain.Service{}).
		Where("id = ?", serviceID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.dbWithContext(ctx).
		First(&service, "id = ?", serviceID).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *Service) AdminDeleteService(ctx context.Context, serviceID uuid.UUID) error {
	return s.dbWithContext(ctx).
		Delete(&domain.Service{}, "id = ?", serviceID).Error
}
