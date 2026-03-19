package config

import (
	"fmt"
	"log"
	"os"

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

	// Environment
	Environment string
}

var AppConfig *Config

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

		// Environment
		Environment: getEnv("ENVIRONMENT", "development"),
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
