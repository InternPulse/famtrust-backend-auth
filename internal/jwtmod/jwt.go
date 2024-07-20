package jwtmod

import (
	"time"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var JwtKey []byte

type JwtClaim struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Email string    `json:"email" binding:"required"`
	Role  Role      `json:"role" binding:"required"`
	jwt.StandardClaims
}

type Role struct {
	ID          string                  `json:"id" binding:"required"`
	Permissions []interfaces.Permission `json:"permissions" binding:"required"`
}

func GenerateJWT(user *interfaces.User) (string, error) {
	// create expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// user claims payload
	Role := Role{
		ID:          user.Role.ID,
		Permissions: user.Role.Permissions,
	}
	claims := JwtClaim{
		ID:    user.ID,
		Email: user.Email,
		Role:  Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
