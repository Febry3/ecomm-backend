package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SellerRepositoryPg struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewSellerRepositoryPg(db *gorm.DB, log *logrus.Logger) repository.SellerRepository {
	return &SellerRepositoryPg{db: db, log: log}
}

func (s *SellerRepositoryPg) CreateSeller(ctx context.Context, seller *entity.Seller) (*entity.Seller, error) {
	db := TxFromContext(ctx, s.db)
	result := db.Create(&seller)
	if result.Error != nil {
		s.log.Errorf("[SellerRepositoryPg] Create Seller Error: %v", result.Error)
		return nil, result.Error
	}
	return seller, nil
}

func (s *SellerRepositoryPg) GetSeller(ctx context.Context, userID int64) (*entity.Seller, error) {
	var seller entity.Seller
	result := s.db.Where("user_id = ?", userID).First(&seller)

	if result.Error != nil {
		s.log.Errorf("[SellerRepositoryPg] Get Seller Error: %v", result.Error)
		return nil, result.Error
	}
	return &seller, nil
}

func (s *SellerRepositoryPg) UpdateSeller(ctx context.Context, seller *entity.Seller) (*entity.Seller, error) {
	result := s.db.Save(&seller)
	if result.Error != nil {
		s.log.Errorf("[SellerRepositoryPg] Update Seller Error: %v", result.Error)
		return nil, result.Error
	}
	return seller, nil
}
