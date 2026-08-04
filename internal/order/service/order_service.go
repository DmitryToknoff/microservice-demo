package order_service

import (
	"context"
	"fmt"

	redisClient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	domain "github.com/DmitryToknoff/microservice-demo/internal/order/domain"
	"github.com/DmitryToknoff/microservice-demo/pkg/logger"
)

type Repository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id int64) (*domain.Order, error)
}

type Cache interface {
	Get(ctx context.Context, id int64) (*domain.Order, error)
	Set(ctx context.Context, order *domain.Order) error
}

type EventProducer interface {
	PublishOrderCreated(ctx context.Context, order *domain.Order) error
}

type OrderService struct {
	repo     Repository
	cache    Cache
	producer EventProducer
	log      *logger.Logger
}

func NewOrderService(repo Repository, cache Cache, producer EventProducer, log *logger.Logger) *OrderService {
	return &OrderService{
		repo:     repo,
		cache:    cache,
		producer: producer,
		log:      log,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID int64, amount float64) (*domain.Order, error) {
	order := &domain.Order{
		UserID: userID,
		Amount: amount,
	}

	if err := order.Validate(); err != nil {
		return nil, fmt.Errorf("validate order: %w", err)
	}

	if err := s.repo.Create(ctx, order); err != nil {
		s.log.Error("failed to create order in db", zap.Error(err))
		return nil, err
	}

	if err := s.producer.PublishOrderCreated(ctx, order); err != nil {
		s.log.Error("failed to publish order created event", zap.Int64("order_id", order.ID), zap.Error(err))
	}

	if err := s.cache.Set(ctx, order); err != nil {
		s.log.Warn("failed to cache new order", zap.Int64("order_id", order.ID), zap.Error(err))
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id int64) (*domain.Order, error) {
	cachedOrder, err := s.cache.Get(ctx, id)
	if err == nil {
		s.log.Debug("cache hit", zap.Int64("order_id", id))
		return cachedOrder, nil
	}

	if err != nil && err != redisClient.Nil {
		s.log.Warn("redis error on get", zap.Error(err))
	}

	s.log.Debug("cache miss", zap.Int64("order_id", id))
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, order); err != nil {
		s.log.Warn("failed to write order to cache", zap.Int64("order_id", id), zap.Error(err))
	}

	return order, nil
}
