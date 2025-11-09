package api

import (
	"log"

	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/GezzyDax/timelith/go-backend/internal/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SettingsHandler handles settings management endpoints
type SettingsHandler struct {
	settingsService *settings.Service
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(settingsService *settings.Service) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
	}
}

// GetAllSettings returns all settings
func (h *SettingsHandler) GetAllSettings(c *fiber.Ctx) error {
	settings, err := h.settingsService.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve settings"})
	}

	// Filter out sensitive values that shouldn't be shown in UI
	for i := range settings {
		if settings[i].Encrypted && !settings[i].Editable {
			settings[i].Value = "***HIDDEN***"
		}
	}

	return c.JSON(settings)
}

// GetSettingsByCategory returns settings filtered by category
func (h *SettingsHandler) GetSettingsByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Category is required"})
	}

	settings, err := h.settingsService.GetByCategory(category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve settings"})
	}

	// Filter out sensitive values that shouldn't be shown in UI
	for i := range settings {
		if settings[i].Encrypted && !settings[i].Editable {
			settings[i].Value = "***HIDDEN***"
		}
	}

	return c.JSON(settings)
}

// UpdateSettingRequest for updating a setting
type UpdateSettingRequest struct {
	Value string `json:"value"`
}

// UpdateSetting updates a single setting
func (h *SettingsHandler) UpdateSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Setting key is required"})
	}

	var req UpdateSettingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get current user ID from JWT token
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Get current setting to check if it's editable and encrypted
	allSettings, err := h.settingsService.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve setting"})
	}

	var currentSetting *models.Setting
	for i := range allSettings {
		if allSettings[i].Key == key {
			currentSetting = &allSettings[i]
			break
		}
	}

	if currentSetting == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Setting not found"})
	}

	if !currentSetting.Editable {
		return c.Status(403).JSON(fiber.Map{"error": "This setting is not editable"})
	}

	// Update the setting
	if err := h.settingsService.Set(key, req.Value, currentSetting.Encrypted, currentSetting.Category, &userID); err != nil {
		log.Printf("Failed to update setting %s: %v", key, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update setting"})
	}

	log.Printf("Setting '%s' updated by user %s", key, userID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Setting updated successfully",
	})
}

// DeleteSetting deletes a setting (only if editable)
func (h *SettingsHandler) DeleteSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Setting key is required"})
	}

	// Get current setting to check if it's editable
	allSettings, err := h.settingsService.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve setting"})
	}

	var currentSetting *models.Setting
	for i := range allSettings {
		if allSettings[i].Key == key {
			currentSetting = &allSettings[i]
			break
		}
	}

	if currentSetting == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Setting not found"})
	}

	if !currentSetting.Editable {
		return c.Status(403).JSON(fiber.Map{"error": "This setting cannot be deleted"})
	}

	if err := h.settingsService.Delete(key); err != nil {
		log.Printf("Failed to delete setting %s: %v", key, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete setting"})
	}

	log.Printf("Setting '%s' deleted", key)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Setting deleted successfully",
	})
}

// CreateSettingRequest for creating a new setting
type CreateSettingRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Encrypted   bool   `json:"encrypted"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// CreateSetting creates a new setting
func (h *SettingsHandler) CreateSetting(c *fiber.Ctx) error {
	var req CreateSettingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Key == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Setting key is required"})
	}
	if req.Category == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Category is required"})
	}

	// Get current user ID from JWT token
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Create the setting
	if err := h.settingsService.Set(req.Key, req.Value, req.Encrypted, req.Category, &userID); err != nil {
		log.Printf("Failed to create setting %s: %v", req.Key, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create setting"})
	}

	log.Printf("Setting '%s' created by user %s", req.Key, userID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Setting created successfully",
	})
}
