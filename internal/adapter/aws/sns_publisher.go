package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type SNSPublisher struct {
	client   *sns.Client
	topicARN string
}

func NewSNSPublisher(ctx context.Context) (*SNSPublisher, error) {
	cfg, err := GetAWSConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS config: %w", err)
	}

	topicARN := os.Getenv("SNS_TOPIC_CATALOG_EVENTS_ARN")
	if topicARN == "" {
		return nil, fmt.Errorf("SNS_TOPIC_CATALOG_EVENTS_ARN not configured")
	}

	return &SNSPublisher{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

func (s *SNSPublisher) publish(ctx context.Context, eventType string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = s.client.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(s.topicARN),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(eventType),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to publish to SNS: %w", err)
	}
	return nil
}

type stockReservedPayload struct {
	Type      string    `json:"type"`
	ProductID uuid.UUID `json:"product_id"`
	OrderID   uuid.UUID `json:"order_id"`
	Quantity  int       `json:"quantity"`
}

type backorderCreatedPayload struct {
	Type      string    `json:"type"`
	ProductID uuid.UUID `json:"product_id"`
	OrderID   uuid.UUID `json:"order_id"`
	Quantity  int       `json:"quantity"`
}

type stockReservationFailedPayload struct {
	Type              string    `json:"type"`
	OrderID           uuid.UUID `json:"order_id"`
	ProductID         uuid.UUID `json:"product_id"`
	QuantityRequested int       `json:"quantity_requested"`
	StockAvailable    int       `json:"stock_available"`
	FailureReason     string    `json:"failure_reason"`
	ErrorDetail       string    `json:"error_detail"`
	OccurredAt        time.Time `json:"occurred_at"`
}

func (s *SNSPublisher) NotifyStockReserved(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error {
	return s.publish(ctx, "STOCK_RESERVED", stockReservedPayload{
		Type:      "STOCK_RESERVED",
		ProductID: productID,
		OrderID:   orderID,
		Quantity:  quantity,
	})
}

func (s *SNSPublisher) NotifyBackorderCreated(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error {
	return s.publish(ctx, "BACKORDER_CREATED", backorderCreatedPayload{
		Type:      "BACKORDER_CREATED",
		ProductID: productID,
		OrderID:   orderID,
		Quantity:  quantity,
	})
}

func (s *SNSPublisher) NotifyStockReservationFailed(ctx context.Context, failure domain.StockReservationFailure) error {
	return s.publish(ctx, "STOCK_RESERVATION_FAILED", stockReservationFailedPayload{
		Type:              "STOCK_RESERVATION_FAILED",
		OrderID:           failure.OrderID,
		ProductID:         failure.ProductID,
		QuantityRequested: failure.QuantityRequested,
		StockAvailable:    failure.StockAvailable,
		FailureReason:     failure.FailureReason,
		ErrorDetail:       failure.ErrorDetail,
		OccurredAt:        failure.OccurredAt,
	})
}

var _ ports.CatalogNotificationPort = (*SNSPublisher)(nil)
