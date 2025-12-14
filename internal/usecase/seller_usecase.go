package usecase

import (
	"context"
	"errors"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/infra/storage"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SellerUsecaseContract interface {
	RegisterSeller(ctx context.Context, req dto.SellerRequest, userID int64, fileData []byte) (*entity.Seller, error)
	UpdateSeller(ctx context.Context, req dto.UpdateSellerRequest, userID int64, fileData []byte) (*entity.Seller, error)
	GetSeller(ctx context.Context, userID int64) (*entity.Seller, error)
}

type SellerUsecase struct {
	repo    repository.SellerRepository
	user    repository.UserRepository
	tx      repository.TxManager
	storage storage.ObjectStorage
	log     *logrus.Logger
}

func NewSellerUsecase(repo repository.SellerRepository, user repository.UserRepository, tx repository.TxManager, log *logrus.Logger, storage storage.ObjectStorage) SellerUsecaseContract {
	return &SellerUsecase{
		repo:    repo,
		user:    user,
		tx:      tx,
		log:     log,
		storage: storage,
	}
}

func (s *SellerUsecase) RegisterSeller(ctx context.Context, request dto.SellerRequest, userID int64, fileData []byte) (*entity.Seller, error) {
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

	var logoURL string
	if fileData != nil {
		fileName := uuid.New().String()
		logoURL, err = s.storage.Upload(ctx, fileName, fileData, "store_logo")
		if err != nil {
			s.log.Errorf("[SellerUsecase] Upload logo error: %v", err.Error())
			return &entity.Seller{}, err
		}
	}

	var seller *entity.Seller

	// Use transaction to ensure atomicity - both seller creation and role update must succeed
	err = s.tx.WithTransaction(ctx, func(txCtx context.Context) error {
		var txErr error
		seller, txErr = s.repo.CreateSeller(txCtx, &entity.Seller{
			UserID:        userID,
			StoreName:     request.StoreName,
			StoreSlug:     request.StoreSlug,
			Description:   request.Description,
			LogoURL:       logoURL,
			BusinessEmail: request.BusinessEmail,
			BusinessPhone: request.BusinessPhone,
		})
		if txErr != nil {
			s.log.Error("[SellerUsecase] Create Seller Error: ", txErr)
			return txErr
		}

		user.Role = "seller"
		if _, txErr = s.user.Update(txCtx, user); txErr != nil {
			s.log.Errorf("[SellerUsecase] Failed to update user role: %v", txErr)
			return txErr
		}

		return nil
	})

	if err != nil {
		s.log.Error("[SellerUsecase] RegisterSeller transaction failed: ", err)
		return nil, err
	}

	return seller, nil
}

func (s *SellerUsecase) UpdateSeller(ctx context.Context, req dto.UpdateSellerRequest, userID int64, fileData []byte) (*entity.Seller, error) {
	if err := validator.New().Struct(req); err != nil {
		s.log.Errorf("[SellerUsecase] Validate Seller Error: %v", err.Error())
		return &entity.Seller{}, err
	}

	existingSeller, err := s.repo.GetSeller(ctx, userID)
	if err != nil {
		s.log.Errorf("[SellerUsecase] Get Seller Error: %v", err.Error())
		return &entity.Seller{}, err
	}

	if fileData != nil {
		fileName := uuid.New().String()
		logoURL, err := s.storage.Upload(ctx, fileName, fileData, "store_logo")
		if err != nil {
			s.log.Errorf("[SellerUsecase] Upload logo error: %v", err.Error())
			return &entity.Seller{}, err
		}
		existingSeller.LogoURL = logoURL
	}

	existingSeller.StoreName = req.StoreName
	existingSeller.Description = req.Description
	existingSeller.BusinessEmail = req.BusinessEmail
	existingSeller.BusinessPhone = req.BusinessPhone

	return s.repo.UpdateSeller(ctx, existingSeller)
}

func (s *SellerUsecase) GetSeller(ctx context.Context, userID int64) (*entity.Seller, error) {
	return s.repo.GetSeller(ctx, userID)
}
