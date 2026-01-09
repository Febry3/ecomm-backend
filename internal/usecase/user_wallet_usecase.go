package usecase

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type UserWalletUsecaseContract interface {
	GetUserWalletByUserID(ctx context.Context, userID int64) (*entity.UserWallet, error)
	CreateOrUpdateUserWallet(ctx context.Context, userID int64, balance float64) error
}

type userWalletUsecase struct {
	userWalletRepository repository.UserWalletRepository
}

func NewUserWalletUsecase(userWalletRepository repository.UserWalletRepository) UserWalletUsecaseContract {
	return &userWalletUsecase{
		userWalletRepository: userWalletRepository,
	}
}

func (u *userWalletUsecase) GetUserWalletByUserID(ctx context.Context, userID int64) (*entity.UserWallet, error) {
	return u.userWalletRepository.GetUserWalletByUserID(ctx, userID)
}

func (u *userWalletUsecase) CreateOrUpdateUserWallet(ctx context.Context, userID int64, balance float64) error {

	userWallet, err := u.userWalletRepository.GetUserWalletByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userWallet := &entity.UserWallet{
				UserID:  userID,
				Balance: balance,
			}
			return u.userWalletRepository.CreateOrUpdateUserWallet(ctx, userWallet)
		}

		return err
	}

	if userWallet != nil {
		userWallet.Balance += balance
		return u.userWalletRepository.CreateOrUpdateUserWallet(ctx, userWallet)
	}

	return nil
}
