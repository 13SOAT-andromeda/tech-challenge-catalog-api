package mocks

import (
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindByID(id string) (*domain.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryRepository) List() ([]*domain.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}
