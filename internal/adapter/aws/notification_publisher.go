package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/ports"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type EmailNotificationPublisher struct {
	client   *sns.Client
	topicARN string
}

type notifRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type notifMessage struct {
	Recipient notifRecipient    `json:"recipient"`
	Data      map[string]string `json:"data"`
}

func NewEmailNotificationPublisher(ctx context.Context) (*EmailNotificationPublisher, error) {
	cfg, err := GetAWSConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS config: %w", err)
	}

	topicARN := os.Getenv("SNS_TOPIC_NOTIFICATION_EVENTS_ARN")
	if topicARN == "" {
		return nil, fmt.Errorf("SNS_TOPIC_NOTIFICATION_EVENTS_ARN not configured")
	}

	return &EmailNotificationPublisher{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

func (p *EmailNotificationPublisher) publish(ctx context.Context, templateType string, recipient notifRecipient, data map[string]string) error {
	body, err := json.Marshal(notifMessage{Recipient: recipient, Data: data})
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	_, err = p.client.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(p.topicARN),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"templateType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(templateType),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to publish notification event [%s]: %w", templateType, err)
	}
	return nil
}

func (p *EmailNotificationPublisher) SendStockReservedEmail(ctx context.Context, recipientEmail, recipientName string, productID uuid.UUID, productName string, orderID uuid.UUID, quantity int) error {
	return p.publish(ctx, "STOCK_RESERVED", notifRecipient{Email: recipientEmail, Name: recipientName}, map[string]string{
		"product_id":   productID.String(),
		"product_name": productName,
		"order_id":     orderID.String(),
		"quantity":     fmt.Sprintf("%d", quantity),
	})
}

func (p *EmailNotificationPublisher) SendBackorderCreatedEmail(ctx context.Context, recipientEmail, recipientName string, productID uuid.UUID, productName string, orderID uuid.UUID, quantity int) error {
	return p.publish(ctx, "BACKORDER_CREATED", notifRecipient{Email: recipientEmail, Name: recipientName}, map[string]string{
		"product_id":   productID.String(),
		"product_name": productName,
		"order_id":     orderID.String(),
		"quantity":     fmt.Sprintf("%d", quantity),
	})
}

func (p *EmailNotificationPublisher) SendStockReservationFailedEmail(ctx context.Context, recipientEmail, recipientName string, failure domain.StockReservationFailure, productName string) error {
	return p.publish(ctx, "STOCK_RESERVATION_FAILED", notifRecipient{Email: recipientEmail, Name: recipientName}, map[string]string{
		"product_id":         failure.ProductID.String(),
		"product_name":       productName,
		"order_id":           failure.OrderID.String(),
		"quantity_requested": fmt.Sprintf("%d", failure.QuantityRequested),
		"stock_available":    fmt.Sprintf("%d", failure.StockAvailable),
		"failure_reason":     failure.FailureReason,
	})
}

var _ ports.EmailNotificationPort = (*EmailNotificationPublisher)(nil)
