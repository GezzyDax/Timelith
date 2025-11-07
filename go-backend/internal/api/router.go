package api

import (
	"github.com/GezzyDax/timelith/go-backend/internal/config"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRouter(cfg *config.Config, db *database.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Timelith API v1.0",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-API-Key",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Store config in locals for middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("jwt_secret", cfg.JWTSecret)
		return c.Next()
	})

	handler := NewHandler(db)

	// Public routes
	api := app.Group("/api")

	api.Get("/health", handler.HealthCheck)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	// Protected routes - require authentication
	protected := api.Group("/", AuthMiddleware(cfg.JWTSecret))

	// Accounts
	accounts := protected.Group("/accounts")
	accounts.Get("/", handler.ListAccounts)
	accounts.Post("/", handler.CreateAccount)
	accounts.Get("/:id", handler.GetAccount)
	accounts.Delete("/:id", handler.DeleteAccount)

	// Templates
	templates := protected.Group("/templates")
	templates.Get("/", handler.ListTemplates)
	templates.Post("/", handler.CreateTemplate)
	templates.Get("/:id", handler.GetTemplate)
	templates.Put("/:id", handler.UpdateTemplate)
	templates.Delete("/:id", handler.DeleteTemplate)

	// Channels
	channels := protected.Group("/channels")
	channels.Get("/", handler.ListChannels)
	channels.Post("/", handler.CreateChannel)
	channels.Get("/:id", handler.GetChannel)
	channels.Put("/:id", handler.UpdateChannel)
	channels.Delete("/:id", handler.DeleteChannel)

	// Schedules
	schedules := protected.Group("/schedules")
	schedules.Get("/", handler.ListSchedules)
	schedules.Post("/", handler.CreateSchedule)
	schedules.Get("/:id", handler.GetSchedule)
	schedules.Put("/:id", handler.UpdateSchedule)
	schedules.Patch("/:id/status", handler.UpdateScheduleStatus)
	schedules.Delete("/:id", handler.DeleteSchedule)
	schedules.Get("/:id/logs", handler.GetScheduleLogs)

	// Logs
	logs := protected.Group("/logs")
	logs.Get("/", handler.GetAllLogs)

	return app
}
