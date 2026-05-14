package categorymodel

import (
	"time"

	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type Category struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Category) ToDomain() *domain.Category {
	return &domain.Category{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromDomain(d *domain.Category) *Category {
	return &Category{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
