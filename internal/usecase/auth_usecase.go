package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/helpers"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthUsecaseContract interface {
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, string, error)
	Register(ctx context.Context, request dto.RegisterRequest) (dto.RegisterResponse, error)
	Logout(ctx context.Context, tokenId string) error
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
	LoginOrRegisterWithGoogle(ctx context.Context, request dto.LoginWithGoogleData) (dto.LoginResponse, string, error)
}

type AuthUsecase struct {
	token        repository.TokenRepository
	user         repository.UserRepository
	authProvider repository.AuthProviderRepository
	log          *logrus.Logger
	jwt          helpers.JwtService
}

func NewAuthUsecase(user repository.UserRepository, log *logrus.Logger, jwt helpers.JwtService, token repository.TokenRepository, authProvider repository.AuthProviderRepository) AuthUsecaseContract {
	return &AuthUsecase{
		token:        token,
		user:         user,
		log:          log,
		jwt:          jwt,
		authProvider: authProvider,
	}
}

func (a *AuthUsecase) Register(ctx context.Context, request dto.RegisterRequest) (dto.RegisterResponse, error) {
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
		Username:    request.Username,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
	}

	if err := a.user.Create(ctx, &user); err != nil {
		a.log.Errorf("[AuthUsecase] Create User Error: %v", err.Error())
		return dto.RegisterResponse{}, err
	}

	authProvider := entity.AuthProvider{
		UserId:     user.ID,
		Password:   sql.NullString{String: hashedPassword, Valid: true},
		Provider:   "email",
		ProviderId: user.Email,
	}

	if err := a.authProvider.Create(ctx, &authProvider); err != nil {
		a.log.Errorf("[AuthUsecase] Create Provider error: %v", err.Error())
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

func (a *AuthUsecase) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, string, error) {
	if err := validator.New().Struct(request); err != nil {
		a.log.Errorf("[AuthUsecase] Validate Register Error: %v", err.Error())
		return dto.LoginResponse{}, "", err
	}

	user, isFound, err := a.user.FindByEmail(ctx, request.Email)
	if err != nil {
		if !isFound {
			a.log.Errorf("[AuthUsecase] Email Not Found: %v", err.Error())
			return dto.LoginResponse{}, "", errorx.ErrInvalidLogin
		}
		a.log.Errorf("[AuthUsecase] FindByEmail Error: %v", err.Error())
		return dto.LoginResponse{}, "", err
	}

	authProvider, err := a.authProvider.FindByUserID(ctx, user.ID)
	if err != nil {
		a.log.Errorf("[AuthUsecase] FindByEmail on AuthProvider Error: %v", err.Error())
		return dto.LoginResponse{}, "", err
	}

	isMatch := helpers.Compare([]byte(authProvider.Password.String), request.Password)
	if !isMatch {
		return dto.LoginResponse{}, "", errorx.ErrInvalidCredentials
	}

	accessToken := a.jwt.IssueAccessToken(dto.JwtPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})

	plainTextRefreshToken := uuid.New().String()

	expiresAt := time.Now().Add(a.jwt.Config.RefreshTTL)
	refreshToken := entity.RefreshToken{
		TokenId:    uuid.New().String(),
		UserId:     user.ID,
		TokenHash:  plainTextRefreshToken,
		Role:       user.Role,
		DeviceInfo: request.DeviceInfo,
		IsRevoked:  false,
		ExpiresAt:  expiresAt,
	}

	if _, err := a.token.CreateOrUpdate(ctx, &refreshToken); err != nil {
		a.log.Errorf("[AuthUsecase] Failed to save refresh token: %v", err)
		return dto.LoginResponse{}, "", err
	}

	return dto.LoginResponse{
		ID:          user.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		AccessToken: accessToken,
	}, plainTextRefreshToken, nil
}

func (a *AuthUsecase) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", errorx.ErrTokenEmpty
	}

	token, err := a.token.FindByAccessToken(ctx, refreshToken)
	if err != nil {
		a.log.Errorf("[AuthUsecase] Find Refresh Token Error: %v", err.Error())
		return "", err
	}

	a.log.Debug("[AuthUsecase] Refresh Token Expired", token.IsExpired(), token.ExpiresAt, time.Now())
	if token.IsExpired() {
		a.log.Errorf("[AuthUsecase] Refresh Token is expired")
		return "", errorx.ErrTokenExpired
	}

	user, err := a.user.FindByID(ctx, token.UserId)

	if err != nil {
		a.log.Errorf("[AuthUsecase] Find User Error: %v", err.Error())
		return "", err
	}

	newAccessToken := a.jwt.IssueAccessToken(dto.JwtPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})

	return newAccessToken, nil
}

func (a *AuthUsecase) LoginOrRegisterWithGoogle(ctx context.Context, request dto.LoginWithGoogleData) (dto.LoginResponse, string, error) {
	authProvider, err := a.authProvider.FindByProviderId(ctx, request.ID, "google")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.log.Errorf("[AuthUsecase] Find Provider Error: %v", err.Error())
		return dto.LoginResponse{}, "", err
	}

	user, err := a.user.FindByID(ctx, authProvider.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tempUser := entity.User{
				Email:      request.Email,
				FirstName:  request.FirstName,
				ProfileUrl: request.PictureUrl,
			}
			err = a.user.Create(ctx, &tempUser)
			if err != nil {
				a.log.Errorf("[AuthUsecase] Create User Error: %v", err.Error())
				return dto.LoginResponse{}, "", err
			}

			authProvider := entity.AuthProvider{
				UserId:     tempUser.ID,
				Password:   sql.NullString{},
				Provider:   "google",
				ProviderId: request.ID,
			}

			if err := a.authProvider.Create(ctx, &authProvider); err != nil {
				a.log.Errorf("[AuthUsecase] Create Provider error: %v", err.Error())
				return dto.LoginResponse{}, "", err
			}
			user = tempUser
		} else {
			a.log.Errorf("[AuthUsecase] Find User Error: %v", err.Error())
			return dto.LoginResponse{}, "", err
		}
	}

	accessToken := a.jwt.IssueAccessToken(dto.JwtPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})

	plainTextRefreshToken := uuid.New().String()

	expiresAt := time.Now().Add(a.jwt.Config.RefreshTTL)
	refreshToken := entity.RefreshToken{
		TokenId:    uuid.New().String(),
		UserId:     user.ID,
		TokenHash:  plainTextRefreshToken,
		Role:       user.Role,
		DeviceInfo: request.DeviceInfo,
		IsRevoked:  false,
		ExpiresAt:  expiresAt,
	}

	if _, err := a.token.CreateOrUpdate(ctx, &refreshToken); err != nil {
		a.log.Errorf("[AuthUsecase] Failed to save refresh token: %v", err)
		return dto.LoginResponse{}, "", err
	}

	return dto.LoginResponse{
		ID:          user.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		AccessToken: accessToken,
		ProfileUrl:  user.ProfileUrl,
	}, plainTextRefreshToken, nil
}

func (a *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errorx.ErrTokenEmpty
	}

	err := a.token.DeleteByAccessToken(ctx, refreshToken)
	if err != nil {
		a.log.Errorf("[AuthUsecase] Delete Refresh Token Error: %v", err.Error())
		return err
	}
	return nil
}
