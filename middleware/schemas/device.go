package schemas

import (
	"time"

	"github.com/google/uuid"
)

// DeviceCommandRequest represents a command sent to a device
type DeviceCommandRequest struct {
	CommandType  string  `json:"command_type" example:"on"`
	CommandValue *string `json:"command_value,omitempty" example:"100"`
}

// AddDeviceRequest represents the request to add a new device
type AddDeviceRequest struct {
	DeviceID         string     `json:"device_id" example:"GPIO_1"`
	Name             string     `json:"name" example:"Main Pump"`
	Type             string     `json:"type" example:"relay"`
	Model            *string    `json:"model,omitempty" example:"ESP32"`
	FirmwareVersion  string     `json:"firmware_version" example:"1.0.0"`
	HardwareID       string     `json:"hardware_id" example:"HW_ABC_123"`
	Location         *string    `json:"location,omitempty" example:"Coop 1 East"`
	CoopID           *uuid.UUID `json:"coop_id,omitempty"`
	IsMainController bool       `json:"is_main_controller" example:"false"`
}

// UpdateDeviceRequest represents the request to update device details
type UpdateDeviceRequest struct {
	Name     *string `json:"name,omitempty" example:"Back Pump"`
	Location *string `json:"location,omitempty" example:"Coop 1 West"`
}

// DeviceStatusResponse represents the real-time status of a device
type DeviceStatusResponse struct {
	DeviceID           uuid.UUID  `json:"device_id"`
	IsOnline           bool       `json:"is_online"`
	LastHeartbeat      *time.Time `json:"last_heartbeat"`
	LastCommandStatus *string    `json:"last_command_status"`
	LastCommandAt     *time.Time `json:"last_command_at"`
	CurrentValue       *float64   `json:"current_value"`
	Unit               *string    `json:"unit"`
}

// CommandEntry represents a single command record in history
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

// BatchResult represents the result of a single command in a batch
type BatchResult struct {
	CommandID uuid.UUID `json:"command_id"`
	DeviceID  uuid.UUID `json:"device_id"`
	Status    string    `json:"status"`
}

// BatchCommandRequest represents a request to issue commands to multiple devices
type BatchCommandRequest struct {
	DeviceIDs   []uuid.UUID `json:"device_ids"`
	CommandType string      `json:"command_type"`
	Parameters  *string     `json:"parameters,omitempty"`
}
