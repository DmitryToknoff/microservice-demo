package orders

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	gatewayClient "github.com/DmitryToknoff/microservice-demo/internal/gateway/client"
	"github.com/DmitryToknoff/microservice-demo/pkg/logger"
)

type CreateHandler struct {
	orderClient *gatewayClient.OrderClient
	log         *logger.Logger
}

func NewCreateHandler(orderClient *gatewayClient.OrderClient, log *logger.Logger) *CreateHandler {
	return &CreateHandler{
		orderClient: orderClient,
		log:         log,
	}
}

type CreateOrderDTO struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	var dto CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"error": "invalid json body"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.orderClient.CreateOrder(ctx, dto.UserID, dto.Amount)
	if err != nil {
		h.log.Error("failed to create order via gRPC", zap.Error(err))
		http.Error(w, `{"error": "failed to create order"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to write response", zap.Error(err))
	}
}
