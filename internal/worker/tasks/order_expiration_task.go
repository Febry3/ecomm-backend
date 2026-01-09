package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const TypeOrderExpiration = "order:expire"

type OrderExpirationPayload struct {
	OrderID     string `json:"order_id"`
	OrderNumber string `json:"order_number"`
	TierID      string `json:"tier_id"`
	UserID      int64  `json:"user_id"`
	PaidAmount  int64  `json:"paid_amount"`
}

// NewOrderExpirationTask creates a new order expiration task
func NewOrderExpirationTask(orderID, orderNumber, tierID string, userID int64, paidAmount int64) (*asynq.Task, error) {
	payload, err := json.Marshal(OrderExpirationPayload{
		OrderID:     orderID,
		OrderNumber: orderNumber,
		TierID:      tierID,
		UserID:      userID,
		PaidAmount:  paidAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order expiration payload: %w", err)
	}
	return asynq.NewTask(TypeOrderExpiration, payload), nil
}

// OrderExpirationHandler defines the interface for handling order expiration
type OrderExpirationHandler interface {
	ExpireOrder(ctx context.Context, orderID string) error
}

// HandleOrderExpirationTask processes the order expiration task
func HandleOrderExpirationTask(handler OrderExpirationHandler) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload OrderExpirationPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal order expiration payload: %w", err)
		}

		return handler.ExpireOrder(ctx, payload.OrderID)
	}
}
