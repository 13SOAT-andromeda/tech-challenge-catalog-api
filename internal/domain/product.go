package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	SerialID              int64          `gorm:"autoIncrement;uniqueIndex;not null" json:"serial_id"`
	ID                    uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name                  string         `gorm:"type:varchar(255);not null" json:"name"`
	Description           string         `gorm:"type:text" json:"description"`
	Price                 float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQuantity         int            `gorm:"not null" json:"stock_quantity"`
	MinThreshold          int            `gorm:"not null" json:"min_threshold"`
	MaxStockLevel         int            `gorm:"not null" json:"max_stock_level"`
	ReplenishmentStrategy string         `gorm:"type:varchar(50)" json:"replenishment_strategy"`
	SupplierSKU           string         `gorm:"type:varchar(100)" json:"supplier_sku"`
	DefaultLeadTimeDays   int            `json:"default_lead_time_days"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
