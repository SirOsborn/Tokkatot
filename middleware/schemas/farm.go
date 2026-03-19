package schemas

import (
	"middleware/models"
	"github.com/google/uuid"
	"time"
)

// FarmWithRole represents a farm with the user's role in it
type FarmWithRole struct {
	models.Farm
	Role      string `json:"role"`
	CoopCount int    `json:"coop_count"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateFarmRequest for creating a new farm
type CreateFarmRequest struct {
	Name        string `json:"name" example:"Sunny Valley Farm"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
}

// UpdateFarmRequest for updating an existing farm
type UpdateFarmRequest struct {
	Name        *string `json:"name,omitempty"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
}

// MemberInfo represents a member of a farm
type MemberInfo struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Role      string    `json:"role"`
	InvitedBy uuid.UUID `json:"invited_by"`
	JoinedAt  string    `json:"joined_at"`
}
