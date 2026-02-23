package main

import (
	"log"
	"os"
	"path/filepath"

	"middleware/api"
	"middleware/config"
	"middleware/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("✅ Configuration loaded - Environment: %s", cfg.Environment)

	// Initialize database (try PostgreSQL, fallback to SQLite for testing)
	_, err := database.InitDB()
	if err != nil {
		log.Printf("⚠️  PostgreSQL connection failed: %v", err)
		log.Println("🔄 Falling back to SQLite for testing...")

		_, sqliteErr := database.InitDBSQLite()
		if sqliteErr != nil {
			log.Fatalf("❌ Failed to initialize SQLite fallback: %v", sqliteErr)
		}

		// Create SQLite schema
		if err := database.CreateSchemaSQLite(); err != nil {
			log.Fatalf("❌ Failed to create SQLite schema: %v", err)
		}
	} else {
		defer database.CloseDB()

		// Create PostgreSQL schema
		if err := database.CreateSchema(); err != nil {
			log.Fatalf("❌ Failed to create database schema: %v", err)
		}
	}

	// Create Fiber app with optimized settings
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		Immutable:     true,
		BodyLimit:     10 * 1024 * 1024, // 10MB for image uploads
	})

	// Get frontend path
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get current working directory:", err)
	}

	var frontendPath string
	if filepath.Base(currentDir) == "middleware" {
		frontendPath = filepath.Join(filepath.Dir(currentDir), "frontend")
	} else {
		frontendPath = filepath.Join(currentDir, "frontend")
	}

	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		log.Printf("⚠️  Frontend directory not found at: %s (continuing without frontend)", frontendPath)
	}

	// Setup routes
	setupRoutes(app, frontendPath)

	// Start server
	log.Printf("✅ Server starting on %s:%s", cfg.ServerHost, cfg.ServerPort)
	if err := app.Listen(cfg.ServerHost + ":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

func setupRoutes(app *fiber.App, frontendPath string) {
	// ===== API ROUTES (v1) =====
	v1 := app.Group("/v1")

	// Authentication routes (no auth required)
	auth := v1.Group("/auth")
	auth.Post("/signup", api.SignupHandler)
	auth.Post("/login", api.LoginHandler)
	auth.Post("/verify", api.VerifyContactHandler) // Email/Phone verification
	auth.Post("/refresh", api.RefreshTokenHandler)
	auth.Post("/logout", api.LogoutHandler)
	auth.Post("/forgot-password", api.ForgotPasswordHandler)
	auth.Post("/reset-password", api.ResetPasswordHandler)

	// Protected routes (require authentication)
	protected := v1.Group("")
	protected.Use(api.AuthMiddleware)
	// Protected endpoints will be added here:
	// protected.Get("/users/me", api.GetCurrentUserHandler)
	// protected.Get("/farms", api.ListFarmsHandler)
	// etc.

	// ===== FRONTEND STATIC ROUTES =====
	// Serve static files
	app.Static("/assets", filepath.Join(frontendPath, "assets"))
	app.Static("/components", filepath.Join(frontendPath, "components"))
	app.Static("/css", filepath.Join(frontendPath, "css"))
	app.Static("/js", filepath.Join(frontendPath, "js"))

	// Static page routes
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "login.html"))
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "signup.html"))
	})

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "not_found",
				"message": "Endpoint not found",
			},
		})
	})
}
