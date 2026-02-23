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
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "tokkatot"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"), // "disable" for local, "require" for production

		// Server configuration
		ServerPort: getEnv("SERVER_PORT", "3000"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		// JWT configuration
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),

		// Environment
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Validate required fields
	if config.JWTSecret == "your-secret-key-change-in-production" && config.Environment == "production" {
		log.Fatal("JWT_SECRET must be set in production environment")
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
