package deliveries

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

type GetStatusHandler struct {
	deliveryClient *gatewayClient.DeliveryClient
	log            *logger.Logger
}

func NewGetStatusHandler(deliveryClient *gatewayClient.DeliveryClient, log *logger.Logger) *GetStatusHandler {
	return &GetStatusHandler{
		deliveryClient: deliveryClient,
		log:            log,
	}
}

func (h *GetStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orderIDStr := r.PathValue("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil || orderID <= 0 {
		http.Error(w, `{"error": "invalid order id"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.deliveryClient.GetDeliveryStatus(ctx, orderID)
	if err != nil {
		h.log.Warn("failed to fetch delivery via gRPC", zap.Int64("order_id", orderID), zap.Error(err))
		http.Error(w, `{"error": "delivery status not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to write response", zap.Int64("order_id", orderID), zap.Error(err))
	}
}
