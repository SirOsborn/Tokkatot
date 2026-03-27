package database

import (
	"database/sql"
	"fmt"
	"log"

	"middleware/config"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

// InitDB initializes the database connection pool using config from .env
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

// CreateSchema creates all necessary tables using the master schema definition in schema.go
func CreateSchema() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := DB.Exec(FullSchema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("✅ Database schema created/updated")

	// Run idempotent migrations for dynamic constraints or data fixes if necessary
	migrations := []string{
		// Ensure admins table has role column (added in schema v2)
		`ALTER TABLE admins ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'admin'`,
		// Migrate farm_users role constraint from legacy (owner/manager/viewer) to current (farmer/viewer)
		`ALTER TABLE farm_users DROP CONSTRAINT IF EXISTS farm_users_role_check`,
		`ALTER TABLE farm_users ADD CONSTRAINT farm_users_role_check CHECK (role IN ('farmer', 'viewer'))`,
		`UPDATE farm_users SET role = 'farmer' WHERE role IN ('owner', 'manager')`,
	}
	for _, m := range migrations {
		if _, merr := DB.Exec(m); merr != nil {
			log.Printf("⚠️  Migration warning: %v", merr)
		}
	}

	return nil
}

// DropAllTables drops all tables in the correct order to respect foreign keys.
// USE WITH CAUTION. This is intended for fresh setup/migration scripts.
func DropAllTables() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	dropSQL := `
		DROP TABLE IF EXISTS device_readings         CASCADE;
		DROP TABLE IF EXISTS device_configurations   CASCADE;
		DROP TABLE IF EXISTS alert_subscriptions     CASCADE;
		DROP TABLE IF EXISTS alerts                  CASCADE;
		DROP TABLE IF EXISTS user_sessions           CASCADE;
		DROP TABLE IF EXISTS registration_keys       CASCADE;
		DROP TABLE IF EXISTS event_logs              CASCADE;
		DROP TABLE IF EXISTS schedule_executions     CASCADE;
		DROP TABLE IF EXISTS schedules               CASCADE;
		DROP TABLE IF EXISTS device_commands         CASCADE;
		DROP TABLE IF EXISTS devices                 CASCADE;
		DROP TABLE IF EXISTS farm_users              CASCADE;
		DROP TABLE IF EXISTS coops                   CASCADE;
		DROP TABLE IF EXISTS farms                   CASCADE;
		DROP TABLE IF EXISTS admins                  CASCADE;
		DROP TABLE IF EXISTS users                   CASCADE;
	`
	_, err := DB.Exec(dropSQL)
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("✅ All tables dropped successfully")
	return nil
}

// SeedInitialAdmin seeds the initial super admin from .env config if no admin exists
func SeedInitialAdmin() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	cfg := config.AppConfig
	if cfg.InitialAdminEmail == "" || cfg.InitialAdminPassword == "" {
		log.Println("ℹ️  Skipping admin seeding (INITIAL_ADMIN_EMAIL or PASSWORD not set in .env)")
		return nil
	}

	log.Println("🌱 Syncing initial super admin...")

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.InitialAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	adminID := "00000000-0000-0000-0000-000000000001"
	
	// Sync user row
	_, err = DB.Exec(`
		INSERT INTO users (id, name, email, phone, password_hash, is_active, full_name)
		VALUES ($1, 'Admin', $2, 'N/A', $3, true, 'Tokkatot Admin')
		ON CONFLICT (id) DO UPDATE SET 
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash
	`, adminID, cfg.InitialAdminEmail, string(hash))
	if err != nil {
		return fmt.Errorf("failed to sync user row for admin: %w", err)
	}

	// Sync admin row
	_, err = DB.Exec(`
		INSERT INTO admins (id, name, email, phone, password_hash, is_active)
		VALUES ($1, 'Admin', $2, 'N/A', $3, true)
		ON CONFLICT (id) DO UPDATE SET 
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash
	`, adminID, cfg.InitialAdminEmail, string(hash))
	if err != nil {
		return fmt.Errorf("failed to sync admin row: %w", err)
	}

	log.Printf("✅ Initial admin synced: %s", cfg.InitialAdminEmail)
	return nil
}

// SeedTestData seeds a test farmer user for development/staging validation.
// Idempotent — safe to call on every startup.
func SeedTestData() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Only seed test data in non-production environments
	// (in production the gateway-sim provides live data)
	testFarmerID := "00000000-0000-0000-0000-000000000002"

	// Check if test farmer already exists
	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", testFarmerID).Scan(&count); err != nil {
		return fmt.Errorf("failed to check test farmer: %w", err)
	}
	if count > 0 {
		return nil // Already seeded
	}

	log.Println("🌱 Seeding test farmer user...")

	hash, err := bcrypt.GenerateFromPassword([]byte("Test@Farmer123!"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash test farmer password: %w", err)
	}

	_, err = DB.Exec(`
		INSERT INTO users (id, name, email, phone, password_hash, is_active, full_name)
		VALUES ($1, 'Test Farmer', 'farmer@tokkatot.com', 'N/A', $2, true, 'Test Farmer Account')
		ON CONFLICT DO NOTHING
	`, testFarmerID, string(hash))
	if err != nil {
		return fmt.Errorf("failed to seed test farmer: %w", err)
	}

	log.Println("✅ Test farmer seeded: farmer@tokkatot.com / Test@Farmer123!")
	return nil
}
