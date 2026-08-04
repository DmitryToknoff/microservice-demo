package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	gatewayClient "github.com/DmitryToknoff/microservice-demo/internal/gateway/client"
	deliveriesHandler "github.com/DmitryToknoff/microservice-demo/internal/gateway/handler/delivery"
	ordersHandler "github.com/DmitryToknoff/microservice-demo/internal/gateway/handler/orders"
	"github.com/DmitryToknoff/microservice-demo/pkg/logger"
)

var (
	httpPort = ":8080"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logCfg := logger.MustNewConfig()
	log := logger.NewLoggerMust(logCfg)
	defer log.Close()

	log.Info("Starting API Gateway...")

	orderSvcAddr := os.Getenv("ORDER_SERVICE_GRPC_HOST")
	if orderSvcAddr == "" {
		orderSvcAddr = "localhost:50051"
	}

	orderClient, err := gatewayClient.NewOrderClient(orderSvcAddr)
	if err != nil {
		log.Fatal("failed to connect to order service", zap.Error(err))
	}
	defer orderClient.Close()

	deliverySvcAddr := os.Getenv("DELIVERY_SERVICE_GRPC_HOST")
	if deliverySvcAddr == "" {
		deliverySvcAddr = "localhost:50052"
	}

	deliveryClient, err := gatewayClient.NewDeliveryClient(deliverySvcAddr)
	if err != nil {
		log.Fatal("failed to connect to delivery service", zap.Error(err))
	}
	defer deliveryClient.Close()

	createOrderH := ordersHandler.NewCreateHandler(orderClient, log)
	getOrderH := ordersHandler.NewGetHandler(orderClient, log)
	getDeliveryH := deliveriesHandler.NewGetStatusHandler(deliveryClient, log)

	mux := http.NewServeMux()
	mux.Handle("POST /api/v1/orders", createOrderH)
	mux.Handle("GET /api/v1/orders/{id}", getOrderH)
	mux.Handle("GET /api/v1/deliveries/{order_id}", getDeliveryH)

	server := &http.Server{
		Addr:         httpPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("API Gateway REST server running", zap.String("port", httpPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server stopped with error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down API Gateway gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("GateWay failed shutdown", zap.Error(err))
		_ = server.Close()
	}

	log.Info("API Gateway stopped successfully")
}
