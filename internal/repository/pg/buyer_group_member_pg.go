package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type BuyerGroupMemberRepositoryPg struct {
	db *gorm.DB
}

func NewBuyerGroupMemberRepositoryPg(db *gorm.DB) repository.BuyerGroupMemberRepository {
	return &BuyerGroupMemberRepositoryPg{db: db}
}

func (b *BuyerGroupMemberRepositoryPg) Create(ctx context.Context, member *entity.BuyerGroupMember) error {
	return b.db.WithContext(ctx).Create(member).Error
}

func (b *BuyerGroupMemberRepositoryPg) Delete(ctx context.Context, memberID string) error {
	return b.db.WithContext(ctx).Delete(&entity.BuyerGroupMember{}, memberID).Error
}

func (b *BuyerGroupMemberRepositoryPg) GetMembersBySessionID(ctx context.Context, sessionID string) ([]entity.BuyerGroupMember, error) {
	var member []entity.BuyerGroupMember
	if err := b.db.WithContext(ctx).Where("session_id = ?", sessionID).Find(&member).Error; err != nil {
		return nil, err
	}
	return member, nil
}
