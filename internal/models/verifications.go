package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationCodes struct {
	DB *gorm.DB
}

func (v *VerificationCodes) GetEmailCodeByUserID(UserID uuid.UUID) (*interfaces.VerCode, error) {
	var profile interfaces.VerCode
	if err := v.DB.Order("created_at DESC").Where("type = ?", "email").Where("user_id = ?", UserID).Limit(1).Find(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (v *VerificationCodes) Get2FACodeByUserID(UserID uuid.UUID) (*interfaces.VerCode, error) {
	var profile interfaces.VerCode
	if err := v.DB.Order("created_at DESC").Where("type = ?", "2fa").Where("user_id = ?", UserID).Limit(1).Find(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (v *VerificationCodes) DeleteEmailCodeByUserID(UserID uuid.UUID) error {
	if err := v.DB.Delete(&interfaces.VerCode{}).Where("type = ?", "email").Where("user_id = ?", UserID).Error; err != nil {
		return err
	}
	return nil
}

func (v *VerificationCodes) Delete2FACodeByUserID(UserID uuid.UUID) error {
	if err := v.DB.Delete(&interfaces.VerCode{}).Where("type = ?", "2fa").Where("user_id = ?", UserID).Error; err != nil {
		return err
	}
	return nil
}

func (v *VerificationCodes) CreateVerificationCode(verCode *interfaces.VerCode) error {
	if err := v.DB.Create(&verCode).Error; err != nil {
		return err
	}
	return nil
}
