package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB is declared in postgres.go - shared across both implementations

// InitDBSQLite initializes SQLite for testing (fallback when PostgreSQL unavailable)
func InitDBSQLite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./tokkatot_test.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	// Configure connection
	db.SetMaxOpenConns(1) // SQLite works best with single connection
	db.SetMaxIdleConns(1)

	log.Println("✅ SQLite database connection established (testing mode)")
	DB = db
	return db, nil
}

// CreateSchemaSQLite creates schema for SQLite (slightly different from PostgreSQL)
func CreateSchemaSQLite() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	schema := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		phone TEXT UNIQUE,
		phone_country_code TEXT,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		language TEXT DEFAULT 'km',
		timezone TEXT DEFAULT 'Asia/Phnom_Penh',
		avatar_url TEXT,
		is_active INTEGER DEFAULT 1,
		contact_verified INTEGER DEFAULT 0,
		verification_type TEXT,
		last_login DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Farms table
	CREATE TABLE IF NOT EXISTS farms (
		id TEXT PRIMARY KEY,
		owner_id TEXT NOT NULL,
		name TEXT NOT NULL,
		location TEXT,
		timezone TEXT DEFAULT 'Asia/Phnom_Penh',
		latitude REAL,
		longitude REAL,
		description TEXT,
		image_url TEXT,
		is_active INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(id)
	);

	-- Farm users (membership table)
	CREATE TABLE IF NOT EXISTS farm_users (
		id TEXT PRIMARY KEY,
		farm_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		role TEXT NOT NULL CHECK (role IN ('owner', 'manager', 'viewer')),
		invited_by TEXT NOT NULL,
		is_active INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (farm_id) REFERENCES farms(id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE(farm_id, user_id)
	);

	-- Devices table
	CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		farm_id TEXT NOT NULL,
		device_id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('gpio', 'relay', 'pwm', 'adc', 'servo', 'sensor')),
		model TEXT,
		firmware_version TEXT NOT NULL,
		hardware_id TEXT NOT NULL UNIQUE,
		location TEXT,
		is_active INTEGER DEFAULT 1,
		is_online INTEGER DEFAULT 0,
		last_heartbeat DATETIME,
		last_command_status TEXT,
		last_command_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (farm_id) REFERENCES farms(id)
	);

	-- Device commands
	CREATE TABLE IF NOT EXISTS device_commands (
		id TEXT PRIMARY KEY,
		device_id TEXT NOT NULL,
		farm_id TEXT NOT NULL,
		issued_by TEXT NOT NULL,
		command_type TEXT NOT NULL,
		command_value TEXT,
		status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'timeout')),
		response TEXT,
		issued_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		executed_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (device_id) REFERENCES devices(id),
		FOREIGN KEY (farm_id) REFERENCES farms(id),
		FOREIGN KEY (issued_by) REFERENCES users(id)
	);

	-- Schedules
	CREATE TABLE IF NOT EXISTS schedules (
		id TEXT PRIMARY KEY,
		farm_id TEXT NOT NULL,
		device_id TEXT NOT NULL,
		name TEXT NOT NULL,
		rule TEXT NOT NULL,
		is_active INTEGER DEFAULT 1,
		created_by TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (farm_id) REFERENCES farms(id),
		FOREIGN KEY (device_id) REFERENCES devices(id),
		FOREIGN KEY (created_by) REFERENCES users(id)
	);

	-- Event logs (audit trail)
	CREATE TABLE IF NOT EXISTS event_logs (
		id TEXT PRIMARY KEY,
		farm_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		resource_id TEXT,
		old_value TEXT,
		new_value TEXT,
		ip_address TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (farm_id) REFERENCES farms(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
	CREATE INDEX IF NOT EXISTS idx_farms_owner_id ON farms(owner_id);
	CREATE INDEX IF NOT EXISTS idx_farm_users_farm_id ON farm_users(farm_id);
	CREATE INDEX IF NOT EXISTS idx_farm_users_user_id ON farm_users(user_id);
	CREATE INDEX IF NOT EXISTS idx_devices_farm_id ON devices(farm_id);
	CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
	CREATE INDEX IF NOT EXISTS idx_schedules_farm_device ON schedules(farm_id, device_id);
	CREATE INDEX IF NOT EXISTS idx_event_logs_farm_user ON event_logs(farm_id, user_id);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create SQLite schema: %w", err)
	}

	log.Println("✅ SQLite database schema created/updated")
	return nil
}
