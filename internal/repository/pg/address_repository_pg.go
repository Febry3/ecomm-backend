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

func NewAddressRepositoryPg(db *gorm.DB) repository.AddressRepository {
	return &AddressProviderPg{
		db: db,
	}
}

func (a *AddressProviderPg) Create(ctx context.Context, address entity.Address) (entity.Address, error) {
	if err := a.db.Create(&address).Error; err != nil {
		return entity.Address{}, err
	}
	return address, nil
}

func (a *AddressProviderPg) Delete(ctx context.Context, id string, userId int64) error {
	if err := a.db.Delete(&entity.Address{}, "address_id = ? AND user_id = ?", id, userId).Error; err != nil {
		return err
	}
	return nil
}

func (a *AddressProviderPg) FindAll(ctx context.Context, userId int64) ([]entity.Address, error) {
	var address []entity.Address
	if err := a.db.Where("user_id = ?", userId).Find(&address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

func (a *AddressProviderPg) Update(ctx context.Context, address entity.Address) (entity.Address, error) {
	if err := a.db.Save(&address).Error; err != nil {
		return entity.Address{}, err
	}
	return address, nil
}

func (a *AddressProviderPg) FindById(ctx context.Context, id string, userId int64) (entity.Address, error) {
	var address entity.Address
	if err := a.db.Where("address_id = ? AND user_id = ?", id, userId).First(&address).Error; err != nil {
		return entity.Address{}, err
	}
	return address, nil
}
