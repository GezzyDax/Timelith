package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupSetupRouter creates a minimal router for initial setup
// This router only serves setup endpoints and doesn't require database
func SetupSetupRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Timelith Setup Wizard",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	handler := NewSetupHandler()

	// Setup API routes
	api := app.Group("/api")
	api.Get("/setup/status", handler.CheckSetupStatus)
	api.Post("/setup", handler.PerformSetup)

	// Health check (minimal)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":       "setup_mode",
			"setup_needed": true,
			"message":      "Application is in setup mode",
		})
	})

	return app
}
