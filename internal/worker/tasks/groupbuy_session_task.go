package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeGroupBuySessionEnd      = "groupbuy:session_end"
	TypeGroupBuySessionEndMail  = "groupbuy:session_end:mail"
	TypeBuyerGroupBuySessionEnd = "groupbuy:buyer_session_end"
)

type GroupBuySessionEndPayload struct {
	SessionID        string `json:"session_id"`
	ProductVariantID string `json:"product_variant_id"`
	SellerID         int64  `json:"seller_id"`
}

type BuyerGroupBuySessionEndPayload struct {
	BuyerSessionID string `json:"session_id"`
}

type GroupBuySessionEndMailPayload struct {
	To      string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewGroupBuySessionEndTask(payload GroupBuySessionEndPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeGroupBuySessionEnd, data), nil
}

func NewGroupBuySessionEndMailTask(payload GroupBuySessionEndMailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeGroupBuySessionEndMail, data), nil
}

func NewBuyerGrupBuySessionEndTask(payload BuyerGroupBuySessionEndPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeBuyerGroupBuySessionEnd, data), nil
}
