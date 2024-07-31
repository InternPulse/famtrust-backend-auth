package handlers

import (
	"time"

	"github.com/google/uuid"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type validateResponse struct {
	StatusCode uint          `json:"statusCode"`
	Status     string        `json:"status"`
	Message    string        `json:"message"`
	Token      string        `json:"token,omitempty"`
	User       cleanUserData `json:"user,omitempty"`
}

type loginResponse struct {
	StatusCode uint   `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Token      string `json:"token,omitempty"`
}

type cleanUserData struct {
	Id           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Role         role      `json:"role"`
	DefaultGroup uuid.UUID `json:"defaultGroup"`
	Has2FA       bool      `json:"has2FA"`
	IsVerified   bool      `json:"isVerified"`
	IsFrozen     bool      `json:"isFrozen"`
	LastLogin    time.Time `json:"lastLogin"`
}

type role struct {
	Id          string   `json:"id" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}
