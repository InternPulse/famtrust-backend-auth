package interfaces

import "github.com/gin-gonic/gin"

type Handlers interface {
	Users() UserHandlers
	AuthMiddleware() gin.HandlerFunc
	Verifications() VerificationHandlers
}

type UserHandlers interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Validate(c *gin.Context)
	GetUserProfileByID(c *gin.Context)
	// ResetPassword(c *gin.Context)
}

type VerificationHandlers interface {
	VerifyEmail(c *gin.Context)
	// Verify2FA(c *gin.Context)
	VerifyNIN(c *gin.Context)
	VerifyBVN(c *gin.Context)
}
