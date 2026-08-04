package order_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	domain "github.com/DmitryToknoff/microservice-demo/internal/order/domain"
	"github.com/DmitryToknoff/microservice-demo/pkg/postgres"
)

type OrderRepository struct {
	pool *postgres.Pool
}

func NewOrderRepository(pool *postgres.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
		INSERT INTO orders (user_id, amount, status) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query, order.UserID, order.Amount, domain.StatusCreated).
		Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return fmt.Errorf("create order in db: %w", err)
	}

	order.Status = domain.StatusCreated
	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
		SELECT id, user_id, amount, status, created_at 
		FROM orders 
		WHERE id = $1
	`

	var order domain.Order
	var statusStr string

	err := r.pool.QueryRow(ctx, query, id).
		Scan(&order.ID, &order.UserID, &order.Amount, &statusStr, &order.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order by id from db: %w", err)
	}

	order.Status = domain.Status(statusStr)
	return &order, nil
}
