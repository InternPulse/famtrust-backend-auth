package interfaces

import (
	"time"

	"gorm.io/gorm"
)

type Models interface {
	User() UserModels
	Role() UserRoles
}

type UserModels interface {
	CreateUser(user *User) error
	GetUserByID(userID uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUserByID(userID uint) error
}

type UserRoles interface {
	GetAllRoles() ([]Role, error)
	GetRoleByID(roleID uint) (*Role, error)
	CreateRole(role *Role) error
	DeleteRoleByID(roleID uint) error
}

type User struct {
	gorm.Model
	First_name          string    `gorm:"not null"`
	Last_name           string    `gorm:"not null"`
	Email               string    `gorm:"not null;unique"`
	Password_hash       string    `gorm:"not null"`
	Bio                 string    `gorm:"not null"`
	NIN                 uint      `gorm:"not null;unique"`
	BVN                 uint      `gorm:"not null;unique"`
	Profile_picture_url string    `gorm:"not null"`
	RoleID              string    `gorm:"not null"`
	Is_freezed          bool      `gorm:"not null"`
	Last_login          time.Time `gorm:"not null"`
	Role                Role      `gorm:"foreignKey:RoleID;references:Name"`
}

type Role struct {
	gorm.Model
	Name                 string `gorm:"not null;unique"`
	CanOperateFamilyAcct bool   `gorm:"not null"`
	CanCreateUser        bool   `gorm:"not null"`
	CanCreateSubAcc      bool   `gorm:"not null"`
	CanDeleteSubAcc      bool   `gorm:"not null"`
	CanEditFamilyAcc     bool   `gorm:"not null"`
	CanSendtoSubAcc      bool   `gorm:"not null"`
	CanSendtoBank        bool   `gorm:"not null"`
	CanFreezeSubAcc      bool   `gorm:"not null"`
}
