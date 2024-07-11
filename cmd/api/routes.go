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

	// Use Logger
	mux.Use(gin.Logger())
	// Use Recovery
	mux.Use(gin.Recovery())
	// Use Cors
	mux.Use(cors.Default())

	// Make api base
	api := mux.Group("/api")

	// Swagger docs
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return mux

}
