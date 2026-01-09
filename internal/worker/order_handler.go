package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/febry3/gamingin/internal/worker/tasks"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	orderUsecase      usecase.OrderUsecaseContract
	groupBuyUsecase   usecase.GroupBuyUsecaseContract
	userWalletUsecase usecase.UserWalletUsecaseContract
	log               *logrus.Logger
}

func NewOrderHandler(orderUsecase usecase.OrderUsecaseContract, groupBuyUsecase usecase.GroupBuyUsecaseContract, userWallet usecase.UserWalletUsecaseContract, log *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		orderUsecase:      orderUsecase,
		groupBuyUsecase:   groupBuyUsecase,
		userWalletUsecase: userWallet,
		log:               log,
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

	order, err := h.orderUsecase.GetOrderByID(ctx, payload.UserID, payload.OrderID)
	if err != nil {
		h.log.Errorf("Failed to get order %s: %v", payload.OrderNumber, err)
		return err
	}

	if order.Status == "paid" {
		var tier *entity.GroupBuyTier
		if payload.TierID != "" {
			var err error
			if tier, err = h.groupBuyUsecase.GetBuyerGroupSessionTier(ctx, payload.TierID); err != nil {
				h.log.Errorf("Failed to get group buy tier %s: %v", payload.TierID, err)
				return err
			}
		}
		if tier != nil {
			err := h.userWalletUsecase.CreateOrUpdateUserWallet(ctx, payload.UserID, (float64(payload.PaidAmount)*tier.DiscountPercentage)/100)
			if err != nil {
				h.log.Errorf("Failed to create or update user wallet for order %s: %v", payload.OrderNumber, err)
				return err
			}
		}
	}

	h.log.Infof("Successfully expired order: %s", payload.OrderNumber)
	return nil
}
