package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/GezzyDax/timelith/go-backend/internal/api"
	"github.com/GezzyDax/timelith/go-backend/internal/config"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/logger"
	"github.com/GezzyDax/timelith/go-backend/internal/scheduler"
	"github.com/GezzyDax/timelith/go-backend/internal/setup"
	"github.com/GezzyDax/timelith/go-backend/internal/telegram"
	"go.uber.org/zap"
)

func main() {
	// Check if setup is needed
	if setup.CheckIfSetupNeeded() {
		fmt.Println("Конфигурация не найдена. Запуск мастера установки...\n")

		// Run interactive setup
		setupConfig, err := setup.RunSetup()
		if err != nil {
			fmt.Printf("Ошибка установки: %v\n", err)
			os.Exit(1)
		}

		// Validate configuration
		if err := setup.ValidateConfig(setupConfig); err != nil {
			fmt.Printf("Ошибка валидации: %v\n", err)
			os.Exit(1)
		}

		// Save configuration to .env
		if err := setup.SaveConfig(setupConfig); err != nil {
			fmt.Printf("Ошибка сохранения конфигурации: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nПодключение к базе данных для создания администратора...")

		// Build database URL
		databaseURL := fmt.Sprintf(
			"postgres://timelith:%s@localhost:5432/timelith?sslmode=disable",
			setupConfig.PostgresPassword,
		)

		// Connect to database
		db, err := database.Connect(databaseURL)
		if err != nil {
			fmt.Printf("Ошибка подключения к БД: %v\n", err)
			fmt.Println("Убедитесь, что PostgreSQL запущен и доступен.")
			os.Exit(1)
		}

		// Run migrations
		if err := db.RunMigrations(); err != nil {
			fmt.Printf("Ошибка миграций БД: %v\n", err)
			db.Close()
			os.Exit(1)
		}

		// Create admin user
		if err := setup.CreateAdminUser(db, setupConfig.AdminUsername, setupConfig.AdminPassword); err != nil {
			fmt.Printf("Ошибка создания администратора: %v\n", err)
			db.Close()
			os.Exit(1)
		}

		db.Close()

		// Show summary
		setup.ShowSummary(setupConfig)

		fmt.Println("Запуск сервера...\n")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Environment); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Log.Info("Starting Timelith backend",
		zap.String("environment", cfg.Environment),
		zap.String("port", cfg.ServerPort))

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Log.Info("Connected to database")

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		logger.Log.Fatal("Failed to run migrations", zap.Error(err))
	}

	logger.Log.Info("Database migrations completed")

	// Initialize Telegram session manager
	sessionManager, err := telegram.NewSessionManager(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to initialize session manager", zap.Error(err))
	}
	defer sessionManager.Close()

	logger.Log.Info("Telegram session manager initialized")

	// Initialize scheduler
	sched := scheduler.NewScheduler(db, sessionManager)
	ctx := context.Background()

	if err := sched.Start(ctx); err != nil {
		logger.Log.Fatal("Failed to start scheduler", zap.Error(err))
	}
	defer sched.Stop()

	logger.Log.Info("Scheduler started")

	// Setup API router
	app := api.SetupRouter(cfg, db)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.ServerPort)
		logger.Log.Info("Starting HTTP server", zap.String("address", addr))

		if err := app.Listen(addr); err != nil {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		logger.Log.Error("Server shutdown error", zap.Error(err))
	}

	logger.Log.Info("Server shutdown complete")
}
