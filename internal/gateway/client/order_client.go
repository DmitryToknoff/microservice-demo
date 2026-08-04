package gateway_client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/DmitryToknoff/microservice-demo/pkg/pb"
)

type OrderClient struct {
	client pb.OrderServiceClient
	conn   *grpc.ClientConn
}

func NewOrderClient(addr string) (*OrderClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service at %s: %w", addr, err)
	}

	return &OrderClient{
		client: pb.NewOrderServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *OrderClient) CreateOrder(ctx context.Context, userID int64, amount float64) (*pb.OrderResponse, error) {
	req := &pb.CreateOrderRequest{
		UserId: userID,
		Amount: amount,
	}

	resp, err := c.client.CreateOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc CreateOrder error: %w", err)
	}

	return resp, nil
}

func (c *OrderClient) GetOrder(ctx context.Context, id int64) (*pb.OrderResponse, error) {
	req := &pb.GetOrderRequest{
		Id: id,
	}

	resp, err := c.client.GetOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc GetOrder error: %w", err)
	}

	return resp, nil
}

func (c *OrderClient) Close() error {
	return c.conn.Close()
}
