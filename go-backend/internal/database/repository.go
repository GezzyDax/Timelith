package database

import (
	"database/sql"
	"fmt"

	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/google/uuid"
)

// Account Repository

func (db *DB) CreateAccount(account *models.Account) error {
	query := `INSERT INTO accounts (id, phone, session_data, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, NOW(), NOW())
			  RETURNING id, created_at, updated_at`

	account.ID = uuid.New()
	return db.QueryRow(query, account.ID, account.Phone, account.SessionData, account.Status).
		Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
}

func (db *DB) GetAccount(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	query := `SELECT * FROM accounts WHERE id = $1`
	err := db.Get(&account, query, id)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (db *DB) GetAccountByPhone(phone string) (*models.Account, error) {
	var account models.Account
	query := `SELECT * FROM accounts WHERE phone = $1`
	err := db.Get(&account, query, phone)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (db *DB) ListAccounts() ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT * FROM accounts ORDER BY created_at DESC`
	err := db.Select(&accounts, query)
	return accounts, err
}

func (db *DB) UpdateAccount(account *models.Account) error {
	query := `UPDATE accounts
			  SET session_data = $1, status = $2, last_login_at = $3,
			      error_message = $4, updated_at = NOW()
			  WHERE id = $5`

	_, err := db.Exec(query, account.SessionData, account.Status,
		account.LastLoginAt, account.ErrorMessage, account.ID)
	return err
}

func (db *DB) DeleteAccount(id uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Template Repository

func (db *DB) CreateTemplate(template *models.Template) error {
	query := `INSERT INTO templates (id, name, content, variables, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, NOW(), NOW())
			  RETURNING id, created_at, updated_at`

	template.ID = uuid.New()
	return db.QueryRow(query, template.ID, template.Name, template.Content, template.Variables).
		Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
}

func (db *DB) GetTemplate(id uuid.UUID) (*models.Template, error) {
	var template models.Template
	query := `SELECT * FROM templates WHERE id = $1`
	err := db.Get(&template, query, id)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (db *DB) ListTemplates() ([]models.Template, error) {
	var templates []models.Template
	query := `SELECT * FROM templates ORDER BY created_at DESC`
	err := db.Select(&templates, query)
	return templates, err
}

func (db *DB) UpdateTemplate(template *models.Template) error {
	query := `UPDATE templates
			  SET name = $1, content = $2, variables = $3, updated_at = NOW()
			  WHERE id = $4`

	_, err := db.Exec(query, template.Name, template.Content, template.Variables, template.ID)
	return err
}

func (db *DB) DeleteTemplate(id uuid.UUID) error {
	query := `DELETE FROM templates WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Channel Repository

func (db *DB) CreateChannel(channel *models.Channel) error {
	query := `INSERT INTO channels (id, name, chat_id, type, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, NOW(), NOW())
			  RETURNING id, created_at, updated_at`

	channel.ID = uuid.New()
	return db.QueryRow(query, channel.ID, channel.Name, channel.ChatID, channel.Type).
		Scan(&channel.ID, &channel.CreatedAt, &channel.UpdatedAt)
}

func (db *DB) GetChannel(id uuid.UUID) (*models.Channel, error) {
	var channel models.Channel
	query := `SELECT * FROM channels WHERE id = $1`
	err := db.Get(&channel, query, id)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (db *DB) ListChannels() ([]models.Channel, error) {
	var channels []models.Channel
	query := `SELECT * FROM channels ORDER BY created_at DESC`
	err := db.Select(&channels, query)
	return channels, err
}

func (db *DB) UpdateChannel(channel *models.Channel) error {
	query := `UPDATE channels
			  SET name = $1, chat_id = $2, type = $3, updated_at = NOW()
			  WHERE id = $4`

	_, err := db.Exec(query, channel.Name, channel.ChatID, channel.Type, channel.ID)
	return err
}

func (db *DB) DeleteChannel(id uuid.UUID) error {
	query := `DELETE FROM channels WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Schedule Repository

func (db *DB) CreateSchedule(schedule *models.Schedule) error {
	query := `INSERT INTO schedules (id, name, account_id, template_id, channel_id,
				cron_expr, timezone, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			  RETURNING id, created_at, updated_at`

	schedule.ID = uuid.New()
	return db.QueryRow(query, schedule.ID, schedule.Name, schedule.AccountID,
		schedule.TemplateID, schedule.ChannelID, schedule.CronExpr,
		schedule.Timezone, schedule.Status).
		Scan(&schedule.ID, &schedule.CreatedAt, &schedule.UpdatedAt)
}

func (db *DB) GetSchedule(id uuid.UUID) (*models.Schedule, error) {
	var schedule models.Schedule
	query := `SELECT * FROM schedules WHERE id = $1`
	err := db.Get(&schedule, query, id)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (db *DB) ListSchedules() ([]models.Schedule, error) {
	var schedules []models.Schedule
	query := `SELECT * FROM schedules ORDER BY created_at DESC`
	err := db.Select(&schedules, query)
	return schedules, err
}

func (db *DB) ListActiveSchedules() ([]models.Schedule, error) {
	var schedules []models.Schedule
	query := `SELECT * FROM schedules WHERE status = 'active' ORDER BY next_run_at ASC`
	err := db.Select(&schedules, query)
	return schedules, err
}

func (db *DB) UpdateSchedule(schedule *models.Schedule) error {
	query := `UPDATE schedules
			  SET name = $1, cron_expr = $2, timezone = $3, status = $4,
			      next_run_at = $5, last_run_at = $6, updated_at = NOW()
			  WHERE id = $7`

	_, err := db.Exec(query, schedule.Name, schedule.CronExpr, schedule.Timezone,
		schedule.Status, schedule.NextRunAt, schedule.LastRunAt, schedule.ID)
	return err
}

func (db *DB) DeleteSchedule(id uuid.UUID) error {
	query := `DELETE FROM schedules WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// JobLog Repository

func (db *DB) CreateJobLog(log *models.JobLog) error {
	query := `INSERT INTO job_logs (id, schedule_id, status, message, error, executed_at, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, NOW())
			  RETURNING id, created_at`

	log.ID = uuid.New()
	return db.QueryRow(query, log.ID, log.ScheduleID, log.Status,
		log.Message, log.Error, log.ExecutedAt).
		Scan(&log.ID, &log.CreatedAt)
}

func (db *DB) GetJobLogs(scheduleID uuid.UUID, limit int) ([]models.JobLog, error) {
	var logs []models.JobLog
	query := `SELECT * FROM job_logs WHERE schedule_id = $1
			  ORDER BY executed_at DESC LIMIT $2`
	err := db.Select(&logs, query, scheduleID, limit)
	return logs, err
}

func (db *DB) GetAllJobLogs(limit int) ([]models.JobLog, error) {
	var logs []models.JobLog
	query := `SELECT * FROM job_logs ORDER BY executed_at DESC LIMIT $1`
	err := db.Select(&logs, query, limit)
	return logs, err
}

// User Repository

func (db *DB) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, username, password_hash, created_at, updated_at)
			  VALUES ($1, $2, $3, NOW(), NOW())
			  RETURNING id, created_at, updated_at`

	user.ID = uuid.New()
	return db.QueryRow(query, user.ID, user.Username, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE username = $1`
	err := db.Get(&user, query, username)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) CountUsers() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := db.Get(&count, query)
	return count, err
}

func (db *DB) ListUsers() ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users ORDER BY created_at DESC`
	err := db.Select(&users, query)
	return users, err
}

func (db *DB) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) UpdateUser(user *models.User) error {
	query := `UPDATE users
			  SET username = $1, password_hash = $2, updated_at = NOW()
			  WHERE id = $3`

	_, err := db.Exec(query, user.Username, user.PasswordHash, user.ID)
	return err
}

func (db *DB) DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Settings Repository

func (db *DB) CreateSetting(setting *models.Setting) error {
	query := `INSERT INTO settings (id, key, value, encrypted, category, description, editable, requires_restart, updated_by, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
			  RETURNING id, created_at, updated_at`

	setting.ID = uuid.New()
	return db.QueryRow(query, setting.ID, setting.Key, setting.Value, setting.Encrypted,
		setting.Category, setting.Description, setting.Editable, setting.RequiresRestart, setting.UpdatedBy).
		Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

func (db *DB) GetSettingByKey(key string) (*models.Setting, error) {
	var setting models.Setting
	query := `SELECT * FROM settings WHERE key = $1`
	err := db.Get(&setting, query, key)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("setting not found: %s", key)
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (db *DB) GetAllSettings() ([]models.Setting, error) {
	var settings []models.Setting
	query := `SELECT * FROM settings ORDER BY category, key`
	err := db.Select(&settings, query)
	return settings, err
}

func (db *DB) GetSettingsByCategory(category string) ([]models.Setting, error) {
	var settings []models.Setting
	query := `SELECT * FROM settings WHERE category = $1 ORDER BY key`
	err := db.Select(&settings, query, category)
	return settings, err
}

func (db *DB) UpsertSetting(key, value string, encrypted bool, category string, updatedBy *uuid.UUID) error {
	query := `INSERT INTO settings (id, key, value, encrypted, category, updated_by, created_at, updated_at)
			  VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW(), NOW())
			  ON CONFLICT (key) DO UPDATE
			  SET value = $2, encrypted = $3, category = $4, updated_by = $5, updated_at = NOW()`

	_, err := db.Exec(query, key, value, encrypted, category, updatedBy)
	return err
}

func (db *DB) UpdateSetting(setting *models.Setting) error {
	query := `UPDATE settings
			  SET value = $1, encrypted = $2, category = $3, description = $4,
			      editable = $5, requires_restart = $6, updated_by = $7, updated_at = NOW()
			  WHERE key = $8`

	_, err := db.Exec(query, setting.Value, setting.Encrypted, setting.Category,
		setting.Description, setting.Editable, setting.RequiresRestart, setting.UpdatedBy, setting.Key)
	return err
}

func (db *DB) DeleteSetting(key string) error {
	query := `DELETE FROM settings WHERE key = $1`
	_, err := db.Exec(query, key)
	return err
}

func (db *DB) CountSettings() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM settings`
	err := db.Get(&count, query)
	return count, err
}
