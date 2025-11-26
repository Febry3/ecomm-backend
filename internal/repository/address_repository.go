package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type AddressRepository interface {
	Create(ctx context.Context, address entity.Address) (entity.Address, error)
	Update(ctx context.Context, address entity.Address) (entity.Address, error)
	FindByUserID(ctx context.Context, id int64) (entity.Address, error)
	Delete(ctx context.Context, id int64) error
}
