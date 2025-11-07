package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/timelith/backend/internal/api"
	"github.com/timelith/backend/internal/config"
	"github.com/timelith/backend/internal/database"
	"github.com/timelith/backend/internal/scheduler"
	"github.com/timelith/backend/internal/telegram"
	"github.com/timelith/backend/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()
	defer log.Sync()

	log.Info("Starting Timelith Backend...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	// Initialize Redis
	redisClient := database.ConnectRedis(cfg.RedisURL)
	defer redisClient.Close()

	// Initialize Telegram manager
	telegramManager := telegram.NewManager(cfg, db, log)
	if err := telegramManager.Initialize(); err != nil {
		log.Error("Failed to initialize Telegram manager", "error", err)
	}

	// Initialize scheduler
	schedulerService := scheduler.New(db, redisClient, telegramManager, log)
	schedulerService.Start()
	defer schedulerService.Stop()

	// Initialize API server
	router := api.NewRouter(cfg, db, redisClient, telegramManager, log)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
