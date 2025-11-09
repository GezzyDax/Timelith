package api

import (
	"github.com/gofiber/fiber/v2"
)

// SetupMiddleware проверяет статус setup и блокирует доступ к API если setup не завершен
func SetupMiddleware(setupRequired bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Разрешенные пути во время setup
		allowedPaths := map[string]bool{
			"/api/setup/status": true,
			"/api/setup":        true,
			"/api/health":       true,
		}

		if setupRequired {
			// Если setup нужен, разрешаем только setup endpoints
			if !allowedPaths[path] {
				return c.Status(503).JSON(fiber.Map{
					"error":          "Setup required",
					"message":        "Please complete the setup process first",
					"setup_required": true,
				})
			}
		} else {
			// Если setup завершен, блокируем setup endpoints
			if path == "/api/setup" {
				return c.Status(403).JSON(fiber.Map{
					"error":   "Setup already completed",
					"message": "Setup can only be performed once",
				})
			}
		}

		return c.Next()
	}
}
