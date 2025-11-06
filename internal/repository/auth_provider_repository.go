package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type AuthProviderRepository interface {
	Create(ctx context.Context, authProvider *entity.AuthProvider) error
	FindByUserID(ctx context.Context, userId int64) (entity.AuthProvider, error)
	FindByProviderId(ctx context.Context, providerId string, provider string) (entity.AuthProvider, error)
}
