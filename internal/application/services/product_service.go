package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type ProductService struct {
	repo             ports.ProductRepository
	notificationPort ports.CatalogNotificationPort
}

func NewProductService(repo ports.ProductRepository, notificationPort ports.CatalogNotificationPort) *ProductService {
	return &ProductService{
		repo:             repo,
		notificationPort: notificationPort,
	}
}

func (s *ProductService) Create(product *domain.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetByID(id uuid.UUID) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) List() ([]*domain.Product, error) {
	return s.repo.List()
}

func (s *ProductService) Update(product *domain.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *ProductService) DecreaseStock(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return err
	}

	if product.StockQuantity < quantity {
		if product.ReplenishmentStrategy == "BACKORDER" {
			// In a real scenario, we would create a backorder record here
			s.notificationPort.NotifyBackorderCreated(ctx, productID, orderID, quantity)
			return nil
		}
		s.notificationPort.NotifyStockInsufficient(ctx, productID, orderID, quantity)
		return domain.ErrInsufficientStock
	}

	err = s.repo.UpdateStock(productID, -quantity)
	if err != nil {
		return err
	}

	return s.notificationPort.NotifyStockReserved(ctx, productID, orderID, quantity)
}
