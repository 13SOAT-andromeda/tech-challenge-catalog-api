package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrProductNotFound     = errors.New("product not found")
	ErrMaintenanceNotFound = errors.New("maintenance not found")
)

const (
	FailureReasonInsufficientStock = "INSUFFICIENT_STOCK"
	FailureReasonProductNotFound   = "PRODUCT_NOT_FOUND"
	FailureReasonDBError           = "DB_ERROR"
)

type StockReservationFailure struct {
	OrderID           uuid.UUID
	ProductID         uuid.UUID
	QuantityRequested int
	StockAvailable    int
	FailureReason     string
	ErrorDetail       string
	OccurredAt        time.Time
}
