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

// ===== DEVICE MANAGEMENT HANDLERS =====

// ListDevicesHandler returns all devices in a farm or coop
// @GET /v1/farms/:farm_id/devices?coop_id=xxx&limit=20&offset=0&filter=online
func ListDevicesHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check user access
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	// Optional coop filter
	var coopID *uuid.UUID
	if coopIDStr := c.Query("coop_id"); coopIDStr != "" {
		parsedID, err := uuid.Parse(coopIDStr)
		if err == nil {
			coopID = &parsedID
		}
	}

	// Pagination
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	filter := c.Query("filter", "all") // all, online, offline

	if limit > 100 {
		limit = 100
	}

	// Build query
	query := `
	SELECT id, farm_id, coop_id, device_id, name, type, model, is_main_controller,
	       firmware_version, hardware_id, location, is_active, is_online,
	       last_heartbeat, last_command_status, last_command_at, created_at, updated_at
	FROM devices
	WHERE farm_id = $1 AND is_active = true
	`
	args := []interface{}{farmID}

	if coopID != nil {
		query += " AND coop_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, coopID)
	}

	if filter == "online" {
		query += " AND is_online = true"
	} else if filter == "offline" {
		query += " AND is_online = false"
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("List devices error: %v", err)
		return utils.InternalError(c, "Failed to fetch devices")
	}
	defer rows.Close()

	devices := []models.Device{}
	for rows.Next() {
		var device models.Device
		err := rows.Scan(
			&device.ID, &device.FarmID, &device.CoopID, &device.DeviceID, &device.Name,
			&device.Type, &device.Model, &device.IsMainController, &device.FirmwareVersion,
			&device.HardwareID, &device.Location, &device.IsActive, &device.IsOnline,
			&device.LastHeartbeat, &device.LastCommandStatus, &device.LastCommandAt,
			&device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan device error: %v", err)
			continue
		}
		devices = append(devices, device)
	}

	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM devices WHERE farm_id = $1 AND is_active = true"
	countArgs := []interface{}{farmID}
	if coopID != nil {
		countQuery += " AND coop_id = $2"
		countArgs = append(countArgs, coopID)
	}
	database.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	return utils.SuccessListResponse(c, devices, total, offset/limit+1, limit)
}

// GetDeviceHandler returns a single device by ID
// @GET /v1/farms/:farm_id/devices/:device_id
func GetDeviceHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	deviceIDStr := c.Params("device_id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	// Check user access
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	var device models.Device
	query := `
	SELECT id, farm_id, coop_id, device_id, name, type, model, is_main_controller,
	       firmware_version, hardware_id, location, is_active, is_online,
	       last_heartbeat, last_command_status, last_command_at, created_at, updated_at
	FROM devices
	WHERE id = $1 AND farm_id = $2
	`

	err = database.DB.QueryRow(query, deviceID, farmID).Scan(
		&device.ID, &device.FarmID, &device.CoopID, &device.DeviceID, &device.Name,
		&device.Type, &device.Model, &device.IsMainController, &device.FirmwareVersion,
		&device.HardwareID, &device.Location, &device.IsActive, &device.IsOnline,
		&device.LastHeartbeat, &device.LastCommandStatus, &device.LastCommandAt,
		&device.CreatedAt, &device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		log.Printf("Get device error: %v", err)
		return utils.InternalError(c, "Failed to fetch device")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, device, "Device fetched successfully")
}

// SendDeviceCommandHandler sends a command to a device
// @POST /v1/farms/:farm_id/devices/:device_id/commands
func SendDeviceCommandHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	deviceIDStr := c.Params("device_id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	// Check user has manager or owner role (device control permission)
	err = checkFarmAccess(userID, farmID, "manager")
	if err != nil {
		return err
	}

	var req struct {
		CommandType  string  `json:"command_type"` // "on", "off", "set_value", "status"
		CommandValue *string `json:"command_value,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.CommandType == "" {
		return utils.BadRequest(c, "missing_command", "Command type is required")
	}

	// Validate command type
	validCommands := map[string]bool{
		"on":        true,
		"off":       true,
		"set_value": true,
		"status":    true,
		"reboot":    true,
	}

	if !validCommands[req.CommandType] {
		return utils.BadRequest(c, "invalid_command", "Invalid command type")
	}

	// Get device info (check it exists and get coop_id)
	var coopID *uuid.UUID
	var deviceIsOnline bool
	err = database.DB.QueryRow(`
		SELECT coop_id, is_online FROM devices 
		WHERE id = $1 AND farm_id = $2 AND is_active = true
	`, deviceID, farmID).Scan(&coopID, &deviceIsOnline)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		log.Printf("Get device info error: %v", err)
		return utils.InternalError(c, "Failed to fetch device")
	}

	// Create device command
	commandID := uuid.New()
	query := `
	INSERT INTO device_commands (id, device_id, farm_id, coop_id, issued_by, command_type, command_value, status, issued_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING id, device_id, farm_id, coop_id, issued_by, command_type, command_value, status, issued_at, created_at
	`

	var cmd models.DeviceCommand
	err = database.DB.QueryRow(
		query,
		commandID, deviceID, farmID, coopID, userID, req.CommandType, req.CommandValue,
	).Scan(
		&cmd.ID, &cmd.DeviceID, &cmd.FarmID, &cmd.CoopID, &cmd.IssuedBy,
		&cmd.CommandType, &cmd.CommandValue, &cmd.Status, &cmd.IssuedAt, &cmd.CreatedAt,
	)

	if err != nil {
		log.Printf("Create device command error: %v", err)
		return utils.InternalError(c, "Failed to send command")
	}

	// Broadcast command update via WebSocket
	BroadcastCommandUpdate(cmd, farmID, coopID)

	// TODO: Publish command to MQTT broker
	// mqtt.Publish("devices/"+deviceID.String()+"/commands", cmd)

	response := fiber.Map{
		"command": cmd,
		"warning": nil,
	}

	if !deviceIsOnline {
		response["warning"] = "Device is currently offline. Command will be executed when device comes online."
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, response, "Command sent successfully")
}

// GetDeviceCommandStatusHandler returns the status of a command
// @GET /v1/farms/:farm_id/devices/:device_id/commands/:command_id
func GetDeviceCommandStatusHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	deviceIDStr := c.Params("device_id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	commandIDStr := c.Params("command_id")
	commandID, err := uuid.Parse(commandIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid command ID")
	}

	// Check user access
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	var cmd models.DeviceCommand
	query := `
	SELECT id, device_id, farm_id, coop_id, issued_by, command_type, command_value, 
	       status, response, issued_at, executed_at, created_at
	FROM device_commands
	WHERE id = $1 AND device_id = $2 AND farm_id = $3
	`

	err = database.DB.QueryRow(query, commandID, deviceID, farmID).Scan(
		&cmd.ID, &cmd.DeviceID, &cmd.FarmID, &cmd.CoopID, &cmd.IssuedBy,
		&cmd.CommandType, &cmd.CommandValue, &cmd.Status, &cmd.Response,
		&cmd.IssuedAt, &cmd.ExecutedAt, &cmd.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Command not found")
	}
	if err != nil {
		log.Printf("Get command status error: %v", err)
		return utils.InternalError(c, "Failed to fetch command status")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cmd, "Command status fetched successfully")
}

// ListDeviceCommandsHandler returns command history for a device
// @GET /v1/farms/:farm_id/devices/:device_id/commands?limit=50&offset=0&status=all
func ListDeviceCommandsHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	deviceIDStr := c.Params("device_id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	// Check user access
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	// Pagination
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	statusFilter := c.Query("status", "all") // all, pending, success, failed, timeout

	if limit > 100 {
		limit = 100
	}

	query := `
	SELECT id, device_id, farm_id, coop_id, issued_by, command_type, command_value,
	       status, response, issued_at, executed_at, created_at
	FROM device_commands
	WHERE device_id = $1 AND farm_id = $2
	`
	args := []interface{}{deviceID, farmID}

	if statusFilter != "all" {
		query += " AND status = $3"
		args = append(args, statusFilter)
	}

	query += " ORDER BY issued_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("List commands error: %v", err)
		return utils.InternalError(c, "Failed to fetch commands")
	}
	defer rows.Close()

	commands := []models.DeviceCommand{}
	for rows.Next() {
		var cmd models.DeviceCommand
		err := rows.Scan(
			&cmd.ID, &cmd.DeviceID, &cmd.FarmID, &cmd.CoopID, &cmd.IssuedBy,
			&cmd.CommandType, &cmd.CommandValue, &cmd.Status, &cmd.Response,
			&cmd.IssuedAt, &cmd.ExecutedAt, &cmd.CreatedAt,
		)
		if err != nil {
			log.Printf("Scan command error: %v", err)
			continue
		}
		commands = append(commands, cmd)
	}

	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM device_commands WHERE device_id = $1 AND farm_id = $2"
	countArgs := []interface{}{deviceID, farmID}
	if statusFilter != "all" {
		countQuery += " AND status = $3"
		countArgs = append(countArgs, statusFilter)
	}
	database.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	return utils.SuccessListResponse(c, commands, total, offset/limit+1, limit)
}

// UpdateDeviceHeartbeatHandler updates device heartbeat (called by devices)
// @POST /v1/devices/:hardware_id/heartbeat
func UpdateDeviceHeartbeatHandler(c *fiber.Ctx) error {
	// This endpoint is called by devices themselves (ESP32/Raspberry Pi)
	// In production, would use device API key authentication
	hardwareID := c.Params("hardware_id")

	if hardwareID == "" {
		return utils.BadRequest(c, "missing_id", "Hardware ID is required")
	}

	now := time.Now()

	// Update last_heartbeat and set is_online = true
	_, err := database.DB.Exec(`
		UPDATE devices 
		SET last_heartbeat = $1, is_online = true, updated_at = $1
		WHERE hardware_id = $2
	`, now, hardwareID)

	if err != nil {
		log.Printf("Update heartbeat error: %v", err)
		return utils.InternalError(c, "Failed to update heartbeat")
	}

	// Fetch device details for WebSocket broadcast
	var device models.Device
	err = database.DB.QueryRow(`
		SELECT id, farm_id, coop_id, device_id, name, type, is_online, last_heartbeat
		FROM devices WHERE hardware_id = $1
	`, hardwareID).Scan(
		&device.ID, &device.FarmID, &device.CoopID, &device.DeviceID,
		&device.Name, &device.Type, &device.IsOnline, &device.LastHeartbeat,
	)

	if err == nil {
		// Broadcast device update via WebSocket
		BroadcastDeviceUpdate(device)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message":   "Heartbeat received",
		"timestamp": now,
	}, "Heartbeat updated")
}
