package pg

import (
	"context"
	"github.com/sirupsen/logrus"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type UserRepositoryPg struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewUserRepositoryPg(db *gorm.DB, log *logrus.Logger) repository.UserRepository {
	return &UserRepositoryPg{db: db, log: log}
}

func (u UserRepositoryPg) Create(ctx context.Context, user *entity.User) error {
	result := u.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		u.log.Errorf("[UserRepositoryPg] Create User Error: %v]", result.Error.Error())
		return result.Error
	}
	return nil
}

func (u UserRepositoryPg) FindByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	result := u.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		u.log.Errorf("[UserRepositoryPg] Find User By ID Error: %v]", result.Error.Error())
		return entity.User{}, result.Error
	}
	return user, nil
}

func (u UserRepositoryPg) FindByEmail(ctx context.Context, email string) (entity.User, bool, error) {
	var user entity.User
	result := u.db.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		u.log.Errorf("[UserRepositoryPg] Find User By Email Error: %v]", result.Error.Error())
		return entity.User{}, false, result.Error
	}

	return user, true, nil
}

func (u UserRepositoryPg) Update(ctx context.Context, user entity.User) (entity.User, error) {
	result := u.db.WithContext(ctx)
	if result.Error != nil {
		u.log.Errorf("[UserRepositoryPg] Update User Error: %v]", result.Error.Error())
		return entity.User{}, result.Error
	}
	return user, nil
}
