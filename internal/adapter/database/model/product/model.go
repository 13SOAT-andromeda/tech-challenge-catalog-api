package productmodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"gorm.io/gorm"
)

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Product struct {
	SerialID              int64     `gorm:"autoIncrement;uniqueIndex;not null"`
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name                  string    `gorm:"type:varchar(255);not null"`
	Description           string    `gorm:"type:text"`
	Price                 float64   `gorm:"type:decimal(10,2);not null"`
	StockQuantity         int       `gorm:"not null"`
	MinThreshold          int       `gorm:"not null"`
	MaxStockLevel         int       `gorm:"not null"`
	ReplenishmentStrategy string    `gorm:"type:varchar(50)"`
	SupplierSKU           string    `gorm:"type:varchar(100)"`
	DefaultLeadTimeDays   int
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

func (m *Product) ToDomain() *domain.Product {
	return &domain.Product{
		SerialID:              m.SerialID,
		ID:                    m.ID,
		Name:                  m.Name,
		Description:           m.Description,
		Price:                 m.Price,
		StockQuantity:         m.StockQuantity,
		MinThreshold:          m.MinThreshold,
		MaxStockLevel:         m.MaxStockLevel,
		ReplenishmentStrategy: m.ReplenishmentStrategy,
		SupplierSKU:           m.SupplierSKU,
		DefaultLeadTimeDays:   m.DefaultLeadTimeDays,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		DeletedAt:             m.DeletedAt,
	}
}

func FromDomain(d *domain.Product) *Product {
	return &Product{
		SerialID:              d.SerialID,
		ID:                    d.ID,
		Name:                  d.Name,
		Description:           d.Description,
		Price:                 d.Price,
		StockQuantity:         d.StockQuantity,
		MinThreshold:          d.MinThreshold,
		MaxStockLevel:         d.MaxStockLevel,
		ReplenishmentStrategy: d.ReplenishmentStrategy,
		SupplierSKU:           d.SupplierSKU,
		DefaultLeadTimeDays:   d.DefaultLeadTimeDays,
		CreatedAt:             d.CreatedAt,
		UpdatedAt:             d.UpdatedAt,
		DeletedAt:             d.DeletedAt,
	}
}
