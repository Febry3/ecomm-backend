package pg

import (
	"context"
	"errors"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type userWalletRepositoryPg struct {
	db *gorm.DB
}

func NewUserWalletRepositoryPg(db *gorm.DB) repository.UserWalletRepository {
	return &userWalletRepositoryPg{
		db: db,
	}
}

func (r *userWalletRepositoryPg) GetUserWalletByUserID(ctx context.Context, userID int64) (*entity.UserWallet, error) {
	var userWallet entity.UserWallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userWallet).Error
	if err != nil {
		return nil, err
	}
	return &userWallet, nil
}

func (r *userWalletRepositoryPg) CreateOrUpdateUserWallet(ctx context.Context, userWallet *entity.UserWallet) error {
	return r.db.WithContext(ctx).Save(userWallet).Error
}

func (r *userWalletRepositoryPg) GetUserBalanceByUserID(ctx context.Context, userID int64) (float64, error) {
	var userWallet entity.UserWallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userWallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return userWallet.Balance, nil
}

func (r *userWalletRepositoryPg) CountUserWallet(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.UserWallet{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
