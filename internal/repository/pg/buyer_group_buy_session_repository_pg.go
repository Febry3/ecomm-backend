package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type BuyerGroupBuySessionRepositoryPg struct {
	db *gorm.DB
}

func NewBuyerGroupBuySessionRepositoryPg(db *gorm.DB) repository.BuyerGroupBuySessionRepository {
	return &BuyerGroupBuySessionRepositoryPg{db: db}
}

func (b *BuyerGroupBuySessionRepositoryPg) Create(ctx context.Context, session *entity.BuyerGroupSession) error {
	return b.db.WithContext(ctx).Create(session).Error
}
func (b *BuyerGroupBuySessionRepositoryPg) Delete(ctx context.Context, sessionID string) error {
	return b.db.WithContext(ctx).Delete(&entity.BuyerGroupSession{}, sessionID).Error
}

func (b *BuyerGroupBuySessionRepositoryPg) GetSessionByCode(ctx context.Context, sessionCode string) (*entity.BuyerGroupSession, error) {
	var session entity.BuyerGroupSession
	if err := b.db.WithContext(ctx).Where("session_code = ?", sessionCode).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (b *BuyerGroupBuySessionRepositoryPg) GetSessionByOrganizerUserID(ctx context.Context, organizerUserID int64) (*entity.BuyerGroupSession, error) {
	var session *entity.BuyerGroupSession

	if err := b.db.WithContext(ctx).Where("organizer_user_id = ? AND status = ?", organizerUserID, "open").First(&session).Error; err != nil {
		return nil, err
	}
	return session, nil
}
