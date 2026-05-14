package services

import (
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type CategoryService struct {
	repo ports.CategoryRepository
}

func NewCategoryService(repo ports.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(name string) (*domain.Category, error) {
	category, err := domain.NewCategory(name)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) GetByID(id string) (*domain.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) List() ([]*domain.Category, error) {
	return s.repo.List()
}

func (s *CategoryService) Update(id, name string) (*domain.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	category.Name = name
	if err := s.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Delete(id string) error {
	return s.repo.Delete(id)
}
