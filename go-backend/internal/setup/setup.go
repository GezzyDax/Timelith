package setup

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/GezzyDax/timelith/go-backend/internal/auth"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"golang.org/x/term"
)

type SetupConfig struct {
	TelegramAppID    string
	TelegramAppHash  string
	ServerPort       string
	PostgresPassword string
	JWTSecret        string
	EncryptionKey    string
	AdminUsername    string
	AdminPassword    string
	Environment      string
}

// CheckIfSetupNeeded проверяет, нужна ли установка
// Теперь проверяет наличие пользователей в БД, а не файл .env
func CheckIfSetupNeeded(db interface{ CountUsers() (int, error) }) bool {
	if db == nil {
		return true
	}

	count, err := db.CountUsers()
	if err != nil {
		// Если БД недоступна или таблицы нет, setup нужен
		return true
	}

	// Если нет пользователей, setup нужен
	return count == 0
}

// RunSetup запускает интерактивную установку
func RunSetup() (*SetupConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	config := &SetupConfig{}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  Установка Timelith")
	fmt.Println("========================================")
	fmt.Println()

	// Telegram API credentials
	fmt.Println("1. Настройки Telegram API")
	fmt.Println("   Получите API credentials на https://my.telegram.org")
	fmt.Println()

	config.TelegramAppID = readInput(reader, "Telegram App ID", "")
	config.TelegramAppHash = readInput(reader, "Telegram App Hash", "")
	fmt.Println()

	// Server settings
	fmt.Println("2. Настройки сервера")
	config.ServerPort = readInput(reader, "Порт сервера", "8080")
	config.Environment = readInput(reader, "Окружение (production/development)", "production")
	fmt.Println()

	// Database settings
	fmt.Println("3. Настройки базы данных PostgreSQL")
	config.PostgresPassword = readInput(reader, "Пароль PostgreSQL", "timelith_password")
	fmt.Println()

	// Security settings - auto-generate
	fmt.Println("4. Генерация ключей безопасности...")
	var err error
	config.JWTSecret, err = GenerateSecret(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
	}
	config.EncryptionKey, err = GenerateSecret(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}
	fmt.Println("   ✓ JWT_SECRET сгенерирован")
	fmt.Println("   ✓ ENCRYPTION_KEY сгенерирован")
	fmt.Println()

	// Admin user
	fmt.Println("5. Создание первого администратора")
	config.AdminUsername = readInput(reader, "Логин администратора", "admin")
	config.AdminPassword = readPassword("Пароль администратора")
	confirmPassword := readPassword("Подтвердите пароль")

	if config.AdminPassword != confirmPassword {
		return nil, fmt.Errorf("пароли не совпадают")
	}
	fmt.Println()

	return config, nil
}

// SaveConfig сохраняет конфигурацию в .env файл
func SaveConfig(config *SetupConfig) error {
	// Determine PostgreSQL host based on environment
	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == "" {
		pgHost = "postgres" // Default to Docker service name
	}

	databaseURL := fmt.Sprintf(
		"postgres://timelith:%s@%s:5432/timelith?sslmode=disable",
		config.PostgresPassword,
		pgHost,
	)

	envContent := fmt.Sprintf(`# Database
DATABASE_URL=%s
POSTGRES_PASSWORD=%s

# Telegram API Credentials
TELEGRAM_APP_ID=%s
TELEGRAM_APP_HASH=%s

# Security (auto-generated)
JWT_SECRET=%s
ENCRYPTION_KEY=%s

# Server
SERVER_PORT=%s
ENVIRONMENT=%s

# API URL (for web frontend)
NEXT_PUBLIC_API_URL=http://localhost:%s
`,
		databaseURL,
		config.PostgresPassword,
		config.TelegramAppID,
		config.TelegramAppHash,
		config.JWTSecret,
		config.EncryptionKey,
		config.ServerPort,
		config.Environment,
		config.ServerPort,
	)

	// Создаем backup существующего .env файла, если он есть
	if err := backupEnvFile(); err != nil {
		// Предупреждение, но не критическая ошибка
		fmt.Printf("⚠ Warning: Failed to backup .env file: %v\n", err)
	}

	err := os.WriteFile(".env", []byte(envContent), 0600)
	if err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	fmt.Println("✓ Конфигурация сохранена в .env")
	return nil
}

// backupEnvFile создает резервную копию .env файла с timestamp
func backupEnvFile() error {
	// Проверяем, существует ли .env файл
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// Файл не существует, backup не нужен
		return nil
	}

	// Создаем имя файла с timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := fmt.Sprintf(".env.backup_%s", timestamp)

	// Открываем исходный файл
	source, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("failed to open .env: %w", err)
	}
	defer source.Close()

	// Создаем backup файл
	backup, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer backup.Close()

	// Копируем содержимое
	if _, err := io.Copy(backup, source); err != nil {
		return fmt.Errorf("failed to copy .env to backup: %w", err)
	}

	fmt.Printf("✓ Backup создан: %s\n", backupPath)
	return nil
}

// CreateAdminUser создает первого администратора в базе данных
func CreateAdminUser(db *database.DB, username, password string) error {
	// Проверяем, есть ли уже пользователи
	existingUser, _ := db.GetUserByUsername(username)
	if existingUser != nil {
		return fmt.Errorf("пользователь '%s' уже существует", username)
	}

	// Хешируем пароль
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя
	user := &models.User{
		Username:     username,
		PasswordHash: passwordHash,
	}

	if err := db.CreateUser(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("✓ Администратор '%s' создан\n", username)
	return nil
}

// GenerateSecret генерирует криптографически стойкий случайный ключ (exported)
func GenerateSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// readInput читает ввод пользователя с значением по умолчанию
func readInput(reader *bufio.Reader, prompt, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultValue != "" {
		return defaultValue
	}
	return input
}

// readPassword читает пароль без отображения на экране
func readPassword(prompt string) string {
	fmt.Printf("%s: ", prompt)
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(password)
}

// ValidateConfig проверяет корректность конфигурации
func ValidateConfig(config *SetupConfig) error {
	// Validate Telegram App ID
	if _, err := strconv.Atoi(config.TelegramAppID); err != nil {
		return fmt.Errorf("telegram App ID должен быть числом")
	}

	// Validate App Hash
	if len(config.TelegramAppHash) < 32 {
		return fmt.Errorf("telegram App Hash слишком короткий")
	}

	// Validate port
	port, err := strconv.Atoi(config.ServerPort)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("некорректный порт сервера")
	}

	// Validate admin credentials
	if len(config.AdminUsername) < 3 {
		return fmt.Errorf("логин администратора должен быть минимум 3 символа")
	}

	if len(config.AdminPassword) < 6 {
		return fmt.Errorf("пароль администратора должен быть минимум 6 символов")
	}

	return nil
}

// ShowSummary показывает итоговую информацию после установки
func ShowSummary(config *SetupConfig) {
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  Установка завершена!")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Printf("Сервер запущен на порту: %s\n", config.ServerPort)
	fmt.Printf("Администратор: %s\n", config.AdminUsername)
	fmt.Println("\nДля входа в систему используйте созданные учетные данные.")
	fmt.Printf("API URL: http://localhost:%s\n", config.ServerPort)
	fmt.Println()
}
