package api

import (
	"time"

	"github.com/GezzyDax/timelith/go-backend/internal/auth"
	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
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
	jwtSecret := c.Locals("jwt_secret").(string)
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

func (h *Handler) CreateAccount(c *fiber.Ctx) error {
	var req CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	account := &models.Account{
		Phone:  req.Phone,
		Status: "inactive",
	}

	if err := h.db.CreateAccount(account); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(account)
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
	Name       string    `json:"name"`
	AccountID  uuid.UUID `json:"account_id"`
	TemplateID uuid.UUID `json:"template_id"`
	ChannelID  uuid.UUID `json:"channel_id"`
	CronExpr   string    `json:"cron_expr"`
	Timezone   string    `json:"timezone"`
}

func (h *Handler) CreateSchedule(c *fiber.Ctx) error {
	var req CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	schedule := &models.Schedule{
		Name:       req.Name,
		AccountID:  req.AccountID,
		TemplateID: req.TemplateID,
		ChannelID:  req.ChannelID,
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
