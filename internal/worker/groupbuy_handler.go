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

type GroupBuySessionHandler struct {
	groupBuyUsecase usecase.GroupBuyUsecaseContract
	log             *logrus.Logger
}

func NewGroupBuySessionHandler(groupBuyUsecase usecase.GroupBuyUsecaseContract, log *logrus.Logger) *GroupBuySessionHandler {
	return &GroupBuySessionHandler{
		groupBuyUsecase: groupBuyUsecase,
		log:             log,
	}
}

func (h *GroupBuySessionHandler) HandleSessionEnd(ctx context.Context, t *asynq.Task) error {
	var payload tasks.GroupBuySessionEndPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	h.log.Infof("Processing group buy session end: SessionID=%s", payload.SessionID)

	err := h.groupBuyUsecase.EndSession(ctx, payload.SessionID, payload.ProductVariantID, payload.SellerID)
	if err != nil {
		h.log.Errorf("failed to end session: %v", err)
		return fmt.Errorf("failed to end session: %w", err)
	}

	return nil
}
