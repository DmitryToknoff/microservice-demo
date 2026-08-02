package order_transport_grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domain "github.com/DmitryToknoff/microservice-demo/internal/order/domain"
	"github.com/DmitryToknoff/microservice-demo/internal/order/service"
	pb "github.com/DmitryToknoff/microservice-demo/pkg/pb"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	service *order_service.OrderService
}

func NewServer(svc *order_service.OrderService) *Server {
	return &Server{service: svc}
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	order, err := s.service.CreateOrder(ctx, req.UserId, req.Amount)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidAmount) || errors.Is(err, domain.ErrInvalidUserID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to create order")
	}

	return &pb.OrderResponse{
		Id:     order.ID,
		UserId: order.UserID,
		Amount: order.Amount,
		Status: string(order.Status),
	}, nil
}

func (s *Server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	order, err := s.service.GetOrder(ctx, req.Id)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, "failed to get order")
	}

	return &pb.OrderResponse{
		Id:     order.ID,
		UserId: order.UserID,
		Amount: order.Amount,
		Status: string(order.Status),
	}, nil
}
