package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// Task types
const (
	TypeEmailDelivery  = "email:deliver"
	TypeWelcomeEmail   = "email:welcome"
	TypeOrderConfirmed = "order:confirmed"
	TypeGroupBuyNotify = "groupbuy:notify"
)

// EmailDeliveryPayload contains the data needed for sending an email
type EmailDeliveryPayload struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
	TemplateID string `json:"template_id,omitempty"`
}

// WelcomeEmailPayload contains data for welcome emails
type WelcomeEmailPayload struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// NewEmailDeliveryTask creates a new email delivery task
func NewEmailDeliveryTask(payload EmailDeliveryPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, data), nil
}

// NewWelcomeEmailTask creates a new welcome email task
func NewWelcomeEmailTask(payload WelcomeEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWelcomeEmail, data), nil
}

// HandleEmailDeliveryTask processes email delivery tasks
func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// TODO: Implement your email sending logic here
	// Example: use an SMTP service, SendGrid, AWS SES, etc.
	fmt.Printf("Sending email to %s: %s\n", payload.Email, payload.Subject)

	return nil
}

// HandleWelcomeEmailTask processes welcome email tasks
func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// TODO: Implement welcome email logic
	fmt.Printf("Sending welcome email to %s (%s)\n", payload.Username, payload.Email)

	return nil
}
