package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusDone      = "done"
	BookingStatusCanceled  = "canceled"
)

const (
	SalesOrderStatusDraft     = "draft"
	SalesOrderStatusConfirmed = "confirmed"
	SalesOrderStatusPaid      = "paid"
	SalesOrderStatusCanceled  = "canceled"
)

const (
	InventoryMovementIn         = "in"
	InventoryMovementOut        = "out"
	InventoryMovementAdjustment = "adjustment"
)

// BaseModel consolida campos comuns de auditoria.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TenantModel adiciona referência obrigatória ao tenant.
type TenantModel struct {
	BaseModel
	TenantID uuid.UUID `gorm:"type:uuid;index" json:"tenant_id"`
}

type Company struct {
	BaseModel
	Name     string            `gorm:"size:140;not null" json:"name"`
	Document string            `gorm:"size:36;uniqueIndex" json:"document"`
	Timezone string            `gorm:"size:60;default:'America/Sao_Paulo'" json:"timezone"`
	Settings datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"settings"`
	Phone    string            `gorm:"size:32" json:"phone"`
	Email    string            `gorm:"size:120" json:"email"`
	Metadata datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	Users    []User            `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	Clients  []Client          `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	Services []Service         `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	Products []Product         `gorm:"foreignKey:TenantID;references:ID" json:"-"`
}

type User struct {
	TenantModel
	Name         string            `gorm:"size:120;not null" json:"name"`
	Email        string            `gorm:"size:160;not null;index:idx_users_email_tenant,unique" json:"email"`
	Phone        string            `gorm:"size:32" json:"phone"`
	Role         string            `gorm:"size:32;not null" json:"role"`
	PasswordHash string            `gorm:"size:255;not null" json:"-"`
	Profile      datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"profile"`
	Active       bool              `gorm:"default:true" json:"active"`
	LastLoginAt  *time.Time        `json:"last_login_at,omitempty"`
}

type Client struct {
	TenantModel
	Name    string            `gorm:"size:160;not null" json:"name"`
	Email   string            `gorm:"size:160" json:"email"`
	Phone   string            `gorm:"size:32" json:"phone"`
	Notes   string            `gorm:"type:text" json:"notes"`
	Tags    datatypes.JSON    `gorm:"type:jsonb" json:"tags"`
	Contact datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"contact"`
}

type Professional struct {
	TenantModel
	UserID       *uuid.UUID         `gorm:"type:uuid" json:"user_id"`
	Name         string             `gorm:"size:160;not null" json:"name"`
	Specialties  datatypes.JSON     `gorm:"type:jsonb" json:"specialties"`
	MaxParallel  int                `gorm:"default:1" json:"max_parallel"`
	Active       bool               `gorm:"default:true" json:"active"`
	Availability []AvailabilityRule `gorm:"foreignKey:ProfessionalID;references:ID"`
}

type AvailabilityRule struct {
	TenantModel
	ProfessionalID uuid.UUID `gorm:"type:uuid;not null;index" json:"professional_id"`
	Weekday        int       `gorm:"not null" json:"weekday"`
	StartTime      string    `gorm:"size:8;not null" json:"start_time"`
	EndTime        string    `gorm:"size:8;not null" json:"end_time"`
}

type Service struct {
	TenantModel
	Name            string            `gorm:"size:160;not null;index:idx_services_name_tenant,unique" json:"name"`
	Category        string            `gorm:"size:80" json:"category"`
	Description     string            `gorm:"type:text" json:"description"`
	DurationMinutes int               `gorm:"not null" json:"duration_minutes"`
	Price           float64           `gorm:"type:numeric(12,2);not null" json:"price"`
	Color           string            `gorm:"size:16" json:"color"`
	Metadata        datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata"`
}

type Product struct {
	TenantModel
	Name        string            `gorm:"size:160;not null" json:"name"`
	SKU         string            `gorm:"size:80;not null;index:idx_products_sku_tenant,unique" json:"sku"`
	Price       float64           `gorm:"type:numeric(12,2);not null" json:"price"`
	Cost        float64           `gorm:"type:numeric(12,2)" json:"cost"`
	StockQty    int               `gorm:"default:0" json:"stock_qty"`
	MinStock    int               `gorm:"default:0" json:"min_stock"`
	Description string            `gorm:"type:text" json:"description"`
	Metadata    datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata"`
}

type InventoryMovement struct {
	TenantModel
	ProductID uuid.UUID  `gorm:"type:uuid;not null;index" json:"product_id"`
	OrderID   *uuid.UUID `gorm:"type:uuid" json:"order_id"`
	Type      string     `gorm:"size:16;not null" json:"type"`
	Quantity  int        `gorm:"not null" json:"quantity"`
	Reason    string     `gorm:"size:160" json:"reason"`
}

type Booking struct {
	TenantModel
	ClientID       uuid.UUID         `gorm:"type:uuid;not null;index" json:"client_id"`
	ProfessionalID uuid.UUID         `gorm:"type:uuid;not null;index" json:"professional_id"`
	ServiceID      uuid.UUID         `gorm:"type:uuid;not null" json:"service_id"`
	Status         string            `gorm:"size:32;not null" json:"status"`
	StartAt        time.Time         `gorm:"not null" json:"start_at"`
	EndAt          time.Time         `gorm:"not null" json:"end_at"`
	Notes          string            `gorm:"type:text" json:"notes"`
	Metadata       datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata"`
}

type SalesOrder struct {
	TenantModel
	ClientID    uuid.UUID   `gorm:"type:uuid;not null;index" json:"client_id"`
	BookingID   *uuid.UUID  `gorm:"type:uuid" json:"booking_id"`
	Status      string      `gorm:"size:32;not null" json:"status"`
	PaymentType string      `gorm:"size:32" json:"payment_method"`
	Total       float64     `gorm:"type:numeric(12,2);default:0" json:"total"`
	Discount    float64     `gorm:"type:numeric(12,2);default:0" json:"discount"`
	Notes       string      `gorm:"type:text" json:"notes"`
	Items       []SalesItem `gorm:"foreignKey:OrderID;references:ID" json:"items"`
}

type SalesItem struct {
	TenantModel
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ItemType  string    `gorm:"size:16;not null" json:"item_type"`
	ItemRefID uuid.UUID `gorm:"type:uuid;not null" json:"item_ref_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	UnitPrice float64   `gorm:"type:numeric(12,2);not null" json:"unit_price"`
}

type Payment struct {
	TenantModel
	OrderID uuid.UUID         `gorm:"type:uuid;not null;index" json:"order_id"`
	Method  string            `gorm:"size:32;not null" json:"method"`
	Amount  float64           `gorm:"type:numeric(12,2);not null" json:"amount"`
	PaidAt  time.Time         `gorm:"not null" json:"paid_at"`
	Details datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"details"`
}

type AuditLog struct {
	TenantModel
	Entity   string            `gorm:"size:64;not null" json:"entity"`
	Action   string            `gorm:"size:64;not null" json:"action"`
	ActorID  uuid.UUID         `gorm:"type:uuid;not null" json:"actor_id"`
	Metadata datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata"`
}
