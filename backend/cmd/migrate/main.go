package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/pkg/database"
)

type migrationFile struct {
	Version string
	Path    string
	Name    string
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.New(database.Config{
		URL:             cfg.DatabaseURL,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		MaxOpenConns:    cfg.DBMaxOpenConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
		LogMode:         0,
	})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("expose db: %v", err)
	}
	defer sqlDB.Close()

	ctx := context.Background()
	if err := applyMigrations(ctx, db); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully.")
}

func applyMigrations(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).
		Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`).
		Error; err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	files, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var migrations []migrationFile
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}
		version := strings.SplitN(file.Name(), "_", 2)[0]
		migrations = append(migrations, migrationFile{
			Version: version,
			Path:    filepath.Join("migrations", file.Name()),
			Name:    file.Name(),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for _, m := range migrations {
		applied, err := isMigrationApplied(ctx, db, m.Version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		log.Printf("Applying migration %s...", m.Name)
		contents, err := os.ReadFile(m.Path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", m.Name, err)
		}

		if err := db.WithContext(ctx).Exec(string(contents)).Error; err != nil {
			return fmt.Errorf("execute migration %s: %w", m.Name, err)
		}

		if err := db.WithContext(ctx).
			Exec("INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)", m.Version, time.Now()).
			Error; err != nil {
			return fmt.Errorf("record migration %s: %w", m.Name, err)
		}
		log.Printf("Migration %s applied.", m.Name)
	}

	return nil
}

func isMigrationApplied(ctx context.Context, db *gorm.DB, version string) (bool, error) {
	var count int64
	if err := db.WithContext(ctx).
		Raw("SELECT COUNT(1) FROM schema_migrations WHERE version = ?", version).
		Scan(&count).Error; err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}
	return count > 0, nil
}
