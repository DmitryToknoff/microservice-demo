package gateway_client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/DmitryToknoff/microservice-demo/pkg/pb"
)

type DeliveryClient struct {
	client pb.DeliveryServiceClient
	conn   *grpc.ClientConn
}

func NewDeliveryClient(addr string) (*DeliveryClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to delivery service at %s: %w", addr, err)
	}

	return &DeliveryClient{
		client: pb.NewDeliveryServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *DeliveryClient) GetDeliveryStatus(ctx context.Context, orderID int64) (*pb.DeliveryStatusResponse, error) {
	req := &pb.GetDeliveryStatusRequest{
		OrderId: orderID,
	}

	resp, err := c.client.GetDeliveryStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc GetDeliveryStatus error: %w", err)
	}

	return resp, nil
}

func (c *DeliveryClient) Close() error {
	return c.conn.Close()
}
