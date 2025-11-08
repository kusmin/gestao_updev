package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository provê acesso centralizado ao banco.
type Repository struct {
	db *gorm.DB
}

// New instancia um repositório baseado no GORM.
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// DB expõe o ponteiro cru para cenários avançados/transactions.
func (r *Repository) DB() *gorm.DB {
	return r.db
}

// WithTenant injeta o filtro de tenant padrão em uma query.
func (r *Repository) WithTenant(tenantID uuid.UUID) *gorm.DB {
	return r.db.Where("tenant_id = ?", tenantID)
}

// Transaction executa uma função dentro de uma transação.
func (r *Repository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
