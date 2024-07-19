package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"gorm.io/gorm"
)

type Models struct {
	user interfaces.UserModels
}

func (m *Models) User() interfaces.UserModels {
	return m.user
}

func NewModel(DB *gorm.DB) interfaces.Models {
	return &Models{
		user: &UserModels{DB: DB},
	}
}
