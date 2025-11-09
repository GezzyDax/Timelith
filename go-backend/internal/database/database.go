package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Connect(databaseURL string) (*DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) RunMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS accounts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			phone VARCHAR(20) UNIQUE NOT NULL,
			session_data BYTEA,
			status VARCHAR(50) NOT NULL DEFAULT 'inactive',
			last_login_at TIMESTAMP,
			error_message TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS templates (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			variables JSONB DEFAULT '[]',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS channels (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			chat_id VARCHAR(255) NOT NULL,
			type VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS schedules (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
			template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
			channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
			cron_expr VARCHAR(255) NOT NULL,
			timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			next_run_at TIMESTAMP,
			last_run_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS job_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
			status VARCHAR(50) NOT NULL,
			message TEXT,
			error TEXT,
			executed_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			key VARCHAR(255) UNIQUE NOT NULL,
			value TEXT NOT NULL,
			encrypted BOOLEAN DEFAULT false,
			category VARCHAR(100) NOT NULL,
			description TEXT,
			editable BOOLEAN DEFAULT true,
			requires_restart BOOLEAN DEFAULT false,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_settings_category ON settings(category)`,
		`CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key)`,
		`CREATE INDEX IF NOT EXISTS idx_settings_editable ON settings(editable)`,
		`INSERT INTO settings (key, value, encrypted, category, description, editable, requires_restart)
		 VALUES
			('setup_completed', 'false', false, 'system', 'Setup wizard completion flag', false, false),
			('server_port', '8080', false, 'system', 'HTTP server port', true, true),
			('environment', 'production', false, 'system', 'Application environment', true, false),
			('log_level', 'info', false, 'system', 'Logging level', true, false)
		 ON CONFLICT (key) DO NOTHING`,
		`CREATE INDEX IF NOT EXISTS idx_accounts_status ON accounts(status)`,
		`CREATE INDEX IF NOT EXISTS idx_schedules_status ON schedules(status)`,
		`CREATE INDEX IF NOT EXISTS idx_schedules_next_run ON schedules(next_run_at)`,
		`CREATE INDEX IF NOT EXISTS idx_job_logs_schedule ON job_logs(schedule_id)`,
		`CREATE INDEX IF NOT EXISTS idx_job_logs_executed_at ON job_logs(executed_at)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
