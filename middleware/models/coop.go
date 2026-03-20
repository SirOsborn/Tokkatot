package models

import (
	"time"

	"github.com/google/uuid"
)

// Coop represents a chicken house/building within a farm
type Coop struct {
	ID           uuid.UUID  `json:"id"`
	FarmID       uuid.UUID  `json:"farm_id"`
	Number       int        `json:"number"`
	Name         string     `json:"name"`
	Capacity     *int       `json:"capacity,omitempty"`
	CurrentCount *int       `json:"current_count,omitempty"`
	ChickenType  *string    `json:"chicken_type,omitempty"`
	MainDeviceID *uuid.UUID `json:"main_device_id,omitempty"`
	TempMin      *float64   `json:"temp_min,omitempty"`
	TempMax      *float64   `json:"temp_max,omitempty"`
	WaterLevelHalfThreshold *float64 `json:"water_level_half_threshold,omitempty"`
	Description  *string    `json:"description,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
