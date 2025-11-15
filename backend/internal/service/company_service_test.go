package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

type fakeCompanyRepo struct {
	listResp      []domain.Company
	listErr       error
	findResp      *domain.Company
	findErr       error
	createErr     error
	updateErr     error
	deleteErr     error
	createdRecord *domain.Company
	updatedRecord *domain.Company
	deletedID     uuid.UUID
}

func (f *fakeCompanyRepo) ListAll(ctx context.Context) ([]domain.Company, error) {
	return f.listResp, f.listErr
}

func (f *fakeCompanyRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	if f.findResp != nil {
		return f.findResp, f.findErr
	}
	return nil, f.findErr
}

func (f *fakeCompanyRepo) Create(ctx context.Context, company *domain.Company) error {
	f.createdRecord = company
	return f.createErr
}

func (f *fakeCompanyRepo) Update(ctx context.Context, company *domain.Company) error {
	f.updatedRecord = company
	return f.updateErr
}

func (f *fakeCompanyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	f.deletedID = id
	return f.deleteErr
}

func TestCompanyServiceCreateCompany(t *testing.T) {
	repo := &fakeCompanyRepo{}
	svc := NewCompanyService(repo)

	result, err := svc.CreateCompany(context.Background(), CreateCompanyInput{
		Name:     "Acme",
		Document: "123",
		Phone:    "9999",
		Email:    "contact@acme.dev",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Acme", repo.createdRecord.Name)
	assert.Equal(t, "contact@acme.dev", repo.createdRecord.Email)

	repo.createErr = errors.New("boom")
	_, err = svc.CreateCompany(context.Background(), CreateCompanyInput{Name: "Fail"})
	assert.EqualError(t, err, "boom")
}

func TestCompanyServiceUpdateCompany(t *testing.T) {
	originalID := uuid.New()
	repo := &fakeCompanyRepo{
		findResp: &domain.Company{
			BaseModel: domain.BaseModel{ID: originalID},
			Name:      "Original",
			Document:  "DOC",
			Phone:     "1111",
			Email:     "old@example.com",
		},
	}
	svc := NewCompanyService(repo)

	updated, err := svc.UpdateCompany(context.Background(), originalID, UpdateCompanyInput{
		Name:     "New Name",
		Document: "",
		Phone:    "2222",
		Email:    "new@example.com",
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "New Name", repo.updatedRecord.Name)
	assert.Equal(t, "DOC", repo.updatedRecord.Document, "empty fields should not override existing value")
	assert.Equal(t, "2222", repo.updatedRecord.Phone)
	assert.Equal(t, "new@example.com", repo.updatedRecord.Email)

	repo.findErr = errors.New("not found")
	_, err = svc.UpdateCompany(context.Background(), uuid.New(), UpdateCompanyInput{Name: "Should fail"})
	assert.EqualError(t, err, "not found")
}

func TestCompanyServiceDeleteCompany(t *testing.T) {
	targetID := uuid.New()
	repo := &fakeCompanyRepo{}
	svc := NewCompanyService(repo)

	require.NoError(t, svc.DeleteCompany(context.Background(), targetID))
	assert.Equal(t, targetID, repo.deletedID)

	repo.deleteErr = errors.New("delete failed")
	err := svc.DeleteCompany(context.Background(), targetID)
	assert.EqualError(t, err, "delete failed")
}
