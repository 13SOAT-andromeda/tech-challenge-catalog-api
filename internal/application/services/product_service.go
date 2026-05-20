package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type ProductService struct {
	repo                 ports.ProductRepository
	backorderRepo        ports.BackorderRepository
	notificationPort     ports.CatalogNotificationPort
	emailPort            ports.EmailNotificationPort
	defaultRecipientEmail string
	defaultRecipientName  string
}

func NewProductService(
	repo ports.ProductRepository,
	backorderRepo ports.BackorderRepository,
	notificationPort ports.CatalogNotificationPort,
	emailPort ports.EmailNotificationPort,
	recipientEmail, recipientName string,
) *ProductService {
	return &ProductService{
		repo:                 repo,
		backorderRepo:        backorderRepo,
		notificationPort:     notificationPort,
		emailPort:            emailPort,
		defaultRecipientEmail: recipientEmail,
		defaultRecipientName:  recipientName,
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
		reason := domain.FailureReasonDBError
		if errors.Is(err, domain.ErrProductNotFound) {
			reason = domain.FailureReasonProductNotFound
		}
		s.notifyFailed(ctx, domain.StockReservationFailure{
			OrderID:           orderID,
			ProductID:         productID,
			QuantityRequested: quantity,
			StockAvailable:    0,
			FailureReason:     reason,
			ErrorDetail:       err.Error(),
		})
		return err
	}

	if product.StockQuantity >= quantity {
		if err := s.repo.UpdateStock(productID, -quantity); err != nil {
			s.notifyFailed(ctx, domain.StockReservationFailure{
				OrderID:           orderID,
				ProductID:         productID,
				QuantityRequested: quantity,
				StockAvailable:    product.StockQuantity,
				FailureReason:     domain.FailureReasonDBError,
				ErrorDetail:       err.Error(),
			})
			return err
		}
		if err := s.notificationPort.NotifyStockReserved(ctx, productID, orderID, quantity); err != nil {
			return err
		}
		if s.emailPort != nil {
			if err := s.emailPort.SendStockReservedEmail(ctx, s.defaultRecipientEmail, s.defaultRecipientName, productID, product.Name, orderID, quantity); err != nil {
				log.Printf("failed to send stock reserved email: %v", err)
			}
		}
		return nil
	}

	switch product.ReplenishmentStrategy {
	case "IMMEDIATE", "BATCH":
		return s.createBackorder(ctx, product, orderID, quantity)
	default:
		s.notifyFailed(ctx, domain.StockReservationFailure{
			OrderID:           orderID,
			ProductID:         productID,
			QuantityRequested: quantity,
			StockAvailable:    product.StockQuantity,
			FailureReason:     domain.FailureReasonInsufficientStock,
			ErrorDetail:       fmt.Sprintf("requested %d units, only %d available", quantity, product.StockQuantity),
		})
		return domain.ErrInsufficientStock
	}
}

func (s *ProductService) createBackorder(ctx context.Context, product *domain.Product, orderID uuid.UUID, quantity int) error {
	backorder := &domain.Backorder{
		ProductID: product.ID,
		OrderID:   orderID,
		Quantity:  quantity,
		Status:    domain.BackorderStatusPending,
	}

	if err := s.backorderRepo.Create(backorder); err != nil {
		s.notifyFailed(ctx, domain.StockReservationFailure{
			OrderID:           orderID,
			ProductID:         product.ID,
			QuantityRequested: quantity,
			StockAvailable:    product.StockQuantity,
			FailureReason:     domain.FailureReasonDBError,
			ErrorDetail:       err.Error(),
		})
		return err
	}

	if err := s.notificationPort.NotifyBackorderCreated(ctx, product.ID, orderID, quantity); err != nil {
		log.Printf("failed to notify backorder created for product %s: %v", product.ID, err)
	}
	if s.emailPort != nil {
		if err := s.emailPort.SendBackorderCreatedEmail(ctx, s.defaultRecipientEmail, s.defaultRecipientName, product.ID, product.Name, orderID, quantity); err != nil {
			log.Printf("failed to send backorder created email: %v", err)
		}
	}
	return nil
}

func (s *ProductService) notifyFailed(ctx context.Context, f domain.StockReservationFailure) {
	if f.OccurredAt.IsZero() {
		f.OccurredAt = time.Now().UTC()
	}
	if err := s.notificationPort.NotifyStockReservationFailed(ctx, f); err != nil {
		log.Printf("failed to publish StockReservationFailed [reason=%s product=%s]: %v", f.FailureReason, f.ProductID, err)
	}
	if s.emailPort != nil {
		product, _ := s.repo.FindByID(f.ProductID)
		productName := f.ProductID.String()
		if product != nil {
			productName = product.Name
		}
		if err := s.emailPort.SendStockReservationFailedEmail(ctx, s.defaultRecipientEmail, s.defaultRecipientName, f, productName); err != nil {
			log.Printf("failed to send stock reservation failed email: %v", err)
		}
	}
}
