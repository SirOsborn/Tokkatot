package schemas

import (
	"time"
	"github.com/google/uuid"
)

// FarmerRow represents a farmer in the admin list
type FarmerRow struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Email            *string    `json:"email,omitempty"`
	Phone            *string    `json:"phone,omitempty"`
	Language         *string    `json:"language,omitempty"`
	IsActive         bool       `json:"is_active"`
	ContactVerified  bool       `json:"contact_verified"`
	CreatedAt        time.Time  `json:"created_at"`
	LastLogin        *time.Time `json:"last_login,omitempty"`
	FarmID           uuid.UUID  `json:"farm_id"`
	FarmName         string     `json:"farm_name"`
	NationalIDNumber *string    `json:"national_id_number,omitempty"`
	FullName         *string    `json:"full_name,omitempty"`
	Sex              *string    `json:"sex,omitempty"`
	Province         *string    `json:"province,omitempty"`
	RegKey           *string    `json:"reg_key,omitempty"`
	RegKeyUsed       *bool      `json:"reg_key_used,omitempty"`
}

// RegKeyRow represents a registration key in the admin list
type RegKeyRow struct {
	ID            uuid.UUID  `json:"id"`
	KeyCode       string     `json:"key_code"`
	FarmName      *string    `json:"farm_name,omitempty"`
	CustomerPhone *string    `json:"customer_phone,omitempty"`
	IsUsed        bool       `json:"is_used"`
	UsedAt        *time.Time `json:"used_at,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedBy     *string    `json:"created_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UsedByName    *string    `json:"used_by_name,omitempty"`
	UsedByPhone   *string    `json:"used_by_phone,omitempty"`
}

// ViewerRow represents a viewer/worker in the admin list
type ViewerRow struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      *string   `json:"email,omitempty"`
	Phone      *string   `json:"phone,omitempty"`
	Language   *string   `json:"language,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	FarmID     uuid.UUID `json:"farm_id"`
	FarmName   string    `json:"farm_name"`
	FarmerName *string   `json:"farmer_name,omitempty"`
}
