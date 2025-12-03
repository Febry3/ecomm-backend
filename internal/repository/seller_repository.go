package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type SellerRepository interface {
	CreateSeller(ctx context.Context, seller *entity.Seller) (*entity.Seller, error)
	UpdateSeller(ctx context.Context, seller *entity.Seller) (*entity.Seller, error)
	GetSeller(ctx context.Context, sellerID int64) (*entity.Seller, error)
}
