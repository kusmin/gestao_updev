package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"github.com/kusmin/gestao_updev/backend/internal/repository"
)

type CompanyService struct {
	repo *repository.CompanyRepository
}

func NewCompanyService(repo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) ListAllCompanies(ctx context.Context) ([]domain.Company, error) {
	return s.repo.ListAll(ctx)
}

func (s *CompanyService) GetCompanyByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	return s.repo.FindByID(ctx, id)
}

type CreateCompanyInput struct {
	Name     string `json:"name" binding:"required"`
	Document string `json:"document"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func (s *CompanyService) CreateCompany(ctx context.Context, input CreateCompanyInput) (*domain.Company, error) {
	company := &domain.Company{
		Name:     input.Name,
		Document: input.Document,
		Phone:    input.Phone,
		Email:    input.Email,
	}

	if err := s.repo.Create(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

type UpdateCompanyInput struct {
	Name     string `json:"name"`
	Document string `json:"document"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func (s *CompanyService) UpdateCompany(ctx context.Context, id uuid.UUID, input UpdateCompanyInput) (*domain.Company, error) {
	company, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		company.Name = input.Name
	}
	if input.Document != "" {
		company.Document = input.Document
	}
	if input.Phone != "" {
		company.Phone = input.Phone
	}
	if input.Email != "" {
		company.Email = input.Email
	}

	if err := s.repo.Update(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
