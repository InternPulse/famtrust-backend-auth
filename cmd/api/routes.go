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
	// signup
	// v1.GET("/signup", app.Handlers.Users().Signup)
	// login
	v1.POST("/login", app.Handlers.Users().Login)

	// // Verification Routes
	v1.GET("/verify-nin", app.Handlers.Verifications().VerifyNIN)
	v1.GET("/verify-bvn", app.Handlers.Verifications().VerifyBVN)
	v1.GET("/verify-email/verify", app.Handlers.Verifications().VerifyEmailToken)
	// v1.GET("/reset-password", app.Handlers.Users().ResetPassword)

	// Protected Routes
	// validate token
	v1.Use(app.Handlers.AuthMiddleware()).GET("/validate", app.Handlers.Users().Validate)
	// verify email
	v1.Use(app.Handlers.AuthMiddleware()).GET("/verify-email", app.Handlers.Verifications().VerifyEmail)
	// v1.GET("/verify-2fa", app.Handlers.Users().Verify2FA)

	// // User & UserProfile Routes
	v1.Use(app.Handlers.AuthMiddleware()).GET("/profile", app.Handlers.Users().GetUserProfileByID)
	// v1.POST("/profile/create", app.Handlers.Users().CreateUser)
	// v1.PUT("/profile/update", app.Handlers.Users().UpdateUserByID)
	// v1.DELETE("/profile/delete", app.Handlers.Users().DeleteUserByID)

	return mux

}
