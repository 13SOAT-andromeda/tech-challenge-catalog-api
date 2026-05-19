package repository

import (
	"github.com/google/uuid"
	backordermodel "github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/model/backorder"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"gorm.io/gorm"
)

type BackorderRepositoryGorm struct {
	db *gorm.DB
}

func NewBackorderRepository(db *gorm.DB) ports.BackorderRepository {
	return &BackorderRepositoryGorm{db: db}
}

func (r *BackorderRepositoryGorm) Create(b *domain.Backorder) error {
	model := backordermodel.FromDomain(b)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	*b = *model.ToDomain()
	return nil
}

func (r *BackorderRepositoryGorm) FindByID(id uuid.UUID) (*domain.Backorder, error) {
	var model backordermodel.Backorder
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *BackorderRepositoryGorm) UpdateStatus(id uuid.UUID, status domain.BackorderStatus) error {
	return r.db.Model(&backordermodel.Backorder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BackorderRepositoryGorm) FindPendingByProductID(productID uuid.UUID) ([]*domain.Backorder, error) {
	var models []backordermodel.Backorder
	if err := r.db.Where("product_id = ? AND status = ?", productID, domain.BackorderStatusPending).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.Backorder, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}
