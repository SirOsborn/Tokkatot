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

	// Start WebSocket hub (for real-time updates)
	go api.WSHub.RunHub()
	log.Println("✅ WebSocket hub started")

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

	// User profile endpoints
	protected.Get("/users/me", api.GetCurrentUserHandler)
	protected.Put("/users/me", api.UpdateProfileHandler)
	protected.Post("/users/me/change-password", api.ChangePasswordHandler)

	// Farm management endpoints
	protected.Get("/farms", api.ListFarmsHandler)
	protected.Post("/farms", api.CreateFarmHandler)
	protected.Get("/farms/:farm_id", api.GetFarmHandler)
	protected.Put("/farms/:farm_id", api.UpdateFarmHandler)
	protected.Delete("/farms/:farm_id", api.DeleteFarmHandler)

	// Coop management endpoints
	protected.Get("/farms/:farm_id/coops", api.ListCoopsHandler)
	protected.Post("/farms/:farm_id/coops", api.CreateCoopHandler)
	protected.Get("/farms/:farm_id/coops/:coop_id", api.GetCoopHandler)
	protected.Put("/farms/:farm_id/coops/:coop_id", api.UpdateCoopHandler)
	protected.Delete("/farms/:farm_id/coops/:coop_id", api.DeleteCoopHandler)

	// Device management endpoints
	protected.Get("/farms/:farm_id/devices", api.ListDevicesHandler)
	protected.Get("/farms/:farm_id/devices/:device_id", api.GetDeviceHandler)
	protected.Post("/farms/:farm_id/devices/:device_id/commands", api.SendDeviceCommandHandler)
	protected.Get("/farms/:farm_id/devices/:device_id/commands/:command_id", api.GetDeviceCommandStatusHandler)
	protected.Get("/farms/:farm_id/devices/:device_id/commands", api.ListDeviceCommandsHandler)

	// Schedule management endpoints
	protected.Post("/farms/:farm_id/schedules", api.CreateScheduleHandler)
	protected.Get("/farms/:farm_id/schedules", api.ListSchedulesHandler)
	protected.Get("/farms/:farm_id/schedules/:schedule_id", api.GetScheduleHandler)
	protected.Put("/farms/:farm_id/schedules/:schedule_id", api.UpdateScheduleHandler)
	protected.Delete("/farms/:farm_id/schedules/:schedule_id", api.DeleteScheduleHandler)
	protected.Get("/farms/:farm_id/schedules/:schedule_id/executions", api.GetScheduleExecutionHistoryHandler)
	protected.Post("/farms/:farm_id/schedules/:schedule_id/execute-now", api.ExecuteScheduleNowHandler)

	// WebSocket for real-time updates (requires authentication)
	protected.Get("/ws", api.WebSocketUpgradeHandler)
	protected.Get("/ws/stats", api.GetWebSocketStatsHandler)

	// Device heartbeat (for IoT devices - no AuthMiddleware, uses device key)
	v1.Post("/devices/:hardware_id/heartbeat", api.UpdateDeviceHeartbeatHandler)

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
