package interfaces

import (
	"time"

	"gorm.io/gorm"
)

type Models interface {
	User() UserModels
}

type UserModels interface {
}

type User struct {
	gorm.Model
	First_name          string
	Last_name           string
	Email               string
	Password_hash       string
	Bio                 string
	Profile_picture_url string
	Role_id             uint
	Is_freezed          bool
	Last_login          time.Time
}

type Role struct {
	Name              string
	CanTopUpFamilyAcc bool
	CanCreateSubAcc   bool
	CanDeleteSubAcc   bool
	CanEditFamilyAcc  bool
	CanSendtoSubAcc   bool
	CanSendtoBank     bool
	gorm.Model
}
