package usecase

import (
	"context"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type GroupBuyUsecaseContract interface {
	CreateGroupBuySession(ctx context.Context, request *dto.GroupBuySessionRequest, sellerID int64) (*dto.GroupBuySessionResponse, error)
	DeleteGroupBuySession(ctx context.Context, sessionID string) error
	FindGroupBuySessionByID(ctx context.Context, sessionID string) (*entity.GroupBuySession, error)
	GetAllGroupBuySessionForSeller(ctx context.Context, sellerID int64) ([]entity.GroupBuySession, error)
	GetAllGroupBuySessionForBuyer(ctx context.Context) ([]entity.GroupBuySession, error)
	ChangeGroupBuySessionStatus(ctx context.Context, sessionID string, status string, sellerID int64) error
}

type GroupBuyUsecase struct {
	groupBuySessionRepo repository.GroupBuySessionRepository
	groupBuyTierRepo    repository.GroupBuyTierRepository
	productRepo         repository.ProductRepository
	productVariantRepo  repository.ProductVariantRepository
	tx                  repository.TxManager
	log                 *logrus.Logger
}

func NewGroupBuyUsecase(groupBuySessionRepo repository.GroupBuySessionRepository, groupBuyTierRepo repository.GroupBuyTierRepository, productRepo repository.ProductRepository, productVariantRepo repository.ProductVariantRepository, tx repository.TxManager, log *logrus.Logger) GroupBuyUsecaseContract {
	return &GroupBuyUsecase{
		groupBuySessionRepo: groupBuySessionRepo,
		groupBuyTierRepo:    groupBuyTierRepo,
		productRepo:         productRepo,
		productVariantRepo:  productVariantRepo,
		tx:                  tx,
		log:                 log,
	}
}

func (g *GroupBuyUsecase) CreateGroupBuySession(ctx context.Context, request *dto.GroupBuySessionRequest, sellerID int64) (*dto.GroupBuySessionResponse, error) {
	var groupBuySession *entity.GroupBuySession
	var tiers []entity.GroupBuyTier
	err := g.tx.WithTransaction(ctx, func(txCtx context.Context) error {
		productVariant, err := g.productVariantRepo.GetProductVariant(ctx, request.ProductVariantID)
		if err != nil {
			g.log.Errorf("failed to get product variant: %v", err)
			return err
		}

		if productVariant.Stock.CurrentStock < request.MaxQuantity {
			g.log.Errorf("product variant stock is not enough")
			return errorx.ErrInsufficientStock
		}

		productVariant.Stock.CurrentStock -= request.MaxQuantity
		productVariant.Stock.ReservedStock += request.MaxQuantity

		if err := g.productVariantRepo.UpdateProductVariant(ctx, productVariant, request.ProductVariantID); err != nil {
			g.log.Errorf("failed to update product variant: %v", err)
			return err
		}

		sessionCode := "JMK" + uuid.New().String()[:8]
		groupBuySession = &entity.GroupBuySession{
			ProductVariantID: request.ProductVariantID,
			SellerID:         sellerID,
			MinParticipants:  request.MinParticipants,
			MaxParticipants:  request.MaxParticipants,
			SessionCode:      sessionCode,
			ExpiresAt:        request.ExpiresAt,
		}

		if err := g.groupBuySessionRepo.Create(txCtx, groupBuySession); err != nil {
			g.log.Errorf("failed to create group buy session: %v", err)
			return err
		}

		for _, val := range request.Tiers {
			tier := &entity.GroupBuyTier{
				GroupBuySessionID:    groupBuySession.ID,
				ParticipantThreshold: val.ParticipantThreshold,
				DiscountPercentage:   float64(val.DiscountPercentage),
			}
			err := g.groupBuyTierRepo.Create(txCtx, tier)
			if err != nil {
				g.log.Errorf("failed to create group buy tier: %v", err)
				return err
			}
			tiers = append(tiers, *tier)
		}
		return nil
	})
	if err != nil {
		g.log.Errorf("failed to create group buy session: %v", err)
		return &dto.GroupBuySessionResponse{}, err
	}

	return &dto.GroupBuySessionResponse{
		ID:               groupBuySession.ID,
		ProductVariantID: groupBuySession.ProductVariantID,
		SellerID:         groupBuySession.SellerID,
		MinParticipants:  groupBuySession.MinParticipants,
		MaxParticipants:  groupBuySession.MaxParticipants,
		SessionCode:      groupBuySession.SessionCode,
		ExpiresAt:        groupBuySession.ExpiresAt,
		Tiers:            tiers,
	}, nil
}

func (g *GroupBuyUsecase) DeleteGroupBuySession(ctx context.Context, sessionID string) error {
	return g.groupBuySessionRepo.Delete(ctx, sessionID)
}

func (g *GroupBuyUsecase) FindGroupBuySessionByID(ctx context.Context, sessionID string) (*entity.GroupBuySession, error) {
	return g.groupBuySessionRepo.FindByID(ctx, sessionID)
}

func (g *GroupBuyUsecase) GetAllGroupBuySessionForSeller(ctx context.Context, sellerID int64) ([]entity.GroupBuySession, error) {
	return g.groupBuySessionRepo.GetAllForSeller(ctx, sellerID)
}

func (g *GroupBuyUsecase) GetAllGroupBuySessionForBuyer(ctx context.Context) ([]entity.GroupBuySession, error) {
	return g.groupBuySessionRepo.GetAllForBuyer(ctx)
}

func (g *GroupBuyUsecase) ChangeGroupBuySessionStatus(ctx context.Context, sessionID string, status string, sellerID int64) error {
	return g.groupBuySessionRepo.ChangeStatus(ctx, sessionID, status, sellerID)
}
