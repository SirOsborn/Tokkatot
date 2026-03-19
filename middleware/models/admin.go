package models

import (
	"time"

	"github.com/google/uuid"
)

// EventLog represents an audit trail entry
type EventLog struct {
	ID         uuid.UUID  `json:"id"`
	FarmID     uuid.UUID  `json:"farm_id"`
	UserID     uuid.UUID  `json:"user_id"`
	EventType  string     `json:"event_type"`
	ResourceID *uuid.UUID `json:"resource_id,omitempty"`
	OldValue   *string    `json:"old_value,omitempty"`
	NewValue   *string    `json:"new_value,omitempty"`
	IPAddress  string     `json:"ip_address"`
	CreatedAt  time.Time  `json:"created_at"`
}

// RegistrationKey for on-site account creation
type RegistrationKey struct {
	ID               uuid.UUID  `json:"id"`
	KeyCode          string     `json:"key_code"`
	FarmName         *string    `json:"farm_name,omitempty"`
	CustomerPhone    *string    `json:"customer_phone,omitempty"`
	IsUsed           bool       `json:"is_used"`
	UsedByUserID     *uuid.UUID `json:"used_by_user_id,omitempty"`
	UsedAt           *time.Time `json:"used_at,omitempty"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	CreatedBy        string     `json:"created_by"`
	NationalIDNumber *string    `json:"national_id_number,omitempty"`
	FullName         *string    `json:"full_name,omitempty"`
	Sex              *string    `json:"sex,omitempty"`
	Province         *string    `json:"province,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}
