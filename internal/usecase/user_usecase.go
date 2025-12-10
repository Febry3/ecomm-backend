package usecase

import (
	"context"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/infra/storage"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserUsecaseContract interface {
	UpdateAvatar(ctx context.Context, fileData []byte, userId int64) (dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userRequest dto.UserRequest) (dto.UserResponse, error)
	GetProfile(ctx context.Context, userId int64) (dto.UserResponse, error)
}

type UserUsecase struct {
	storage storage.ObjectStorage
	user    repository.UserRepository
	seller  repository.SellerRepository
	log     *logrus.Logger
}

func NewUserUsecase(user repository.UserRepository, log *logrus.Logger, storage storage.ObjectStorage, seller repository.SellerRepository) UserUsecaseContract {
	return &UserUsecase{
		user:    user,
		log:     log,
		storage: storage,
		seller:  seller,
	}
}

func (u *UserUsecase) GetProfile(ctx context.Context, userId int64) (dto.UserResponse, error) {
	user, err := u.user.FindByID(ctx, userId)
	if err != nil {
		u.log.Error("[UserUsecase] FindByID Error", err.Error())
		return dto.UserResponse{}, err
	}

	var seller *entity.Seller
	if user.Role == "seller" {
		seller, _ = u.seller.GetSeller(ctx, user.ID)
	}

	if seller == nil {
		seller = &entity.Seller{ID: 0}
	}
	return dto.UserResponse{
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		ID:          user.ID,
		ProfileUrl:  user.ProfileUrl,
		Role:        user.Role,
		SellerID:    seller.ID,
	}, nil
}

func (u *UserUsecase) UpdateProfile(ctx context.Context, userRequest dto.UserRequest) (dto.UserResponse, error) {
	user, err := u.user.FindByID(ctx, userRequest.UserID)
	if err != nil {
		u.log.Error("[UserUsecase] FindByID Error", err.Error())
		return dto.UserResponse{}, err
	}

	userRequest.UpdateEntity(&user)
	updatedUser, err := u.user.Update(ctx, user)
	if err != nil {
		u.log.Error("[UserUsecase] UpdateUser Error", err.Error())
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:          updatedUser.ID,
		ProfileUrl:  updatedUser.ProfileUrl,
		Username:    updatedUser.Username,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		PhoneNumber: updatedUser.PhoneNumber,
		Role:        updatedUser.Role,
	}, nil
}

func (u *UserUsecase) UpdateAvatar(ctx context.Context, fileData []byte, userId int64) (dto.UserResponse, error) {
	fileName := uuid.New().String()

	user, err := u.user.FindByID(ctx, userId)
	if err != nil {
		u.log.Error("[UserUsecase] FindByID Error", err.Error())
		return dto.UserResponse{}, err
	}
	newUrl, err := u.storage.Upload(ctx, fileName, fileData, "avatar")
	if err != nil {
		u.log.Error("[UserUsecase] Upload Error", err.Error())
		return dto.UserResponse{}, err
	}

	user.ProfileUrl = newUrl
	updatedUser, err := u.user.Update(ctx, user)
	if err != nil {
		u.log.Error("[UserUsecase] UpdateUser Error", err.Error())
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:          updatedUser.ID,
		ProfileUrl:  updatedUser.ProfileUrl,
		Username:    updatedUser.Username,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		PhoneNumber: updatedUser.PhoneNumber,
	}, nil
}
