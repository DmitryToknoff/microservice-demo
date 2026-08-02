package delivery_service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	domain "github.com/DmitryToknoff/microservice-demo/internal/delivery/domain"
	"github.com/DmitryToknoff/microservice-demo/pkg/logger"
)

type Repository interface {
	Create(ctx context.Context, delivery *domain.Delivery) error
	GetByOrderID(ctx context.Context, orderID int64) (*domain.Delivery, error)
}

type DeliveryService struct {
	repo Repository
	log  *logger.Logger
}

func NewDeliveryService(repo Repository, log *logger.Logger) *DeliveryService {
	return &DeliveryService{
		repo: repo,
		log:  log,
	}
}

func (s *DeliveryService) ProcessOrderCreated(ctx context.Context, orderID int64) error {
	delivery := &domain.Delivery{
		OrderID: orderID,
		Address: "defualt address",
	}

	if err := delivery.Validate(); err != nil {
		return fmt.Errorf("validate delivery: %w", err)
	}

	if err := s.repo.Create(ctx, delivery); err != nil {
		s.log.Error("failed to create delivery", zap.Int64("order_id", orderID), zap.Error(err))
		return err
	}

	s.log.Info("delivery created successfully for order", zap.Int64("order_id", orderID))
	return nil
}

func (s *DeliveryService) GetDeliveryStatus(ctx context.Context, orderID int64) (*domain.Delivery, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}
