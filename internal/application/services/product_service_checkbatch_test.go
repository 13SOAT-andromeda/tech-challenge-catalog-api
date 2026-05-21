package services_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/services"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func newProductServiceForBatch(mockRepo *mocks.MockProductRepository) *services.ProductService {
	mockBackorder := new(mocks.MockBackorderRepository)
	mockNotif := new(mocks.MockCatalogNotificationPort)
	return services.NewProductService(mockRepo, mockBackorder, mockNotif, nil)
}

func TestProductService_CheckBatch_AllAvailable(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := newProductServiceForBatch(mockRepo)

	products := []*domain.Product{
		{SerialID: 1, ID: uuid.New(), Name: "Óleo", Price: 50.0, StockQuantity: 10},
		{SerialID: 2, ID: uuid.New(), Name: "Filtro", Price: 30.0, StockQuantity: 5},
	}
	mockRepo.On("FindBatchBySerialIDs", []int64{1, 2}).Return(products, nil)

	items := []domain.CheckBatchItem{
		{ID: 1, Qty: 3},
		{ID: 2, Qty: 2},
	}
	available, result, err := svc.CheckBatch(items)

	assert.NoError(t, err)
	assert.True(t, available)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CheckBatch_InsufficientStock(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := newProductServiceForBatch(mockRepo)

	products := []*domain.Product{
		{SerialID: 1, ID: uuid.New(), Name: "Óleo", Price: 50.0, StockQuantity: 2},
	}
	mockRepo.On("FindBatchBySerialIDs", []int64{1}).Return(products, nil)

	items := []domain.CheckBatchItem{{ID: 1, Qty: 5}}
	available, result, err := svc.CheckBatch(items)

	assert.NoError(t, err)
	assert.False(t, available)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CheckBatch_ProductNotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := newProductServiceForBatch(mockRepo)

	mockRepo.On("FindBatchBySerialIDs", []int64{99}).Return([]*domain.Product{}, nil)

	items := []domain.CheckBatchItem{{ID: 99, Qty: 1}}
	available, result, err := svc.CheckBatch(items)

	assert.NoError(t, err)
	assert.False(t, available)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CheckBatch_Empty(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := newProductServiceForBatch(mockRepo)

	mockRepo.On("FindBatchBySerialIDs", []int64{}).Return([]*domain.Product{}, nil)

	available, result, err := svc.CheckBatch([]domain.CheckBatchItem{})

	assert.NoError(t, err)
	assert.True(t, available)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}
