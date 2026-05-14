package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackorderStatus string

const (
	BackorderStatusPending   BackorderStatus = "PENDING"
	BackorderStatusConfirmed BackorderStatus = "CONFIRMED"
	BackorderStatusCancelled BackorderStatus = "CANCELLED"
	BackorderStatusFulfilled BackorderStatus = "FULFILLED"
)

type Backorder struct {
	ID               uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID        uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	OrderID          uuid.UUID       `gorm:"type:uuid;not null" json:"order_id"`
	Quantity         int             `gorm:"not null" json:"quantity"`
	Status           BackorderStatus `gorm:"type:varchar(20);not null" json:"status"`
	EstimatedArrival *time.Time      `json:"estimated_arrival"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (b *Backorder) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
