package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type GroupBuySessionRepositoryPg struct {
	db *gorm.DB
}

func NewGroupBuySessionRepositoryPg(db *gorm.DB) repository.GroupBuySessionRepository {
	return &GroupBuySessionRepositoryPg{db: db}
}

func (g *GroupBuySessionRepositoryPg) GetAllForSeller(ctx context.Context, sellerID int64) ([]entity.GroupBuySession, error) {
	var sessions []entity.GroupBuySession
	if err := g.db.Preload("Seller").Preload("Product").Preload("GroupBuyTiers").Where("seller_id = ?", sellerID).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (g *GroupBuySessionRepositoryPg) GetAllForBuyer(ctx context.Context) ([]entity.GroupBuySession, error) {
	var sessions []entity.GroupBuySession
	if err := g.db.Preload("Seller").Preload("Product").Preload("GroupBuyTiers").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (g *GroupBuySessionRepositoryPg) Create(ctx context.Context, session *entity.GroupBuySession) error {
	return g.db.Create(session).Error
}

func (g *GroupBuySessionRepositoryPg) Delete(ctx context.Context, sessionID string) error {
	return g.db.Delete(&entity.GroupBuySession{}, sessionID).Error
}

func (g *GroupBuySessionRepositoryPg) FindByID(ctx context.Context, sessionID string) (*entity.GroupBuySession, error) {
	var tier entity.GroupBuySession
	if err := g.db.Preload("Seller").Preload("Product").Preload("GroupBuyTiers").First(&tier, sessionID).Error; err != nil {
		return nil, err
	}
	return &tier, nil
}
