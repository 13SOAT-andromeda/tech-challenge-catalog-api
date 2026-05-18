package maintenancemodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"gorm.io/gorm"
)

func (m *Maintenance) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type Maintenance struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	Description       string    `gorm:"type:text;not null"`
	BasePrice         float64   `gorm:"type:decimal(10,2);not null"`
	EstimatedDuration int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (m *Maintenance) ToDomain() *domain.Maintenance {
	return &domain.Maintenance{
		ID:                m.ID,
		Description:       m.Description,
		BasePrice:         m.BasePrice,
		EstimatedDuration: m.EstimatedDuration,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		DeletedAt:         m.DeletedAt,
	}
}

func FromDomain(d *domain.Maintenance) *Maintenance {
	return &Maintenance{
		ID:                d.ID,
		Description:       d.Description,
		BasePrice:         d.BasePrice,
		EstimatedDuration: d.EstimatedDuration,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
		DeletedAt:         d.DeletedAt,
	}
}
