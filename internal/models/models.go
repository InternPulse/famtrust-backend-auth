package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"gorm.io/gorm"
)

type Models struct {
	user interfaces.UserModels
	role interfaces.UserRoles
}

func (m *Models) User() interfaces.UserModels {
	return m.user
}

func (m *Models) Role() interfaces.UserRoles {
	return m.role
}

func NewModel(DB *gorm.DB) interfaces.Models {
	return &Models{
		user: &UserModels{DB: DB},
		role: &UserRoles{DB: DB},
	}
}
