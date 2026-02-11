package main

import (
	"log"
	"os"

	"wallet-api/internal/delivery/http"
	"wallet-api/internal/delivery/http/handler"
	"wallet-api/internal/pkg/database"

	"github.com/joho/godotenv"
)

// @title Wallet API
// @version 1.0
// @description RESTful API for E-Wallet Application
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (using system environment variables if available)")
	}

	// 2. Initialize Database
	db := database.InitMySQL()
	defer db.Close()

	// 3. Initialize Handlers (Dependency Injection)
	handlers := handler.InitHandlers(db)

	// 4. Setup Router with Handlers
	app := http.SetupRouter(handlers)

	// 5. Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Default port if not specified
	}

	log.Printf("Server is starting on port %s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
