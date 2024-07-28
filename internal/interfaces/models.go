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
	CreateUserProfile(profile *UserProfile) error
	UpdateUserProfile(profile *UserProfile) error
	GetUserProfileByID(userID uuid.UUID) (*UserProfile, error)
	UpdateUser(user *User) error
	DeleteUserByID(userID uuid.UUID) error
	PasswordMatches(passswordHash string, plainText string) (bool, error)
	GetUserByBVN(bvn int) (*User, error)
	GetUserByNIN(nin int) (*User, error)
	SetIsVerified(userID uuid.UUID, value bool) error
	GetUsersByDefaultGroup(groupID uuid.UUID) (*[]User, error)
	GetUserByDefaultGroup(userID uuid.UUID, groupID uuid.UUID) (*User, error)
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
	Get2FACodeByUserID(codeID uuid.UUID) (*VerCode, error)
	GetResetCodeByID(codeID uuid.UUID) (*VerCode, error)
	DeleteEmailCodeByUserID(userID uuid.UUID) error
	DeleteResetCodeByUserID(userID uuid.UUID) error
	Delete2FACodeByUserID(userID uuid.UUID) error
}

// Create uuid model.
type UUIDModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
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
	Email        string      `json:"email" gorm:"not null;unique"`
	PasswordHash string      `json:"_" gorm:"not null"`
	RoleID       string      `json:"roleId" gorm:"not null"`
	DefaultGroup uuid.UUID   `json:"defaultGroup"`
	Has2FA       bool        `json:"has2FA" gorm:"column:has_2fa;not null"`
	IsVerified   bool        `json:"isVerified" gorm:"not null"`
	IsFrozen     bool        `json:"isFreezed" gorm:"not null"`
	LastLogin    time.Time   `json:"lastLogin" gorm:"not null"`
	Role         Role        `json:"role" gorm:"foreignKey:RoleID;references:ID"`
	UserProfile  UserProfile `json:"userProfile" gorm:"foreignKey:UserID;references:ID;constraints:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserProfile struct {
	UUIDModel
	UserID            uuid.UUID `json:"userId"`
	FirstName         string    `json:"firstName" gorm:"not null"`
	LastName          string    `json:"lastName" gorm:"not null"`
	Bio               string    `json:"bio" gorm:"not null"`
	NIN               uint      `json:"nin" gorm:"not null"`
	BVN               uint      `json:"bvn" gorm:"not null"`
	ProfilePictureUrl string    `json:"profilePictureUrl" gorm:"not null"`
}

type Role struct {
	ID          string         `json:"Id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Permissions []Permission   `json:"permissions" gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Permission struct {
	ID        string `json:"Id" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Roles     []Role         `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// VerCode Type can either be 'email' or '2fa' or 'password'
// TODO: Implement enums
type VerCode struct {
	UUIDModel
	UserID uuid.UUID `json:"userId" gorm:"not null"`
	Type   string    `json:"type" gorm:"not null"`
}
