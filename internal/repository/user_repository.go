package repository

import (
	"context"
	"github.com/febry3/gamingin/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int64) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, bool, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
}
