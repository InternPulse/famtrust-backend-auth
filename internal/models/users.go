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
	if err := u.DB.Preload("Role").Preload("Role.Permissions").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModels) GetUserByEmail(email string) (*interfaces.User, error) {
	var user interfaces.User
	if err := u.DB.Where("email = ?", email).Preload("Role").Preload("Role.Permissions").First(&user).Error; err != nil {
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

func (u *UserModels) PasswordMatches(passwordHash string, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainText))
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

func (u *UserModels) CreateUserProfile(profile *interfaces.UserProfile) error {
	if err := u.DB.Create(&profile).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserModels) UpdateUserProfile(profile *interfaces.UserProfile) error {
	if err := u.DB.Model(&interfaces.UserProfile{}).Where("user_id = ?", profile.UserID).Updates(&profile).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserModels) GetUserProfileByID(userID uuid.UUID) (*interfaces.UserProfile, error) {
	var profile interfaces.UserProfile
	if err := u.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (u *UserModels) GetUserByNIN(nin int) (*interfaces.User, error) {
	var user interfaces.User
	if err := u.DB.Where("nin = ?", nin).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModels) GetUserByBVN(bvn int) (*interfaces.User, error) {
	var user interfaces.User
	if err := u.DB.Where("nin = ?", bvn).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModels) SetIsVerified(userID uuid.UUID, value bool) error {
	// Update the `Active` field of a user with a specific ID
	if err := u.DB.Model(&interfaces.User{}).Where("id = ?", userID).Update("is_verified", value).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserModels) GetUsersByDefaultGroup(groupID uuid.UUID) (*[]interfaces.User, error) {
	var users []interfaces.User
	if err := u.DB.Preload("UserProfile").
		Joins("UserProfile").
		Where(`"UserProfile".default_group = ?`, groupID).
		Omit("password_hash").
		Find(&users).Error; err != nil {

		return nil, err
	}
	return &users, nil
}
