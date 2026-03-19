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

// @title Tokkatot API
// @version 2.0
// @description Backend API for Tokkatot Agri-Tech Platform
// @contact.name Tokkatot Support
// @contact.email info@tokkatot.com
// @license.name MIT
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("✅ Configuration loaded - Environment: %s", cfg.Environment)

	// Initialize PostgreSQL database
	_, err := database.InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}
	defer database.CloseDB()

	// Create schema
	if err := database.CreateSchema(); err != nil {
		log.Fatalf("❌ Failed to create database schema: %v", err)
	}

	// Seed initial admin if needed
	if err := database.SeedInitialAdmin(); err != nil {
		log.Printf("⚠️  Admin seeding warning: %v", err)
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
	absPath, _ := filepath.Abs(currentDir)
	if filepath.Base(absPath) == "middleware" {
		frontendPath = filepath.Join(filepath.Dir(absPath), "frontend")
	} else {
		frontendPath = filepath.Join(absPath, "frontend")
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
	// ===== FRONTEND STATIC ROUTES =====
	// Serve static files
	app.Static("/assets", filepath.Join(frontendPath, "assets"))
	app.Static("/components", filepath.Join(frontendPath, "components"))
	app.Static("/css", filepath.Join(frontendPath, "css"))
	app.Static("/js", filepath.Join(frontendPath, "js"))
	app.Static("/", frontendPath) // Fallback for any other frontend assets (manifest, robots, etc.)

	// Static page routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "index.html"))
	})
	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "admin.html"))
	})
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "login.html"))
	})
	app.Get("/register", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "signup.html"))
	})
	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Redirect("/register")
	})
	app.Get("/index.html", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "index.html"))
	})
	app.Get("/profile", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "profile.html"))
	})
	app.Get("/settings", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "settings.html"))
	})
	app.Get("/disease-detection", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "disease-detection.html"))
	})
	app.Get("/monitoring", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "monitoring.html"))
	})
	app.Get("/schedules", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "schedules.html"))
	})
	app.Get("/alerts", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(frontendPath, "pages", "alerts.html"))
	})

	// v1 API Group
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

	// User session & activity endpoints
	protected.Get("/users/sessions", api.GetUserSessionsHandler)
	protected.Delete("/users/sessions/:session_id", api.RevokeUserSessionHandler)
	protected.Get("/users/activity-log", api.GetUserActivityLogHandler)

	// Alert subscription endpoints (user-scoped)
	protected.Post("/users/alert-subscriptions", api.CreateAlertSubscriptionHandler)
	protected.Get("/users/alert-subscriptions", api.GetAlertSubscriptionsHandler)
	protected.Put("/users/alert-subscriptions/:subscription_id", api.UpdateAlertSubscriptionHandler)
	protected.Delete("/users/alert-subscriptions/:subscription_id", api.DeleteAlertSubscriptionHandler)

	// Farm management endpoints
	protected.Get("/farms", api.ListFarmsHandler)
	protected.Post("/farms", api.CreateFarmHandler)
	protected.Get("/farms/:farm_id", api.GetFarmHandler)
	protected.Put("/farms/:farm_id", api.UpdateFarmHandler)
	protected.Delete("/farms/:farm_id", api.DeleteFarmHandler)

	// Farm member endpoints
	protected.Get("/farms/:farm_id/members", api.GetFarmMembersHandler)
	protected.Post("/farms/:farm_id/members", api.InviteFarmMemberHandler)
	protected.Put("/farms/:farm_id/members/:user_id", api.UpdateFarmMemberRoleHandler)
	protected.Delete("/farms/:farm_id/members/:user_id", api.RemoveFarmMemberHandler)

	// Coop management endpoints
	protected.Get("/farms/:farm_id/coops", api.ListCoopsHandler)
	protected.Post("/farms/:farm_id/coops", api.CreateCoopHandler)
	protected.Get("/farms/:farm_id/coops/:coop_id", api.GetCoopHandler)
	protected.Put("/farms/:farm_id/coops/:coop_id", api.UpdateCoopHandler)
	protected.Delete("/farms/:farm_id/coops/:coop_id", api.DeleteCoopHandler)
	protected.Get("/farms/:farm_id/coops/:coop_id/temperature-timeline", api.TemperatureTimelineHandler)

	// Device management endpoints
	protected.Get("/farms/:farm_id/devices", api.ListDevicesHandler)
	protected.Post("/farms/:farm_id/devices", api.AddDeviceHandler)
	protected.Get("/farms/:farm_id/devices/:device_id", api.GetDeviceHandler)
	protected.Put("/farms/:farm_id/devices/:device_id", api.UpdateDeviceHandler)
	protected.Delete("/farms/:farm_id/devices/:device_id", api.DeleteDeviceHandler)

	// Device advanced endpoints
	protected.Get("/farms/:farm_id/devices/:device_id/history", api.GetDeviceHistoryHandler)
	protected.Get("/farms/:farm_id/devices/:device_id/status", api.GetDeviceStatusHandler)
	protected.Get("/farms/:farm_id/devices/:device_id/config", api.GetDeviceConfigHandler)
	protected.Put("/farms/:farm_id/devices/:device_id/config", api.UpdateDeviceConfigHandler)
	protected.Post("/farms/:farm_id/devices/:device_id/calibrate", api.CalibrateDeviceHandler)

	// Device command endpoints
	protected.Post("/farms/:farm_id/devices/:device_id/commands", api.SendDeviceCommandHandler)
	protected.Get("/farms/:farm_id/devices/:device_id/commands/:command_id", api.GetDeviceCommandStatusHandler)
	protected.Get("/farms/:farm_id/devices/:device_id/commands", api.ListDeviceCommandsHandler)
	protected.Delete("/farms/:farm_id/devices/:device_id/commands/:command_id", api.CancelCommandHandler)

	// Farm-level device control
	protected.Get("/farms/:farm_id/commands", api.GetFarmCommandHistoryHandler)
	protected.Post("/farms/:farm_id/emergency-stop", api.EmergencyStopHandler)
	protected.Post("/farms/:farm_id/devices/batch-command", api.BatchDeviceCommandHandler)

	// Schedule management endpoints
	protected.Post("/farms/:farm_id/schedules", api.CreateScheduleHandler)
	protected.Get("/farms/:farm_id/schedules", api.ListSchedulesHandler)
	protected.Get("/farms/:farm_id/schedules/:schedule_id", api.GetScheduleHandler)
	protected.Put("/farms/:farm_id/schedules/:schedule_id", api.UpdateScheduleHandler)
	protected.Delete("/farms/:farm_id/schedules/:schedule_id", api.DeleteScheduleHandler)
	protected.Get("/farms/:farm_id/schedules/:schedule_id/executions", api.GetScheduleExecutionHistoryHandler)
	protected.Post("/farms/:farm_id/schedules/:schedule_id/execute-now", api.ExecuteScheduleNowHandler)

	// Alert endpoints
	protected.Get("/farms/:farm_id/alerts/history", api.GetAlertHistoryHandler)
	protected.Get("/farms/:farm_id/alerts", api.GetFarmAlertsHandler)
	protected.Get("/farms/:farm_id/alerts/:alert_id", api.GetAlertHandler)
	protected.Put("/farms/:farm_id/alerts/:alert_id/acknowledge", api.AcknowledgeAlertHandler)

	// Analytics & reporting endpoints
	protected.Get("/farms/:farm_id/dashboard", api.GetFarmDashboardHandler)
	protected.Get("/farms/:farm_id/reports/device-metrics", api.GetDeviceMetricsReportHandler)
	protected.Get("/farms/:farm_id/reports/device-usage", api.GetDeviceUsageReportHandler)
	protected.Get("/farms/:farm_id/reports/farm-performance", api.GetFarmPerformanceReportHandler)
	protected.Get("/farms/:farm_id/reports/export", api.ExportReportHandler)
	protected.Get("/farms/:farm_id/events", api.GetFarmEventLogHandler)

	// WebSocket for real-time updates (requires authentication)
	protected.Get("/ws", api.WebSocketUpgradeHandler)
	protected.Get("/ws/stats", api.GetWebSocketStatsHandler)

	// Device heartbeat (for IoT devices - no AuthMiddleware, uses device key)
	v1.Post("/devices/:hardware_id/heartbeat", api.UpdateDeviceHeartbeatHandler)

	// ===== ADMIN ROUTES (role="admin" required) =====
	admin := v1.Group("/admin")
	admin.Use(api.AuthMiddleware, api.AdminMiddleware)
	admin.Get("/stats", api.GetAdminStatsHandler)
	admin.Get("/farmers", api.ListFarmersHandler)
	admin.Post("/farmers", api.RegisterFarmerHandler)
	admin.Delete("/farmers/:user_id", api.DeactivateFarmerHandler)
	admin.Get("/farmers/:user_id/profile", api.GetFarmerProfileHandler)
	admin.Get("/viewers", api.ListViewersHandler)
	admin.Get("/reg-keys", api.ListRegKeysHandler)
	admin.Post("/reg-keys", api.RegisterFarmerHandler)
	admin.Put("/profile", api.UpdateAdminProfileHandler)


	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		// If requesting HTML (browser), show 404 page
		if c.Accepts("text/html") != "" {
			return c.Status(fiber.StatusNotFound).SendFile(filepath.Join(frontendPath, "pages", "404.html"))
		}

		// Otherwise return JSON
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "not_found",
				"message": "Endpoint not found",
			},
		})
	})
}
