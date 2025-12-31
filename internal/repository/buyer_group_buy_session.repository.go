package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type BuyerGroupBuySessionRepository interface {
	Create(ctx context.Context, session *entity.BuyerGroupSession) error
	GetSessionByCode(ctx context.Context, sessionCode string) (*entity.BuyerGroupSession, error)
	GetSessionByOrganizerUserID(ctx context.Context, organizerUserID int64) (*entity.BuyerGroupSession, error)
	Delete(ctx context.Context, sessionID string) error
	AddMember(ctx context.Context, buyer_session *entity.BuyerGroupSession) error
	ChangeBuyerSessionStatus(ctx context.Context, buyerSessionID string, status string) error
	GetSessionByID(ctx context.Context, buyerSessionID string) (*entity.BuyerGroupSession, error)
}
