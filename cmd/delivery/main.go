package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	deliveryRepoPg "github.com/DmitryToknoff/microservice-demo/internal/delivery/repository/postgres"
	deliverySvc "github.com/DmitryToknoff/microservice-demo/internal/delivery/service"
	deliveryGrpc "github.com/DmitryToknoff/microservice-demo/internal/delivery/transport/grpc"
	deliveryKafka "github.com/DmitryToknoff/microservice-demo/internal/delivery/transport/kafka"

	"github.com/DmitryToknoff/microservice-demo/pkg/logger"
	"github.com/DmitryToknoff/microservice-demo/pkg/pb"
	"github.com/DmitryToknoff/microservice-demo/pkg/postgres"
)

var (
	kafkaBrokers = []string{"localhost:9092", "localhost:9093", "localhost:9094"}
	kafkaTopic = "order.created"
	groupID = "delivery-service-group"

	grpcPort = ":50052"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logCfg := logger.MustNewConfig()
	log := logger.NewLoggerMust(logCfg)
	defer log.Close()

	log.Info("Starting Delivery Service...")

	pgCfg := postgres.NewConfigMust()
	pgPool, err := postgres.NewPool(ctx, pgCfg)
	if err != nil {
		log.Fatal("failed to initialize postgres pool", zap.Error(err))
	}

	defer pgPool.Close()
	log.Info("Connected to PostgreSQL for Delivery Service", zap.String("addr", pgCfg.Addr))

	repo := deliveryRepoPg.NewDeliveryRepository(pgPool)
	service := deliverySvc.NewDeliveryService(repo, log)
	grpcHandler := deliveryGrpc.NewServer(service)

	consumer := deliveryKafka.NewConsumer(kafkaBrokers, kafkaTopic, groupID, service, log)
	defer consumer.Close()

	go consumer.Start(ctx)

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal("failed to listen tcp port", zap.String("port", grpcPort), zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(grpcServer, grpcHandler)

	go func() {
		log.Info("Delivery Service gRPC running", zap.String("port", grpcPort))
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("gRPC server stopped with error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down Delivery Service gracefully...")

	grpcServer.GracefulStop()
	log.Info("Delivery Service stopped successfully")
}
