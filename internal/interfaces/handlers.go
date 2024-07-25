package interfaces

import "github.com/gin-gonic/gin"

type Handlers interface {
	Users() UserHandlers
	AuthMiddleware() gin.HandlerFunc
	Verifications() VerificationHandlers
}

type UserHandlers interface {
	Signup(c *gin.Context)
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
	Validate(c *gin.Context)
	GetPermissions(permissions []Permission) []string
	// ResetPassword(c *gin.Context)

	// User Profiles
	GetUserProfileByID(c *gin.Context)
	CreateUserProfile(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	GetProfilePicture(c *gin.Context)

	// Get Users By...
	GetUsersByDefaultGroup(c *gin.Context)
	GetUserByDefaultGroup(c *gin.Context)
}

type VerificationHandlers interface {
	VerifyEmail(c *gin.Context)
	VerifyEmailToken(c *gin.Context)
	VerifyNIN(c *gin.Context)
	VerifyBVN(c *gin.Context)
}
