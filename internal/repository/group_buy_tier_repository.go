package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type GroupBuyTierRepository interface {
	Create(ctx context.Context, tier *entity.GroupBuyTier) error
	FindByID(ctx context.Context, tierID string) (*entity.GroupBuyTier, error)
}
