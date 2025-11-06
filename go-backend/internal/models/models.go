package models

import (
	"time"
)

type TelegramAccount struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PhoneNumber    string    `gorm:"uniqueIndex;not null" json:"phone_number"`
	SessionData    string    `gorm:"size:10000" json:"session_data,omitempty"`
	Status         string    `gorm:"default:'pending';not null" json:"status"`
	FirstName      string    `json:"first_name,omitempty"`
	LastName       string    `json:"last_name,omitempty"`
	Username       string    `json:"username,omitempty"`
	TelegramUserID int64     `gorm:"uniqueIndex" json:"telegram_user_id,omitempty"`
	LastActiveAt   *time.Time `json:"last_active_at,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type MessageTemplate struct {
	ID                     uint   `gorm:"primaryKey" json:"id"`
	Name                   string `gorm:"not null" json:"name"`
	Content                string `gorm:"not null;type:text" json:"content"`
	MediaType              string `json:"media_type,omitempty"`
	MediaURL               string `json:"media_url,omitempty"`
	ParseModeEnabled       bool   `gorm:"default:false" json:"parse_mode_enabled"`
	ParseMode              string `json:"parse_mode,omitempty"`
	DisableWebPagePreview  bool   `gorm:"default:false" json:"disable_web_page_preview"`
	Buttons                string `gorm:"type:json" json:"buttons,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type Channel struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	ChannelType   string    `gorm:"not null" json:"channel_type"`
	TelegramID    int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username      string    `json:"username,omitempty"`
	Title         string    `json:"title,omitempty"`
	MembersCount  int       `json:"members_count,omitempty"`
	Description   string    `gorm:"type:text" json:"description,omitempty"`
	LastSyncedAt  *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Schedule struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Name               string    `gorm:"not null" json:"name"`
	TelegramAccountID  uint      `gorm:"not null" json:"telegram_account_id"`
	MessageTemplateID  uint      `gorm:"not null" json:"message_template_id"`
	ScheduleType       string    `gorm:"not null" json:"schedule_type"`
	CronExpression     string    `json:"cron_expression,omitempty"`
	IntervalMinutes    int       `json:"interval_minutes,omitempty"`
	ScheduledAt        *time.Time `json:"scheduled_at,omitempty"`
	Timezone           string    `gorm:"default:'UTC'" json:"timezone"`
	Active             bool      `gorm:"default:false" json:"active"`
	NextRunAt          *time.Time `json:"next_run_at,omitempty"`
	LastRunAt          *time.Time `json:"last_run_at,omitempty"`
	TotalRuns          int       `gorm:"default:0" json:"total_runs"`
	SuccessfulRuns     int       `gorm:"default:0" json:"successful_runs"`
	FailedRuns         int       `gorm:"default:0" json:"failed_runs"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	TelegramAccount    TelegramAccount `gorm:"foreignKey:TelegramAccountID" json:"telegram_account"`
	MessageTemplate    MessageTemplate `gorm:"foreignKey:MessageTemplateID" json:"message_template"`
	ScheduleChannels   []ScheduleChannel `gorm:"foreignKey:ScheduleID" json:"schedule_channels"`
}

type ScheduleChannel struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ScheduleID uint      `gorm:"not null" json:"schedule_id"`
	ChannelID  uint      `gorm:"not null" json:"channel_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Channel    Channel   `gorm:"foreignKey:ChannelID" json:"channel"`
}

type SendLog struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ScheduleID        uint      `gorm:"not null" json:"schedule_id"`
	TelegramAccountID uint      `gorm:"not null" json:"telegram_account_id"`
	ChannelID         uint      `gorm:"not null" json:"channel_id"`
	Status            string    `gorm:"not null" json:"status"`
	MessageContent    string    `gorm:"type:text" json:"message_content"`
	TelegramMessageID int64     `json:"telegram_message_id,omitempty"`
	ErrorMessage      string    `gorm:"type:text" json:"error_message,omitempty"`
	SentAt            *time.Time `json:"sent_at,omitempty"`
	RetryCount        int       `gorm:"default:0" json:"retry_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (TelegramAccount) TableName() string {
	return "telegram_accounts"
}

func (MessageTemplate) TableName() string {
	return "message_templates"
}

func (Channel) TableName() string {
	return "channels"
}

func (Schedule) TableName() string {
	return "schedules"
}

func (ScheduleChannel) TableName() string {
	return "schedule_channels"
}

func (SendLog) TableName() string {
	return "send_logs"
}
