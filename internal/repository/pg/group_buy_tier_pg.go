package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type GroupBuyTierRepositoryPg struct {
	db *gorm.DB
}

func NewGroupBuyTierRepositoryPg(db *gorm.DB) repository.GroupBuyTierRepository {
	return &GroupBuyTierRepositoryPg{db: db}
}

func (g *GroupBuyTierRepositoryPg) Create(ctx context.Context, tier *entity.GroupBuyTier) error {
	return g.db.Create(tier).Error
}

func (g *GroupBuyTierRepositoryPg) FindByID(ctx context.Context, tierID string) (*entity.GroupBuyTier, error) {
	var tier entity.GroupBuyTier
	if err := g.db.Where("id = ?", tierID).First(&tier).Error; err != nil {
		return nil, err
	}
	return &tier, nil
}
