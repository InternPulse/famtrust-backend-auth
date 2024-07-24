package interfaces

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Models interface {
	Users() UserModels
	Roles() UserRoles
	Permissions() UserPermissions
	VerCodes() VerCodeModels
}

type UserModels interface {
	CreateUser(user *User) error
	GetUserByID(userID uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserProfileByID(userID uuid.UUID) (*UserProfile, error)
	UpdateUser(user *User) error
	DeleteUserByID(userID uuid.UUID) error
	PasswordMatches(passswordHash string, plainText string) (bool, error)
	GetUserByBVN(bvn int) (*User, error)
	GetUserByNIN(nin int) (*User, error)
	SetIsVerified(userID uuid.UUID, value bool) error
}

type UserRoles interface {
	GetAllRoles() ([]Role, error)
	GetRoleByID(roleID string) (*Role, error)
	CreateRole(role *Role) error
	UpdateRoleByID(role *Role) error
	DeleteRoleByID(roleID string) error
}

type UserPermissions interface {
	GetAllPermissions() ([]Permission, error)
	GetPermission(perm *Permission) (*Permission, error)
	CreatePermission(perm *Permission) error
	UpdatePermission(perm *Permission) error
	DeletePermission(perm Permission) error
}

type VerCodeModels interface {
	CreateVerificationCode(verCode *VerCode) error
	GetEmailCodeByID(codeID uuid.UUID) (*VerCode, error)
	Get2FACodeByID(codeID uuid.UUID) (*VerCode, error)
	DeleteEmailCodeByID(codeID uuid.UUID) error
	Delete2FACodeByID(codeID uuid.UUID) error
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
	Email        string      `gorm:"not null;unique"`
	PasswordHash string      `gorm:"not null"`
	RoleID       string      `gorm:"not null"`
	Has2FA       bool        `gorm:"column:has_2fa;not null"`
	IsVerified   bool        `gorm:"not null"`
	IsFreezed    bool        `gorm:"not null"`
	DefaultGroup uuid.UUID   `gorm:"not null"`
	LastLogin    time.Time   `gorm:"not null"`
	Role         Role        `gorm:"foreignKey:RoleID;references:ID"`
	UserProfile  UserProfile `gorm:"foreignKey:UserID;references:ID;constraints:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserProfile struct {
	UUIDModel
	UserID              uuid.UUID
	FirstName           string `gorm:"not null"`
	LastName            string `gorm:"not null"`
	Bio                 string `gorm:"not null"`
	NIN                 uint   `gorm:"unique"`
	BVN                 uint   `gorm:"unique"`
	Profile_picture_url string `gorm:"not null"`
}

type Role struct {
	ID          string `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Permissions []Permission   `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Permission struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Roles     []Role         `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// VerCode Type can either be 'email' or '2fa'
// TODO: Implement enums
type VerCode struct {
	UUIDModel
	UserID uuid.UUID `gorm:"not null"`
	Type   string    `gorm:"not null"`
}
