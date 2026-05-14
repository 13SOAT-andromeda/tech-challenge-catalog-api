package services_test

import (
	"testing"

	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/services"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryService_Create(t *testing.T) {
	mockRepo := new(mocks.MockCategoryRepository)
	svc := services.NewCategoryService(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*domain.Category")).Return(nil)

	category, err := svc.Create("Test Category")
	assert.NoError(t, err)
	assert.Equal(t, "Test Category", category.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID(t *testing.T) {
	mockRepo := new(mocks.MockCategoryRepository)
	svc := services.NewCategoryService(mockRepo)

	expectedCategory := &domain.Category{ID: "1", Name: "Test Category"}
	mockRepo.On("FindByID", "1").Return(expectedCategory, nil)

	category, err := svc.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, category)
	mockRepo.AssertExpectations(t)
}
