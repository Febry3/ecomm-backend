package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenRepositoryPg struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewTokenRepositoryPg(db *gorm.DB, log *logrus.Logger) repository.TokenRepository {
	return &TokenRepositoryPg{
		db:  db,
		log: log,
	}
}

func (t *TokenRepositoryPg) CreateOrUpdate(ctx context.Context, token *entity.RefreshToken) (*entity.RefreshToken, error) {
	if token == nil {
		return nil, errorx.ErrTokenEmpty
	}

	err := t.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "device_info"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_revoked", "role", "token_hash", "expires_at", "created_at"}),
		}).
		Create(token).Error

	if err != nil {
		t.log.WithError(err).Error("failed to upsert refresh token")
		return nil, err
	}
	return token, nil
}

func (t *TokenRepositoryPg) FindByAccessToken(ctx context.Context, accessToken string) (entity.RefreshToken, error) {
	var token entity.RefreshToken
	err := t.db.WithContext(ctx).First(&token, "token_hash = ?", accessToken).Error
	if err != nil {
		return token, err
	}

	return token, nil
}

func (t *TokenRepositoryPg) DeleteByUserID(ctx context.Context, id int) error {
	err := t.db.WithContext(ctx).Where("user_id = ?", id).Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return err
	}
	return nil
}
