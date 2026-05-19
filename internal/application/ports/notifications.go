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
