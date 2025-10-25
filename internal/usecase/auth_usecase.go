package usecase

import (
	"context"
	"errors"
	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/helpers"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthUsecaseContract interface {
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	Register(ctx context.Context, request dto.RegisterRequest) (dto.RegisterResponse, error)
	Logout(ctx context.Context, tokenId string) error
}

type AuthUsecase struct {
	user repository.UserRepository
	log  *logrus.Logger
	jwt  helpers.JwtService
}

func NewAuthUsecase(user repository.UserRepository, log *logrus.Logger, jwt helpers.JwtService) AuthUsecaseContract {
	return &AuthUsecase{
		user: user,
		log:  log,
		jwt:  jwt,
	}
}

func (a AuthUsecase) Register(ctx context.Context, request dto.RegisterRequest) (dto.RegisterResponse, error) {
	if err := validator.New().Struct(request); err != nil {
		a.log.Errorf("[AuthUsecase] Validate Register Error: %v", err.Error())
		return dto.RegisterResponse{}, err
	}

	_, isTaken, err := a.user.FindByEmail(ctx, request.Email)
	if isTaken {
		a.log.Errorf("[AuthUsecase] Email Taken")
		return dto.RegisterResponse{}, errorx.ErrEmailTaken
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.log.Errorf("[AuthUsecase] Register Error: %v", err.Error())
		return dto.RegisterResponse{}, err
	}

	hashedPassword, err := helpers.Hash(request.Password)
	if err != nil {
		a.log.Errorf("[AuthUsecase] Hash Password Error: %v", err.Error())
		return dto.RegisterResponse{}, err
	}

	user := entity.User{
		Email:       request.Email,
		Password:    hashedPassword,
		Username:    request.Username,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
	}

	if err := a.user.Create(ctx, &user); err != nil {
		a.log.Errorf("[AuthUsecase] Create User Error: %v", err.Error())
		return dto.RegisterResponse{}, err
	}

	return dto.RegisterResponse{
		ID:          user.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}, nil
}

func (a AuthUsecase) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	if err := validator.New().Struct(request); err != nil {
		a.log.Errorf("[AuthUsecase] Validate Register Error: %v", err.Error())
		return dto.LoginResponse{}, err
	}

	user, isFound, err := a.user.FindByEmail(ctx, request.Email)
	if err != nil {
		if !isFound {
			a.log.Errorf("[AuthUsecase] Email Not Found: %v", err.Error())
			return dto.LoginResponse{}, errorx.ErrInvalidLogin
		}
		a.log.Errorf("[AuthUsecase] FindByEmail Error: %v", err.Error())
		return dto.LoginResponse{}, err
	}

	isMatch := helpers.Compare([]byte(user.Password), request.Password)

	if !isMatch {
		return dto.LoginResponse{}, errorx.ErrInvalidCredentials
	}

	accessToken := a.jwt.IssueJwt(dto.JwtPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})

	return dto.LoginResponse{
		ID:          user.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}

func (a AuthUsecase) Logout(ctx context.Context, tokenId string) error {
	//TODO implement me
	panic("implement me")
}
