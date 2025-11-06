package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	"github.com/timelith/backend/internal/config"
	"github.com/timelith/backend/internal/telegram"
	"github.com/timelith/backend/pkg/logger"
	"gorm.io/gorm"
)

type Server struct {
	config          *config.Config
	db              *gorm.DB
	redis           *redis.Client
	telegramManager *telegram.Manager
	log             *logger.Logger
}

func NewRouter(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, tm *telegram.Manager, log *logger.Logger) http.Handler {
	s := &Server{
		config:          cfg,
		db:              db,
		redis:           redisClient,
		telegramManager: tm,
		log:             log,
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", s.handleHealth)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(s.apiKeyMiddleware)

		// Telegram auth endpoints
		r.Post("/telegram/auth/send_code", s.handleSendCode)
		r.Post("/telegram/auth/verify_code", s.handleVerifyCode)
		r.Post("/telegram/disconnect", s.handleDisconnect)
		r.Get("/telegram/status/{account_id}", s.handleStatus)

		// Channels
		r.Post("/telegram/channels/sync", s.handleSyncChannels)

		// Schedules
		r.Get("/schedules", s.handleGetSchedules)

		// Send logs
		r.Post("/send_logs", s.handleCreateSendLog)
	})

	return r
}

func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != s.config.APIKey {
			s.respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
