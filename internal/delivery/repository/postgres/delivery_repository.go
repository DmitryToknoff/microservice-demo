package delivery_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	domain "github.com/DmitryToknoff/microservice-demo/internal/delivery/domain"
	"github.com/DmitryToknoff/microservice-demo/pkg/postgres"
)

type DeliveryRepository struct {
	pool *postgres.Pool
}

func NewDeliveryRepository(pool *postgres.Pool) *DeliveryRepository {
	return &DeliveryRepository{pool: pool}
}

func (r *DeliveryRepository) Create(ctx context.Context, delivery *domain.Delivery) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
		INSERT INTO deliveries (order_id, address, status) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (order_id) DO NOTHING
		RETURNING id, updated_at
	`

	err := r.pool.QueryRow(ctx, query, delivery.OrderID, delivery.Address, domain.StatusProcessing).
		Scan(&delivery.ID, &delivery.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("create delivery in db: %w", err)
	}

	delivery.Status = domain.StatusProcessing
	return nil
}

func (r *DeliveryRepository) GetByOrderID(ctx context.Context, orderID int64) (*domain.Delivery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
		SELECT id, order_id, address, status, updated_at 
		FROM deliveries 
		WHERE order_id = $1
	`

	var delivery domain.Delivery
	var statusStr string

	err := r.pool.QueryRow(ctx, query, orderID).
		Scan(&delivery.ID, &delivery.OrderID, &delivery.Address, &statusStr, &delivery.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDeliveryNotFound
		}
		return nil, fmt.Errorf("get delivery by order_id from db: %w", err)
	}

	delivery.Status = domain.Status(statusStr)
	return &delivery, nil
}
