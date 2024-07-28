package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationCodes struct {
	DB *gorm.DB
}

func (v *VerificationCodes) GetEmailCodeByID(codeID uuid.UUID) (*interfaces.VerCode, error) {
	var profile interfaces.VerCode
	if err := v.DB.Order("created_at DESC").Where("type = ?", "email").Where("id = ?", codeID).Limit(1).Find(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (v *VerificationCodes) GetResetCodeByID(codeID uuid.UUID) (*interfaces.VerCode, error) {
	var profile interfaces.VerCode
	if err := v.DB.Order("created_at DESC").Where("type = ?", "password").Where("id = ?", codeID).Limit(1).Find(&profile).Error; err != nil {
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
	if err := v.DB.Where("type = ?", "email").Where("user_id = ?", UserID).Delete(&interfaces.VerCode{}).Error; err != nil {
		return err
	}
	return nil
}

func (v *VerificationCodes) DeleteResetCodeByUserID(UserID uuid.UUID) error {
	if err := v.DB.Where("type = ?", "password").Where("user_id = ?", UserID).Delete(&interfaces.VerCode{}).Error; err != nil {
		return err
	}
	return nil
}

func (v *VerificationCodes) Delete2FACodeByUserID(UserID uuid.UUID) error {
	if err := v.DB.Where("type = ?", "2fa").Where("user_id = ?", UserID).Delete(&interfaces.VerCode{}).Error; err != nil {
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
