package main

import (
	"log"
	"os"

	"github.com/InternPulse/famtrust-backend-auth/internal/db"
	"github.com/InternPulse/famtrust-backend-auth/internal/handlers"
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/InternPulse/famtrust-backend-auth/internal/jwtmod"
	"github.com/InternPulse/famtrust-backend-auth/internal/mailer"
	"github.com/InternPulse/famtrust-backend-auth/internal/models"
	"github.com/joho/godotenv"

	_ "github.com/InternPulse/famtrust-backend-auth/docs"
)

const webPort = ":8001"

type Config struct {
	Handlers interfaces.Handlers
}

// @title						FamTrust API Backend - Auth
// @version					1.0
// @description				This is the Authentication and Authrization micro-service for the FamTrust Web App.
// @BasePath					/api/v1/
// @securityDefinitions.apiKey	BearerAuth
// @name						Authorization
// @in							header
func main() {
	// load env vars
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	// init jwt
	jwtmod.JwtKey = []byte(os.Getenv("JWTKEY"))

	// new postgres instance
	postgresDB := db.NewPostgresDB()

	// new model instance
	models := models.NewModel(postgresDB)

	// new mailer instance
	mailer := mailer.NewMailer()

	// new app instance
	app := Config{
		Handlers: handlers.NewHandler(models, mailer),
	}

	// Run app
	err := app.routes().Run(webPort)
	if err != nil {
		log.Fatalf("Failed to start web api; %v", err)
	}
}
