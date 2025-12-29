package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type BuyerGroupMemberRepository interface {
	Create(ctx context.Context, member *entity.BuyerGroupMember) error
	Delete(ctx context.Context, memberID string) error
	GetMembersBySessionID(ctx context.Context, sessionID string) ([]entity.BuyerGroupMember, error)
}
