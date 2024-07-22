package models

import (
	"errors"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (u *UserModels) GetUserByID(userID uuid.UUID) (*interfaces.User, error) {
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

func (u *UserModels) DeleteUserByID(userID uuid.UUID) error {
	if err := u.DB.Delete(&interfaces.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *UserModels) PasswordMatches(user *interfaces.User, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *UserModels) GetUserProfileByID(userID uuid.UUID) (*interfaces.UserProfile, error) {
	var profile interfaces.UserProfile
	if err := u.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}
