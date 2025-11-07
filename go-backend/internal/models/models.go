package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Account represents a Telegram account
type Account struct {
	ID           uuid.UUID      `db:"id" json:"id"`
	Phone        string         `db:"phone" json:"phone"`
	SessionData  []byte         `db:"session_data" json:"-"` // Encrypted session
	Status       string         `db:"status" json:"status"`  // active, inactive, error
	LastLoginAt  sql.NullTime   `db:"last_login_at" json:"last_login_at"`
	ErrorMessage sql.NullString `db:"error_message" json:"error_message,omitempty"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}

// Template represents a message template
type Template struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Content   string    `db:"content" json:"content"`
	Variables []string  `db:"variables" json:"variables"` // JSON array of variable names
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Channel represents a Telegram channel/chat target
type Channel struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	ChatID    string    `db:"chat_id" json:"chat_id"` // Telegram chat ID or username
	Type      string    `db:"type" json:"type"`       // channel, group, user
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Schedule represents a scheduled message job
type Schedule struct {
	ID         uuid.UUID    `db:"id" json:"id"`
	Name       string       `db:"name" json:"name"`
	AccountID  uuid.UUID    `db:"account_id" json:"account_id"`
	TemplateID uuid.UUID    `db:"template_id" json:"template_id"`
	ChannelID  uuid.UUID    `db:"channel_id" json:"channel_id"`
	CronExpr   string       `db:"cron_expr" json:"cron_expr"`
	Timezone   string       `db:"timezone" json:"timezone"` // e.g., "Europe/Moscow"
	Status     string       `db:"status" json:"status"`     // active, paused, completed
	NextRunAt  sql.NullTime `db:"next_run_at" json:"next_run_at"`
	LastRunAt  sql.NullTime `db:"last_run_at" json:"last_run_at"`
	CreatedAt  time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at" json:"updated_at"`
}

// JobLog represents execution history
type JobLog struct {
	ID         uuid.UUID      `db:"id" json:"id"`
	ScheduleID uuid.UUID      `db:"schedule_id" json:"schedule_id"`
	Status     string         `db:"status" json:"status"` // success, failed, retry
	Message    sql.NullString `db:"message" json:"message,omitempty"`
	Error      sql.NullString `db:"error" json:"error,omitempty"`
	ExecutedAt time.Time      `db:"executed_at" json:"executed_at"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
}

// User represents admin user for authentication
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
