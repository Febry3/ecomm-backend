package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type AddressRepository interface {
	Create(ctx context.Context, address entity.Address) (entity.Address, error)
	Update(ctx context.Context, address entity.Address) (entity.Address, error)
	FindAll(ctx context.Context, userId int64) ([]entity.Address, error)
	FindById(ctx context.Context, id string, userId int64) (entity.Address, error)
	Delete(ctx context.Context, id string, userId int64) error
}
