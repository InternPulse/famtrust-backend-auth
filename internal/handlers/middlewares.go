package handlers

import (
	"net/http"
	"strings"

	"github.com/InternPulse/famtrust-backend-auth/internal/jwtmod"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (h *Handlers) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, loginResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "error",
				Message:    "Invalid Token",
			})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix to get the actual token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := jwtmod.JwtClaim{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtmod.JwtKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, loginResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "error",
				Message:    "Invalid Token",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, loginResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "error",
				Message:    "Invalid Token",
			})
			c.Abort()
			return
		}

		c.Set("token", token.Raw)
		c.Set("UserID", claims.ID)
		c.Next()
	}
}
