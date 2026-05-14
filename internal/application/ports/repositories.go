package ports

import "github.com/juliovaz/tech-challenge-catalog-api/internal/domain"

type CategoryRepository interface {
	Create(category *domain.Category) error
	FindByID(id string) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id string) error
	List() ([]*domain.Category, error)
}
