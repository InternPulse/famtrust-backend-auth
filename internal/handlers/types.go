package handlers

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/google/uuid"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	StatusCode    uint                     `json:"status_code" binding:"required"`
	StatusMessage string                   `json:"status_message" binding:"required"`
	Token         string                   `json:"token,omitempty"`
	Data          map[string]loginUserData `json:"data,omitempty"`
}

type loginUserData struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Email string    `json:"email" binding:"required"`
	Role  Role      `json:"role" binding:"required"`
}

type Role struct {
	ID          string                  `json:"id" binding:"required"`
	Permissions []interfaces.Permission `json:"permissions" binding:"required"`
}
