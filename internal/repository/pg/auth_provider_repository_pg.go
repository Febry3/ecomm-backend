package pg

import (
	"context"
	"errors"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"gorm.io/gorm"
)

type AuthProviderPg struct {
	db *gorm.DB
}

func NewAuthProvider(db *gorm.DB) *AuthProviderPg {
	return &AuthProviderPg{
		db: db,
	}
}

func (a *AuthProviderPg) Create(ctx context.Context, authProvider *entity.AuthProvider) error {
	err := a.db.WithContext(ctx).Create(authProvider).Error
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthProviderPg) FindByUserID(ctx context.Context, userId int64) (entity.AuthProvider, error) {
	authProvider := entity.AuthProvider{}
	err := a.db.WithContext(ctx).First(&authProvider, "user_id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AuthProvider{}, errorx.ErrTokenInvalid
		}
		return entity.AuthProvider{}, err
	}
	return authProvider, nil
}

func (a *AuthProviderPg) FindByProviderId(ctx context.Context, authProviderId string, provider string) (entity.AuthProvider, error) {
	authProvider := entity.AuthProvider{}
	err := a.db.WithContext(ctx).First(&authProvider, "provider_id = ? AND provider = ?", authProviderId, provider).Error
	if err != nil {
		return entity.AuthProvider{}, err
	}
	return authProvider, nil
}
