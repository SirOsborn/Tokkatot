package api

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"middleware/database"
	"middleware/models"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== COOP MANAGEMENT HANDLERS =====

// ListCoopsHandler returns all coops in a farm
// @GET /v1/farms/:farm_id/coops?limit=20&offset=0
func ListCoopsHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check user access to farm
	err = checkFarmAccess(userID, farmID, "viewer") // Minimum role: viewer
	if err != nil {
		return err
	}

	// Pagination
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	// Get coops with device count
	query := `
	SELECT c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count, 
	       c.chicken_type, c.main_device_id, c.description, c.is_active, 
	       c.created_at, c.updated_at,
	       COUNT(d.id) AS device_count
	FROM coops c
	LEFT JOIN devices d ON c.id = d.coop_id AND d.is_active = true
	WHERE c.farm_id = $1 AND c.is_active = true
	GROUP BY c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count,
	         c.chicken_type, c.main_device_id, c.description, c.is_active,
	         c.created_at, c.updated_at
	ORDER BY c.number ASC
	LIMIT $2 OFFSET $3
	`

	rows, err := database.DB.Query(query, farmID, limit, offset)
	if err != nil {
		log.Printf("List coops error: %v", err)
		return utils.InternalError(c, "Failed to fetch coops")
	}
	defer rows.Close()

	type CoopWithDevices struct {
		models.Coop
		DeviceCount int `json:"device_count"`
	}

	coops := []CoopWithDevices{}
	for rows.Next() {
		var coop CoopWithDevices
		err := rows.Scan(
			&coop.ID, &coop.FarmID, &coop.Number, &coop.Name, &coop.Capacity,
			&coop.CurrentCount, &coop.ChickenType, &coop.MainDeviceID, &coop.Description,
			&coop.IsActive, &coop.CreatedAt, &coop.UpdatedAt, &coop.DeviceCount,
		)
		if err != nil {
			log.Printf("Scan coop error: %v", err)
			continue
		}
		coops = append(coops, coop)
	}

	// Get total count
	var total int64
	database.DB.QueryRow("SELECT COUNT(*) FROM coops WHERE farm_id = $1 AND is_active = true", farmID).Scan(&total)

	return utils.SuccessListResponse(c, coops, total, offset/limit+1, limit)
}

// GetCoopHandler returns a single coop by ID
// @GET /v1/farms/:farm_id/coops/:coop_id
func GetCoopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	coopIDStr := c.Params("coop_id")
	coopID, err := uuid.Parse(coopIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	// Check user access
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	// Get coop with device count
	var coop struct {
		models.Coop
		DeviceCount int `json:"device_count"`
	}

	query := `
	SELECT c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count,
	       c.chicken_type, c.main_device_id, c.description, c.is_active,
	       c.created_at, c.updated_at,
	       COUNT(d.id) AS device_count
	FROM coops c
	LEFT JOIN devices d ON c.id = d.coop_id AND d.is_active = true
	WHERE c.id = $1 AND c.farm_id = $2
	GROUP BY c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count,
	         c.chicken_type, c.main_device_id, c.description, c.is_active,
	         c.created_at, c.updated_at
	`

	err = database.DB.QueryRow(query, coopID, farmID).Scan(
		&coop.ID, &coop.FarmID, &coop.Number, &coop.Name, &coop.Capacity,
		&coop.CurrentCount, &coop.ChickenType, &coop.MainDeviceID, &coop.Description,
		&coop.IsActive, &coop.CreatedAt, &coop.UpdatedAt, &coop.DeviceCount,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Coop not found")
	}
	if err != nil {
		log.Printf("Get coop error: %v", err)
		return utils.InternalError(c, "Failed to fetch coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, coop, "Coop fetched successfully")
}

// CreateCoopHandler creates a new coop (manager/owner only)
// @POST /v1/farms/:farm_id/coops
func CreateCoopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check user has farmer role
	err = checkFarmAccess(userID, farmID, "farmer")
	if err != nil {
		return err
	}

	var req struct {
		Number       int     `json:"number"`
		Name         string  `json:"name"`
		Capacity     *int    `json:"capacity,omitempty"`
		CurrentCount *int    `json:"current_count,omitempty"`
		ChickenType  *string `json:"chicken_type,omitempty"`
		Description  *string `json:"description,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.Number < 1 {
		return utils.BadRequest(c, "invalid_number", "Coop number must be at least 1")
	}

	if req.Name == "" {
		return utils.BadRequest(c, "missing_name", "Coop name is required")
	}

	// Check if coop number already exists for this farm
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM coops WHERE farm_id = $1 AND number = $2)", farmID, req.Number).Scan(&exists)
	if err != nil {
		log.Printf("Coop number check error: %v", err)
		return utils.InternalError(c, "Failed to validate coop number")
	}
	if exists {
		return utils.Conflict(c, "coop_exists", "A coop with this number already exists in this farm")
	}

	// Create coop
	coopID := uuid.New()
	query := `
	INSERT INTO coops (id, farm_id, number, name, capacity, current_count, chicken_type, description, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING id, farm_id, number, name, capacity, current_count, chicken_type, main_device_id, description, is_active, created_at, updated_at
	`

	var coop models.Coop
	err = database.DB.QueryRow(
		query,
		coopID, farmID, req.Number, req.Name, req.Capacity, req.CurrentCount, req.ChickenType, req.Description,
	).Scan(
		&coop.ID, &coop.FarmID, &coop.Number, &coop.Name, &coop.Capacity,
		&coop.CurrentCount, &coop.ChickenType, &coop.MainDeviceID, &coop.Description,
		&coop.IsActive, &coop.CreatedAt, &coop.UpdatedAt,
	)

	if err != nil {
		log.Printf("Create coop error: %v", err)
		return utils.InternalError(c, "Failed to create coop")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, coop, "Coop created successfully")
}

// UpdateCoopHandler updates coop details (manager/owner only)
// @PUT /v1/farms/:farm_id/coops/:coop_id
func UpdateCoopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	coopIDStr := c.Params("coop_id")
	coopID, err := uuid.Parse(coopIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	// Check user has farmer role
	err = checkFarmAccess(userID, farmID, "farmer")
	if err != nil {
		return err
	}

	var req struct {
		Name         *string `json:"name,omitempty"`
		Capacity     *int    `json:"capacity,omitempty"`
		CurrentCount *int    `json:"current_count,omitempty"`
		ChickenType  *string `json:"chicken_type,omitempty"`
		Description  *string `json:"description,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	query := `
	UPDATE coops SET
		name = COALESCE($1, name),
		capacity = COALESCE($2, capacity),
		current_count = COALESCE($3, current_count),
		chicken_type = COALESCE($4, chicken_type),
		description = COALESCE($5, description),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $6 AND farm_id = $7
	RETURNING id, farm_id, number, name, capacity, current_count, chicken_type, main_device_id, description, is_active, created_at, updated_at
	`

	var coop models.Coop
	err = database.DB.QueryRow(
		query,
		req.Name, req.Capacity, req.CurrentCount, req.ChickenType, req.Description, coopID, farmID,
	).Scan(
		&coop.ID, &coop.FarmID, &coop.Number, &coop.Name, &coop.Capacity,
		&coop.CurrentCount, &coop.ChickenType, &coop.MainDeviceID, &coop.Description,
		&coop.IsActive, &coop.CreatedAt, &coop.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Coop not found")
	}
	if err != nil {
		log.Printf("Update coop error: %v", err)
		return utils.InternalError(c, "Failed to update coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, coop, "Coop updated successfully")
}

// DeleteCoopHandler soft-deletes a coop (owner only)
// @DELETE /v1/farms/:farm_id/coops/:coop_id
func DeleteCoopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	coopIDStr := c.Params("coop_id")
	coopID, err := uuid.Parse(coopIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	// Check user is farmer
	err = checkFarmAccess(userID, farmID, "farmer")
	if err != nil {
		return err
	}

	// Soft delete coop
	_, err = database.DB.Exec("UPDATE coops SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND farm_id = $2", coopID, farmID)
	if err != nil {
		log.Printf("Delete coop error: %v", err)
		return utils.InternalError(c, "Failed to delete coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Coop deleted successfully",
	}, "Coop deleted")
}

// ===== HELPER FUNCTIONS =====

// TemperatureTimelineHandler returns Apple Weather-style temperature data for a coop
// @GET /v1/farms/:farm_id/coops/:coop_id/temperature-timeline?days=7
func TemperatureTimelineHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	coopID, err := uuid.Parse(c.Params("coop_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	days, _ := strconv.Atoi(c.Query("days", "7"))
	if days > 30 {
		days = 30
	}
	if days < 1 {
		days = 1
	}

	// Find coop name
	var coopName string
	database.DB.QueryRow(`SELECT name FROM coops WHERE id = $1`, coopID).Scan(&coopName)

	// Find the temperature sensor device in this coop
	var deviceID uuid.UUID
	err = database.DB.QueryRow(`
		SELECT id FROM devices
		WHERE coop_id = $1 AND is_active = true AND type = 'sensor'
		ORDER BY is_main_controller DESC, created_at ASC LIMIT 1
	`, coopID).Scan(&deviceID)

	if err == sql.ErrNoRows {
		return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
			"coop_id":      coopID,
			"coop_name":    coopName,
			"sensor_found": false,
			"current_temp": nil,
			"bg_hint":      "warm",
			"today":        nil,
			"history":      []interface{}{},
		}, "No temperature sensor found in this coop")
	}
	if err != nil {
		log.Printf("Temperature timeline device lookup: %v", err)
		return utils.InternalError(c, "Failed to find sensor device")
	}

	// ── 1. Current temperature (latest reading) ──────────────────────────────
	var currentTemp *float64
	database.DB.QueryRow(`
		SELECT ROUND(value::numeric, 1)::float FROM device_readings
		WHERE device_id = $1 AND sensor_type = 'temperature'
		ORDER BY timestamp DESC LIMIT 1
	`, deviceID).Scan(&currentTemp)

	// ── 2. Today's hourly averages ────────────────────────────────────────────
	type HourlyPoint struct {
		Hour string  `json:"hour"` // "14:00"
		Temp float64 `json:"temp"`
	}
	hourlyData := []HourlyPoint{}
	hRows, hErr := database.DB.Query(`
		SELECT
			TO_CHAR(DATE_TRUNC('hour', timestamp), 'HH24:MI') AS hour_label,
			ROUND(AVG(value)::numeric, 1)::float AS avg_temp
		FROM device_readings
		WHERE device_id = $1 AND sensor_type = 'temperature'
		  AND timestamp >= CURRENT_DATE
		  AND timestamp < CURRENT_DATE + INTERVAL '1 day'
		GROUP BY DATE_TRUNC('hour', timestamp)
		ORDER BY DATE_TRUNC('hour', timestamp) ASC
	`, deviceID)
	if hErr == nil {
		defer hRows.Close()
		for hRows.Next() {
			var hp HourlyPoint
			if err := hRows.Scan(&hp.Hour, &hp.Temp); err == nil {
				hourlyData = append(hourlyData, hp)
			}
		}
	}

	// ── 3. Today's peak high and low with timestamps ─────────────────────────
	type TempPeak struct {
		Temp float64 `json:"temp"`
		Time string  `json:"time"` // "14:00"
	}
	var todayHigh, todayLow TempPeak
	database.DB.QueryRow(`
		SELECT ROUND(value::numeric, 1)::float, TO_CHAR(timestamp, 'HH24:MI')
		FROM device_readings
		WHERE device_id = $1 AND sensor_type = 'temperature'
		  AND timestamp >= CURRENT_DATE AND timestamp < CURRENT_DATE + INTERVAL '1 day'
		ORDER BY value DESC LIMIT 1
	`, deviceID).Scan(&todayHigh.Temp, &todayHigh.Time)
	database.DB.QueryRow(`
		SELECT ROUND(value::numeric, 1)::float, TO_CHAR(timestamp, 'HH24:MI')
		FROM device_readings
		WHERE device_id = $1 AND sensor_type = 'temperature'
		  AND timestamp >= CURRENT_DATE AND timestamp < CURRENT_DATE + INTERVAL '1 day'
		ORDER BY value ASC LIMIT 1
	`, deviceID).Scan(&todayLow.Temp, &todayLow.Time)

	// ── 4. Daily history (past N days including today) ────────────────────────
	type DayEntry struct {
		Date  string   `json:"date"`  // "2026-02-23"
		Label string   `json:"label"` // "Yesterday", "Mon"
		High  TempPeak `json:"high"`
		Low   TempPeak `json:"low"`
	}
	history := []DayEntry{}

	todayStr := time.Now().Format("2006-01-02")
	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	dRows, dErr := database.DB.Query(`
		SELECT
			TO_CHAR(DATE(timestamp), 'YYYY-MM-DD') AS day,
			ROUND(MAX(value)::numeric, 1)::float    AS high_temp,
			ROUND(MIN(value)::numeric, 1)::float    AS low_temp
		FROM device_readings
		WHERE device_id = $1 AND sensor_type = 'temperature'
		  AND timestamp >= CURRENT_DATE - ($2 * INTERVAL '1 day')
		  AND timestamp < CURRENT_DATE + INTERVAL '1 day'
		GROUP BY DATE(timestamp)
		ORDER BY DATE(timestamp) DESC
	`, deviceID, days)

	if dErr == nil {
		defer dRows.Close()
		for dRows.Next() {
			var entry DayEntry
			if err := dRows.Scan(&entry.Date, &entry.High.Temp, &entry.Low.Temp); err != nil {
				continue
			}
			// Get exact time of the high and low for this day
			database.DB.QueryRow(`
				SELECT TO_CHAR(timestamp, 'HH24:MI') FROM device_readings
				WHERE device_id = $1 AND sensor_type = 'temperature' AND DATE(timestamp) = $2::date
				ORDER BY value DESC LIMIT 1
			`, deviceID, entry.Date).Scan(&entry.High.Time)
			database.DB.QueryRow(`
				SELECT TO_CHAR(timestamp, 'HH24:MI') FROM device_readings
				WHERE device_id = $1 AND sensor_type = 'temperature' AND DATE(timestamp) = $2::date
				ORDER BY value ASC LIMIT 1
			`, deviceID, entry.Date).Scan(&entry.Low.Time)

			switch entry.Date {
			case todayStr:
				entry.Label = "Today"
			case yesterdayStr:
				entry.Label = "Yesterday"
			default:
				t, _ := time.Parse("2006-01-02", entry.Date)
				entry.Label = t.Format("Mon") // "Mon", "Tue", etc.
			}
			history = append(history, entry)
		}
	}

	// ── 5. Background colour hint ─────────────────────────────────────────────
	bgHint := "warm"
	if currentTemp != nil {
		switch {
		case *currentTemp >= 35:
			bgHint = "scorching"
		case *currentTemp >= 32:
			bgHint = "hot"
		case *currentTemp >= 28:
			bgHint = "warm"
		case *currentTemp >= 24:
			bgHint = "neutral"
		case *currentTemp >= 20:
			bgHint = "cool"
		default:
			bgHint = "cold"
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"coop_id":      coopID,
		"coop_name":    coopName,
		"farm_id":      farmID,
		"device_id":    deviceID,
		"sensor_found": true,
		"current_temp": currentTemp,
		"bg_hint":      bgHint,
		"today": fiber.Map{
			"date":   todayStr,
			"hourly": hourlyData,
			"high":   todayHigh,
			"low":    todayLow,
		},
		"history": history,
	}, "Temperature timeline fetched")
}

// checkFarmAccess verifies user has at least minimum role for farm
func checkFarmAccess(userID, farmID uuid.UUID, minRole string) error {
	var userRole string
	err := database.DB.QueryRow(`
		SELECT role FROM farm_users 
		WHERE farm_id = $1 AND user_id = $2 AND is_active = true
	`, farmID, userID).Scan(&userRole)

	if err == sql.ErrNoRows {
		return fiber.NewError(fiber.StatusForbidden, "You do not have access to this farm")
	}
	if err != nil {
		log.Printf("Check farm access error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to verify access")
	}

	// Check role hierarchy: farmer > viewer
	roleHierarchy := map[string]int{
		"farmer": 2,
		"viewer": 1,
	}

	if roleHierarchy[userRole] < roleHierarchy[minRole] {
		return fiber.NewError(fiber.StatusForbidden, "Insufficient permissions")
	}

	return nil
}
