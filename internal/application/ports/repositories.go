package ports

import (
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type CategoryRepository interface {
	Create(category *domain.Category) error
	FindByID(id string) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id string) error
	List() ([]*domain.Category, error)
}

type ProductRepository interface {
	Create(product *domain.Product) error
	FindByID(id uuid.UUID) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uuid.UUID) error
	List() ([]*domain.Product, error)
	UpdateStock(id uuid.UUID, quantity int) error
}

type MaintenanceRepository interface {
	Create(maintenance *domain.Maintenance) error
	FindByID(id uuid.UUID) (*domain.Maintenance, error)
	Update(maintenance *domain.Maintenance) error
	Delete(id uuid.UUID) error
	List() ([]*domain.Maintenance, error)
}

type BackorderRepository interface {
	Create(backorder *domain.Backorder) error
	FindByID(id uuid.UUID) (*domain.Backorder, error)
	UpdateStatus(id uuid.UUID, status domain.BackorderStatus) error
	FindPendingByProductID(productID uuid.UUID) ([]*domain.Backorder, error)
}
