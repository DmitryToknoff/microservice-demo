package order_domain

import (
	"errors"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidAmount = errors.New("order amount must be greater than zero")
	ErrInvalidUserID = errors.New("user_id is required")
	ErrInvalidStatus = errors.New("invalid status transition")
)

type Status string

const (
	StatusCreated   Status = "CREATED"
	StatusPaid      Status = "PAID"
	StatusCancelled Status = "CANCELLED"
)

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Amount    float64   `json:"amount"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (o *Order) Validate() error {
	if o.UserID <= 0 {
		return ErrInvalidUserID
	}
	if o.Amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}
