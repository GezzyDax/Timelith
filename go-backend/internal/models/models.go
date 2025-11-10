package models

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a Telegram account
type Account struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	Phone             string     `db:"phone" json:"phone"`
	SessionData       []byte     `db:"session_data" json:"-"` // Encrypted session
	Status            string     `db:"status" json:"status"`  // active, inactive, error
	ProxyEnabled      bool       `db:"proxy_enabled" json:"proxy_enabled"`
	ProxyHost         NullString `db:"proxy_host" json:"proxy_host"`
	ProxyPort         NullInt64  `db:"proxy_port" json:"proxy_port"`
	ProxyUsername     NullString `db:"proxy_username" json:"proxy_username"`
	ProxyPassword     NullString `db:"proxy_password" json:"proxy_password"` // Encrypted
	MessagesSent      int        `db:"messages_sent" json:"messages_sent"`
	LastUsedAt        NullTime   `db:"last_used_at" json:"last_used_at"`
	LastLoginAt       NullTime   `db:"last_login_at" json:"last_login_at"`
	ErrorMessage      NullString `db:"error_message" json:"error_message,omitempty"`
	PhoneCodeHash     NullString `db:"phone_code_hash" json:"-"`
	LoginCodeSentAt   NullTime   `db:"login_code_sent_at" json:"login_code_sent_at"`
	TwoFactorRequired bool       `db:"two_factor_required" json:"two_factor_required"`
	TwoFactorHint     NullString `db:"two_factor_hint" json:"two_factor_hint"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
}

// Template represents a message template
type Template struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	Name              string     `db:"name" json:"name"`
	Content           string     `db:"content" json:"content"`
	Variables         []string   `db:"variables" json:"variables"`                       // JSON array of variable names
	MediaType         NullString `db:"media_type" json:"media_type"`                     // photo, video, document, album
	MediaUrls         []string   `db:"media_urls" json:"media_urls"`                     // JSON array of media URLs
	CopyFromChatID    NullString `db:"copy_from_chat_id" json:"copy_from_chat_id"`       // Source chat for copying
	CopyFromMessageID NullInt64  `db:"copy_from_message_id" json:"copy_from_message_id"` // Source message ID
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
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
	ID              uuid.UUID  `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	AccountID       uuid.UUID  `db:"account_id" json:"account_id"`
	TemplateID      uuid.UUID  `db:"template_id" json:"template_id"`
	ChannelIDs      []string   `db:"channel_ids" json:"channel_ids"` // JSON array of channel UUIDs
	CronExpr        string     `db:"cron_expr" json:"cron_expr"`
	Timezone        string     `db:"timezone" json:"timezone"`                   // e.g., "Europe/Moscow"
	DayFilter       NullString `db:"day_filter" json:"day_filter"`               // all, weekdays, weekends, custom
	CustomDays      []int      `db:"custom_days" json:"custom_days"`             // [1,3,5] for Mon,Wed,Fri (0=Sunday)
	DelayMinSeconds int        `db:"delay_min_seconds" json:"delay_min_seconds"` // Min delay between messages
	DelayMaxSeconds int        `db:"delay_max_seconds" json:"delay_max_seconds"` // Max delay between messages
	LoadBalance     bool       `db:"load_balance" json:"load_balance"`           // Use account rotation
	Status          string     `db:"status" json:"status"`                       // active, paused, completed
	NextRunAt       NullTime   `db:"next_run_at" json:"next_run_at"`
	LastRunAt       NullTime   `db:"last_run_at" json:"last_run_at"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

// JobLog represents execution history
type JobLog struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	ScheduleID uuid.UUID  `db:"schedule_id" json:"schedule_id"`
	Status     string     `db:"status" json:"status"` // success, failed, retry
	Message    NullString `db:"message" json:"message,omitempty"`
	Error      NullString `db:"error" json:"error,omitempty"`
	ExecutedAt time.Time  `db:"executed_at" json:"executed_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

// User represents admin user for authentication
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// Setting represents a system configuration setting
type Setting struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Key             string     `db:"key" json:"key"`
	Value           string     `db:"value" json:"value"`
	Encrypted       bool       `db:"encrypted" json:"encrypted"`
	Category        string     `db:"category" json:"category"`
	Description     NullString `db:"description" json:"description,omitempty"`
	Editable        bool       `db:"editable" json:"editable"`
	RequiresRestart bool       `db:"requires_restart" json:"requires_restart"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	UpdatedBy       *uuid.UUID `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}
