package order_tranport_kafka

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/DmitryToknoff/microservice-demo/internal/order/domain"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

type OrderCreatedEvent struct {
	OrderID   int64   `json:"order_id"`
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
}

func (p *Producer) PublishOrderCreated(ctx context.Context, order *domain.Order) error {
	event := OrderCreatedEvent{
		OrderID:   order.ID,
		UserID:    order.UserID,
		Amount:    order.Amount,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal order created event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", order.ID)),
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
