package jwtmod

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var JwtKey []byte

type JwtClaim struct {
	ID uuid.UUID `json:"id" binding:"required"`
	jwt.StandardClaims
}

func GenerateJWT(userID uuid.UUID) (string, error) {
	// create expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// user claims payload
	claims := JwtClaim{
		ID: userID,
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
