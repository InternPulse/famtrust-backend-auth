package interfaces

import (
	"time"

	"github.com/google/uuid"
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

// Create uuid model.
type UUIDModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Create uuid BeforeCreate hook.
func (u *UUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type User struct {
	UUIDModel
	Email         string      `gorm:"not null;unique"`
	Password_hash string      `gorm:"not null"`
	RoleID        string      `gorm:"not null"`
	Is_verified   bool        `gorm:"not null"`
	Is_freezed    bool        `gorm:"not null"`
	Last_login    time.Time   `gorm:"not null"`
	Role          Role        `gorm:"foreignKey:RoleID;references:ID"`
	UserProfile   UserProfile `gorm:"foreignKey:UserID;references:ID;constraints:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserProfile struct {
	UUIDModel
	UserID              uuid.UUID
	First_name          string `gorm:"not null"`
	Last_name           string `gorm:"not null"`
	Bio                 string `gorm:"not null"`
	NIN                 uint   `gorm:"not null;unique"`
	BVN                 uint   `gorm:"not null;unique"`
	Profile_picture_url string `gorm:"not null"`
}

type Role struct {
	ID         string `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Permission []Permission   `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Permission struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Role      []Role         `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
