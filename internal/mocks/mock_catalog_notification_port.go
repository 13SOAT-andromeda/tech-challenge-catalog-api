package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockCatalogNotificationPort struct {
	mock.Mock
}

func (m *MockCatalogNotificationPort) NotifyStockReserved(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error {
	args := m.Called(ctx, productID, orderID, quantity)
	return args.Error(0)
}

func (m *MockCatalogNotificationPort) NotifyBackorderCreated(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error {
	args := m.Called(ctx, productID, orderID, quantity)
	return args.Error(0)
}

func (m *MockCatalogNotificationPort) NotifyStockReservationFailed(ctx context.Context, failure domain.StockReservationFailure) error {
	args := m.Called(ctx, failure)
	return args.Error(0)
}
