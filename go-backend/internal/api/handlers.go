package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/GezzyDax/timelith/go-backend/internal/auth"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/GezzyDax/timelith/go-backend/internal/telegram"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	db             *database.DB
	sessionManager *telegram.SessionManager
}

func NewHandler(db *database.DB, sessionManager *telegram.SessionManager) *Handler {
	return &Handler{db: db, sessionManager: sessionManager}
}

func (h *Handler) requireSessionManager(c *fiber.Ctx) bool {
	if h.sessionManager == nil {
		c.Status(503).JSON(fiber.Map{
			"error":   "Telegram integration not ready",
			"message": "Complete initial setup before linking accounts",
		})
		return false
	}
	return true
}

// Health check
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

// Auth handlers

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Get user
	user, err := h.db.GetUserByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT
	jwtSecret, ok := c.Locals("jwt_secret").(string)
	if !ok || jwtSecret == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Server configuration error: JWT secret not available"})
	}
	token, err := auth.GenerateToken(user.ID.String(), user.Username, jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(&LoginResponse{
		Token: token,
		User:  user,
	})
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	if err := h.db.CreateUser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Username already exists"})
	}

	return c.Status(201).JSON(user)
}

// Account handlers

func (h *Handler) ListAccounts(c *fiber.Ctx) error {
	accounts, err := h.db.ListAccounts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(accounts)
}

func (h *Handler) GetAccount(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	account, err := h.db.GetAccount(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account not found"})
	}

	return c.JSON(account)
}

type CreateAccountRequest struct {
	Phone string `json:"phone"`
}

type VerifyAccountCodeRequest struct {
	Code string `json:"code"`
}

type VerifyAccountPasswordRequest struct {
	Password string `json:"password"`
}

func (h *Handler) CreateAccount(c *fiber.Ctx) error {
	if !h.requireSessionManager(c) {
		return nil
	}

	var req CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Phone == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phone is required"})
	}

	ctx := context.Background()

	account, err := h.db.GetAccountByPhone(req.Phone)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to lookup account"})
		}

		account = &models.Account{
			Phone:  req.Phone,
			Status: "pending",
		}

		if err := h.db.CreateAccount(account); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	} else if account.Status == "active" {
		return c.Status(400).JSON(fiber.Map{"error": "Account already active"})
	}

	phoneCodeHash, pendingSession, err := h.sessionManager.AuthenticatePhone(ctx, req.Phone)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to send code: %v", err)})
	}

	if err := h.db.UpdateAccountCodeState(account.ID, phoneCodeHash, "code_sent", pendingSession); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to persist code state"})
	}

	account.Status = "code_sent"
	account.LoginCodeSentAt = models.NewNullTime(time.Now())
	account.TwoFactorRequired = false
	account.TwoFactorHint = models.NullString{}
	account.SessionData = pendingSession

	return c.Status(202).JSON(account)
}

func (h *Handler) VerifyAccountCode(c *fiber.Ctx) error {
	if !h.requireSessionManager(c) {
		return nil
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid account ID"})
	}

	var req VerifyAccountCodeRequest
	if err := c.BodyParser(&req); err != nil || req.Code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Verification code is required"})
	}

	account, err := h.db.GetAccount(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(404).JSON(fiber.Map{"error": "Account not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load account"})
	}

	if !account.PhoneCodeHash.Valid {
		return c.Status(400).JSON(fiber.Map{"error": "Account is not awaiting verification"})
	}

	ctx := context.Background()
	if len(account.SessionData) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Session expired, request a new code"})
	}

	rawPending, err := h.sessionManager.DecryptSession(account.SessionData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Stored session corrupted; request a new code"})
	}

	finalSession, pendingSession, requiresPassword, passwordHint, err := h.sessionManager.VerifyCode(ctx, account.Phone, req.Code, account.PhoneCodeHash.String, rawPending)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to verify code: %v", err)})
	}

	if requiresPassword {
		hint := models.NullString{}
		if passwordHint != "" {
			hint = models.NewNullString(passwordHint)
		}
		if err := h.db.MarkAccountPasswordRequired(account.ID, hint, pendingSession); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update account state"})
		}

		account.Status = "password_required"
		account.TwoFactorRequired = true
		account.TwoFactorHint = hint
		account.PhoneCodeHash = models.NullString{}
		account.SessionData = pendingSession

		return c.JSON(fiber.Map{
			"requires_password": true,
			"password_hint":     passwordHint,
			"account":           account,
		})
	}

	if err := h.db.SaveAccountSession(account.ID, finalSession); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to store session"})
	}

	account.Status = "active"
	account.TwoFactorRequired = false
	account.TwoFactorHint = models.NullString{}
	account.PhoneCodeHash = models.NullString{}
	account.LastLoginAt = models.NewNullTime(time.Now())
	account.SessionData = finalSession

	return c.JSON(fiber.Map{
		"requires_password": false,
		"account":           account,
	})
}

func (h *Handler) VerifyAccountPassword(c *fiber.Ctx) error {
	if !h.requireSessionManager(c) {
		return nil
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid account ID"})
	}

	var req VerifyAccountPasswordRequest
	if err := c.BodyParser(&req); err != nil || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}

	account, err := h.db.GetAccount(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(404).JSON(fiber.Map{"error": "Account not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load account"})
	}

	if !account.TwoFactorRequired {
		return c.Status(400).JSON(fiber.Map{"error": "Password verification not required"})
	}

	if len(account.SessionData) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Pending session not found; request a new code"})
	}

	ctx := context.Background()
	sessionData, err := h.sessionManager.VerifyPassword(ctx, account.Phone, account.SessionData, req.Password)
	if err != nil {
		if errors.Is(err, telegram.ErrInvalidPassword) {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid password"})
		}
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to verify password: %v", err)})
	}

	if err := h.db.SaveAccountSession(account.ID, sessionData); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to store session"})
	}

	account.Status = "active"
	account.TwoFactorRequired = false
	account.TwoFactorHint = models.NullString{}
	account.LastLoginAt = models.NewNullTime(time.Now())
	account.SessionData = sessionData

	return c.JSON(fiber.Map{
		"account": account,
	})
}

func (h *Handler) DeleteAccount(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.db.DeleteAccount(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

// Template handlers

func (h *Handler) ListTemplates(c *fiber.Ctx) error {
	templates, err := h.db.ListTemplates()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(templates)
}

func (h *Handler) GetTemplate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	template, err := h.db.GetTemplate(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Template not found"})
	}

	return c.JSON(template)
}

type CreateTemplateRequest struct {
	Name      string   `json:"name"`
	Content   string   `json:"content"`
	Variables []string `json:"variables"`
}

func (h *Handler) CreateTemplate(c *fiber.Ctx) error {
	var req CreateTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	template := &models.Template{
		Name:      req.Name,
		Content:   req.Content,
		Variables: req.Variables,
	}

	if err := h.db.CreateTemplate(template); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(template)
}

func (h *Handler) UpdateTemplate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req CreateTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	template := &models.Template{
		ID:        id,
		Name:      req.Name,
		Content:   req.Content,
		Variables: req.Variables,
	}

	if err := h.db.UpdateTemplate(template); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(template)
}

func (h *Handler) DeleteTemplate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.db.DeleteTemplate(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

// Channel handlers

func (h *Handler) ListChannels(c *fiber.Ctx) error {
	channels, err := h.db.ListChannels()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(channels)
}

func (h *Handler) GetChannel(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	channel, err := h.db.GetChannel(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Channel not found"})
	}

	return c.JSON(channel)
}

type CreateChannelRequest struct {
	Name   string `json:"name"`
	ChatID string `json:"chat_id"`
	Type   string `json:"type"`
}

func (h *Handler) CreateChannel(c *fiber.Ctx) error {
	var req CreateChannelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	channel := &models.Channel{
		Name:   req.Name,
		ChatID: req.ChatID,
		Type:   req.Type,
	}

	if err := h.db.CreateChannel(channel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(channel)
}

func (h *Handler) UpdateChannel(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req CreateChannelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	channel := &models.Channel{
		ID:     id,
		Name:   req.Name,
		ChatID: req.ChatID,
		Type:   req.Type,
	}

	if err := h.db.UpdateChannel(channel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(channel)
}

func (h *Handler) DeleteChannel(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.db.DeleteChannel(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

// Schedule handlers

func (h *Handler) ListSchedules(c *fiber.Ctx) error {
	schedules, err := h.db.ListSchedules()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(schedules)
}

func (h *Handler) GetSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	schedule, err := h.db.GetSchedule(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Schedule not found"})
	}

	return c.JSON(schedule)
}

type CreateScheduleRequest struct {
	Name       string      `json:"name"`
	AccountID  uuid.UUID   `json:"account_id"`
	TemplateID uuid.UUID   `json:"template_id"`
	ChannelIDs []uuid.UUID `json:"channel_ids"`
	CronExpr   string      `json:"cron_expr"`
	Timezone   string      `json:"timezone"`
}

func formatUUIDs(ids []uuid.UUID) []string {
	if len(ids) == 0 {
		return []string{}
	}

	formatted := make([]string, 0, len(ids))
	for _, id := range ids {
		formatted = append(formatted, id.String())
	}
	return formatted
}

func (h *Handler) CreateSchedule(c *fiber.Ctx) error {
	var req CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if len(req.ChannelIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "At least one channel_id is required"})
	}

	schedule := &models.Schedule{
		Name:       req.Name,
		AccountID:  req.AccountID,
		TemplateID: req.TemplateID,
		ChannelIDs: formatUUIDs(req.ChannelIDs),
		CronExpr:   req.CronExpr,
		Timezone:   req.Timezone,
		Status:     "active",
	}

	if err := h.db.CreateSchedule(schedule); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(schedule)
}

func (h *Handler) UpdateSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Get existing schedule
	schedule, err := h.db.GetSchedule(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Schedule not found"})
	}

	var req CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	schedule.Name = req.Name
	if len(req.ChannelIDs) > 0 {
		schedule.ChannelIDs = formatUUIDs(req.ChannelIDs)
	}
	schedule.CronExpr = req.CronExpr
	schedule.Timezone = req.Timezone

	if err := h.db.UpdateSchedule(schedule); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(schedule)
}

type UpdateScheduleStatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateScheduleStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req UpdateScheduleStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	schedule, err := h.db.GetSchedule(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Schedule not found"})
	}

	schedule.Status = req.Status

	if err := h.db.UpdateSchedule(schedule); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(schedule)
}

func (h *Handler) DeleteSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.db.DeleteSchedule(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

// Job log handlers

func (h *Handler) GetScheduleLogs(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	logs, err := h.db.GetJobLogs(id, 50)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(logs)
}

func (h *Handler) GetAllLogs(c *fiber.Ctx) error {
	logs, err := h.db.GetAllJobLogs(100)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(logs)
}
