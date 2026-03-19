package models

import (
	"time"

	"github.com/google/uuid"
)

// Schedule represents an automated task
type Schedule struct {
	ID             uuid.UUID      `json:"id"`
	FarmID         uuid.UUID      `json:"farm_id"`
	CoopID         *uuid.UUID     `json:"coop_id,omitempty"`
	DeviceID       uuid.UUID      `json:"device_id"`
	Name           string         `json:"name"`
	ScheduleType   string         `json:"schedule_type"`
	CronExpression *string        `json:"cron_expression,omitempty"`
	OnDuration     *int           `json:"on_duration,omitempty"`
	OffDuration    *int           `json:"off_duration,omitempty"`
	ConditionJSON  *string        `json:"condition_json,omitempty"`
	Action         string         `json:"action"`
	ActionValue    *string        `json:"action_value,omitempty"`
	ActionDuration *int           `json:"action_duration,omitempty"`
	ActionSequence NullRawMessage `json:"action_sequence,omitempty"`
	Priority       int            `json:"priority"`
	IsActive       bool           `json:"is_active"`
	NextExecution  *time.Time     `json:"next_execution,omitempty"`
	LastExecution  *time.Time     `json:"last_execution,omitempty"`
	ExecutionCount int            `json:"execution_count"`
	CreatedBy      uuid.UUID      `json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// ScheduleExecution represents a log of schedule execution
type ScheduleExecution struct {
	ID                  uuid.UUID  `json:"id"`
	ScheduleID          uuid.UUID  `json:"schedule_id"`
	DeviceID            uuid.UUID  `json:"device_id"`
	ScheduledTime       time.Time  `json:"scheduled_time"`
	ActualExecutionTime *time.Time `json:"actual_execution_time,omitempty"`
	Status              string     `json:"status"`
	ExecutionDurationMs *int       `json:"execution_duration_ms,omitempty"`
	DeviceResponse      *string    `json:"device_response,omitempty"`
	ErrorMessage        *string    `json:"error_message,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}
