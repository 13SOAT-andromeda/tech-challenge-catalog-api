package mocks

import (
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockMaintenanceRepository struct {
	mock.Mock
}

func (m *MockMaintenanceRepository) Create(maintenance *domain.Maintenance) error {
	args := m.Called(maintenance)
	return args.Error(0)
}

func (m *MockMaintenanceRepository) FindByID(id uuid.UUID) (*domain.Maintenance, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Maintenance), args.Error(1)
}

func (m *MockMaintenanceRepository) Update(maintenance *domain.Maintenance) error {
	args := m.Called(maintenance)
	return args.Error(0)
}

func (m *MockMaintenanceRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMaintenanceRepository) List() ([]*domain.Maintenance, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Maintenance), args.Error(1)
}

func (m *MockMaintenanceRepository) FindBatchBySerialIDs(ids []int64) ([]*domain.Maintenance, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Maintenance), args.Error(1)
}
