package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type CatalogNotificationPort interface {
	NotifyStockReserved(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error
	NotifyBackorderCreated(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error
	NotifyStockReservationFailed(ctx context.Context, failure domain.StockReservationFailure) error
}

type EmailNotificationPort interface {
	SendStockReservedEmail(ctx context.Context, recipientEmail, recipientName string, productID uuid.UUID, productName string, orderID uuid.UUID, quantity int) error
	SendBackorderCreatedEmail(ctx context.Context, recipientEmail, recipientName string, productID uuid.UUID, productName string, orderID uuid.UUID, quantity int) error
	SendStockReservationFailedEmail(ctx context.Context, recipientEmail, recipientName string, failure domain.StockReservationFailure, productName string) error
}
