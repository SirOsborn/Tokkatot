package models

import (
	"time"

	"github.com/google/uuid"
)

// Alert represents a monitoring alert
type Alert struct {
	ID             uuid.UUID  `json:"id"`
	FarmID         uuid.UUID  `json:"farm_id"`
	DeviceID       *uuid.UUID `json:"device_id,omitempty"`
	CoopID         *uuid.UUID `json:"coop_id,omitempty"`
	CoopName       string     `json:"coop_name,omitempty"`
	AlertType      string     `json:"alert_type"`
	Severity       string     `json:"severity"`
	Message        string     `json:"message"`
	ThresholdValue *float64   `json:"threshold_value,omitempty"`
	ActualValue    *float64   `json:"actual_value,omitempty"`
	IsActive       bool       `json:"is_active"`
	IsAcknowledged bool       `json:"is_acknowledged"`
	TriggeredAt    time.Time  `json:"triggered_at"`
	AcknowledgedBy *uuid.UUID `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// AlertSubscription represents user's alert notification preferences
type AlertSubscription struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	AlertType       string    `json:"alert_type"`
	Channel         string    `json:"channel"`
	IsEnabled       bool      `json:"is_enabled"`
	QuietHoursStart *string   `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd   *string   `json:"quiet_hours_end,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
