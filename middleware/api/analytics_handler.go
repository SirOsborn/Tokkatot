package api

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"middleware/database"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== ANALYTICS & REPORTING HANDLERS =====

// GetFarmDashboardHandler returns overview stats for a farm
// @GET /v1/farms/:farm_id/dashboard
func GetFarmDashboardHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	// Farm basic info
	var farmName string
	database.DB.QueryRow("SELECT name FROM farms WHERE id = $1", farmID).Scan(&farmName)

	// Device stats
	var totalDevices, onlineDevices, offlineDevices int64
	database.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE farm_id = $1 AND is_active = true", farmID).Scan(&totalDevices)
	database.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE farm_id = $1 AND is_active = true AND is_online = true", farmID).Scan(&onlineDevices)
	offlineDevices = totalDevices - onlineDevices

	// Alert stats
	var activeAlerts, criticalAlerts, warningAlerts int64
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true", farmID).Scan(&activeAlerts)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true AND severity = 'critical'", farmID).Scan(&criticalAlerts)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true AND severity = 'warning'", farmID).Scan(&warningAlerts)

	// Commands last 24h
	var last24hCommands int64
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM device_commands
		WHERE farm_id = $1 AND created_at > CURRENT_TIMESTAMP - INTERVAL '24 hours'
	`, farmID).Scan(&last24hCommands)

	// Alerts last 24h
	var last24hAlerts int64
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM alerts
		WHERE farm_id = $1 AND triggered_at > CURRENT_TIMESTAMP - INTERVAL '24 hours'
	`, farmID).Scan(&last24hAlerts)

	// Recent events (last 10)
	rows, err := database.DB.Query(`
		SELECT id, event_type, resource_id, ip_address, created_at
		FROM event_logs WHERE farm_id = $1
		ORDER BY created_at DESC LIMIT 10
	`, farmID)

	type EventEntry struct {
		ID         string  `json:"id"`
		EventType  string  `json:"event_type"`
		ResourceID *string `json:"resource_id,omitempty"`
		IPAddress  *string `json:"ip_address,omitempty"`
		Timestamp  string  `json:"timestamp"`
	}
	recentEvents := []EventEntry{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e EventEntry
			rows.Scan(&e.ID, &e.EventType, &e.ResourceID, &e.IPAddress, &e.Timestamp)
			recentEvents = append(recentEvents, e)
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"farm": fiber.Map{
			"id":   farmID,
			"name": farmName,
		},
		"device_status": fiber.Map{
			"total":   totalDevices,
			"online":  onlineDevices,
			"offline": offlineDevices,
		},
		"alerts": fiber.Map{
			"active":   activeAlerts,
			"critical": criticalAlerts,
			"warning":  warningAlerts,
		},
		"quick_stats": fiber.Map{
			"last_24h_commands": last24hCommands,
			"last_24h_alerts":   last24hAlerts,
		},
		"recent_events": recentEvents,
	}, "Dashboard loaded")
}

// GetDeviceMetricsReportHandler returns aggregated sensor data for a device
// @GET /v1/farms/:farm_id/reports/device-metrics?device_id=uuid&from=2026-02-01&to=2026-02-19&metric=temperature
func GetDeviceMetricsReportHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	deviceIDStr := c.Query("device_id")
	if deviceIDStr == "" {
		return utils.BadRequest(c, "missing_device_id", "device_id query parameter is required")
	}
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device_id")
	}

	from := c.Query("from", "")
	to := c.Query("to", "")
	metric := c.Query("metric", "all")

	query := `
		SELECT
			DATE_TRUNC('day', timestamp) AS day,
			sensor_type,
			AVG(value) AS avg_val,
			MIN(value) AS min_val,
			MAX(value) AS max_val,
			unit
		FROM device_readings
		WHERE device_id = $1
	`
	args := []interface{}{deviceID}

	if from != "" {
		query += " AND timestamp >= $" + strconv.Itoa(len(args)+1) + "::date"
		args = append(args, from)
	}
	if to != "" {
		query += " AND timestamp <= ($" + strconv.Itoa(len(args)+1) + "::date + INTERVAL '1 day')"
		args = append(args, to)
	}
	if metric != "all" {
		query += " AND sensor_type = $" + strconv.Itoa(len(args)+1)
		args = append(args, metric)
	}

	query += " GROUP BY day, sensor_type, unit ORDER BY day DESC LIMIT 1000"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Device metrics report error: %v", err)
		return utils.InternalError(c, "Failed to fetch metrics report")
	}
	defer rows.Close()

	type DataPoint struct {
		Timestamp  string  `json:"timestamp"`
		SensorType string  `json:"sensor_type"`
		Avg        float64 `json:"avg"`
		Min        float64 `json:"min"`
		Max        float64 `json:"max"`
		Unit       string  `json:"unit"`
	}

	dataPoints := []DataPoint{}
	for rows.Next() {
		var dp DataPoint
		if err := rows.Scan(&dp.Timestamp, &dp.SensorType, &dp.Avg, &dp.Min, &dp.Max, &dp.Unit); err != nil {
			continue
		}
		dataPoints = append(dataPoints, dp)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"device_id":   deviceID,
		"metric":      metric,
		"from_date":   from,
		"to_date":     to,
		"data_points": dataPoints,
	}, "Device metrics report generated")
}

// GetDeviceUsageReportHandler returns usage statistics for a device
// @GET /v1/farms/:farm_id/reports/device-usage?device_id=uuid&days=30
func GetDeviceUsageReportHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	deviceIDStr := c.Query("device_id")
	if deviceIDStr == "" {
		return utils.BadRequest(c, "missing_device_id", "device_id is required")
	}
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device_id")
	}

	days, _ := strconv.Atoi(c.Query("days", "30"))
	if days > 365 {
		days = 365
	}

	var deviceName string
	database.DB.QueryRow("SELECT name FROM devices WHERE id = $1", deviceID).Scan(&deviceName)

	// Count total commands in period
	var totalCycles, successCycles int64
	database.DB.QueryRow(`
		SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'success')
		FROM device_commands
		WHERE device_id = $1 AND created_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')
	`, deviceID, days).Scan(&totalCycles, &successCycles)

	reliability := 0.0
	if totalCycles > 0 {
		reliability = float64(successCycles) / float64(totalCycles) * 10
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"device_id":         deviceID,
		"device_name":       deviceName,
		"days":              days,
		"total_cycles":      totalCycles,
		"successful_cycles": successCycles,
		"reliability_score": reliability,
	}, "Device usage report generated")
}

// GetFarmPerformanceReportHandler returns farm-level performance summary
// @GET /v1/farms/:farm_id/reports/farm-performance?days=30
func GetFarmPerformanceReportHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	days, _ := strconv.Atoi(c.Query("days", "30"))
	if days > 365 {
		days = 365
	}

	var totalDevices, onlineDevices int64
	database.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE farm_id = $1 AND is_active = true", farmID).Scan(&totalDevices)
	database.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE farm_id = $1 AND is_active = true AND is_online = true", farmID).Scan(&onlineDevices)

	var totalAlerts, criticalAlerts int64
	database.DB.QueryRow(`
		SELECT COUNT(*), COUNT(*) FILTER (WHERE severity = 'critical')
		FROM alerts WHERE farm_id = $1 AND triggered_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')
	`, farmID, days).Scan(&totalAlerts, &criticalAlerts)

	var scheduledCmds, successCmds int64
	database.DB.QueryRow(`
		SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'success')
		FROM device_commands
		WHERE farm_id = $1 AND created_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')
	`, farmID, days).Scan(&scheduledCmds, &successCmds)

	successRate := 0.0
	if scheduledCmds > 0 {
		successRate = float64(successCmds) / float64(scheduledCmds) * 100
	}

	uptimePct := 0.0
	if totalDevices > 0 {
		uptimePct = float64(onlineDevices) / float64(totalDevices) * 100
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"farm_id": farmID,
		"days":    days,
		"device_health": fiber.Map{
			"total":   totalDevices,
			"online":  onlineDevices,
			"offline": totalDevices - onlineDevices,
		},
		"alerts_triggered": totalAlerts,
		"critical_alerts":  criticalAlerts,
		"automation_efficiency": fiber.Map{
			"total_commands": scheduledCmds,
			"successful":     successCmds,
			"success_rate":   fmt.Sprintf("%.1f%%", successRate),
		},
		"uptime_percent": fmt.Sprintf("%.1f%%", uptimePct),
	}, "Farm performance report generated")
}

// GetFarmEventLogHandler returns the event audit log for a farm
// @GET /v1/farms/:farm_id/events?limit=100&offset=0&event_type=all&days=30
func GetFarmEventLogHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	days, _ := strconv.Atoi(c.Query("days", "30"))
	eventType := c.Query("event_type", "all")
	if limit > 500 {
		limit = 500
	}
	if days > 365 {
		days = 365
	}

	query := `
		SELECT id, event_type, user_id, resource_id, ip_address, created_at
		FROM event_logs
		WHERE farm_id = $1 AND created_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')
	`
	args := []interface{}{farmID, days}

	if eventType != "all" {
		query += " AND event_type = $" + strconv.Itoa(len(args)+1)
		args = append(args, eventType)
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Get event log error: %v", err)
		return utils.InternalError(c, "Failed to fetch event log")
	}
	defer rows.Close()

	type EventEntry struct {
		ID         string  `json:"id"`
		EventType  string  `json:"event_type"`
		UserID     string  `json:"user_id"`
		ResourceID *string `json:"resource_id,omitempty"`
		IPAddress  *string `json:"ip_address,omitempty"`
		Timestamp  string  `json:"timestamp"`
	}

	events := []EventEntry{}
	for rows.Next() {
		var e EventEntry
		if err := rows.Scan(&e.ID, &e.EventType, &e.UserID, &e.ResourceID, &e.IPAddress, &e.Timestamp); err != nil {
			continue
		}
		events = append(events, e)
	}

	var total int64
	database.DB.QueryRow("SELECT COUNT(*) FROM event_logs WHERE farm_id = $1", farmID).Scan(&total)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"events": events,
		"total":  total,
	}, "Event log fetched successfully")
}

// ExportReportHandler exports report data as CSV (MVP: CSV only)
// @GET /v1/farms/:farm_id/reports/export?type=device_metrics&device_id=uuid&from=2026-02-01&to=2026-02-19
func ExportReportHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	reportType := c.Query("type", "device_metrics")
	from := c.Query("from", "")
	to := c.Query("to", "")
	deviceIDStr := c.Query("device_id", "")

	var csvRows []string

	switch reportType {
	case "device_metrics":
		if deviceIDStr == "" {
			return utils.BadRequest(c, "missing_device_id", "device_id is required for device_metrics export")
		}
		deviceID, err := uuid.Parse(deviceIDStr)
		if err != nil {
			return utils.BadRequest(c, "invalid_id", "Invalid device_id")
		}

		query := `
			SELECT timestamp, sensor_type, value, unit, quality
			FROM device_readings WHERE device_id = $1
		`
		args := []interface{}{deviceID}
		if from != "" {
			query += " AND timestamp >= $" + strconv.Itoa(len(args)+1) + "::date"
			args = append(args, from)
		}
		if to != "" {
			query += " AND timestamp <= ($" + strconv.Itoa(len(args)+1) + "::date + INTERVAL '1 day')"
			args = append(args, to)
		}
		query += " ORDER BY timestamp DESC LIMIT 50000"

		rows, err := database.DB.Query(query, args...)
		if err != nil {
			return utils.InternalError(c, "Failed to fetch data for export")
		}
		defer rows.Close()

		csvRows = append(csvRows, "timestamp,sensor_type,value,unit,quality")
		for rows.Next() {
			var ts, stype, unit, quality string
			var val float64
			if err := rows.Scan(&ts, &stype, &val, &unit, &quality); err != nil {
				continue
			}
			csvRows = append(csvRows, fmt.Sprintf("%s,%s,%.4f,%s,%s", ts, stype, val, unit, quality))
		}

	case "farm_events":
		rows, err := database.DB.Query(`
			SELECT created_at, event_type, user_id, ip_address
			FROM event_logs WHERE farm_id = $1 ORDER BY created_at DESC LIMIT 50000
		`, farmID)
		if err != nil {
			return utils.InternalError(c, "Failed to fetch data for export")
		}
		defer rows.Close()

		csvRows = append(csvRows, "timestamp,event_type,user_id,ip_address")
		for rows.Next() {
			var ts, evType, uid string
			var ip *string
			if err := rows.Scan(&ts, &evType, &uid, &ip); err != nil {
				continue
			}
			ipStr := ""
			if ip != nil {
				ipStr = *ip
			}
			csvRows = append(csvRows, fmt.Sprintf("%s,%s,%s,%s", ts, evType, uid, ipStr))
		}

	default:
		return utils.BadRequest(c, "invalid_type", "type must be 'device_metrics' or 'farm_events'")
	}

	csvContent := strings.Join(csvRows, "\n")
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_%s.csv\"", reportType, farmID))
	return c.SendString(csvContent)
}
