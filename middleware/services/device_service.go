package services

import (
	"database/sql"
	"errors"
	"middleware/database"
	"middleware/models"
	"middleware/schemas"
	"middleware/utils"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDeviceNotFound  = errors.New("device_not_found")
	ErrCommandNotFound = errors.New("command_not_found")
)

// DeviceService handles all business logic related to device management
type DeviceService struct {
	farmService *FarmService
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		farmService: NewFarmService(),
	}
}

// ListDevices returns all devices for a farm with optional coop filter
func (s *DeviceService) ListDevices(userID, farmID uuid.UUID, coopID *uuid.UUID, filter string, limit, offset int) ([]models.Device, int64, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, farm_id, coop_id, device_id, name, type, firmware_version, hardware_id, location, is_active, is_online, last_heartbeat, created_at FROM devices WHERE farm_id = $1 AND is_active = true`
	args := []interface{}{farmID}
	nextArg := 2

	if coopID != nil {
		query += " AND coop_id = $" + strconv.Itoa(nextArg)
		args = append(args, *coopID)
		nextArg++
	}

	switch filter {
case "online":
		query += " AND is_online = true"
	case "offline":
		query += " AND is_online = false"
	}

	query += " ORDER BY name ASC LIMIT $" + strconv.Itoa(nextArg) + " OFFSET $" + strconv.Itoa(nextArg+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var d models.Device
		if err := rows.Scan(&d.ID, &d.FarmID, &d.CoopID, &d.DeviceID, &d.Name, &d.Type, &d.FirmwareVersion, &d.HardwareID, &d.Location, &d.IsActive, &d.IsOnline, &d.LastHeartbeat, &d.CreatedAt); err != nil {
			continue
		}
		devices = append(devices, d)
	}

	var total int64
	database.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE farm_id = $1 AND is_active = true", farmID).Scan(&total)

	return devices, total, nil
}

// GetDevice retrieves a single device by ID
func (s *DeviceService) GetDevice(userID, farmID, deviceID uuid.UUID) (models.Device, error) {
	var d models.Device
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return d, err
	}

	err := database.DB.QueryRow(`
		SELECT id, farm_id, coop_id, device_id, name, type, firmware_version, hardware_id, location, is_active, is_online, last_heartbeat, created_at 
		FROM devices 
		WHERE id = $1 AND farm_id = $2
	`, deviceID, farmID).Scan(&d.ID, &d.FarmID, &d.CoopID, &d.DeviceID, &d.Name, &d.Type, &d.FirmwareVersion, &d.HardwareID, &d.Location, &d.IsActive, &d.IsOnline, &d.LastHeartbeat, &d.CreatedAt)

	if err == sql.ErrNoRows {
		return d, ErrDeviceNotFound
	}
	return d, err
}

// AddDevice registers a new device
func (s *DeviceService) AddDevice(userID, farmID uuid.UUID, req schemas.AddDeviceRequest) (*models.Device, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}

	if req.CoopID != nil {
		var exists bool
		if err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM coops WHERE id = $1 AND farm_id = $2)", *req.CoopID, farmID).Scan(&exists); err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrCoopNotFound
		}
	}

	now := time.Now()
	deviceID := uuid.New()

	_, err := database.DB.Exec(`
		INSERT INTO devices (id, farm_id, coop_id, device_id, name, type, model, is_main_controller, firmware_version, hardware_id, location, is_active, is_online, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true, false, $12, $12)
	`, deviceID, farmID, req.CoopID, req.DeviceID, req.Name, req.Type, req.Model, req.IsMainController, req.FirmwareVersion, req.HardwareID, req.Location, now)
	if err != nil {
		return nil, err
	}

	return &models.Device{
		ID:               deviceID,
		FarmID:           farmID,
		CoopID:           req.CoopID,
		DeviceID:         req.DeviceID,
		Name:             req.Name,
		Type:             req.Type,
		Model:            req.Model,
		IsMainController: req.IsMainController,
		FirmwareVersion:  req.FirmwareVersion,
		HardwareID:       req.HardwareID,
		Location:         req.Location,
		IsActive:         true,
		IsOnline:         false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// UpdateDevice updates device metadata
func (s *DeviceService) UpdateDevice(userID, farmID, deviceID uuid.UUID, req schemas.UpdateDeviceRequest) (*models.Device, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}

	var d models.Device
	err := database.DB.QueryRow(`
		UPDATE devices SET
			name = COALESCE($1, name),
			location = COALESCE($2, location),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND farm_id = $4
		RETURNING id, farm_id, coop_id, device_id, name, type, model, is_main_controller, firmware_version, hardware_id, location, is_active, is_online, last_heartbeat, last_command_status, last_command_at, created_at, updated_at
	`, req.Name, req.Location, deviceID, farmID).Scan(
		&d.ID, &d.FarmID, &d.CoopID, &d.DeviceID, &d.Name, &d.Type, &d.Model, &d.IsMainController, &d.FirmwareVersion, &d.HardwareID,
		&d.Location, &d.IsActive, &d.IsOnline, &d.LastHeartbeat, &d.LastCommandStatus, &d.LastCommandAt, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrDeviceNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// DeleteDevice deactivates a device
func (s *DeviceService) DeleteDevice(userID, farmID, deviceID uuid.UUID) error {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return err
	}
	res, err := database.DB.Exec(`
		UPDATE devices SET is_active = false, is_online = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND farm_id = $2
	`, deviceID, farmID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

// GetDeviceStatus returns live status info
func (s *DeviceService) GetDeviceStatus(userID, farmID, deviceID uuid.UUID) (*schemas.DeviceStatusResponse, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	var status schemas.DeviceStatusResponse
	err := database.DB.QueryRow(`
		SELECT id, is_online, last_heartbeat, last_command_status, last_command_at
		FROM devices WHERE id = $1 AND farm_id = $2
	`, deviceID, farmID).Scan(&status.DeviceID, &status.IsOnline, &status.LastHeartbeat, &status.LastCommandStatus, &status.LastCommandAt)
	if err == sql.ErrNoRows {
		return nil, ErrDeviceNotFound
	}
	if err != nil {
		return nil, err
	}

	// Latest reading (optional)
	var value sql.NullFloat64
	var unit sql.NullString
	_ = database.DB.QueryRow(`
		SELECT value, unit FROM device_readings
		WHERE device_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`, deviceID).Scan(&value, &unit)
	if value.Valid {
		status.CurrentValue = &value.Float64
	}
	if unit.Valid {
		status.Unit = &unit.String
	}

	return &status, nil
}

// GetDeviceHistory returns sensor readings
func (s *DeviceService) GetDeviceHistory(userID, farmID, deviceID uuid.UUID, sensorType string, limit, offset int) ([]models.DeviceReading, int64, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, 0, err
	}

	var exists bool
	if err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM devices WHERE id = $1 AND farm_id = $2)", deviceID, farmID).Scan(&exists); err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, ErrDeviceNotFound
	}

	query := `SELECT id, device_id, sensor_type, value, unit, quality, timestamp FROM device_readings WHERE device_id = $1`
	args := []interface{}{deviceID}
	if sensorType != "" {
		query += " AND sensor_type = $2"
		args = append(args, sensorType)
	}
	query += " ORDER BY timestamp DESC LIMIT $3 OFFSET $4"
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var readings []models.DeviceReading
	for rows.Next() {
		var r models.DeviceReading
		if err := rows.Scan(&r.ID, &r.DeviceID, &r.SensorType, &r.Value, &r.Unit, &r.Quality, &r.Timestamp); err != nil {
			continue
		}
		readings = append(readings, r)
	}

	var total int64
	if sensorType == "" {
		database.DB.QueryRow("SELECT COUNT(*) FROM device_readings WHERE device_id = $1", deviceID).Scan(&total)
	} else {
		database.DB.QueryRow("SELECT COUNT(*) FROM device_readings WHERE device_id = $1 AND sensor_type = $2", deviceID, sensorType).Scan(&total)
	}

	return readings, total, nil
}

// GetDeviceConfig returns configuration parameters for a device
func (s *DeviceService) GetDeviceConfig(userID, farmID, deviceID uuid.UUID) ([]models.DeviceConfiguration, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	// Ensure device belongs to farm
	var exists bool
	if err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM devices WHERE id = $1 AND farm_id = $2)", deviceID, farmID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrDeviceNotFound
	}

	rows, err := database.DB.Query(`
		SELECT id, device_id, parameter_name, parameter_value, unit, min_value, max_value, is_calibrated, calibrated_at, created_at, updated_at
		FROM device_configurations
		WHERE device_id = $1
		ORDER BY parameter_name ASC
	`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cfgs []models.DeviceConfiguration
	for rows.Next() {
		var c models.DeviceConfiguration
		if err := rows.Scan(&c.ID, &c.DeviceID, &c.ParameterName, &c.ParameterValue, &c.Unit, &c.MinValue, &c.MaxValue, &c.IsCalibrated, &c.CalibratedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
			continue
		}
		cfgs = append(cfgs, c)
	}
	return cfgs, nil
}

// UpdateDeviceConfig upserts a configuration parameter
func (s *DeviceService) UpdateDeviceConfig(userID, farmID, deviceID uuid.UUID, name, value string, unit *string, minValue, maxValue *float64) (*models.DeviceConfiguration, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return nil, err
	}

	// Ensure device belongs to farm
	var exists bool
	if err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM devices WHERE id = $1 AND farm_id = $2)", deviceID, farmID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrDeviceNotFound
	}

	var cfg models.DeviceConfiguration
	err := database.DB.QueryRow(`
		INSERT INTO device_configurations (id, device_id, parameter_name, parameter_value, unit, min_value, max_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (device_id, parameter_name)
		DO UPDATE SET parameter_value = $4, unit = $5, min_value = $6, max_value = $7, updated_at = CURRENT_TIMESTAMP
		RETURNING id, device_id, parameter_name, parameter_value, unit, min_value, max_value, is_calibrated, calibrated_at, created_at, updated_at
	`, uuid.New(), deviceID, name, value, unit, minValue, maxValue).Scan(
		&cfg.ID, &cfg.DeviceID, &cfg.ParameterName, &cfg.ParameterValue, &cfg.Unit, &cfg.MinValue, &cfg.MaxValue, &cfg.IsCalibrated, &cfg.CalibratedAt, &cfg.CreatedAt, &cfg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// CalibrateDevice marks a parameter as calibrated
func (s *DeviceService) CalibrateDevice(userID, farmID, deviceID uuid.UUID, name string) (*models.DeviceConfiguration, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return nil, err
	}

	var cfg models.DeviceConfiguration
	err := database.DB.QueryRow(`
		UPDATE device_configurations
		SET is_calibrated = true, calibrated_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE device_id = $1 AND parameter_name = $2
		RETURNING id, device_id, parameter_name, parameter_value, unit, min_value, max_value, is_calibrated, calibrated_at, created_at, updated_at
	`, deviceID, name).Scan(&cfg.ID, &cfg.DeviceID, &cfg.ParameterName, &cfg.ParameterValue, &cfg.Unit, &cfg.MinValue, &cfg.MaxValue, &cfg.IsCalibrated, &cfg.CalibratedAt, &cfg.CreatedAt, &cfg.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrDeviceNotFound
	}
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// IssueCommand sends a command to a device
func (s *DeviceService) IssueCommand(userID, farmID, deviceID uuid.UUID, commandType string, commandValue *string, actionDuration *int) (*models.DeviceCommand, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return nil, err
	}

	cmdID := uuid.New()
	now := time.Now()
	cmd := &models.DeviceCommand{
		ID:             cmdID,
		FarmID:       farmID,
		DeviceID:     deviceID,
		IssuedBy:     userID,
		CommandType:  commandType,
		CommandValue: commandValue,
		ActionDuration: actionDuration,
		Status:       "pending",
		IssuedAt:     now,
		CreatedAt:    now,
	}

	_, err := database.DB.Exec(`
		INSERT INTO device_commands (id, farm_id, device_id, issued_by, command_type, command_value, action_duration, status, issued_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, cmd.ID, cmd.FarmID, cmd.DeviceID, cmd.IssuedBy, cmd.CommandType, cmd.CommandValue, cmd.ActionDuration, cmd.Status, cmd.IssuedAt, cmd.CreatedAt)
	if err != nil {
		return nil, err
	}

	_, _ = database.DB.Exec(`
		UPDATE devices SET last_command_status = $1, last_command_at = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, cmd.Status, cmd.IssuedAt, cmd.DeviceID)

	return cmd, err
}

// GetCommandStatus returns a single command
func (s *DeviceService) GetCommandStatus(userID, farmID, commandID uuid.UUID) (*models.DeviceCommand, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	var c models.DeviceCommand
	err := database.DB.QueryRow(`
		SELECT id, device_id, farm_id, coop_id, issued_by, command_type, command_value, action_duration, status, response, issued_at, executed_at, created_at
		FROM device_commands
		WHERE id = $1 AND farm_id = $2
	`, commandID, farmID).Scan(&c.ID, &c.DeviceID, &c.FarmID, &c.CoopID, &c.IssuedBy, &c.CommandType, &c.CommandValue, &c.ActionDuration, &c.Status, &c.Response, &c.IssuedAt, &c.ExecutedAt, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrCommandNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// CancelCommand marks a pending command as failed
func (s *DeviceService) CancelCommand(userID, farmID, commandID uuid.UUID) error {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return err
	}
	res, err := database.DB.Exec(`
		UPDATE device_commands
		SET status = 'failed', response = 'cancelled', executed_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND farm_id = $2 AND status = 'pending'
	`, commandID, farmID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrCommandNotFound
	}
	return nil
}

// GetFarmCommandHistory returns recent commands for a farm
func (s *DeviceService) GetFarmCommandHistory(userID, farmID uuid.UUID, limit int) ([]schemas.CommandEntry, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(`
		SELECT dc.id, dc.device_id, d.name, dc.command_type, dc.status, dc.issued_by, dc.issued_at, dc.executed_at
		FROM device_commands dc
		JOIN devices d ON dc.device_id = d.id
		WHERE dc.farm_id = $1
		ORDER BY dc.created_at DESC
		LIMIT $2
	`, farmID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []schemas.CommandEntry
	for rows.Next() {
		var e schemas.CommandEntry
		if err := rows.Scan(&e.ID, &e.DeviceID, &e.DeviceName, &e.CommandType, &e.Status, &e.IssuedBy, &e.IssuedAt, &e.ExecutedAt); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// BatchCommands issues a command to multiple devices
func (s *DeviceService) BatchCommands(userID, farmID uuid.UUID, req schemas.BatchCommandRequest) ([]schemas.BatchResult, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return nil, err
	}

	var results []schemas.BatchResult
	for _, dID := range req.DeviceIDs {
		cmd, err := s.IssueCommand(userID, farmID, dID, req.CommandType, req.Parameters, req.ActionDuration)
		if err != nil {
			results = append(results, schemas.BatchResult{CommandID: uuid.Nil, DeviceID: dID, Status: "failed"})
			continue
		}
		results = append(results, schemas.BatchResult{CommandID: cmd.ID, DeviceID: dID, Status: cmd.Status})
	}
	return results, nil
}

// EmergencyStop issues an "off" command to all active devices in a farm
func (s *DeviceService) EmergencyStop(userID, farmID uuid.UUID) (int, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return 0, err
	}

	rows, err := database.DB.Query(`
		SELECT id FROM devices WHERE farm_id = $1 AND is_active = true
	`, farmID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var dID uuid.UUID
		if err := rows.Scan(&dID); err != nil {
			continue
		}
		_, _ = s.IssueCommand(userID, farmID, dID, "off", nil, nil)
		count++
	}
	return count, nil
}

// UpdateHeartbeat updates a device heartbeat by hardware_id
func (s *DeviceService) UpdateHeartbeat(hardwareID string, status string, response *string, ipAddress string) error {
	res, err := database.DB.Exec(`
		UPDATE devices
		SET is_online = true,
			last_heartbeat = CURRENT_TIMESTAMP,
			last_command_status = $1::TEXT,
			last_command_at = CASE WHEN $1::TEXT IS NOT NULL THEN CURRENT_TIMESTAMP ELSE last_command_at END,
			response = $2::TEXT
		WHERE hardware_id = $3::TEXT
	`, status, response, hardwareID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Device not assigned to any farm yet. Register as unassigned for discovery.
		_, err = database.DB.Exec(`
			INSERT INTO unassigned_gateways (hardware_id, ip_address, last_seen)
			VALUES ($1::TEXT, $2::TEXT, CURRENT_TIMESTAMP)
			ON CONFLICT (hardware_id) DO UPDATE 
			SET ip_address = $2::TEXT, last_seen = CURRENT_TIMESTAMP
		`, hardwareID, ipAddress)
		return err
	}
	return nil
}

// GetUnassignedGateways returns a list of gateways that have checked in but aren't assigned
func (s *DeviceService) GetUnassignedGateways() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT hardware_id, ip_address, last_seen, created_at
		FROM unassigned_gateways
		WHERE hardware_id NOT IN (SELECT hardware_id FROM devices)
		ORDER BY last_seen DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var hID, ip string
		var lastSeen, createdAt time.Time
		if err := rows.Scan(&hID, &ip, &lastSeen, &createdAt); err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"hardware_id": hID,
			"ip_address":  ip,
			"last_seen":   lastSeen,
			"created_at":  createdAt,
		})
	}
	return result, nil
}

// AssignGateway links a hardware ID to a farm and coop
func (s *DeviceService) AssignGateway(hardwareID string, farmID uuid.UUID, coopID *uuid.UUID, name string) (string, *uuid.UUID, error) {
	// If caller didn't provide a coop, default to the first active coop of the farm.
	// This makes the gateway visible in coop-scoped UIs and aligns with telemetry routes.
	if coopID == nil || *coopID == uuid.Nil {
		var firstCoopID uuid.UUID
		if err := database.DB.QueryRow(`
			SELECT id
			FROM coops
			WHERE farm_id = $1 AND is_active = true
			ORDER BY number ASC, created_at ASC
			LIMIT 1
		`, farmID).Scan(&firstCoopID); err == nil && firstCoopID != uuid.Nil {
			coopID = &firstCoopID
		}
	}

	// Check if a gateway device row already exists for this hardware_id.
	var existingDeviceID uuid.UUID
	var exists bool
	_ = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM devices WHERE hardware_id = $1)", hardwareID).Scan(&exists)

	deviceID := uuid.New()
	if exists {
		// Re-assignment / re-tokenization path: keep the same device id.
		_ = database.DB.QueryRow("SELECT id FROM devices WHERE hardware_id = $1 LIMIT 1", hardwareID).Scan(&existingDeviceID)
		if existingDeviceID != uuid.Nil {
			deviceID = existingDeviceID
		}

		_, err := database.DB.Exec(`
			UPDATE devices
			SET farm_id = $1,
				coop_id = $2,
				name = $3,
				is_active = true,
				is_main_controller = true,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $4
		`, farmID, coopID, name, deviceID)
		if err != nil {
			return "", nil, err
		}
	} else {
		// Create the device record
		// NOTE: Gateway acts as a 'sensor' or 'relay' container but we usually register the RPi as the main controller.
		_, err := database.DB.Exec(`
			INSERT INTO devices (id, farm_id, coop_id, device_id, hardware_id, name, type, firmware_version, is_main_controller)
			VALUES ($1, $2, $3, $4, $5, $6, 'sensor', '1.0.0', true)
		`, deviceID, farmID, coopID, hardwareID, hardwareID, name)
		if err != nil {
			return "", nil, err
		}
	}

	// Remove from unassigned
	_, _ = database.DB.Exec("DELETE FROM unassigned_gateways WHERE hardware_id = $1", hardwareID)

	// Create (or rotate) a gateway token so the Pi can authenticate with `X-Gateway-Token`.
	// Token is returned once to the caller to be put on the Pi.
	var ownerID uuid.UUID
	if err := database.DB.QueryRow("SELECT owner_id FROM farms WHERE id = $1", farmID).Scan(&ownerID); err != nil {
		return "", nil, err
	}

	// Deactivate any existing tokens for this device (token rotation).
	_, _ = database.DB.Exec(`UPDATE gateway_tokens SET is_active = false WHERE device_id = $1`, deviceID)

	rawToken := strings.ReplaceAll(uuid.New().String(), "-", "")
	tokenHash := utils.HashToken(rawToken)

	_, err := database.DB.Exec(`
		INSERT INTO gateway_tokens (farm_id, device_id, user_id, token_hash, name, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
	`, farmID, deviceID, ownerID, tokenHash, "Gateway ("+hardwareID+")")
	if err != nil {
		return "", nil, err
	}

	return rawToken, coopID, nil
}

// GetDeviceCommands returns last commands for a device
func (s *DeviceService) GetDeviceCommands(userID, farmID, deviceID uuid.UUID, limit int) ([]models.DeviceCommand, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(`
		SELECT id, farm_id, device_id, issued_by, command_type, command_value, action_duration, status, response, issued_at, executed_at, created_at
		FROM device_commands
		WHERE device_id = $1 AND farm_id = $2
		ORDER BY created_at DESC
		LIMIT $3
	`, deviceID, farmID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []models.DeviceCommand
	for rows.Next() {
		var c models.DeviceCommand
		if err := rows.Scan(&c.ID, &c.FarmID, &c.DeviceID, &c.IssuedBy, &c.CommandType, &c.CommandValue, &c.ActionDuration, &c.Status, &c.Response, &c.IssuedAt, &c.ExecutedAt, &c.CreatedAt); err != nil {
			continue
		}
		commands = append(commands, c)
	}
	return commands, nil
}
// GetPendingCommands returns all pending commands for a specific gateway/hardware
func (s *DeviceService) GetPendingCommands(hardwareID string) ([]models.DeviceCommand, error) {
	rows, err := database.DB.Query(`
		SELECT dc.id, dc.device_id, dc.farm_id, dc.coop_id, dc.issued_by, dc.command_type, dc.command_value, dc.action_duration, dc.status, dc.response, dc.issued_at, dc.executed_at, dc.created_at
		FROM device_commands dc
		JOIN devices d ON dc.device_id = d.id
		WHERE d.hardware_id = $1 AND dc.status = 'pending'
		ORDER BY dc.created_at ASC
	`, hardwareID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []models.DeviceCommand
	for rows.Next() {
		var c models.DeviceCommand
		if err := rows.Scan(&c.ID, &c.DeviceID, &c.FarmID, &c.CoopID, &c.IssuedBy, &c.CommandType, &c.CommandValue, &c.ActionDuration, &c.Status, &c.Response, &c.IssuedAt, &c.ExecutedAt, &c.CreatedAt); err != nil {
			continue
		}
		commands = append(commands, c)
	}
	return commands, nil
}

// UpdateCommandStatus updates the status and response of a command
func (s *DeviceService) UpdateCommandStatus(commandID uuid.UUID, status, response string) error {
	now := time.Now()
	res, err := database.DB.Exec(`
		UPDATE device_commands
		SET status = $1, response = $2, executed_at = $3, updated_at = $3
		WHERE id = $4
	`, status, response, now, commandID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrCommandNotFound
	}

	// Also update the device's last command status for quick status checks
	var deviceID uuid.UUID
	_ = database.DB.QueryRow("SELECT device_id FROM device_commands WHERE id = $1", commandID).Scan(&deviceID)
	if deviceID != uuid.Nil {
		_, _ = database.DB.Exec(`
			UPDATE devices SET last_command_status = $1, last_command_at = $2, updated_at = $2
			WHERE id = $3
		`, status, now, deviceID)
	}

	return nil
}
