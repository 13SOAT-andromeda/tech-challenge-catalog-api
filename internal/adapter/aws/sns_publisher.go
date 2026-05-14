package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
)

type SNSPublisher struct {
	client *sns.Client
	topics map[string]string
}

func NewSNSPublisher(ctx context.Context) (*SNSPublisher, error) {
	cfg, err := GetAWSConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS config: %w", err)
	}

	topics := map[string]string{
		"stock_reserved":     os.Getenv("SNS_TOPIC_STOCK_RESERVED_ARN"),
		"backorder_created":  os.Getenv("SNS_TOPIC_BACKORDER_CREATED_ARN"),
		"stock_insufficient": os.Getenv("SNS_TOPIC_STOCK_INSUFFICIENT_ARN"),
	}

	return &SNSPublisher{
		client: sns.NewFromConfig(cfg),
		topics: topics,
	}, nil
}

type stockEvent struct {
	ProductID uuid.UUID `json:"product_id"`
	OrderID   uuid.UUID `json:"order_id"`
	Quantity  int       `json:"quantity"`
	Type      string    `json:"type"`
}

func (s *SNSPublisher) publish(ctx context.Context, topicKey string, event stockEvent) error {
	topicARN, ok := s.topics[topicKey]
	if !ok || topicARN == "" {
		return fmt.Errorf("topic ARN for %s not configured", topicKey)
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = s.client.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(topicARN),
	})

	if err != nil {
		return fmt.Errorf("failed to publish to SNS: %w", err)
	}

	return nil
}

func (s *SNSPublisher) NotifyStockReserved(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error {
	return s.publish(ctx, "stock_reserved", stockEvent{
		ProductID: productID,
		OrderID:   orderID,
		Quantity:  quantity,
		Type:      "STOCK_RESERVED",
	})
}

func (s *SNSPublisher) NotifyBackorderCreated(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error {
	return s.publish(ctx, "backorder_created", stockEvent{
		ProductID: productID,
		OrderID:   orderID,
		Quantity:  quantity,
		Type:      "BACKORDER_CREATED",
	})
}

func (s *SNSPublisher) NotifyStockInsufficient(ctx context.Context, productID uuid.UUID, orderID uuid.UUID, quantity int) error {
	return s.publish(ctx, "stock_insufficient", stockEvent{
		ProductID: productID,
		OrderID:   orderID,
		Quantity:  quantity,
		Type:      "STOCK_INSUFFICIENT",
	})
}

// Ensure SNSPublisher implements ports.CatalogNotificationPort
var _ ports.CatalogNotificationPort = (*SNSPublisher)(nil)
