package repository

import (
	"context"
	"github.com/febry3/gamingin/internal/entity"
)

type TokenRepository interface {
	CreateOrUpdate(ctx context.Context, token *entity.RefreshToken) (*entity.RefreshToken, error)
	FindByUserID(ctx context.Context, id int) (entity.RefreshToken, error)
	DeleteByUserID(ctx context.Context, id int) error
}
