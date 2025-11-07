package api

import (
	"strings"

	"github.com/GezzyDax/timelith/go-backend/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		token := parts[1]

		// Validate token
		claims, err := auth.ValidateToken(token, jwtSecret)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Store user info in locals
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}

func APIKeyMiddleware(apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip auth endpoints
		if strings.HasPrefix(c.Path(), "/api/auth") || c.Path() == "/api/health" {
			return c.Next()
		}

		// Check API key or JWT
		key := c.Get("X-API-Key")
		if key != "" && key == apiKey {
			return c.Next()
		}

		// Otherwise require JWT
		return AuthMiddleware(c.Locals("jwt_secret").(string))(c)
	}
}
