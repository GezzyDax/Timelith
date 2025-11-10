package api

import (
	"fmt"
	"log"

	"github.com/GezzyDax/timelith/go-backend/internal/auth"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/GezzyDax/timelith/go-backend/internal/settings"
	"github.com/GezzyDax/timelith/go-backend/internal/setup"
	"github.com/gofiber/fiber/v2"
)

type SetupHandler struct {
	db              *database.DB
	settingsService *settings.Service
}

func NewSetupHandler(db *database.DB, settingsService *settings.Service) *SetupHandler {
	return &SetupHandler{
		db:              db,
		settingsService: settingsService,
	}
}

// SetupStatusResponse indicates whether setup is needed
type SetupStatusResponse struct {
	SetupRequired bool `json:"setup_required"`
}

// CheckSetupStatus returns whether the application needs setup
func (h *SetupHandler) CheckSetupStatus(c *fiber.Ctx) error {
	setupRequired := !h.settingsService.IsSetupCompleted()

	// Fallback: also check if users exist (for backward compatibility)
	if !setupRequired {
		setupRequired = setup.CheckIfSetupNeeded(h.db)
	}

	return c.JSON(SetupStatusResponse{
		SetupRequired: setupRequired,
	})
}

// DatabaseSetupRequest for Stage 1: Database configuration
type DatabaseSetupRequest struct {
	UseDockerDatabase bool   `json:"use_docker_database"` // true = use Docker, false = external
	DatabaseURL       string `json:"database_url"`        // Only if use_docker_database = false
}

// AdminSetupRequest for Stage 2: Admin user creation
type AdminSetupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ApiKeysSetupRequest for Stage 3: API keys
type ApiKeysSetupRequest struct {
	TelegramAppID   string `json:"telegram_app_id"`
	TelegramAppHash string `json:"telegram_app_hash"`
}

// SetupRequest contains all configuration data from the web UI (legacy)
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
	log.Println("Starting setup process...")

	// Проверяем, не завершен ли уже setup
	if !setup.CheckIfSetupNeeded(h.db) {
		log.Println("Setup already completed")
		return c.Status(403).JSON(fiber.Map{
			"error":   "Setup already completed",
			"message": "Setup can only be performed once",
		})
	}

	var req SetupRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Println("Request parsed successfully")

	// Validate required fields
	if req.TelegramAppID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Telegram App ID is required"})
	}
	if req.TelegramAppHash == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Telegram App Hash is required"})
	}
	if req.AdminUsername == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Admin username is required"})
	}
	if req.AdminPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Admin password is required"})
	}
	if len(req.AdminPassword) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Admin password must be at least 6 characters long"})
	}

	log.Println("Input validation passed")

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

	log.Println("Generating security keys...")

	// Generate security keys
	var err error
	setupConfig.JWTSecret, err = setup.GenerateSecret(32)
	if err != nil {
		log.Printf("Failed to generate JWT secret: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate JWT secret"})
	}

	setupConfig.EncryptionKey, err = setup.GenerateSecret(32)
	if err != nil {
		log.Printf("Failed to generate encryption key: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate encryption key"})
	}

	log.Println("Validating configuration...")

	// Validate configuration
	if err := setup.ValidateConfig(setupConfig); err != nil {
		log.Printf("Configuration validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	log.Println("Saving configuration to .env...")

	// Save configuration to .env
	if err := setup.SaveConfig(setupConfig); err != nil {
		log.Printf("Failed to save configuration: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save configuration"})
	}

	log.Println("Creating admin user...")

	// Create admin user
	passwordHash, err := auth.HashPassword(setupConfig.AdminPassword)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := &models.User{
		Username:     setupConfig.AdminUsername,
		PasswordHash: passwordHash,
	}

	if err := h.db.CreateUser(user); err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create admin user: %v", err)})
	}

	log.Println("Setup completed successfully!")

	return c.JSON(SetupResponse{
		Success: true,
		Message: "Setup completed successfully! Please refresh the page to continue.",
	})
}

// SetupDatabase handles Stage 1: Database configuration
func (h *SetupHandler) SetupDatabase(c *fiber.Ctx) error {
	if h.settingsService.IsSetupCompleted() {
		return c.Status(403).JSON(fiber.Map{"error": "Setup already completed"})
	}

	var req DatabaseSetupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var databaseURL string
	if req.UseDockerDatabase {
		// Generate random password for Docker database
		password, err := setup.GenerateSecret(16)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate database password"})
		}
		databaseURL = fmt.Sprintf("postgres://timelith:%s@postgres:5432/timelith?sslmode=disable", password)

		// Save postgres password to settings
		if err := h.settingsService.Set("postgres_password", password, true, "database", nil); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save database password"})
		}
	} else {
		if req.DatabaseURL == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Database URL is required"})
		}
		databaseURL = req.DatabaseURL
	}

	// Save database URL to settings (encrypted)
	if err := h.settingsService.Set("database_url", databaseURL, true, "database", nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save database configuration"})
	}

	log.Println("Database configuration saved successfully")

	return c.JSON(SetupResponse{
		Success: true,
		Message: "Database configured successfully",
	})
}

// SetupAdmin handles Stage 2: Admin user creation
func (h *SetupHandler) SetupAdmin(c *fiber.Ctx) error {
	if h.settingsService.IsSetupCompleted() {
		return c.Status(403).JSON(fiber.Map{"error": "Setup already completed"})
	}

	var req AdminSetupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate input
	if req.Username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username is required"})
	}
	if req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}
	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters long"})
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Create admin user
	user := &models.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	if err := h.db.CreateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create admin user: %v", err)})
	}

	log.Printf("Admin user '%s' created successfully", req.Username)

	return c.JSON(SetupResponse{
		Success: true,
		Message: "Admin user created successfully",
	})
}

// SetupComplete handles Stage 3: API keys and completes setup
func (h *SetupHandler) SetupComplete(c *fiber.Ctx) error {
	if h.settingsService.IsSetupCompleted() {
		return c.Status(403).JSON(fiber.Map{"error": "Setup already completed"})
	}

	var req ApiKeysSetupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate Telegram credentials
	if req.TelegramAppID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Telegram App ID is required"})
	}
	if req.TelegramAppHash == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Telegram App Hash is required"})
	}

	// Save Telegram credentials to settings (encrypted)
	if err := h.settingsService.Set("telegram_app_id", req.TelegramAppID, true, "telegram", nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save Telegram App ID"})
	}

	if err := h.settingsService.Set("telegram_app_hash", req.TelegramAppHash, true, "telegram", nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save Telegram App Hash"})
	}

	// Generate JWT secret if not exists
	jwtSecret, err := setup.GenerateSecret(32)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate JWT secret"})
	}
	if err := h.settingsService.Set("jwt_secret", jwtSecret, true, "security", nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save JWT secret"})
	}

	// Mark setup as completed
	if err := h.settingsService.MarkSetupCompleted(nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to mark setup as completed"})
	}

	log.Println("Setup completed successfully!")

	return c.JSON(SetupResponse{
		Success: true,
		Message: "Setup completed successfully! You can now log in with your admin credentials.",
	})
}

// SetupAutomatic handles automatic all-in-one setup - just username and password needed
func (h *SetupHandler) SetupAutomatic(c *fiber.Ctx) error {
	if h.settingsService.IsSetupCompleted() {
		return c.Status(403).JSON(fiber.Map{"error": "Setup already completed"})
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate input
	if req.Username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username is required"})
	}
	if len(req.Username) < 3 {
		return c.Status(400).JSON(fiber.Map{"error": "Username must be at least 3 characters long"})
	}
	if req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}
	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters long"})
	}

	log.Println("Starting automatic setup...")

	// Step 1: Auto-configure Docker database
	password, err := setup.GenerateSecret(16)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate database password"})
	}
	databaseURL := fmt.Sprintf("postgres://timelith:%s@postgres:5432/timelith?sslmode=disable", password)

	if err := h.settingsService.Set("postgres_password", password, true, "database", nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save database password"})
	}

	if err := h.settingsService.Set("database_url", databaseURL, true, "database", nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save database configuration"})
	}

	log.Println("✓ Database configured automatically")

	// Step 2: Create admin user
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	if err := h.db.CreateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create admin user: %v", err)})
	}

	log.Printf("✓ Admin user '%s' created", req.Username)

	// Step 3: Generate JWT secret
	jwtSecret, err := setup.GenerateSecret(32)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate JWT secret"})
	}
	if err := h.settingsService.Set("jwt_secret", jwtSecret, true, "security", nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save JWT secret"})
	}

	log.Println("✓ Security keys generated")

	// Step 4: Mark setup as completed
	if err := h.settingsService.MarkSetupCompleted(nil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to mark setup as completed"})
	}

	log.Println("✓ Automatic setup completed successfully!")

	return c.JSON(SetupResponse{
		Success: true,
		Message: "Setup completed automatically! You can now log in with your admin credentials.",
	})
}

// ExportedGenerateSecret is an exported wrapper for setup.GenerateSecret
func ExportedGenerateSecret(length int) (string, error) {
	return setup.GenerateSecret(length)
}
