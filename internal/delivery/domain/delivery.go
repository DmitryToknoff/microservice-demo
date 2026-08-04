package delivery_domain

import (
	"errors"
	"time"
)

var (
	ErrDeliveryNotFound = errors.New("delivery not found")
	ErrInvalidOrderID   = errors.New("order_id is required")
)

type Status string

const (
	StatusProcessing Status = "PROCESSING"
	StatusInTransit  Status = "IN_TRANSIT"
	StatusDelivered  Status = "DELIVERED"
	StatusFailed     Status = "FAILED"
)

type Delivery struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	Address   string    `json:"address"`
	Status    Status    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Delivery) Validate() error {
	if d.OrderID <= 0 {
		return ErrInvalidOrderID
	}
	return nil
}
