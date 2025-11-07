package main

import (
	"log"
	"os"
	"path/filepath"

	"middleware/api"
	"middleware/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	// Check if we're in tokkatot directory, then look for .env in middleware subdirectory
	currentDir, _ := os.Getwd()
	var envPath string

	if filepath.Base(currentDir) == "middleware" {
		envPath = ".env"
	} else {
		envPath = "middleware/.env"
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("No .env file found at %s, using system environment variables", envPath)
	} else {
		log.Printf("Loaded environment variables from %s", envPath)
	}

	// Get absolute paths for frontend files
	// First get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get current working directory:", err)
	}

	log.Println("Current working directory:", currentDir)

	// Check if we're in the middleware directory or the parent tokkatot directory
	var frontendPath string
	if filepath.Base(currentDir) == "middleware" {
		// We're in middleware, go up one level to tokkatot, then into frontend
		frontendPath = filepath.Join(filepath.Dir(currentDir), "frontend")
	} else {
		// We're likely in the tokkatot directory already, just add frontend
		frontendPath = filepath.Join(currentDir, "frontend")
	}

	log.Println("Calculated frontend path:", frontendPath)

	// Verify the path exists
	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		log.Fatal("Frontend directory not found at: ", frontendPath)
	}

	log.Println("Frontend path resolved successfully to:", frontendPath)

	app := fiber.New()

	// Initialize database
	api.DB = database.InitDB()

	// Serve static files with absolute paths
	app.Static("/assets", filepath.Join(frontendPath, "assets"))
	app.Static("/components", filepath.Join(frontendPath, "components"))
	app.Static("/css", filepath.Join(frontendPath, "css"))
	app.Static("/js", filepath.Join(frontendPath, "js"))

	// Authentication and static page routes
	app.Get("/login", func(c *fiber.Ctx) error {
		if api.ValidateCookie(c) == nil {
			return c.Redirect("/")
		}
		return c.SendFile(filepath.Join(frontendPath, "pages", "login.html"))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		if api.ValidateCookie(c) != nil {
			return c.Redirect("/login")
		}
		return c.SendFile(filepath.Join(frontendPath, "pages", "index.html"))
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		if api.ValidateCookie(c) == nil {
			return c.Redirect("/")
		}
		return c.SendFile(filepath.Join(frontendPath, "pages", "signup.html"))
	})

	app.Get("/profile", func(c *fiber.Ctx) error {
		if api.ValidateCookie(c) != nil {
			return c.Redirect("/login")
		}
		return c.SendFile(filepath.Join(frontendPath, "pages", "profile.html"))
	})

	app.Get("/settings", func(c *fiber.Ctx) error {
		if api.ValidateCookie(c) != nil {
			return c.Redirect("/login")
		}
		return c.SendFile(filepath.Join(frontendPath, "pages", "settings.html"))
	})

	app.Get("/disease-detection", func(c *fiber.Ctx) error {
		if api.ValidateCookie(c) != nil {
			return c.Redirect("/login")
		}
		return c.SendFile(filepath.Join(frontendPath, "pages", "disease-detection.html"))
	})

	// User authentication routes
	app.Post("/register", api.RegisterHandler)
	app.Post("/login", api.LoginHandler)

	// API routes (protected by authentication)
	apiRoutes := app.Group("/api", func(c *fiber.Ctx) error {
		if api.ValidateCookie(c) != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
		}
		return c.Next()
	})

	// Profile routes
	apiRoutes.Get("/profile", api.GetProfileHandler)
	apiRoutes.Post("/profile", api.UpdateProfileHandler)

	// Poultry system sensor data retrieval
	apiRoutes.Get("/get-initial-state", api.GetInitialStateHandler)
	apiRoutes.Get("/get-current-data", api.GetCurrentDataHandler)
	apiRoutes.Get("/get-historical-data", api.GetHistoricalDataHandler)

	// Poultry system control routes
	apiRoutes.Get("/toggle-auto", api.ToggleAutoHandler)
	apiRoutes.Get("/toggle-belt", api.ToggleBeltHandler)
	apiRoutes.Get("/toggle-fan", api.ToggleFanHandler)
	apiRoutes.Get("/toggle-bulb", api.ToggleBulbHandler)
	apiRoutes.Get("/toggle-feeder", api.ToggleFeederHandler)
	apiRoutes.Get("/toggle-pump", api.TogglePumpHandler)

	// AI Disease Detection routes
	apiRoutes.Get("/ai/health", api.AIHealthCheckHandler)
	apiRoutes.Post("/ai/predict-disease", api.PredictDiseaseHandler)
	apiRoutes.Get("/ai/disease-info", api.GetDiseaseInfoHandler)

	// Schedule management routes
	/* apiRoutes.Post("/schedule", api.SaveScheduleHandler)      // Save schedule
	apiRoutes.Get("/toggle-schedule", api.GetScheduleHandler) // Retrieve schedule
	*/

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendFile("../frontend/pages/404.html")
	})

	// Close the database connection when the application shuts down
	defer api.DB.Close()

	// Start the server
	log.Println("Server is running on port 443")
	log.Fatal(app.ListenTLS(":443", os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")))
}
