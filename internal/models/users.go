package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"gorm.io/gorm"
)

type UserModels struct {
	DB *gorm.DB
}

func (u *UserModels) CreateUser(user *interfaces.User) error {
	if err := u.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserModels) GetUserByID(userID uint) (*interfaces.User, error) {
	var user interfaces.User
	if err := u.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModels) GetUserByEmail(email string) (*interfaces.User, error) {
	var user interfaces.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModels) UpdateUser(user *interfaces.User) error {
	if err := u.DB.Model(&interfaces.User{}).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserModels) DeleteUserByID(userID uint) error {
	if err := u.DB.Delete(&interfaces.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}
