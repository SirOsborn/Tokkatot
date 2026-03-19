package services

import (
	"database/sql"
	"middleware/database"
	"middleware/models"
	"strconv"
	"time"

	"github.com/google/uuid"
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

	query := `SELECT id, farm_id, coop_id, device_id, name, type, firmware_version, hardware_id, location, is_active, is_online, last_heartbeat, created_at FROM devices WHERE farm_id = $1`
	args := []interface{}{farmID}
	nextArg := 2

	if coopID != nil {
		query += " AND coop_id = $" + strconv.Itoa(nextArg)
		args = append(args, *coopID)
		nextArg++
	}

	if filter == "online" {
		query += " AND is_online = true"
	} else if filter == "offline" {
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
	database.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE farm_id = $1", farmID).Scan(&total)

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
		return d, sql.ErrNoRows
	}
	return d, err
}

// IssueCommand sends a command to a device
func (s *DeviceService) IssueCommand(userID, farmID, deviceID uuid.UUID, commandType string, commandValue *string) (*models.DeviceCommand, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return nil, err
	}

	cmdID := uuid.New()
	now := time.Now()
	cmd := &models.DeviceCommand{
		ID:           cmdID,
		FarmID:       farmID,
		DeviceID:     deviceID,
		IssuedBy:     userID,
		CommandType:  commandType,
		CommandValue: commandValue,
		Status:       "pending",
		IssuedAt:     now,
		CreatedAt:    now,
	}

	_, err := database.DB.Exec(`
		INSERT INTO device_commands (id, farm_id, device_id, issued_by, command_type, command_value, status, issued_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, cmd.ID, cmd.FarmID, cmd.DeviceID, cmd.IssuedBy, cmd.CommandType, cmd.CommandValue, cmd.Status, cmd.IssuedAt, cmd.CreatedAt)

	return cmd, err
}

// GetDeviceCommands returns last commands for a device
func (s *DeviceService) GetDeviceCommands(userID, farmID, deviceID uuid.UUID, limit int) ([]models.DeviceCommand, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(`
		SELECT id, farm_id, device_id, issued_by, command_type, command_value, status, response, issued_at, executed_at, created_at
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
		if err := rows.Scan(&c.ID, &c.FarmID, &c.DeviceID, &c.IssuedBy, &c.CommandType, &c.CommandValue, &c.Status, &c.Response, &c.IssuedAt, &c.ExecutedAt, &c.CreatedAt); err != nil {
			continue
		}
		commands = append(commands, c)
	}
	return commands, nil
}
