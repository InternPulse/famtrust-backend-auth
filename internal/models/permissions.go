package models

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"gorm.io/gorm"
)

type UserPermissions struct {
	DB *gorm.DB
}

func (p *UserPermissions) GetAllPermissions() ([]interfaces.Permission, error) {
	var perms []interfaces.Permission
	if err := p.DB.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (p *UserPermissions) GetPermission(perm *interfaces.Permission) (*interfaces.Permission, error) {
	var permission interfaces.Permission
	if err := p.DB.First(&permission, perm).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (p *UserPermissions) CreatePermission(perm *interfaces.Permission) error {
	if err := p.DB.Create(&perm).Error; err != nil {
		return err
	}
	return nil
}

func (p *UserPermissions) UpdatePermission(perm *interfaces.Permission) error {
	if err := p.DB.Model(&interfaces.Permission{}).Updates(&perm).Error; err != nil {
		return err
	}
	return nil
}

func (p *UserPermissions) DeletePermission(perm interfaces.Permission) error {
	if err := p.DB.Delete(&interfaces.Permission{}, perm).Error; err != nil {
		return err
	}
	return nil
}
