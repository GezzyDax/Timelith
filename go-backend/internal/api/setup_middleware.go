package api

import (
	"strings"

	"github.com/GezzyDax/timelith/go-backend/internal/settings"

	"github.com/gofiber/fiber/v2"
)

// SetupMiddleware проверяет статус setup и блокирует доступ к API если setup не завершен
// Теперь проверяет статус динамически через settings service вместо статичного boolean
func SetupMiddleware(settingsService *settings.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Разрешенные пути во время setup
		allowedPaths := map[string]bool{
			"/api/setup/status":    true,
			"/api/setup":           true,
			"/api/setup/automatic": true,
			"/api/setup/database":  true,
			"/api/setup/admin":     true,
			"/api/setup/complete":  true,
			"/api/health":          true,
		}

		// Динамически проверяем статус setup
		setupRequired := !settingsService.IsSetupCompleted()

		isSetupPath := strings.HasPrefix(path, "/api/setup")
		isStatusPath := path == "/api/setup/status"

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
			// Если setup завершен, блокируем любые setup endpoints кроме статуса
			if isSetupPath && !isStatusPath {
				return c.Status(403).JSON(fiber.Map{
					"error":   "Setup already completed",
					"message": "Setup can only be performed once",
				})
			}
		}

		return c.Next()
	}
}
