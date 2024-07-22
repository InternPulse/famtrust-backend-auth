package handlers

import (
	"time"

	"github.com/google/uuid"
)

type loginSampleResponse200 struct {
	StatusCode uint   `example:"200"`
	Status     string `example:"success"`
	Message    string `example:"User Logged in successfully"`
	Token      string `example:"b6d4a7e1d2d841a1afe874a2a5c15d8b"`
}

type loginSampleResponseError401 struct {
	StatusCode uint   `example:"401"`
	Status     string `example:"error"`
	Message    string `example:"Invalid Credentials"`
}

type loginSampleResponseError500 struct {
	StatusCode uint   `example:"500"`
	Status     string `example:"error"`
	Message    string `example:"An error occured"`
}

type validateSampleResponse200 struct {
	StatusCode uint   `example:"200"`
	Status     string `example:"success"`
	Message    string `example:"Request successful"`
	Token      string `example:"b6d4a7e1d2d841a1afe874a2a5c15d8b"`
	Data       struct {
		User validateSampleResponse200User
	}
}

type profileSampleResponse200 struct {
	StatusCode uint   `example:"200"`
	Status     string `example:"success"`
	Message    string `example:"Request successful"`
	Token      string `example:"b6d4a7e1d2d841a1afe874a2a5c15d8b"`
	Data       struct {
		Profile struct {
			ID                  uuid.UUID
			UserID              uuid.UUID
			FirstName           string
			LastName            string
			Bio                 string
			NIN                 uint
			BVN                 uint
			Profile_picture_url string
			CreatedAt           time.Time
			UpdatedAt           time.Time
		}
	}
}

type validateSampleResponseRole struct {
	Id          string   `example:"admin"`
	Permissions []string `example:"canTransact, canWithdraw"`
}

type validateSampleResponse200User struct {
	Email      string `example:"user@example.com"`
	Has2FA     bool   `example:"true"`
	ID         string `example:"user-id"`
	IsFreezed  bool   `example:"true"`
	IsVerified bool   `example:"true"`
	LastLogin  string `example:"2024-07-22T14:30:00Z"`
	Role       validateSampleResponseRole
}
