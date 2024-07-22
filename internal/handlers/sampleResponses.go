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
			ID                  uuid.UUID `example:"a5c9f82e-6b7a-4a53-a81c-82b1e2f453a6"`
			UserID              uuid.UUID `example:"d38f91b2-dc3b-4f9d-aeb4-7b95c91e9d08"`
			FirstName           string    `example:"Famtrust"`
			LastName            string    `example:"Guru"`
			Bio                 string    `example:"The best FamTrust user of all time."`
			NIN                 uint      `example:"35473745433"`
			BVN                 uint      `example:"35473783473"`
			Profile_picture_url string    `example:"https://image.famtrust.biz/dkkjieikdjfoej.jpg"`
			CreatedAt           time.Time `example:"2024-07-22T14:30:00Z"`
			UpdatedAt           time.Time `example:"2024-07-22T14:30:00Z"`
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
	ID         string `example:"d38f91b2-dc3b-4f9d-aeb4-7b95c91e9d08"`
	IsFreezed  bool   `example:"true"`
	IsVerified bool   `example:"true"`
	LastLogin  string `example:"2024-07-22T14:30:00Z"`
	Role       validateSampleResponseRole
}
