package services_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/services"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMaintenanceService_GetBatch_ReturnsList(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	svc := services.NewMaintenanceService(mockRepo)

	maintenances := []*domain.Maintenance{
		{SerialID: 1, ID: uuid.New(), Description: "Troca de óleo", BasePrice: 150.0},
		{SerialID: 2, ID: uuid.New(), Description: "Alinhamento", BasePrice: 80.0},
	}
	mockRepo.On("FindBatchBySerialIDs", []int64{1, 2}).Return(maintenances, nil)

	result, err := svc.GetBatch([]int64{1, 2})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].SerialID)
	assert.Equal(t, "Troca de óleo", result[0].Description)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_GetBatch_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	svc := services.NewMaintenanceService(mockRepo)

	mockRepo.On("FindBatchBySerialIDs", []int64{99}).Return([]*domain.Maintenance{}, nil)

	result, err := svc.GetBatch([]int64{99})

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_GetBatch_Empty(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	svc := services.NewMaintenanceService(mockRepo)

	mockRepo.On("FindBatchBySerialIDs", []int64{}).Return([]*domain.Maintenance{}, nil)

	result, err := svc.GetBatch([]int64{})

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}
