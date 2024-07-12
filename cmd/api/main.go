package main

import (
	"log"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
)

const webPort = ":8001"

type Config struct {
	Handlers interfaces.Handlers
}

// @title			FamTrust API Backend - Auth
// @version			1.0
// @description		This is the Authentication and Authrization micro-service for the FamTrust Web App.
// @host			localhost:8001
// @BasePath		/api/v1/
func main() {

	app := Config{}

	// Run api
	err := app.routes().Run(webPort)
	if err != nil {
		log.Fatalf("Failed to start web api; %v", err)
	}
}
