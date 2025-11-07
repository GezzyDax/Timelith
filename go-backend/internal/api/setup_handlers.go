package api

import (
	"fmt"

	"github.com/GezzyDax/timelith/go-backend/internal/auth"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/GezzyDax/timelith/go-backend/internal/setup"
	"github.com/gofiber/fiber/v2"
)

type SetupHandler struct {
	// No database initially since it's not configured yet
}

func NewSetupHandler() *SetupHandler {
	return &SetupHandler{}
}

// SetupStatusResponse indicates whether setup is needed
type SetupStatusResponse struct {
	SetupRequired bool `json:"setup_required"`
}

// CheckSetupStatus returns whether the application needs setup
func (h *SetupHandler) CheckSetupStatus(c *fiber.Ctx) error {
	return c.JSON(SetupStatusResponse{
		SetupRequired: setup.CheckIfSetupNeeded(),
	})
}

// SetupRequest contains all configuration data from the web UI
type SetupRequest struct {
	TelegramAppID    string `json:"telegram_app_id"`
	TelegramAppHash  string `json:"telegram_app_hash"`
	ServerPort       string `json:"server_port"`
	PostgresPassword string `json:"postgres_password"`
	AdminUsername    string `json:"admin_username"`
	AdminPassword    string `json:"admin_password"`
	Environment      string `json:"environment"`
}

// SetupResponse contains the result of setup process
type SetupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PerformSetup handles the complete setup process via web UI
func (h *SetupHandler) PerformSetup(c *fiber.Ctx) error {
	var req SetupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Convert to internal setup config
	setupConfig := &setup.SetupConfig{
		TelegramAppID:    req.TelegramAppID,
		TelegramAppHash:  req.TelegramAppHash,
		ServerPort:       req.ServerPort,
		PostgresPassword: req.PostgresPassword,
		AdminUsername:    req.AdminUsername,
		AdminPassword:    req.AdminPassword,
		Environment:      req.Environment,
	}

	// Generate security keys
	var err error
	setupConfig.JWTSecret, err = setup.GenerateSecret(32)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate JWT secret"})
	}

	setupConfig.EncryptionKey, err = setup.GenerateSecret(32)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate encryption key"})
	}

	// Validate configuration
	if err := setup.ValidateConfig(setupConfig); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Save configuration to .env
	if err := setup.SaveConfig(setupConfig); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save configuration"})
	}

	// Connect to database
	databaseURL := fmt.Sprintf(
		"postgres://timelith:%s@localhost:5432/timelith?sslmode=disable",
		setupConfig.PostgresPassword,
	)

	db, err := database.Connect(databaseURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to connect to database. Make sure PostgreSQL is running.",
		})
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to run database migrations"})
	}

	// Create admin user
	passwordHash, err := auth.HashPassword(setupConfig.AdminPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := &models.User{
		Username:     setupConfig.AdminUsername,
		PasswordHash: passwordHash,
	}

	if err := db.CreateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create admin user"})
	}

	return c.JSON(SetupResponse{
		Success: true,
		Message: "Setup completed successfully! The server will restart automatically.",
	})
}

// ExportedGenerateSecret is an exported wrapper for setup.GenerateSecret
func ExportedGenerateSecret(length int) (string, error) {
	return setup.GenerateSecret(length)
}
