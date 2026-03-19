package schemas

import (
	"encoding/json"
	"middleware/models"

	"github.com/google/uuid"
)

// CreateScheduleRequest represents the request to create an automated schedule
type CreateScheduleRequest struct {
	DeviceID       uuid.UUID       `json:"device_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CoopID         *uuid.UUID      `json:"coop_id,omitempty"`
	Name           string          `json:"name" example:"Daily Watering"`
	ScheduleType   string          `json:"schedule_type" example:"time_based"` // "time_based", "duration_based", "condition_based"
	CronExpression *string         `json:"cron_expression,omitempty" example:"0 8 * * *"`
	OnDuration     *int            `json:"on_duration,omitempty" example:"3600"`
	OffDuration    *int            `json:"off_duration,omitempty"`
	ConditionJSON  *string         `json:"condition_json,omitempty"`
	Action         string          `json:"action" example:"on"`
	ActionValue    *string         `json:"action_value,omitempty"`
	ActionDuration *int            `json:"action_duration,omitempty"`
	ActionSequence json.RawMessage `json:"action_sequence,omitempty"`
	Priority       *int            `json:"priority,omitempty"`
	IsActive       *bool           `json:"is_active,omitempty"`
}

// UpdateScheduleRequest represents the request to update an automated schedule
type UpdateScheduleRequest struct {
	Name           *string         `json:"name,omitempty"`
	ScheduleType   *string         `json:"schedule_type,omitempty"`
	CronExpression *string         `json:"cron_expression,omitempty"`
	OnDuration     *int            `json:"on_duration,omitempty"`
	OffDuration    *int            `json:"off_duration,omitempty"`
	ConditionJSON  *string         `json:"condition_json,omitempty"`
	Action         *string         `json:"action,omitempty"`
	ActionValue    *string         `json:"action_value,omitempty"`
	ActionDuration *int            `json:"action_duration,omitempty"`
	ActionSequence json.RawMessage `json:"action_sequence,omitempty"`
	Priority       *int            `json:"priority,omitempty"`
	IsActive       *bool           `json:"is_active,omitempty"`
}

// ScheduleWithDevice represents a schedule including the associated device name
type ScheduleWithDevice struct {
	models.Schedule
	DeviceName *string `json:"device_name,omitempty"`
}
