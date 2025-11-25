package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type AddressProviderPg struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) repository.AddressRepository {
	return &AddressProviderPg{
		db: db,
	}
}

// Create implements repository.AddressRepository.
func (a *AddressProviderPg) Create(ctx context.Context, address entity.Address) (entity.Address, error) {
	panic("unimplemented")
}

// Delete implements repository.AddressRepository.
func (a *AddressProviderPg) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// FindByUserID implements repository.AddressRepository.
func (a *AddressProviderPg) FindByUserID(ctx context.Context, id int64) (entity.Address, error) {
	panic("unimplemented")
}

// Update implements repository.AddressRepository.
func (a *AddressProviderPg) Update(ctx context.Context, address entity.Address, id int64) (entity.Address, error) {
	panic("unimplemented")
}
