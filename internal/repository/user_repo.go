package repository

import (
	"test-backend-1-d1ma11/internal/entity"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetByEmail(user *entity.User, email string) error {
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) Create(user *entity.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
