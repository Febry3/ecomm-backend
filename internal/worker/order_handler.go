package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/febry3/gamingin/internal/usecase"
	"github.com/febry3/gamingin/internal/worker/tasks"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecaseContract
	log          *logrus.Logger
}

func NewOrderHandler(orderUsecase usecase.OrderUsecaseContract, log *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
		log:          log,
	}
}

func (h *OrderHandler) HandleOrderExpiration(ctx context.Context, task *asynq.Task) error {
	var payload tasks.OrderExpirationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal order expiration payload: %w", err)
	}

	h.log.Infof("Processing order expiration for order: %s", payload.OrderNumber)

	if err := h.orderUsecase.ExpireOrder(ctx, payload.OrderID); err != nil {
		h.log.Errorf("Failed to expire order %s: %v", payload.OrderNumber, err)
		return err
	}

	h.log.Infof("Successfully expired order: %s", payload.OrderNumber)
	return nil
}
