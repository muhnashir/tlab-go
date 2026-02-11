package http

import (
	"wallet-api/internal/delivery/http/handler"
	"wallet-api/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "wallet-api/docs" // Import generated docs

	"github.com/gofiber/swagger"
)

// SetupRouter initializes the Fiber app and registers all routes
func SetupRouter(handlers *handler.AllHandlers) *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())  // Request logging
	app.Use(recover.New()) // Panic recovery
	app.Use(cors.New())    // CORS support

	// Swagger Route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API Group
	api := app.Group("/api")

	// Auth Routes (Public)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", handlers.AuthHandler.Register)
	authGroup.Post("/login", handlers.AuthHandler.Login)

	// Protected Routes Group
	protected := api.Group("/", middleware.JWTProtected())

	// User Routes
	userGroup := protected.Group("/users")
	userGroup.Get("/profile", handlers.AuthHandler.GetProfile)

	// Wallet Routes
	walletGroup := protected.Group("/wallets")
	// Assuming balance endpoint logic exists or will be added to handler
	walletGroup.Post("/topup", handlers.WalletHandler.TopUp)
	walletGroup.Get("/balance", handlers.WalletHandler.GetBalance)

	// Transaction Routes
	transactionGroup := protected.Group("/transactions")
	transactionGroup.Post("/transfer", handlers.WalletHandler.Transfer)
	transactionGroup.Get("/history", handlers.WalletHandler.GetHistory) // Can act as history

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "system is healthy",
		})
	})

	return app
}
