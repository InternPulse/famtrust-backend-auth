package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"gorm.io/gorm"
)

type UserRoles struct {
	DB *gorm.DB
}

func (r *UserRoles) GetAllRoles() ([]interfaces.Role, error) {
	var roles []interfaces.Role
	if err := r.DB.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRoles) GetRoleByID(roleID uint) (*interfaces.Role, error) {
	var role interfaces.Role
	if err := r.DB.First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *UserRoles) CreateRole(role *interfaces.Role) error {
	if err := r.DB.Create(&role).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRoles) DeleteRoleByID(roleID uint) error {
	if err := r.DB.Delete(&interfaces.Role{}, roleID).Error; err != nil {
		return err
	}
	return nil
}
