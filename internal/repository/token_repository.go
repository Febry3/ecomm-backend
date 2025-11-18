package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type TokenRepository interface {
	CreateOrUpdate(ctx context.Context, token *entity.RefreshToken) (*entity.RefreshToken, error)
	FindByAccessToken(ctx context.Context, accessToken string) (entity.RefreshToken, error)
	DeleteByUserID(ctx context.Context, id int) error
	DeleteByAccessToken(ctx context.Context, accessToken string) error
}
