package repository

import (
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/model/category"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepositoryGorm struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return &CategoryRepositoryGorm{db: db}
}

func (r *CategoryRepositoryGorm) Create(categoryDomain *domain.Category) error {
	model := categorymodel.FromDomain(categoryDomain)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	*categoryDomain = *model.ToDomain()
	return nil
}

func (r *CategoryRepositoryGorm) FindByID(id string) (*domain.Category, error) {
	var model categorymodel.Category
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *CategoryRepositoryGorm) Update(categoryDomain *domain.Category) error {
	model := categorymodel.FromDomain(categoryDomain)
	return r.db.Save(model).Error
}

func (r *CategoryRepositoryGorm) Delete(id string) error {
	return r.db.Delete(&categorymodel.Category{}, "id = ?", id).Error
}

func (r *CategoryRepositoryGorm) List() ([]*domain.Category, error) {
	var models []categorymodel.Category
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	categories := make([]*domain.Category, len(models))
	for i, m := range models {
		categories[i] = m.ToDomain()
	}
	return categories, nil
}
