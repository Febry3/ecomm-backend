package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/febry3/gamingin/internal/usecase"
	"github.com/febry3/gamingin/internal/worker/tasks"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

type GroupBuySessionHandler struct {
	groupBuyUsecase usecase.GroupBuyUsecaseContract
	asynqClient     *asynq.Client
	email           *mail.Dialer
	log             *logrus.Logger
}

func NewGroupBuySessionHandler(groupBuyUsecase usecase.GroupBuyUsecaseContract, asynqClient *asynq.Client, email *mail.Dialer, log *logrus.Logger) *GroupBuySessionHandler {
	return &GroupBuySessionHandler{
		groupBuyUsecase: groupBuyUsecase,
		asynqClient:     asynqClient,
		email:           email,
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

	// Chain: Enqueue email notification after session ends successfully
	emailTask, err := tasks.NewGroupBuySessionEndMailTask(tasks.GroupBuySessionEndMailPayload{
		To:      "kadallaut96@gmail.com", // TODO: Get from participants
		Subject: "Group Buy Session Ended",
		Body:    fmt.Sprintf("Your group buy session %s has ended!", payload.SessionID),
	})
	if err != nil {
		h.log.Errorf("failed to create email task: %v", err)
	} else {
		_, err = h.asynqClient.Enqueue(emailTask, asynq.Queue("low"))
		if err != nil {
			h.log.Errorf("failed to enqueue email task: %v", err)
		} else {
			h.log.Infof("Enqueued email notification for session %s", payload.SessionID)
		}
	}

	return nil
}

func (h *GroupBuySessionHandler) HandleSessionEndMail(ctx context.Context, t *asynq.Task) error {
	var payload tasks.GroupBuySessionEndMailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", h.email.Username)
	m.SetHeader("To", payload.To)
	m.SetHeader("Subject", payload.Subject)
	m.SetBody("text/html", payload.Body)

	if err := h.email.DialAndSend(m); err != nil {
		h.log.Errorf("Failed to send email to %s: %v", payload.To, err)
		return err
	}

	h.log.Infof("Email sent successfully to: %s", payload.To)

	return nil
}

func (h *GroupBuySessionHandler) HandleBuyerSessionEnd(ctx context.Context, t *asynq.Task) error {
	var payload tasks.BuyerGroupBuySessionEndPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	err := h.groupBuyUsecase.ChangeBuyerSessionStatus(ctx, payload.BuyerSessionID, "expired")
	if err != nil {
		h.log.Errorf("failed to change buyer session status: %v", err)
		return fmt.Errorf("failed to change buyer session status: %w", err)
	}
	return nil
}
