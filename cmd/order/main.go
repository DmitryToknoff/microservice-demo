package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	orderRepoPg "github.com/DmitryToknoff/microservice-demo/internal/order/repository/postgres"
	orderRepoRedis "github.com/DmitryToknoff/microservice-demo/internal/order/repository/redis"
	orderSvc "github.com/DmitryToknoff/microservice-demo/internal/order/service"
	orderGrpc "github.com/DmitryToknoff/microservice-demo/internal/order/transport/grpc"
	orderKafka "github.com/DmitryToknoff/microservice-demo/internal/order/transport/kafka"

	"github.com/DmitryToknoff/microservice-demo/pkg/logger"
	"github.com/DmitryToknoff/microservice-demo/pkg/pb"
	"github.com/DmitryToknoff/microservice-demo/pkg/postgres"
	"github.com/DmitryToknoff/microservice-demo/pkg/redis"
)

var (
	kafkaBrokers = []string{"kafka-1:9092", "kafka-2:9093", "kafka-3:9094"}
	grpcPort     = ":50051"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logCfg := logger.MustNewConfig()
	log := logger.NewLoggerMust(logCfg)
	defer log.Close()

	log.Info("Starting Order Service...")

	pgCfg := postgres.NewConfigMust()
	pgPool, err := postgres.NewPool(ctx, pgCfg)
	if err != nil {
		log.Fatal("failed to initialize postgres pool", zap.Error(err))
	}
	defer pgPool.Close()
	log.Info("Connected to PostgreSQL", zap.String("addr", pgCfg.Addr))

	redisCfg := redis.MustNewConfig()
	redisClient, err := redis.NewClient(ctx, redisCfg)
	if err != nil {
		log.Fatal("failed to initialize redis client", zap.Error(err))
	}
	defer redisClient.Close()
	log.Info("Connected to Redis", zap.String("addr", redisCfg.Addr))

	kafkaTopic := "order.created"
	producer := orderKafka.NewProducer(kafkaBrokers, kafkaTopic)
	defer producer.Close()
	log.Info("Initialized Kafka Producer", zap.Strings("brokers", kafkaBrokers), zap.String("topic", kafkaTopic))

	pgRepository := orderRepoPg.NewOrderRepository(pgPool)
	redisCache := orderRepoRedis.NewOrderCache(redisClient, 5*time.Minute)

	service := orderSvc.NewOrderService(pgRepository, redisCache, producer, log)
	grpcHandler := orderGrpc.NewServer(service)

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal("failed to listen tcp port", zap.String("port", grpcPort), zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, grpcHandler)

	go func() {
		log.Info("gRPC Server running", zap.String("port", grpcPort))
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("gRPC server stopped with error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down Order Service gracefully...")

	grpcServer.GracefulStop()
	log.Info("Order Service stopped successfully")
}
