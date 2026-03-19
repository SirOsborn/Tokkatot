package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a farmer or team member
type User struct {
	ID               uuid.UUID  `json:"id"`
	Email            *string    `json:"email,omitempty"`
	Phone            *string    `json:"phone,omitempty"`
	PhoneCountryCode *string    `json:"phone_country_code,omitempty"`
	Name             string     `json:"name"`
	PasswordHash     string     `json:"-"`
	IsActive         bool       `json:"is_active"`
	LastLogin        *time.Time `json:"last_login,omitempty"`
	NationalIDNumber *string    `json:"national_id_number,omitempty"`
	Sex              *string    `json:"sex,omitempty"`
	Province         *string    `json:"province,omitempty"`
	FullName         *string    `json:"full_name,omitempty"`
	FarmID           *uuid.UUID `json:"farm_id,omitempty"`
	FarmName         *string    `json:"farm_name,omitempty"`
	Role             string     `json:"role"` // Add role for convenience
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// UserSession represents an active authenticated session
type UserSession struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	DeviceName   *string   `json:"device_name,omitempty"`
	IPAddress    *string   `json:"ip_address,omitempty"`
	UserAgent    *string   `json:"user_agent,omitempty"`
	RefreshToken string    `json:"-"`
	LastActivity time.Time `json:"last_activity"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
