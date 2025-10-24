package api

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	dataProvider = "https://10.0.0.2"
	httpClient   *http.Client
)

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
}

// ====== DATA HANDLERS ====== //
func getDataHandler(c **fiber.Ctx, endpoint string) error {
	resp, err := httpClient.Get(dataProvider + endpoint)
	if err != nil || resp.StatusCode != http.StatusOK {
		return (*c).Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get data from data provider"})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return (*c).Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Forward raw data from the data provider directly to the client
	return (*c).JSON(fiber.Map{"data": string(body)})
}

func GetInitialStateHandler(c *fiber.Ctx) error {
	return getDataHandler(&c, "/get-initial-state")
}

func GetCurrentDataHandler(c *fiber.Ctx) error {
	return getDataHandler(&c, "/get-current-data")
}

func GetHistoricalDataHandler(c *fiber.Ctx) error {
	return getDataHandler(&c, "/get-historical-data")
}

// ====== TOGGLE HANDLERS ====== //
func toggleHandler(c **fiber.Ctx, endpoint string) error {
	resp, err := httpClient.Get(dataProvider + endpoint)
	if err != nil || resp.StatusCode != http.StatusOK {
		return (*c).Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to toggle device"})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return (*c).Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Verify step removed: return the provider response as the new state
	return (*c).Status(fiber.StatusOK).JSON(fiber.Map{"state": string(body)})
}

func ToggleAutoHandler(c *fiber.Ctx) error {
	return toggleHandler(&c, "/toggle-auto")
}

func ToggleBeltHandler(c *fiber.Ctx) error {
	return toggleHandler(&c, "/toggle-belt")
}

func ToggleFanHandler(c *fiber.Ctx) error {
	return toggleHandler(&c, "/toggle-fan")
}

func ToggleBulbHandler(c *fiber.Ctx) error {
	return toggleHandler(&c, "/toggle-bulb")
}

func ToggleFeederHandler(c *fiber.Ctx) error {
	return toggleHandler(&c, "/toggle-feeder")
}

func ToggleWaterHandler(c *fiber.Ctx) error {
	return toggleHandler(&c, "/toggle-water")
}

// ====== SCHEDULE HANDLERS ====== //
/* func SaveScheduleHandler(c *fiber.Ctx) error {
	var schedule database.Schedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid schedule data"})
	}

	// Save schedule to database
	if err := database.SaveSchedule(DB, schedule); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save schedule"})
	}

	return c.JSON(fiber.Map{"message": "Schedule saved successfully"})
}

func GetScheduleHandler(c *fiber.Ctx) error {
	schedule, err := database.GetSchedule(DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve schedule"})
	}

	return c.JSON(schedule)
} */
