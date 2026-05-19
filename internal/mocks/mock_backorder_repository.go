package mocks

import (
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockBackorderRepository struct {
	mock.Mock
}

func (m *MockBackorderRepository) Create(backorder *domain.Backorder) error {
	args := m.Called(backorder)
	return args.Error(0)
}

func (m *MockBackorderRepository) FindByID(id uuid.UUID) (*domain.Backorder, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Backorder), args.Error(1)
}

func (m *MockBackorderRepository) UpdateStatus(id uuid.UUID, status domain.BackorderStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockBackorderRepository) FindPendingByProductID(productID uuid.UUID) ([]*domain.Backorder, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Backorder), args.Error(1)
}
