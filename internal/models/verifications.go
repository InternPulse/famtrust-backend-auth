package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationCodes struct {
	DB *gorm.DB
}

func (v *VerificationCodes) GetCodeByID(codeID uuid.UUID) (*interfaces.VerCode, error) {
	var profile interfaces.VerCode
	if err := v.DB.First(&profile, codeID).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (v *VerificationCodes) DeleteCodeByID(codeID uuid.UUID) error {
	if err := v.DB.Delete(&interfaces.VerCode{}, codeID).Error; err != nil {
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
