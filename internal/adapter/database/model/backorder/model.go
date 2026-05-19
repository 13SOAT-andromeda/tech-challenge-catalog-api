package backordermodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"gorm.io/gorm"
)

type Backorder struct {
	ID               uuid.UUID              `gorm:"type:uuid;primaryKey"`
	ProductID        uuid.UUID              `gorm:"type:uuid;not null;index"`
	OrderID          uuid.UUID              `gorm:"type:uuid;not null;index"`
	Quantity         int                    `gorm:"not null"`
	Status           domain.BackorderStatus `gorm:"type:varchar(20);not null"`
	EstimatedArrival *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (b *Backorder) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (m *Backorder) ToDomain() *domain.Backorder {
	return &domain.Backorder{
		ID:               m.ID,
		ProductID:        m.ProductID,
		OrderID:          m.OrderID,
		Quantity:         m.Quantity,
		Status:           m.Status,
		EstimatedArrival: m.EstimatedArrival,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		DeletedAt:        m.DeletedAt,
	}
}

func FromDomain(d *domain.Backorder) *Backorder {
	return &Backorder{
		ID:               d.ID,
		ProductID:        d.ProductID,
		OrderID:          d.OrderID,
		Quantity:         d.Quantity,
		Status:           d.Status,
		EstimatedArrival: d.EstimatedArrival,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
		DeletedAt:        d.DeletedAt,
	}
}
