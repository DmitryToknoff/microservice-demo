package order_repository_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domain "github.com/DmitryToknoff/microservice-demo/internal/order/domain"
	"github.com/DmitryToknoff/microservice-demo/pkg/redis"
)

type OrderCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewOrderCache(client *redis.Client, ttl time.Duration) *OrderCache {
	return &OrderCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *OrderCache) Get(ctx context.Context, id int64) (*domain.Order, error) {
	key := fmt.Sprintf("order:%d", id)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return nil, fmt.Errorf("unmarshal cached order: %w", err)
	}

	return &order, nil
}

func (c *OrderCache) Set(ctx context.Context, order *domain.Order) error {
	key := fmt.Sprintf("order:%d", order.ID)

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshal order for cache: %w", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}
