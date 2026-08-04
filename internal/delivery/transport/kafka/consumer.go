package delivery_transport_kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	delivery_service "github.com/DmitryToknoff/microservice-demo/internal/delivery/service"
	"github.com/DmitryToknoff/microservice-demo/pkg/logger"
)

type Consumer struct {
	reader  *kafka.Reader
	service *delivery_service.DeliveryService
	log     *logger.Logger
}

func NewConsumer(brokers []string, topic, groupID string, svc *delivery_service.DeliveryService, log *logger.Logger) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       10 * 1024,
			MaxBytes:       10 * 1024 * 1024,
			CommitInterval: 0,
		}),
		service: svc,
		log:     log,
	}
}

type OrderCreatedEvent struct {
	OrderID int64 `json:"order_id"`
}

func (c *Consumer) Start(ctx context.Context) {
	c.log.Info("Kafka Consumer started listening for order.created...")

	for {
		select {
		case <-ctx.Done():
			c.log.Info("Kafka Consumer stopping...")
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				c.log.Error("error fetching message from kafka", zap.Error(err))
				continue
			}

			var event OrderCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				c.log.Error("failed to unmarshal kafka event, skipping message", zap.Error(err))
				_ = c.reader.CommitMessages(ctx, msg)
				continue
			}

			c.log.Debug("received order.created event", zap.Int64("order_id", event.OrderID))

			if err := c.service.ProcessOrderCreated(ctx, event.OrderID); err != nil {
				c.log.Error("failed to process delivery creation",
					zap.Int64("order_id", event.OrderID),
					zap.Error(err),
				)
				time.Sleep(1 * time.Second)
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.log.Error("failed to commit message offset to kafka", zap.Error(err))
			} else {
				c.log.Info("successfully processed and committed message", zap.Int64("order_id", event.OrderID))
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
