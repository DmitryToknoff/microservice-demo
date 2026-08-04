package orders

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	gatewayClient "github.com/DmitryToknoff/microservice-demo/internal/gateway/client"
	"github.com/DmitryToknoff/microservice-demo/pkg/logger"
)

type GetHandler struct {
	orderClient *gatewayClient.OrderClient
	log         *logger.Logger
}

func NewGetHandler(orderClient *gatewayClient.OrderClient, log *logger.Logger) *GetHandler {
	return &GetHandler{
		orderClient: orderClient,
		log:         log,
	}
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, `{"error": "invalid order id"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.orderClient.GetOrder(ctx, id)
	if err != nil {
		h.log.Warn("failed to fetch order via gRPC", zap.Int64("id", id), zap.Error(err))
		http.Error(w, `{"error": "order not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Warn("failed to write response", zap.Int64("id", id), zap.Error(err))
	}
}
