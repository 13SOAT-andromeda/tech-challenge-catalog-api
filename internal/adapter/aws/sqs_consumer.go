package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/services"
)

type SQSConsumer struct {
	client     *sqs.Client
	queueURL   string
	productSvc *services.ProductService
}

// ordersEventEnvelope is the outer wrapper published by the orders service for all events.
type ordersEventEnvelope struct {
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

type orderApprovedEvent struct {
	OrderID    string      `json:"order_id"`
	Items      []orderItem `json:"items"`
	ApprovedAt time.Time   `json:"approved_at"`
}

type orderItem struct {
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	Name           string `json:"name"`
	Quantity       int    `json:"qty"`
	UnitPriceCents int64  `json:"unit_price_cents"`
}

// snsEnvelope unwraps messages forwarded from SNS → SQS subscriptions.
type snsEnvelope struct {
	Message string `json:"Message"`
}

func NewSQSConsumer(ctx context.Context, productSvc *services.ProductService) (*SQSConsumer, error) {
	cfg, err := GetAWSConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS config: %w", err)
	}

	queueURL := os.Getenv("SQS_ORDERS_APPROVED_QUEUE_URL")
	if queueURL == "" {
		return nil, fmt.Errorf("SQS_ORDERS_APPROVED_QUEUE_URL not configured")
	}

	return &SQSConsumer{
		client:     sqs.NewFromConfig(cfg),
		queueURL:   queueURL,
		productSvc: productSvc,
	}, nil
}

func (c *SQSConsumer) Start(ctx context.Context) {
	log.Printf("SQS consumer started for queue: %s", c.queueURL)
	for {
		select {
		case <-ctx.Done():
			log.Println("SQS consumer stopped")
			return
		default:
			c.poll(ctx)
		}
	}
}

func (c *SQSConsumer) poll(ctx context.Context) {
	output, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		log.Printf("error receiving SQS messages: %v", err)
		time.Sleep(5 * time.Second)
		return
	}

	for _, msg := range output.Messages {
		if err := c.process(ctx, msg); err != nil {
			log.Printf("error processing message %s: %v", aws.ToString(msg.MessageId), err)
			continue
		}
		c.deleteMessage(ctx, msg)
	}
}

func (c *SQSConsumer) process(ctx context.Context, msg types.Message) error {
	body := aws.ToString(msg.Body)

	// Unwrap SNS → SQS notification envelope
	var snsEnv snsEnvelope
	if err := json.Unmarshal([]byte(body), &snsEnv); err == nil && snsEnv.Message != "" {
		body = snsEnv.Message
	}

	// Unwrap orders service event envelope
	var outerEnv ordersEventEnvelope
	if err := json.Unmarshal([]byte(body), &outerEnv); err != nil {
		return fmt.Errorf("unmarshal orders event envelope: %w", err)
	}
	if outerEnv.EventType != "order.approved" {
		log.Printf("unexpected event type %q — skipping", outerEnv.EventType)
		return nil
	}

	var event orderApprovedEvent
	if err := json.Unmarshal(outerEnv.Data, &event); err != nil {
		return fmt.Errorf("unmarshal OrderApproved data: %w", err)
	}

	orderID, err := uuid.Parse(event.OrderID)
	if err != nil {
		return fmt.Errorf("invalid order_id %q: %w", event.OrderID, err)
	}

	for _, item := range event.Items {
		if item.Kind != "product" {
			continue
		}
		productID, err := uuid.Parse(item.ID)
		if err != nil {
			log.Printf("invalid product id %q in order %s — skipping", item.ID, event.OrderID)
			continue
		}
		if err := c.productSvc.DecreaseStock(ctx, productID, orderID, item.Quantity); err != nil {
			log.Printf("DecreaseStock failed product=%s order=%s: %v", item.ID, event.OrderID, err)
		}
	}
	return nil
}

func (c *SQSConsumer) deleteMessage(ctx context.Context, msg types.Message) {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Printf("error deleting message %s: %v", aws.ToString(msg.MessageId), err)
	}
}
