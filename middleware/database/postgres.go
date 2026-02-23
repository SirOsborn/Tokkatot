package database

import (
	"database/sql"
	"fmt"
	"log"

	"middleware/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection pool
func InitDB() (*sql.DB, error) {
	dbURL := config.GetDatabaseURL()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool for performance (low-end device optimization)
	db.SetMaxOpenConns(25)        // Max 25 concurrent connections
	db.SetMaxIdleConns(5)         // Keep 5 connections ready
	db.SetConnMaxLifetime(5 * 60) // Recycle connections every 5 minutes

	log.Println("✅ Database connection established")
	DB = db
	return db, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// CreateSchema creates all necessary tables
func CreateSchema() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	schema := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email TEXT UNIQUE,
		phone TEXT UNIQUE,
		phone_country_code VARCHAR(5),
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		language VARCHAR(10) DEFAULT 'km',
		timezone VARCHAR(40) DEFAULT 'Asia/Phnom_Penh',
		avatar_url TEXT,
		is_active BOOLEAN DEFAULT true,
		contact_verified BOOLEAN DEFAULT false,
		verification_type VARCHAR(10),
		last_login TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Farms table
	CREATE TABLE IF NOT EXISTS farms (
		id UUID PRIMARY KEY,
		owner_id UUID NOT NULL REFERENCES users(id),
		name TEXT NOT NULL,
		location TEXT,
		timezone VARCHAR(40) DEFAULT 'Asia/Phnom_Penh',
		latitude DECIMAL(10,8),
		longitude DECIMAL(11,8),
		description TEXT,
		image_url TEXT,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Coops table (chicken houses within a farm)
	CREATE TABLE IF NOT EXISTS coops (
		id UUID PRIMARY KEY,
		farm_id UUID NOT NULL REFERENCES farms(id),
		number INTEGER NOT NULL,
		name TEXT NOT NULL,
		capacity INTEGER,
		current_count INTEGER,
		chicken_type VARCHAR(20),
		main_device_id UUID,
		description TEXT,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(farm_id, number)
	);

	-- Farm users (membership table)
	CREATE TABLE IF NOT EXISTS farm_users (
		id UUID PRIMARY KEY,
		farm_id UUID NOT NULL REFERENCES farms(id),
		user_id UUID NOT NULL REFERENCES users(id),
		role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'manager', 'viewer')),
		invited_by UUID NOT NULL REFERENCES users(id),
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(farm_id, user_id)
	);

	-- Devices table
	CREATE TABLE IF NOT EXISTS devices (
		id UUID PRIMARY KEY,
		farm_id UUID NOT NULL REFERENCES farms(id),
		coop_id UUID REFERENCES coops(id),
		device_id VARCHAR(50) NOT NULL UNIQUE,
		name TEXT NOT NULL,
		type VARCHAR(50) NOT NULL CHECK (type IN ('gpio', 'relay', 'pwm', 'adc', 'servo', 'sensor')),
		model TEXT,
		is_main_controller BOOLEAN DEFAULT false,
		firmware_version VARCHAR(20) NOT NULL,
		hardware_id TEXT NOT NULL UNIQUE,
		location TEXT,
		is_active BOOLEAN DEFAULT true,
		is_online BOOLEAN DEFAULT false,
		last_heartbeat TIMESTAMP,
		last_command_status VARCHAR(50),
		last_command_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Device commands
	CREATE TABLE IF NOT EXISTS device_commands (
		id UUID PRIMARY KEY,
		device_id UUID NOT NULL REFERENCES devices(id),
		farm_id UUID NOT NULL REFERENCES farms(id),
		coop_id UUID REFERENCES coops(id),
		issued_by UUID NOT NULL REFERENCES users(id),
		command_type VARCHAR(50) NOT NULL,
		command_value TEXT,
		status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'timeout')),
		response TEXT,
		issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		executed_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Schedules
	CREATE TABLE IF NOT EXISTS schedules (
		id UUID PRIMARY KEY,
		farm_id UUID NOT NULL REFERENCES farms(id),
		coop_id UUID REFERENCES coops(id),
		device_id UUID NOT NULL REFERENCES devices(id),
		name TEXT NOT NULL,
		schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('time_based', 'duration_based', 'condition_based')),
		cron_expression TEXT,
		on_duration INTEGER,
		off_duration INTEGER,
		condition_json JSONB,
		action VARCHAR(20) NOT NULL CHECK (action IN ('on', 'off', 'set_value')),
		action_value TEXT,
		action_duration INTEGER,
		action_sequence JSONB,
		priority INTEGER DEFAULT 0,
		is_active BOOLEAN DEFAULT true,
		next_execution TIMESTAMP,
		last_execution TIMESTAMP,
		execution_count INTEGER DEFAULT 0,
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Schedule executions (execution history)
	CREATE TABLE IF NOT EXISTS schedule_executions (
		id UUID PRIMARY KEY,
		schedule_id UUID NOT NULL REFERENCES schedules(id),
		device_id UUID NOT NULL REFERENCES devices(id),
		scheduled_time TIMESTAMP NOT NULL,
		actual_execution_time TIMESTAMP,
		status VARCHAR(20) NOT NULL CHECK (status IN ('executed', 'failed', 'skipped')),
		execution_duration_ms INTEGER,
		device_response JSONB,
		error_message TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Event logs (audit trail)
	CREATE TABLE IF NOT EXISTS event_logs (
		id UUID PRIMARY KEY,
		farm_id UUID NOT NULL REFERENCES farms(id),
		user_id UUID NOT NULL REFERENCES users(id),
		event_type VARCHAR(50) NOT NULL,
		resource_id UUID,
		old_value JSONB,
		new_value JSONB,
		ip_address VARCHAR(45),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Registration keys (for on-site account creation by staff)
	CREATE TABLE IF NOT EXISTS registration_keys (
		id UUID PRIMARY KEY,
		key_code VARCHAR(50) UNIQUE NOT NULL,
		farm_name VARCHAR(255),
		farm_location TEXT,
		customer_name VARCHAR(255),
		customer_phone VARCHAR(20),
		is_used BOOLEAN DEFAULT false,
		used_by_user_id UUID REFERENCES users(id),
		used_at TIMESTAMP,
		expires_at TIMESTAMP,
		created_by VARCHAR(100),
		notes TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- ===== PERFORMANCE INDEXES =====

	-- Users indexes
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
	CREATE INDEX IF NOT EXISTS idx_users_is_active_created ON users(is_active, created_at);

	-- Farms indexes
	CREATE INDEX IF NOT EXISTS idx_farms_owner_id ON farms(owner_id);
	CREATE INDEX IF NOT EXISTS idx_farms_is_active ON farms(is_active);
	CREATE INDEX IF NOT EXISTS idx_farms_created_at ON farms(created_at DESC);

	-- Coops indexes
	CREATE INDEX IF NOT EXISTS idx_coops_farm_id ON coops(farm_id);
	CREATE INDEX IF NOT EXISTS idx_coops_farm_number ON coops(farm_id, number);
	CREATE INDEX IF NOT EXISTS idx_coops_is_active ON coops(is_active);
	CREATE INDEX IF NOT EXISTS idx_coops_main_device ON coops(main_device_id);

	-- Farm users indexes
	CREATE INDEX IF NOT EXISTS idx_farm_users_farm_id ON farm_users(farm_id);
	CREATE INDEX IF NOT EXISTS idx_farm_users_user_id ON farm_users(user_id);
	CREATE INDEX IF NOT EXISTS idx_farm_users_active ON farm_users(is_active);

	-- Devices indexes
	CREATE INDEX IF NOT EXISTS idx_devices_farm_id_active ON devices(farm_id, is_active);
	CREATE INDEX IF NOT EXISTS idx_devices_coop_id ON devices(coop_id);
	CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
	CREATE INDEX IF NOT EXISTS idx_devices_is_online ON devices(is_online);
	CREATE INDEX IF NOT EXISTS idx_devices_main_controller ON devices(coop_id, is_main_controller);

	-- Device commands indexes
	CREATE INDEX IF NOT EXISTS idx_device_commands_device_id ON device_commands(device_id);
	CREATE INDEX IF NOT EXISTS idx_device_commands_farm_id ON device_commands(farm_id);
	CREATE INDEX IF NOT EXISTS idx_device_commands_coop_id ON device_commands(coop_id);
	CREATE INDEX IF NOT EXISTS idx_device_commands_status ON device_commands(status);
	CREATE INDEX IF NOT EXISTS idx_device_commands_created ON device_commands(created_at DESC);

	-- Schedules indexes
	CREATE INDEX IF NOT EXISTS idx_schedules_farm_device ON schedules(farm_id, device_id);
	CREATE INDEX IF NOT EXISTS idx_schedules_coop_id ON schedules(coop_id);
	CREATE INDEX IF NOT EXISTS idx_schedules_is_active ON schedules(is_active);
	CREATE INDEX IF NOT EXISTS idx_schedules_next_execution ON schedules(next_execution);

	-- Schedule executions indexes
	CREATE INDEX IF NOT EXISTS idx_schedule_executions_schedule_id ON schedule_executions(schedule_id);
	CREATE INDEX IF NOT EXISTS idx_schedule_executions_time ON schedule_executions(scheduled_time DESC);
	CREATE INDEX IF NOT EXISTS idx_schedule_executions_status ON schedule_executions(status);

	-- Event logs indexes
	CREATE INDEX IF NOT EXISTS idx_event_logs_farm_user ON event_logs(farm_id, user_id);
	CREATE INDEX IF NOT EXISTS idx_event_logs_created ON event_logs(created_at DESC);

	-- Registration keys indexes
	CREATE INDEX IF NOT EXISTS idx_registration_keys_code ON registration_keys(key_code);
	CREATE INDEX IF NOT EXISTS idx_registration_keys_is_used ON registration_keys(is_used);
	CREATE INDEX IF NOT EXISTS idx_registration_keys_expires ON registration_keys(expires_at);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("✅ Database schema created/updated")
	return nil
}
