package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server
	ServerPort string
	ServerHost string

	// JWT
	JWTSecret string

	// Initial Admin
	InitialAdminEmail    string
	InitialAdminPassword string

	// Test/Local Farmer (seeded in development only)
	TestFarmerEmail    string
	TestFarmerPassword string
	DemoFarmName       string
	DemoCoopName       string

	// Environment
	Environment string

	// Telemetry retention
	TelemetryRetentionDays int

	// Web Push (VAPID)
	VapidPublicKey  string
	VapidPrivateKey string
	VapidSubject    string
}

var AppConfig *Config

func init() {
	LoadConfig()
}

func LoadConfig() *Config {
	// Load .env file (Overload forces .env values to override system environment variables)
	if err := godotenv.Overload(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	config := &Config{
		// Database configuration
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Server configuration
		ServerPort: getEnv("SERVER_PORT", "3000"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		// JWT configuration
		JWTSecret: getEnv("JWT_SECRET", ""),

		// Initial Admin (Must be set in .env)
		InitialAdminEmail:    getEnv("INITIAL_ADMIN_EMAIL", ""),
		InitialAdminPassword: getEnv("INITIAL_ADMIN_PASSWORD", ""),

		// Test/Staging Farmer (never set in production .env)
		TestFarmerEmail:    getEnv("TEST_FARMER_EMAIL", ""),
		TestFarmerPassword: getEnv("TEST_FARMER_PASSWORD", ""),
		DemoFarmName:       getEnv("DEMO_FARM_NAME", "Demo Farm"),
		DemoCoopName:       getEnv("DEMO_COOP_NAME", "Coop 1"),

		// Environment
		Environment: getEnv("ENVIRONMENT", "development"),

		// Telemetry retention (days)
		TelemetryRetentionDays: getEnvInt("TELEMETRY_RETENTION_DAYS", 7),

		// Web Push Configuration
		VapidPublicKey:  getEnv("VAPID_PUBLIC_KEY", ""),
		VapidPrivateKey: getEnv("VAPID_PRIVATE_KEY", ""),
		VapidSubject:    getEnv("VAPID_SUBJECT", "mailto:admin@tokkatot.com"),
	}

	// Validate required fields
	if config.JWTSecret == "" || config.InitialAdminEmail == "" || config.InitialAdminPassword == "" || config.DBHost == "" || config.DBName == "" {
		if config.Environment == "production" {
			log.Fatal("CRITICAL: All database, JWT, and admin variables MUST be set in .env for production.")
		} else {
			log.Println("⚠️  Warning: Some environment variables (DB_HOST, JWT_SECRET, etc.) are missing. The app may fail to connect.")
		}
	}

	AppConfig = config
	return config
}

func GetDatabaseURL() string {
	config := AppConfig
	if config == nil {
		LoadConfig()
		config = AppConfig
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
