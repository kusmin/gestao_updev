package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kusmin/gestao_updev/backend/internal/domain"
)

func TestNew(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	repo := New(db)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	repo := New(db)
	assert.Equal(t, db, repo.DB())
}

func TestWithTenant(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	repo := New(db)
	tenantID := uuid.New()

	var user domain.User
	stmt := repo.WithTenant(tenantID).Session(&gorm.Session{DryRun: true}).First(&user).Statement

	whereClauses := stmt.Clauses["where"].Expression.(clause.Where)
	found := false
	for _, expr := range whereClauses.Exprs {
		if eq, ok := expr.(clause.Eq); ok {
			if col, ok := eq.Column.(clause.Column); ok {
				if col.Name == "tenant_id" {
					found = true
				}
			}
		}
	}
	assert.True(t, found)
}

func TestTransaction(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	repo := New(db)

	t.Run("should commit transaction", func(t *testing.T) {
		err := repo.Transaction(context.Background(), func(tx *gorm.DB) error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("should rollback transaction", func(t *testing.T) {
		err := repo.Transaction(context.Background(), func(tx *gorm.DB) error {
			return errors.New("some error")
		})
		assert.Error(t, err)
	})
}
