//go:build ignore

package main

import (
	"log"
	"middleware/config"
	"middleware/database"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Change working directory to middleware to load .env correctly if run from root
	wd, _ := os.Getwd()
	if filepath.Base(wd) != "middleware" {
		if _, err := os.Stat("middleware"); err == nil {
			os.Chdir("middleware")
		}
	}

	// 1. Load configuration from .env
	cfg := config.LoadConfig()
	log.Printf("✅ Loaded config for environment: %s", cfg.Environment)

	// 2. Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer database.CloseDB()

	// 3. Drop all tables (Clean slate)
	log.Println("🔄 Dropping all tables...")
	if err := database.DropAllTables(); err != nil {
		log.Fatal("❌ Failed to drop tables:", err)
	}

	// 4. Create fresh schema from single source of truth
	log.Println("🏗️ Creating fresh schema...")
	if err := database.CreateSchema(); err != nil {
		log.Fatal("❌ Failed to create fresh schema:", err)
	}

	// ── 5. Seed Initial Data (Super Admin & Test Data) ─────────────────────────
	log.Println("🌱 Seeding initial data...")

	// Super Admin (Configurable from .env)
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.InitialAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash password:", err)
	}
	
	adminID := "00000000-0000-0000-0000-000000000001"
	_, err = db.Exec(`
		INSERT INTO users (id, name, email, phone, password_hash, is_active, full_name)
		VALUES ($1, 'Admin', $2, 'N/A', $3, true, 'Tokkatot Admin')
	`, adminID, cfg.InitialAdminEmail, string(hash))
	if err != nil {
		log.Fatal("seed super admin:", err)
	}

	// Also add to admins table (Tokkatot internal staff)
	_, err = db.Exec(`
		INSERT INTO admins (id, name, email, phone, password_hash, is_active)
		VALUES ($1, 'Admin', $2, 'N/A', $3, true)
	`, adminID, cfg.InitialAdminEmail, string(hash))
	if err != nil {
		log.Fatal("seed admin table:", err)
	}

	log.Printf("✅ Database reset and seeded successfully with Admin: %s", cfg.InitialAdminEmail)
}
