package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]entity.Category, error)
}
