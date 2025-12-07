package user

import (
	"belajarGo2/service/user"
	"context"

	"gorm.io/gorm"
)

type (
	GormRepository struct {
		*gorm.DB
	}
)

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db.Table("bg_users"),
	}
}

func (r *GormRepository) Create(user user.User) (err error) {
	return r.DB.WithContext(context.Background()).Create(&user).Error
}

func (r *GormRepository) GetByEmail(email string) (user user.User, err error) {
	r.DB.WithContext(context.Background()).First(&user, "email = ?", email)
	return
}

func (r *GormRepository) UpdateEmailVerification(user user.User) (err error) {
	err = r.DB.WithContext(context.Background()).Updates(&user).Error
	return
}
