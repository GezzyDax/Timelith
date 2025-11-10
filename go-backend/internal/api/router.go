package api

import (
	"github.com/GezzyDax/timelith/go-backend/internal/config"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/settings"
	"github.com/GezzyDax/timelith/go-backend/internal/telegram"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRouter(cfg *config.Config, db *database.DB, settingsService *settings.Service, sessionManager *telegram.SessionManager) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Timelith API v1.0",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8080, http://127.0.0.1:3000, http://127.0.0.1:8080",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-API-Key",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Store settings service and config in locals for handlers
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("settings", settingsService)
		if cfg.JWTSecret != "" {
			c.Locals("jwt_secret", cfg.JWTSecret)
		}
		return c.Next()
	})

	handler := NewHandler(db, sessionManager)
	setupHandler := NewSetupHandler(db, settingsService)
	settingsHandler := NewSettingsHandler(settingsService)
	usersHandler := NewUsersHandler(db)

	// Public routes
	api := app.Group("/api")

	// Apply setup middleware - теперь проверяет статус динамически
	api.Use(SetupMiddleware(settingsService))

	api.Get("/health", handler.HealthCheck)

	// Setup routes
	api.Get("/setup/status", setupHandler.CheckSetupStatus)
	api.Post("/setup", setupHandler.PerformSetup)             // Legacy endpoint
	api.Post("/setup/automatic", setupHandler.SetupAutomatic) // New automatic setup - all in one
	api.Post("/setup/database", setupHandler.SetupDatabase)
	api.Post("/setup/admin", setupHandler.SetupAdmin)
	api.Post("/setup/complete", setupHandler.SetupComplete)

	// Auth routes - middleware will handle access control
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
	accounts.Post("/:id/verify-code", handler.VerifyAccountCode)
	accounts.Post("/:id/verify-password", handler.VerifyAccountPassword)
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

	// Settings management
	settings := protected.Group("/settings")
	settings.Get("/", settingsHandler.GetAllSettings)
	settings.Get("/category/:category", settingsHandler.GetSettingsByCategory)
	settings.Post("/", settingsHandler.CreateSetting)
	settings.Put("/:key", settingsHandler.UpdateSetting)
	settings.Delete("/:key", settingsHandler.DeleteSetting)

	// User management
	users := protected.Group("/users")
	users.Get("/", usersHandler.ListUsers)
	users.Get("/:id", usersHandler.GetUser)
	users.Post("/", usersHandler.CreateUser)
	users.Put("/:id", usersHandler.UpdateUser)
	users.Delete("/:id", usersHandler.DeleteUser)

	return app
}
