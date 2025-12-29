package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/febry3/gamingin/internal/worker/tasks"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GroupBuyUsecaseContract interface {
	CreateGroupBuySession(ctx context.Context, request *dto.GroupBuySessionRequest, sellerID int64) (*dto.GroupBuySessionResponse, error)
	DeleteGroupBuySession(ctx context.Context, sessionID string) error
	FindGroupBuySessionByID(ctx context.Context, sessionID string) (*entity.GroupBuySession, error)
	GetAllGroupBuySessionForSeller(ctx context.Context, sellerID int64) ([]entity.GroupBuySession, error)
	GetAllGroupBuySessionForBuyer(ctx context.Context) ([]entity.GroupBuySession, error)
	ChangeGroupBuySessionStatus(ctx context.Context, sessionID string, status string, sellerID int64) error
	EndSession(ctx context.Context, sessionID string, productVariantID string, sellerID int64) error
	CreateBuyerSession(ctx context.Context, request *dto.CreateBuyerGroupSessionRequest) (string, error)
}

type GroupBuyUsecase struct {
	groupBuySessionRepo   repository.GroupBuySessionRepository
	groupBuyTierRepo      repository.GroupBuyTierRepository
	productRepo           repository.ProductRepository
	productVariantRepo    repository.ProductVariantRepository
	buyerGroupSessionRepo repository.BuyerGroupBuySessionRepository
	tx                    repository.TxManager
	log                   *logrus.Logger
	asynqClient           *asynq.Client
}

func NewGroupBuyUsecase(groupBuySessionRepo repository.GroupBuySessionRepository, groupBuyTierRepo repository.GroupBuyTierRepository, productRepo repository.ProductRepository, productVariantRepo repository.ProductVariantRepository, buyerGroupSessionRepo repository.BuyerGroupBuySessionRepository, tx repository.TxManager, log *logrus.Logger, asynqClient *asynq.Client) GroupBuyUsecaseContract {
	return &GroupBuyUsecase{
		groupBuySessionRepo:   groupBuySessionRepo,
		groupBuyTierRepo:      groupBuyTierRepo,
		productRepo:           productRepo,
		productVariantRepo:    productVariantRepo,
		buyerGroupSessionRepo: buyerGroupSessionRepo,
		tx:                    tx,
		log:                   log,
		asynqClient:           asynqClient,
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

		groupBuySession = &entity.GroupBuySession{
			ProductVariantID: request.ProductVariantID,
			SellerID:         sellerID,
			MinParticipants:  request.MinParticipants,
			MaxParticipants:  request.MaxParticipants,
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

	task, err := tasks.NewGroupBuySessionEndTask(tasks.GroupBuySessionEndPayload{
		SessionID:        groupBuySession.ID,
		ProductVariantID: groupBuySession.ProductVariantID,
		SellerID:         sellerID,
	})

	if err != nil {
		g.log.Errorf("failed to create session end task: %v", err)
	} else {
		_, err = g.asynqClient.Enqueue(task, asynq.ProcessAt(groupBuySession.ExpiresAt))
		if err != nil {
			g.log.Errorf("failed to enqueue session end task: %v", err)
		} else {
			g.log.Infof("Scheduled session end task for session %s at %v", groupBuySession.ID, groupBuySession.ExpiresAt)
		}
	}

	return &dto.GroupBuySessionResponse{
		ID:               groupBuySession.ID,
		ProductVariantID: groupBuySession.ProductVariantID,
		SellerID:         groupBuySession.SellerID,
		MinParticipants:  groupBuySession.MinParticipants,
		MaxParticipants:  groupBuySession.MaxParticipants,
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

func (g *GroupBuyUsecase) EndSession(ctx context.Context, sessionID string, productVariantID string, sellerID int64) error {
	session, err := g.groupBuySessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		g.log.Errorf("failed to find session: %v", err)
		return err
	}

	if session.Status != "active" {
		g.log.Infof("Session %s is already %s, skipping", sessionID, session.Status)
		return nil
	}

	productVariant, err := g.productVariantRepo.GetProductVariant(ctx, productVariantID)
	if err != nil {
		g.log.Errorf("failed to get product variant: %v", err)
		return err
	}

	err = g.groupBuySessionRepo.ChangeStatus(ctx, sessionID, "completed", sellerID)
	if err != nil {
		g.log.Errorf("failed to change session status: %v", err)
		return err
	}

	g.log.Infof("Group buy session %s completed successfully. Product: %s", sessionID, productVariant.ID)

	return nil
}

func (g *GroupBuyUsecase) CreateBuyerSession(ctx context.Context, request *dto.CreateBuyerGroupSessionRequest) (string, error) {
	session, err := g.buyerGroupSessionRepo.GetSessionByOrganizerUserID(ctx, request.OrganizerUserID)
	g.log.Infof("[GroupBuyUsecase] Session: %v", session)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		g.log.Errorf("[GroupBuyUsecase] Failed to get session: %v", err)
		return "", err
	}

	if session != nil {
		g.log.Infof("[GroupBuyUsecase] You already started a session")
		return "", errorx.ErrSessionAlreadyStarted
	}

	if err := g.groupBuySessionRepo.FindByProductVariantID(ctx, request.ProductVariantID); err != nil {
		g.log.Infof("[GroupBuyUsecase] Group buy session not found for product variant %s", request.ProductVariantID)
		return "", errorx.ErrGroupBuySessionNotFound
	}

	buyerGroupSession := &entity.BuyerGroupSession{
		ProductVariantID:    request.ProductVariantID,
		OrganizerUserID:     request.OrganizerUserID,
		Title:               request.Title,
		SessionCode:         "LBX" + uuid.New().String()[:8],
		ExpiresAt:           time.Now().Add(time.Hour * 1),
		CurrentParticipants: 1,
		Status:              "open",
	}

	if err := g.buyerGroupSessionRepo.Create(ctx, buyerGroupSession); err != nil {
		g.log.Errorf("[GroupBuyUsecase] Failed to create session: %v", err)
		return "", err
	}

	return buyerGroupSession.SessionCode, nil
}

func (g *GroupBuyUsecase) GetSessionForBuyerByCode(ctx context.Context, sessionCode string) (*entity.BuyerGroupSession, error) {
	return g.buyerGroupSessionRepo.GetSessionByCode(ctx, sessionCode)
}
