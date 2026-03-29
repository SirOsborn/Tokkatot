package models

import (
	"time"

	"github.com/google/uuid"
)

// Device represents an IoT device
type Device struct {
	ID                uuid.UUID  `json:"id"`
	FarmID            uuid.UUID  `json:"farm_id"`
	CoopID            *uuid.UUID `json:"coop_id,omitempty"`
	DeviceID          string     `json:"device_id"`
	Name              string     `json:"name"`
	Type              string     `json:"type"`
	Model             *string    `json:"model,omitempty"`
	IsMainController  bool       `json:"is_main_controller"`
	FirmwareVersion   string     `json:"firmware_version"`
	HardwareID        string     `json:"hardware_id"`
	Location          *string    `json:"location,omitempty"`
	IsActive          bool       `json:"is_active"`
	IsOnline          bool       `json:"is_online"`
	LastHeartbeat     *time.Time `json:"last_heartbeat,omitempty"`
	LastCommandStatus *string    `json:"last_command_status,omitempty"`
	LastCommandAt     *time.Time `json:"last_command_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// DeviceCommand represents a command sent to a device
type DeviceCommand struct {
	ID           uuid.UUID  `json:"id"`
	DeviceID     uuid.UUID  `json:"device_id"`
	FarmID       uuid.UUID  `json:"farm_id"`
	CoopID       *uuid.UUID `json:"coop_id,omitempty"`
	IssuedBy     uuid.UUID  `json:"issued_by"`
	CommandType  string     `json:"command_type"`
	CommandValue *string    `json:"command_value,omitempty"`
	ActionDuration *int    `json:"action_duration,omitempty"`
	Status       string     `json:"status"`
	Response     *string    `json:"response,omitempty"`
	IssuedAt     time.Time  `json:"issued_at"`
	ExecutedAt   *time.Time `json:"executed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	DeviceModel  *string    `json:"device_model,omitempty"`
}

// DeviceConfiguration represents a device parameter setting
type DeviceConfiguration struct {
	ID             uuid.UUID  `json:"id"`
	DeviceID       uuid.UUID  `json:"device_id"`
	ParameterName  string     `json:"parameter_name"`
	ParameterValue string     `json:"parameter_value"`
	Unit           *string    `json:"unit,omitempty"`
	MinValue       *float64   `json:"min_value,omitempty"`
	MaxValue       *float64   `json:"max_value,omitempty"`
	IsCalibrated   bool       `json:"is_calibrated"`
	CalibratedAt   *time.Time `json:"calibrated_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// DeviceReading represents a time-series sensor reading
type DeviceReading struct {
	ID         uuid.UUID `json:"id"`
	DeviceID   uuid.UUID `json:"device_id"`
	SensorType string    `json:"sensor_type"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Quality    string    `json:"quality"`
	Timestamp  time.Time `json:"timestamp"`
}
