package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Maintenance struct {
	SerialID          int64          `gorm:"autoIncrement;uniqueIndex;not null" json:"serial_id"`
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Description       string         `gorm:"type:text;not null" json:"description"`
	BasePrice         float64        `gorm:"type:decimal(10,2);not null" json:"base_price"`
	EstimatedDuration int            `json:"estimated_duration"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Maintenance) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
