package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type GroupBuySessionRepository interface {
	Create(ctx context.Context, session *entity.GroupBuySession) error
	FindByID(ctx context.Context, sessionID string) (*entity.GroupBuySession, error)
	FindByProductVariantID(ctx context.Context, productVariantID string) (*entity.GroupBuySession, error)
	Delete(ctx context.Context, sessionID string) error
	GetAllForSeller(ctx context.Context, sellerID int64) ([]entity.GroupBuySession, error)
	GetAllForBuyer(ctx context.Context) ([]entity.GroupBuySession, error)
	ChangeStatus(ctx context.Context, sessionID string, status string, sellerID int64) error
}
