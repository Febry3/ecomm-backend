package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type CategoryRepositoryPg struct {
	db *gorm.DB
}

func NewCategoryRepositoryPg(db *gorm.DB) repository.CategoryRepository {
	return &CategoryRepositoryPg{db: db}
}

func (c *CategoryRepositoryPg) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	err := c.db.WithContext(ctx).Select("id, name, slug").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
