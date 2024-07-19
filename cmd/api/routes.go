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
	v1.GET("/signup", app.Handlers.Users().Signup)
	v1.GET("/login", app.Handlers.Users().Login)
	v1.GET("/validate", app.Handlers.Users().Validate)

	return mux

}
