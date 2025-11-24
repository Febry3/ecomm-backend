package usecase

import (
	"context"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/sirupsen/logrus"
)

type UserUsecaseContract interface {
	UpdateProfile(ctx context.Context, userRequest dto.UserRequest) (dto.UserResponse, error)
	GetProfile(ctx context.Context, userId int64) (dto.UserResponse, error)
}

type UserUsecase struct {
	user repository.UserRepository
	log  *logrus.Logger
}

func NewUserUsecase(user repository.UserRepository, log *logrus.Logger) UserUsecaseContract {
	return &UserUsecase{
		user: user,
		log:  log,
	}
}

func (u *UserUsecase) GetProfile(ctx context.Context, userId int64) (dto.UserResponse, error) {
	user, err := u.user.FindByID(ctx, userId)
	if err != nil {
		u.log.Error("")
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		ID:          user.ID,
		ProfileUrl:  user.ProfileUrl,
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
		ProfileUrl:  updatedUser.ProfileUrl,
		Username:    updatedUser.Username,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		PhoneNumber: updatedUser.PhoneNumber,
	}, nil
}
