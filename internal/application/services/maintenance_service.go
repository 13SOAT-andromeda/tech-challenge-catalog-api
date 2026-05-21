package services

import (
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type MaintenanceService struct {
	repo ports.MaintenanceRepository
}

func NewMaintenanceService(repo ports.MaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{repo: repo}
}

func (s *MaintenanceService) Create(maintenance *domain.Maintenance) error {
	return s.repo.Create(maintenance)
}

func (s *MaintenanceService) GetByID(id uuid.UUID) (*domain.Maintenance, error) {
	return s.repo.FindByID(id)
}

func (s *MaintenanceService) List() ([]*domain.Maintenance, error) {
	return s.repo.List()
}

func (s *MaintenanceService) Update(maintenance *domain.Maintenance) error {
	return s.repo.Update(maintenance)
}

func (s *MaintenanceService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *MaintenanceService) GetBatch(ids []int64) ([]*domain.Maintenance, error) {
	return s.repo.FindBatchBySerialIDs(ids)
}
