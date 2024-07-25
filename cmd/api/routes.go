package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *Config) routes() *gin.Engine {
	// Make router.
	mux := gin.New()

	// Use Middlewares
	mux.Use(gin.Logger())   //Logger
	mux.Use(gin.Recovery()) //Recovery
	mux.Use(cors.Default()) //Cors

	// Make api base
	api := mux.Group("/api")

	// Swagger docs
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Make API v1
	v1 := api.Group("/v1")

	// Auth Routes
	v1.POST("/signup", app.Handlers.Users().Signup)
	v1.POST("/login", app.Handlers.Users().Login)
	// v1.POST("/delete-user", app.Handlers.Users().DeleteUser)
	// v1.GET("/reset-password", app.Handlers.Users().ResetPassword)

	// Get User Profile Picture
	v1.GET("/images/profile-pic/:imageName", app.Handlers.Users().GetProfilePicture)

	// User Routes
	Users := v1.Group("/users").Use(app.Handlers.AuthMiddleware())
	Users.GET("/", app.Handlers.Users().GetUsersByDefaultGroup)
	Users.POST("/", app.Handlers.Users().CreateUser)
	Users.GET("/:userID", app.Handlers.Users().GetUserByDefaultGroup)

	// // Verification Routes
	v1.GET("/verify-nin", app.Handlers.Verifications().VerifyNIN)
	v1.GET("/verify-bvn", app.Handlers.Verifications().VerifyBVN)
	v1.GET("/verify-email/verify", app.Handlers.Verifications().VerifyEmailToken)

	// Protected Routes
	v1.GET("/validate", app.Handlers.AuthMiddleware(), app.Handlers.Users().Validate)
	v1.GET("/verify-email", app.Handlers.AuthMiddleware(), app.Handlers.Verifications().VerifyEmail)

	// UserProfile Routes [Protected]
	profile := v1.Group("/profile").Use(app.Handlers.AuthMiddleware())
	profile.GET("/", app.Handlers.Users().GetUserProfileByID)
	profile.POST("/create", app.Handlers.Users().CreateUserProfile)
	profile.PUT("/update", app.Handlers.Users().UpdateUserProfile)

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	mux.MaxMultipartMemory = 16 << 20 // 16 MiB

	return mux

}
