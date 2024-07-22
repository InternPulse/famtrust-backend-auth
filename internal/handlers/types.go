package handlers

import (
	"time"

	"github.com/google/uuid"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	StatusCode uint
	Status     string
	Message    string
	Token      string                   `json:"token,omitempty"`
	Data       map[string]loginUserData `json:"data,omitempty"`
}

type loginUserData struct {
	Id         uuid.UUID
	Email      string
	Role       role
	Has2FA     bool
	IsVerified bool
	IsFreezed  bool
	LastLogin  time.Time
}

type role struct {
	Id          string   `json:"id" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}
