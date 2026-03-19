package schemas

import (
	"github.com/google/uuid"
)

// UserInfo for response objects
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    *string   `json:"email,omitempty"`
	Phone    *string   `json:"phone,omitempty"`
	Role     string    `json:"role"` // In context of a farm
	Language string    `json:"language"`
}

// SignupRequest for user registration
type SignupRequest struct {
	Email           *string `json:"email,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	Name            string  `json:"name"`
	Password         string  `json:"password"`
	RegistrationKey  *string `json:"registration_key,omitempty"` // For farmers: system-issued reg key
	FarmerID         *string `json:"farmer_id,omitempty"`        // For workers/viewers: the farmer's user ID
	NationalIDNumber *string `json:"national_id_number,omitempty"`
	Sex              *string `json:"sex,omitempty"`
	Province         *string `json:"province,omitempty"`
}

// LoginRequest for user login
type LoginRequest struct {
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Password string  `json:"password"`
}

// TokenResponse for JWT token response
type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"` // seconds
	User         UserInfo `json:"user"`
}

// PaginatedResponse for list endpoints
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// ErrorResponse for error responses
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// JSONResponse for simple success/message responses
type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
