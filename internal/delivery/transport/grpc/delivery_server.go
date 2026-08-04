package delivery_transport_grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domain "github.com/DmitryToknoff/microservice-demo/internal/delivery/domain"
	"github.com/DmitryToknoff/microservice-demo/internal/delivery/service"
	pb "github.com/DmitryToknoff/microservice-demo/pkg/pb"
)

type Server struct {
	pb.UnimplementedDeliveryServiceServer
	service *delivery_service.DeliveryService
}

func NewServer(svc *delivery_service.DeliveryService) *Server {
	return &Server{service: svc}
}

func (s *Server) GetDeliveryStatus(ctx context.Context, req *pb.GetDeliveryStatusRequest) (*pb.DeliveryStatusResponse, error) {
	delivery, err := s.service.GetDeliveryStatus(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, domain.ErrDeliveryNotFound) {
			return nil, status.Error(codes.NotFound, "delivery not found")
		}
		return nil, status.Error(codes.Internal, "failed to get delivery status")
	}

	return &pb.DeliveryStatusResponse{
		OrderId: delivery.OrderID,
		Status:  string(delivery.Status),
		Address: delivery.Address,
	}, nil
}
