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

// ===== DEVICE CRUD (admin-only add/delete, manager update) =====

// AddDeviceHandler adds a device to a farm (Tokkatot admin team only)
// @POST /v1/farms/:farm_id/devices
func AddDeviceHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "owner"); err != nil {
		return err
	}

	var req struct {
		DeviceID         string     `json:"device_id"`
		Name             string     `json:"name"`
		Type             string     `json:"type"`
		Model            *string    `json:"model,omitempty"`
		FirmwareVersion  string     `json:"firmware_version"`
		HardwareID       string     `json:"hardware_id"`
		Location         *string    `json:"location,omitempty"`
		CoopID           *uuid.UUID `json:"coop_id,omitempty"`
		IsMainController bool       `json:"is_main_controller"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if req.DeviceID == "" || req.Name == "" || req.Type == "" || req.FirmwareVersion == "" || req.HardwareID == "" {
		return utils.BadRequest(c, "missing_fields", "device_id, name, type, firmware_version and hardware_id are required")
	}

	validTypes := map[string]bool{"gpio": true, "relay": true, "pwm": true, "adc": true, "servo": true, "sensor": true}
	if !validTypes[req.Type] {
		return utils.BadRequest(c, "invalid_type", "type must be one of: gpio, relay, pwm, adc, servo, sensor")
	}

	id := uuid.New()
	var device models.Device
	err = database.DB.QueryRow(`
		INSERT INTO devices (id, farm_id, coop_id, device_id, name, type, model, is_main_controller,
		                     firmware_version, hardware_id, location, is_active, is_online,
		                     created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,true,false,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)
		RETURNING id, farm_id, coop_id, device_id, name, type, model, is_main_controller,
		          firmware_version, hardware_id, location, is_active, is_online,
		          last_heartbeat, last_command_status, last_command_at, created_at, updated_at
	`, id, farmID, req.CoopID, req.DeviceID, req.Name, req.Type, req.Model, req.IsMainController,
		req.FirmwareVersion, req.HardwareID, req.Location,
	).Scan(
		&device.ID, &device.FarmID, &device.CoopID, &device.DeviceID, &device.Name,
		&device.Type, &device.Model, &device.IsMainController, &device.FirmwareVersion,
		&device.HardwareID, &device.Location, &device.IsActive, &device.IsOnline,
		&device.LastHeartbeat, &device.LastCommandStatus, &device.LastCommandAt,
		&device.CreatedAt, &device.UpdatedAt,
	)
	if err != nil {
		log.Printf("Add device error: %v", err)
		return utils.InternalError(c, "Failed to add device (hardware_id or device_id may already exist)")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, device, "Device added successfully")
}

// UpdateDeviceHandler updates device name/location (manager+)
// @PUT /v1/farms/:farm_id/devices/:device_id
func UpdateDeviceHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	if err := checkFarmAccess(userID, farmID, "manager"); err != nil {
		return err
	}

	var req struct {
		Name     *string `json:"name,omitempty"`
		Location *string `json:"location,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	_, err = database.DB.Exec(`
		UPDATE devices SET
			name     = COALESCE($1, name),
			location = COALESCE($2, location),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND farm_id = $4 AND is_active = true
	`, req.Name, req.Location, deviceID, farmID)
	if err != nil {
		log.Printf("Update device error: %v", err)
		return utils.InternalError(c, "Failed to update device")
	}

	var device models.Device
	database.DB.QueryRow(`
		SELECT id, farm_id, coop_id, device_id, name, type, model, is_main_controller,
		       firmware_version, hardware_id, location, is_active, is_online,
		       last_heartbeat, last_command_status, last_command_at, created_at, updated_at
		FROM devices WHERE id = $1
	`, deviceID).Scan(
		&device.ID, &device.FarmID, &device.CoopID, &device.DeviceID, &device.Name,
		&device.Type, &device.Model, &device.IsMainController, &device.FirmwareVersion,
		&device.HardwareID, &device.Location, &device.IsActive, &device.IsOnline,
		&device.LastHeartbeat, &device.LastCommandStatus, &device.LastCommandAt,
		&device.CreatedAt, &device.UpdatedAt,
	)

	return utils.SuccessResponse(c, fiber.StatusOK, device, "Device updated successfully")
}

// DeleteDeviceHandler soft-deletes a device (owner only)
// @DELETE /v1/farms/:farm_id/devices/:device_id
func DeleteDeviceHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	if err := checkFarmAccess(userID, farmID, "owner"); err != nil {
		return err
	}

	result, err := database.DB.Exec(
		"UPDATE devices SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND farm_id = $2",
		deviceID, farmID,
	)
	if err != nil {
		log.Printf("Delete device error: %v", err)
		return utils.InternalError(c, "Failed to delete device")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NotFound(c, "Device not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Device deleted successfully"}, "Device deleted")
}

// ===== DEVICE ADVANCED HANDLERS =====

// GetDeviceHistoryHandler returns sensor readings for a device
// @GET /v1/farms/:farm_id/devices/:device_id/history?hours=24&limit=1000&metric=all
func GetDeviceHistoryHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	hours, _ := strconv.Atoi(c.Query("hours", "24"))
	limit, _ := strconv.Atoi(c.Query("limit", "1000"))
	metric := c.Query("metric", "all")
	if hours > 720 {
		hours = 720
	}
	if limit > 10000 {
		limit = 10000
	}

	query := `
		SELECT id, device_id, sensor_type, value, unit, quality, timestamp
		FROM device_readings
		WHERE device_id = $1 AND timestamp > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 hour')
	`
	args := []interface{}{deviceID, hours}

	if metric != "all" {
		query += " AND sensor_type = $3"
		args = append(args, metric)
	}
	query += " ORDER BY timestamp DESC LIMIT $" + strconv.Itoa(len(args)+1)
	args = append(args, limit)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Get device history error: %v", err)
		return utils.InternalError(c, "Failed to fetch device history")
	}
	defer rows.Close()

	readings := []models.DeviceReading{}
	for rows.Next() {
		var r models.DeviceReading
		if err := rows.Scan(&r.ID, &r.DeviceID, &r.SensorType, &r.Value, &r.Unit, &r.Quality, &r.Timestamp); err != nil {
			continue
		}
		readings = append(readings, r)
	}

	var total int64
	database.DB.QueryRow("SELECT COUNT(*) FROM device_readings WHERE device_id = $1", deviceID).Scan(&total)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"device_id": deviceID,
		"readings":  readings,
		"total":     total,
	}, "Device history fetched")
}

// GetDeviceStatusHandler returns real-time status of a device
// @GET /v1/farms/:farm_id/devices/:device_id/status
func GetDeviceStatusHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	var device models.Device
	err = database.DB.QueryRow(`
		SELECT id, is_online, last_heartbeat, last_command_status, last_command_at
		FROM devices WHERE id = $1 AND farm_id = $2 AND is_active = true
	`, deviceID, farmID).Scan(
		&device.ID, &device.IsOnline, &device.LastHeartbeat,
		&device.LastCommandStatus, &device.LastCommandAt,
	)
	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch device status")
	}

	// Get latest reading
	var latestValue *float64
	var latestUnit *string
	database.DB.QueryRow(`
		SELECT value, unit FROM device_readings WHERE device_id = $1 ORDER BY timestamp DESC LIMIT 1
	`, deviceID).Scan(&latestValue, &latestUnit)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"device_id":           deviceID,
		"is_online":           device.IsOnline,
		"last_heartbeat":      device.LastHeartbeat,
		"last_command_status": device.LastCommandStatus,
		"last_command_at":     device.LastCommandAt,
		"current_value":       latestValue,
		"unit":                latestUnit,
	}, "Device status fetched")
}

// GetDeviceConfigHandler returns configuration parameters for a device
// @GET /v1/farms/:farm_id/devices/:device_id/config
func GetDeviceConfigHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	rows, err := database.DB.Query(`
		SELECT id, device_id, parameter_name, parameter_value, unit,
		       min_value, max_value, is_calibrated, calibrated_at, created_at, updated_at
		FROM device_configurations WHERE device_id = $1
		ORDER BY parameter_name
	`, deviceID)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch device config")
	}
	defer rows.Close()

	configs := []models.DeviceConfiguration{}
	for rows.Next() {
		var cfg models.DeviceConfiguration
		if err := rows.Scan(
			&cfg.ID, &cfg.DeviceID, &cfg.ParameterName, &cfg.ParameterValue,
			&cfg.Unit, &cfg.MinValue, &cfg.MaxValue, &cfg.IsCalibrated,
			&cfg.CalibratedAt, &cfg.CreatedAt, &cfg.UpdatedAt,
		); err != nil {
			continue
		}
		configs = append(configs, cfg)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"device_id":      deviceID,
		"configurations": configs,
	}, "Config fetched successfully")
}

// UpdateDeviceConfigHandler upserts configuration parameters for a device
// @PUT /v1/farms/:farm_id/devices/:device_id/config
func UpdateDeviceConfigHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	if err := checkFarmAccess(userID, farmID, "manager"); err != nil {
		return err
	}

	var req struct {
		Configurations []struct {
			ParameterName  string  `json:"parameter_name"`
			ParameterValue string  `json:"parameter_value"`
			Unit           *string `json:"unit,omitempty"`
		} `json:"configurations"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if len(req.Configurations) == 0 {
		return utils.BadRequest(c, "missing_configs", "At least one configuration is required")
	}

	for _, cfg := range req.Configurations {
		id := uuid.New()
		_, err = database.DB.Exec(`
			INSERT INTO device_configurations (id, device_id, parameter_name, parameter_value, unit, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			ON CONFLICT (device_id, parameter_name) DO UPDATE
				SET parameter_value = $4, unit = COALESCE($5, device_configurations.unit), updated_at = CURRENT_TIMESTAMP
		`, id, deviceID, cfg.ParameterName, cfg.ParameterValue, cfg.Unit)
		if err != nil {
			log.Printf("Update config error: %v", err)
			return utils.InternalError(c, "Failed to update configuration")
		}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"device_id": deviceID,
		"updated":   len(req.Configurations),
	}, "Configuration updated successfully")
}

// CalibrateDeviceHandler sets a calibration value for a device parameter
// @POST /v1/farms/:farm_id/devices/:device_id/calibrate
func CalibrateDeviceHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	if err := checkFarmAccess(userID, farmID, "manager"); err != nil {
		return err
	}

	var req struct {
		ParameterName  string  `json:"parameter_name"`
		ParameterValue float64 `json:"parameter_value"`
		Unit           *string `json:"unit,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if req.ParameterName == "" {
		return utils.BadRequest(c, "missing_fields", "parameter_name is required")
	}

	now := time.Now()
	id := uuid.New()
	valueStr := strconv.FormatFloat(req.ParameterValue, 'f', -1, 64)

	_, err = database.DB.Exec(`
		INSERT INTO device_configurations (id, device_id, parameter_name, parameter_value, unit, is_calibrated, calibrated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (device_id, parameter_name) DO UPDATE
			SET parameter_value = $4, unit = COALESCE($5, device_configurations.unit),
			    is_calibrated = true, calibrated_at = $6, updated_at = CURRENT_TIMESTAMP
	`, id, deviceID, req.ParameterName, valueStr, req.Unit, now)
	if err != nil {
		log.Printf("Calibrate device error: %v", err)
		return utils.InternalError(c, "Failed to calibrate device")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"parameter_name":  req.ParameterName,
		"parameter_value": req.ParameterValue,
		"unit":            req.Unit,
		"is_calibrated":   true,
		"calibrated_at":   now,
	}, "Device calibrated successfully")
}

// CancelCommandHandler cancels a queued command
// @DELETE /v1/farms/:farm_id/devices/:device_id/commands/:command_id
func CancelCommandHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}
	commandID, err := uuid.Parse(c.Params("command_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid command ID")
	}

	if err := checkFarmAccess(userID, farmID, "manager"); err != nil {
		return err
	}

	result, err := database.DB.Exec(`
		UPDATE device_commands SET status = 'failed', executed_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND device_id = $2 AND farm_id = $3 AND status = 'pending'
	`, commandID, deviceID, farmID)
	if err != nil {
		log.Printf("Cancel command error: %v", err)
		return utils.InternalError(c, "Failed to cancel command")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.Conflict(c, "cannot_cancel", "Command not found or already sent to device")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Command cancelled"}, "Command cancelled")
}

// GetFarmCommandHistoryHandler returns command history across all devices in a farm
// @GET /v1/farms/:farm_id/commands?hours=24&device_id=optional&limit=100
func GetFarmCommandHistoryHandler(c *fiber.Ctx) error {
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

	hours, _ := strconv.Atoi(c.Query("hours", "24"))
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	deviceFilter := c.Query("device_id", "")
	if limit > 500 {
		limit = 500
	}

	query := `
		SELECT dc.id, dc.device_id, d.name AS device_name, dc.command_type, dc.status,
		       dc.issued_by, dc.issued_at, dc.executed_at
		FROM device_commands dc
		INNER JOIN devices d ON dc.device_id = d.id
		WHERE dc.farm_id = $1 AND dc.created_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 hour')
	`
	args := []interface{}{farmID, hours}

	if deviceFilter != "" {
		if did, err := uuid.Parse(deviceFilter); err == nil {
			query += " AND dc.device_id = $" + strconv.Itoa(len(args)+1)
			args = append(args, did)
		}
	}

	query += " ORDER BY dc.issued_at DESC LIMIT $" + strconv.Itoa(len(args)+1)
	args = append(args, limit)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Get command history error: %v", err)
		return utils.InternalError(c, "Failed to fetch command history")
	}
	defer rows.Close()

	type CommandEntry struct {
		ID          uuid.UUID  `json:"command_id"`
		DeviceID    uuid.UUID  `json:"device_id"`
		DeviceName  string     `json:"device_name"`
		CommandType string     `json:"command_type"`
		Status      string     `json:"status"`
		IssuedBy    uuid.UUID  `json:"issued_by"`
		IssuedAt    time.Time  `json:"issued_at"`
		ExecutedAt  *time.Time `json:"executed_at,omitempty"`
	}

	commands := []CommandEntry{}
	for rows.Next() {
		var e CommandEntry
		if err := rows.Scan(&e.ID, &e.DeviceID, &e.DeviceName, &e.CommandType, &e.Status, &e.IssuedBy, &e.IssuedAt, &e.ExecutedAt); err != nil {
			continue
		}
		commands = append(commands, e)
	}

	var total int64
	database.DB.QueryRow("SELECT COUNT(*) FROM device_commands WHERE farm_id = $1", farmID).Scan(&total)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"commands": commands,
		"total":    total,
	}, "Command history fetched")
}

// EmergencyStopHandler immediately stops all active devices on a farm
// @POST /v1/farms/:farm_id/emergency-stop
func EmergencyStopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "manager"); err != nil {
		return err
	}

	// Get all online devices
	rows, err := database.DB.Query(
		"SELECT id FROM devices WHERE farm_id = $1 AND is_active = true AND is_online = true",
		farmID,
	)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch devices")
	}
	defer rows.Close()

	stopped := 0
	for rows.Next() {
		var devID uuid.UUID
		if err := rows.Scan(&devID); err != nil {
			continue
		}
		cmdID := uuid.New()
		database.DB.Exec(`
			INSERT INTO device_commands (id, device_id, farm_id, issued_by, command_type, status, issued_at, created_at)
			VALUES ($1, $2, $3, $4, 'off', 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, cmdID, devID, farmID, userID)
		stopped++
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"stopped_devices": stopped,
		"message":         "Emergency stop issued to all online devices",
	}, "Emergency stop issued")
}

// BatchDeviceCommandHandler sends the same command to multiple devices
// @POST /v1/farms/:farm_id/devices/batch-command
func BatchDeviceCommandHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "manager"); err != nil {
		return err
	}

	var req struct {
		DeviceIDs   []uuid.UUID `json:"device_ids"`
		CommandType string      `json:"command_type"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if len(req.DeviceIDs) == 0 || req.CommandType == "" {
		return utils.BadRequest(c, "missing_fields", "device_ids and command_type are required")
	}

	type BatchResult struct {
		CommandID uuid.UUID `json:"command_id"`
		DeviceID  uuid.UUID `json:"device_id"`
		Status    string    `json:"status"`
	}

	batchID := uuid.New()
	results := []BatchResult{}
	for _, devID := range req.DeviceIDs {
		cmdID := uuid.New()
		_, err := database.DB.Exec(`
			INSERT INTO device_commands (id, device_id, farm_id, issued_by, command_type, status, issued_at, created_at)
			VALUES ($1, $2, $3, $4, $5, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, cmdID, devID, farmID, userID, req.CommandType)
		status := "queued"
		if err != nil {
			status = "error"
		}
		results = append(results, BatchResult{CommandID: cmdID, DeviceID: devID, Status: status})
	}

	return utils.SuccessResponse(c, fiber.StatusAccepted, fiber.Map{
		"command_batch_id": batchID,
		"commands":         results,
		"total":            len(results),
	}, "Batch command issued")
}
