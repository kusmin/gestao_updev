package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kusmin/gestao_updev/backend/internal/config"
	"github.com/kusmin/gestao_updev/backend/internal/domain"
	"github.com/kusmin/gestao_updev/backend/pkg/database"
)

const (
	demoDocument = "00000000000191"
	demoEmail    = "admin@updev.demo"
)

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
	})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("expose db: %v", err)
	}
	defer sqlDB.Close()

	adminPassword := resolveSeedPassword()
	if err := runSeed(context.Background(), db, adminPassword); err != nil {
		log.Fatalf("seed data: %v", err)
	}

	log.Println("Seed executado com sucesso.")
}

func runSeed(ctx context.Context, db *gorm.DB, adminPassword string) error {
	var company domain.Company
	err := db.WithContext(ctx).Where("document = ?", demoDocument).First(&company).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			company = domain.Company{
				Name:     "UpDev Demo",
				Document: demoDocument,
				Timezone: "America/Sao_Paulo",
			}
			if err := db.WithContext(ctx).Create(&company).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := ensureAdminUser(ctx, db, company, adminPassword); err != nil {
		return err
	}
	if err := ensureClient(ctx, db, company); err != nil {
		return err
	}
	if err := ensureService(ctx, db, company); err != nil {
		return err
	}
	if err := ensureProduct(ctx, db, company); err != nil {
		return err
	}
	return nil
}

func ensureAdminUser(ctx context.Context, db *gorm.DB, company domain.Company, password string) error {
	var user domain.User
	err := db.WithContext(ctx).Where("tenant_id = ? AND email = ?", company.ID, demoEmail).First(&user).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user = domain.User{
		TenantModel:  domain.TenantModel{TenantID: company.ID},
		Name:         "Admin Demo",
		Email:        demoEmail,
		Role:         "admin",
		Phone:        "+55 11 99999-0000",
		PasswordHash: string(hash),
		Active:       true,
	}
	return db.WithContext(ctx).Create(&user).Error
}

func ensureClient(ctx context.Context, db *gorm.DB, company domain.Company) error {
	client := domain.Client{
		TenantModel: domain.TenantModel{TenantID: company.ID},
		Name:        "Cliente Demo",
		Email:       "cliente@demo.com",
		Phone:       "+55 11 98888-0000",
	}
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "email"}},
			UpdateAll: false,
		}).
		Create(&client).Error
}

func ensureService(ctx context.Context, db *gorm.DB, company domain.Company) error {
	service := domain.Service{
		TenantModel:     domain.TenantModel{TenantID: company.ID},
		Name:            "Corte Premium",
		Category:        "Barbearia",
		Description:     "Corte com acabamento e massagem.",
		DurationMinutes: 45,
		Price:           60,
		Color:           "#ff9900",
	}
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "name"}},
			UpdateAll: false,
		}).
		Create(&service).Error
}

func ensureProduct(ctx context.Context, db *gorm.DB, company domain.Company) error {
	product := domain.Product{
		TenantModel: domain.TenantModel{TenantID: company.ID},
		Name:        "Pomada Modeladora",
		SKU:         "POM-001",
		Price:       35,
		Cost:        18,
		StockQty:    50,
		MinStock:    5,
		Description: "Pomada de fixação média.",
	}
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "sku"}},
			UpdateAll: false,
		}).
		Create(&product).Error
}

func resolveSeedPassword() string {
	if value := strings.TrimSpace(os.Getenv("SEED_ADMIN_PASSWORD")); value != "" {
		return value
	}
	const size = 16
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		log.Printf("failed to generate random seed password, falling back to timestamp-based value: %v", err)
		return base64.RawURLEncoding.EncodeToString([]byte("gestao-demo-password"))
	}
	password := base64.RawURLEncoding.EncodeToString(buf)
	log.Printf("SEED_ADMIN_PASSWORD not set; generated random admin password for demo data.")
	return password
}
