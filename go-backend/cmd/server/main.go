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
	"github.com/GezzyDax/timelith/go-backend/internal/encryption"
	"github.com/GezzyDax/timelith/go-backend/internal/logger"
	"github.com/GezzyDax/timelith/go-backend/internal/scheduler"
	"github.com/GezzyDax/timelith/go-backend/internal/settings"
	"github.com/GezzyDax/timelith/go-backend/internal/setup"
	"github.com/GezzyDax/timelith/go-backend/internal/telegram"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("  Starting Timelith")
	fmt.Println("========================================")
	fmt.Println()

	// Initialize encryption master key
	fmt.Println("🔐 Initializing encryption...")
	if err := encryption.InitMasterKey(); err != nil {
		fmt.Printf("Failed to initialize encryption: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("⚠️  Configuration not fully loaded: %v\n", err)
		fmt.Println("📋 Running in setup mode...")
		cfg = &config.Config{
			ServerPort:  "8080",
			Environment: "production",
		}
	}

	// Initialize logger
	if err := logger.Init(cfg.Environment); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Log.Info("Timelith backend starting",
		zap.String("environment", cfg.Environment),
		zap.String("port", cfg.ServerPort))

	// Connect to database
	// Try to connect with provided config, or use defaults
	databaseURL := cfg.DatabaseURL
	if databaseURL == "" {
		// Use default connection for setup
		pgHost := os.Getenv("POSTGRES_HOST")
		if pgHost == "" {
			pgHost = "postgres"
		}
		pgPassword := os.Getenv("POSTGRES_PASSWORD")
		if pgPassword == "" {
			pgPassword = "timelith_password"
		}
		databaseURL = fmt.Sprintf("postgres://timelith:%s@%s:5432/timelith?sslmode=disable", pgPassword, pgHost)
	}

	db, err := database.Connect(databaseURL)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Log.Info("Connected to database")

	// Run migrations
	logger.Log.Info("Running database migrations...")
	if err := db.RunMigrations(); err != nil {
		logger.Log.Fatal("Failed to run migrations", zap.Error(err))
	}
	logger.Log.Info("Database migrations completed")

	// Initialize Settings Service
	logger.Log.Info("Initializing settings service...")
	settingsService, err := settings.NewService(db)
	if err != nil {
		logger.Log.Fatal("Failed to initialize settings service", zap.Error(err))
	}
	defer settingsService.Stop()
	logger.Log.Info("Settings service initialized")

	// Check setup status from settings
	setupRequired := !settingsService.IsSetupCompleted()

	// Fallback: also check if users exist (for backward compatibility)
	if !setupRequired {
		setupRequired = setup.CheckIfSetupNeeded(db)
	}

	if setupRequired {
		logger.Log.Info("⚙️  Setup required - setup not completed")
		fmt.Println()
		fmt.Println("📋 Setup Required")
		fmt.Println("=" + "=====================================")
		fmt.Println()
		fmt.Println("Please open your browser and navigate to:")
		fmt.Println()
		fmt.Println("  http://localhost:3000/setup")
		fmt.Println()
		fmt.Println("The web-based setup wizard will guide you through")
		fmt.Println("the initial configuration process.")
		fmt.Println()
		fmt.Println("Note: Make sure the web-ui service is running.")
		fmt.Println()
	} else {
		logger.Log.Info("✅ Setup completed - application ready")
	}

	// Setup API router with settings service
	app := api.SetupRouter(cfg, db, settingsService, setupRequired)

	// Initialize other services only if setup is complete
	var sched *scheduler.Scheduler
	if !setupRequired {
		// Initialize Telegram session manager
		sessionManager, err := telegram.NewSessionManager(cfg)
		if err != nil {
			logger.Log.Warn("Failed to initialize session manager", zap.Error(err))
		} else {
			defer sessionManager.Close()
			logger.Log.Info("Telegram session manager initialized")

			// Initialize scheduler
			sched = scheduler.NewScheduler(db, sessionManager)
			ctx := context.Background()

			if err := sched.Start(ctx); err != nil {
				logger.Log.Error("Failed to start scheduler", zap.Error(err))
			} else {
				defer sched.Stop()
				logger.Log.Info("Scheduler started")
			}
		}
	}

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
