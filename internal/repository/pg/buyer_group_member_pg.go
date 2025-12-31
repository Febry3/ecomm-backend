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
	db := TxFromContext(ctx, b.db)
	return db.Create(member).Error
}

func (b *BuyerGroupMemberRepositoryPg) Delete(ctx context.Context, memberID string) error {
	return b.db.WithContext(ctx).Delete(&entity.BuyerGroupMember{}, memberID).Error
}

func (b *BuyerGroupMemberRepositoryPg) GetMembersBySessionID(ctx context.Context, sessionID string) ([]entity.BuyerGroupMember, error) {
	var member []entity.BuyerGroupMember
	if err := b.db.WithContext(ctx).Where("session_id = ?", sessionID).Preload("User").Find(&member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

func (b *BuyerGroupBuySessionRepositoryPg) AddMember(ctx context.Context, buyer_session *entity.BuyerGroupSession) error {
	return b.db.WithContext(ctx).Model(buyer_session).Update("current_participants", buyer_session.CurrentParticipants+1).Error
}

func (b *BuyerGroupBuySessionRepositoryPg) ChangeBuyerSessionStatus(ctx context.Context, buyerSessionID string, status string) error {
	return b.db.WithContext(ctx).Model(&entity.BuyerGroupSession{}).Where("id = ?", buyerSessionID).Update("status", status).Error
}
