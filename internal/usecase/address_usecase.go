package usecase

import (
	"context"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AddressUsecaseContract interface {
	Create(ctx context.Context, address dto.AddressRequest, userId int64) (entity.Address, error)
	Update(ctx context.Context, address dto.AddressRequest, id string, userId int64) (entity.Address, error)
	Delete(ctx context.Context, id string, userId int64) error
	GetAll(ctx context.Context, userId int64) ([]entity.Address, error)
}

type AddressUsecase struct {
	address repository.AddressRepository
	user    repository.UserRepository
	log     *logrus.Logger
}

func NewAddressUsecase(address repository.AddressRepository, user repository.UserRepository, log *logrus.Logger) AddressUsecaseContract {
	return &AddressUsecase{
		address: address,
		user:    user,
		log:     log,
	}
}

func (a *AddressUsecase) GetAll(ctx context.Context, userId int64) ([]entity.Address, error) {
	return a.address.FindAll(ctx, userId)
}

func (a *AddressUsecase) Create(ctx context.Context, request dto.AddressRequest, userId int64) (entity.Address, error) {
	if err := validator.New().Struct(request); err != nil {
		a.log.Errorf("[AddressUsecase] Validate Address Error: %v", err.Error())
		return entity.Address{}, err
	}

	user, err := a.user.FindByID(ctx, userId)
	if err != nil {
		a.log.Errorf("[AddressUsecase] Find User Error: %v", err.Error())
		return entity.Address{}, err
	}

	address, err := a.address.Create(ctx, entity.Address{
		UserID:        user.ID,
		StreetAddress: request.StreetAddress,
		RT:            request.RT,
		RW:            request.RW,
		Village:       request.Village,
		District:      request.District,
		City:          request.City,
		Province:      request.Province,
		PostalCode:    request.PostalCode,
		Notes:         request.Notes,
		IsDefault:     true,
	})

	if err != nil {
		a.log.Errorf("[AddressUsecase] Create Address Error: %v", err.Error())
		return entity.Address{}, err
	}

	return address, err
}

func (a *AddressUsecase) Delete(ctx context.Context, id string, userId int64) error {
	return a.address.Delete(ctx, id, userId)
}

func (a *AddressUsecase) Update(ctx context.Context, request dto.AddressRequest, id string, userId int64) (entity.Address, error) {
	address, err := a.address.FindById(ctx, id, userId)
	if err != nil {
		a.log.Errorf("[AddressUsecase] Find Address Error: %v", err.Error())
		return entity.Address{}, err
	}
	request.UpdateEntity(&address)
	updatedAddress, err := a.address.Update(ctx, address)
	if err != nil {
		a.log.Errorf("[AddressUsecase] Update Address Error: %v", err.Error())
		return entity.Address{}, err
	}

	return updatedAddress, err
}
