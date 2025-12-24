package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeGroupBuySessionEnd = "groupbuy:session_end"
)

type GroupBuySessionEndPayload struct {
	SessionID        string `json:"session_id"`
	ProductVariantID string `json:"product_variant_id"`
	SellerID         int64  `json:"seller_id"`
}

func NewGroupBuySessionEndTask(payload GroupBuySessionEndPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeGroupBuySessionEnd, data), nil
}
