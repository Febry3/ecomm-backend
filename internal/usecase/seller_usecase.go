package usecase

import (
	"context"
	"errors"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SellerUsecaseContract interface {
	RegisterSeller(ctx context.Context, req dto.SellerRequest, userID int64) (*entity.Seller, error)
	UpdateSeller(ctx context.Context, req dto.SellerRequest, userID int64) (*entity.Seller, error)
	GetSeller(ctx context.Context, userID int64) (*entity.Seller, error)
}

type SellerUsecase struct {
	repo repository.SellerRepository
	user repository.UserRepository
	log  *logrus.Logger
}

func NewSellerUsecase(repo repository.SellerRepository, user repository.UserRepository, log *logrus.Logger) SellerUsecaseContract {
	return &SellerUsecase{
		repo: repo,
		user: user,
		log:  log,
	}
}

func (s *SellerUsecase) RegisterSeller(ctx context.Context, request dto.SellerRequest, userID int64) (*entity.Seller, error) {
	if err := validator.New().Struct(request); err != nil {
		s.log.Errorf("[SellerUsecase] Validate Seller Error: %v", err.Error())
		return &entity.Seller{}, err
	}

	user, err := s.user.FindByID(ctx, userID)
	if err != nil {
		s.log.Errorf("[SellerUsecase] User not found")
		return &entity.Seller{}, errorx.ErrUserNotFound
	}

	checkSeller, err := s.repo.GetSeller(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.log.Errorf("[SellerUsecase] Find seller error")
		return &entity.Seller{}, err
	}

	if checkSeller != nil {
		s.log.Errorf("[SellerUsecase] Seller already exists")
		return &entity.Seller{}, errorx.ErrSellerAlreadyExists
	}

	seller, err := s.repo.CreateSeller(ctx, &entity.Seller{
		UserID:        userID,
		StoreName:     request.StoreName,
		StoreSlug:     request.StoreSlug,
		Description:   request.Description,
		LogoURL:       request.LogoURL,
		BusinessEmail: request.BusinessEmail,
		BusinessPhone: request.BusinessPhone,
	})

	if err != nil {
		s.log.Error("[SellerUsecase] Register Seller Error: ", err)
		return nil, err
	}

	user.Role = "seller"
	if _, err := s.user.Update(ctx, user); err != nil {
		s.log.Errorf("[SellerUsecase] Failed to update user role: %v", err)
		// Note: In a production app, you might want to rollback the seller creation here
		// or use a transaction to ensure atomicity.
	}

	return seller, nil
}

func (s *SellerUsecase) UpdateSeller(ctx context.Context, req dto.SellerRequest, userID int64) (*entity.Seller, error) {
	if err := validator.New().Struct(req); err != nil {
		s.log.Errorf("[SellerUsecase] Validate Seller Error: %v", err.Error())
		return &entity.Seller{}, err
	}

	return s.repo.UpdateSeller(ctx, &entity.Seller{
		StoreName:     req.StoreName,
		StoreSlug:     req.StoreSlug,
		Description:   req.Description,
		LogoURL:       req.LogoURL,
		BusinessEmail: req.BusinessEmail,
		BusinessPhone: req.BusinessPhone,
	})
}

func (s *SellerUsecase) GetSeller(ctx context.Context, userID int64) (*entity.Seller, error) {
	return s.repo.GetSeller(ctx, userID)
}
