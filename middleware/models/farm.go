package models

import (
	"time"

	"github.com/google/uuid"
)

// Farm represents a poultry farm
type Farm struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Location    *string   `json:"location,omitempty"`
	Province    *string   `json:"province,omitempty"`
	Timezone    string    `json:"timezone"`
	Latitude    *float64  `json:"latitude,omitempty"`
	Longitude   *float64  `json:"longitude,omitempty"`
	Description *string   `json:"description,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FarmUser represents a user's membership in a farm with a role
type FarmUser struct {
	ID        uuid.UUID `json:"id"`
	FarmID    uuid.UUID `json:"farm_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
	InvitedBy uuid.UUID `json:"invited_by"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
