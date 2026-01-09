package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type UserWalletRepository interface {
	GetUserWalletByUserID(ctx context.Context, userID int64) (*entity.UserWallet, error)
	GetUserBalanceByUserID(ctx context.Context, userID int64) (float64, error)
	CountUserWallet(ctx context.Context, userID int64) (int64, error)
	CreateOrUpdateUserWallet(ctx context.Context, userWallet *entity.UserWallet) error
}
