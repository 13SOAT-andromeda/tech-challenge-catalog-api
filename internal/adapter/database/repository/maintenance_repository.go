package repository

import (
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/database/model/maintenance"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
	"gorm.io/gorm"
)

type MaintenanceRepositoryGorm struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) ports.MaintenanceRepository {
	return &MaintenanceRepositoryGorm{db: db}
}

func (r *MaintenanceRepositoryGorm) Create(maintenanceDomain *domain.Maintenance) error {
	model := maintenancemodel.FromDomain(maintenanceDomain)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	*maintenanceDomain = *model.ToDomain()
	return nil
}

func (r *MaintenanceRepositoryGorm) FindByID(id uuid.UUID) (*domain.Maintenance, error) {
	var model maintenancemodel.Maintenance
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *MaintenanceRepositoryGorm) Update(maintenanceDomain *domain.Maintenance) error {
	model := maintenancemodel.FromDomain(maintenanceDomain)
	return r.db.Save(model).Error
}

func (r *MaintenanceRepositoryGorm) Delete(id uuid.UUID) error {
	return r.db.Delete(&maintenancemodel.Maintenance{}, "id = ?", id).Error
}

func (r *MaintenanceRepositoryGorm) List() ([]*domain.Maintenance, error) {
	var models []maintenancemodel.Maintenance
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	maintenances := make([]*domain.Maintenance, len(models))
	for i, m := range models {
		maintenances[i] = m.ToDomain()
	}
	return maintenances, nil
}
