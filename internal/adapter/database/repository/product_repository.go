package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/model/product"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepositoryGorm struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &ProductRepositoryGorm{db: db}
}

func (r *ProductRepositoryGorm) Create(productDomain *domain.Product) error {
	model := productmodel.FromDomain(productDomain)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	*productDomain = *model.ToDomain()
	return nil
}

func (r *ProductRepositoryGorm) FindByID(id uuid.UUID) (*domain.Product, error) {
	var model productmodel.Product
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ProductRepositoryGorm) Update(productDomain *domain.Product) error {
	model := productmodel.FromDomain(productDomain)
	return r.db.Save(model).Error
}

func (r *ProductRepositoryGorm) Delete(id uuid.UUID) error {
	return r.db.Delete(&productmodel.Product{}, "id = ?", id).Error
}

func (r *ProductRepositoryGorm) List() ([]*domain.Product, error) {
	var models []productmodel.Product
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	products := make([]*domain.Product, len(models))
	for i, m := range models {
		products[i] = m.ToDomain()
	}
	return products, nil
}

func (r *ProductRepositoryGorm) UpdateStock(id uuid.UUID, quantity int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var model productmodel.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&model, "id = ?", id).Error; err != nil {
			return err
		}

		model.StockQuantity += quantity
		return tx.Save(&model).Error
	})
}

func (r *ProductRepositoryGorm) FindBatchBySerialIDs(ids []int64) ([]*domain.Product, error) {
	var models []productmodel.Product
	if err := r.db.Where("serial_id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}
	products := make([]*domain.Product, len(models))
	for i, m := range models {
		products[i] = m.ToDomain()
	}
	return products, nil
}
