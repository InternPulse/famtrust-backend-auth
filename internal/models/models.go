package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"gorm.io/gorm"
)

type Models struct {
	users       interfaces.UserModels
	roles       interfaces.UserRoles
	permissions interfaces.UserPermissions
	verCodes    interfaces.VerCodeModels
}

func (m *Models) Users() interfaces.UserModels {
	return m.users
}

func (m *Models) Roles() interfaces.UserRoles {
	return m.roles
}

func (m *Models) Permissions() interfaces.UserPermissions {
	return m.permissions
}

func (m *Models) VerCodes() interfaces.VerCodeModels {
	return m.verCodes
}

func NewModel(DB *gorm.DB) interfaces.Models {
	return &Models{
		users:       &UserModels{DB: DB},
		roles:       &UserRoles{DB: DB},
		permissions: &UserPermissions{DB: DB},
		verCodes:    &VerificationCodes{DB: DB},
	}
}
