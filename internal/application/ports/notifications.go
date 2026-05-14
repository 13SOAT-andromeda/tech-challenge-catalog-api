package ports

import (
	"context"

	"github.com/google/uuid"
)

type CatalogNotificationPort interface {
	NotifyStockReserved(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error
	NotifyBackorderCreated(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error
	NotifyStockInsufficient(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error
}
