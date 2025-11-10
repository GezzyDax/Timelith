package api

import (
	"log"

	"github.com/GezzyDax/timelith/go-backend/internal/auth"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UsersHandler handles user management endpoints
type UsersHandler struct {
	db *database.DB
}

// NewUsersHandler creates a new users handler
func NewUsersHandler(db *database.DB) *UsersHandler {
	return &UsersHandler{db: db}
}

// ListUsers returns all users
func (h *UsersHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.db.ListUsers()
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	return c.JSON(users)
}

// GetUser returns a single user by ID
func (h *UsersHandler) GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.db.GetUserByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// CreateUserRequest for creating a new user
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUser creates a new user
func (h *UsersHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
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

	// Check if username already exists
	existingUser, err := h.db.GetUserByUsername(req.Username)
	if err == nil && existingUser != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	if err := h.db.CreateUser(user); err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	log.Printf("User '%s' created successfully", req.Username)

	return c.Status(201).JSON(user)
}

// UpdateUserRequest for updating a user
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// UpdateUser updates an existing user
func (h *UsersHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get existing user
	user, err := h.db.GetUserByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Update username if provided
	if req.Username != "" {
		user.Username = req.Username
	}

	// Update password if provided
	if req.Password != "" {
		if len(req.Password) < 6 {
			return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters long"})
		}

		passwordHash, err := auth.HashPassword(req.Password)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		user.PasswordHash = passwordHash
	}

	// Update in database
	if err := h.db.UpdateUser(user); err != nil {
		log.Printf("Failed to update user: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	log.Printf("User '%s' updated successfully", user.Username)

	return c.JSON(user)
}

// DeleteUser deletes a user
func (h *UsersHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Check if user exists
	user, err := h.db.GetUserByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if this is the last user
	count, err := h.db.CountUsers()
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to verify user count"})
	}

	if count <= 1 {
		return c.Status(403).JSON(fiber.Map{"error": "Cannot delete the last user"})
	}

	// Delete user
	if err := h.db.DeleteUser(id); err != nil {
		log.Printf("Failed to delete user: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	log.Printf("User '%s' deleted successfully", user.Username)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}
