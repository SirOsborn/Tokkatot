package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// NullRawMessage is a json.RawMessage that can be scanned from a nullable SQL JSONB column.
// It handles NULL values gracefully and serializes as inline JSON (not a string).
type NullRawMessage []byte

func (n *NullRawMessage) Scan(src interface{}) error {
	if src == nil {
		*n = nil
		return nil
	}
	switch v := src.(type) {
	case []byte:
		data := make([]byte, len(v))
		copy(data, v)
		*n = data
	case string:
		*n = []byte(v)
	default:
		return fmt.Errorf("NullRawMessage: unsupported type %T", src)
	}
	return nil
}

func (n NullRawMessage) MarshalJSON() ([]byte, error) {
	if len(n) == 0 {
		return []byte("null"), nil
	}
	return json.RawMessage(n).MarshalJSON()
}

// User represents a farmer or team member
type User struct {
	ID               uuid.UUID  `json:"id"`
	Email            *string    `json:"email,omitempty"`
	Phone            *string    `json:"phone,omitempty"`
	PhoneCountryCode *string    `json:"phone_country_code,omitempty"` // "+855" for Cambodia
	PasswordHash     string     `json:"-"`
	Name             string     `json:"name"`
	Language         string     `json:"language"` // "km" or "en"
	Timezone         string     `json:"timezone"` // "Asia/Phnom_Penh"
	AvatarURL        *string    `json:"avatar_url,omitempty"`
	IsActive         bool       `json:"is_active"`
	ContactVerified  bool       `json:"contact_verified"`
	VerificationType *string    `json:"verification_type,omitempty"` // "email" or "phone"
	LastLogin        *time.Time `json:"last_login,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Farm represents a poultry farm
type Farm struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Location    *string   `json:"location,omitempty"`
	Timezone    string    `json:"timezone"`
	Latitude    *float64  `json:"latitude,omitempty"`
	Longitude   *float64  `json:"longitude,omitempty"`
	Description *string   `json:"description,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Coop represents a chicken house/building within a farm
type Coop struct {
	ID           uuid.UUID  `json:"id"`
	FarmID       uuid.UUID  `json:"farm_id"`
	Number       int        `json:"number"`                   // Coop 1, Coop 2, Coop 3, etc.
	Name         string     `json:"name"`                     // "Coop 1", "Layer House", etc.
	Capacity     *int       `json:"capacity,omitempty"`       // Max number of chickens
	CurrentCount *int       `json:"current_count,omitempty"`  // Current chicken count
	ChickenType  *string    `json:"chicken_type,omitempty"`   // "layer", "broiler", "mixed"
	MainDeviceID *uuid.UUID `json:"main_device_id,omitempty"` // Primary Raspberry Pi controller
	Description  *string    `json:"description,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// FarmUser represents a user's membership in a farm with a role
type FarmUser struct {
	ID        uuid.UUID `json:"id"`
	FarmID    uuid.UUID `json:"farm_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"` // "owner", "manager", "viewer"
	InvitedBy uuid.UUID `json:"invited_by"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Device represents an IoT device (ESP32, sensor, etc.)
type Device struct {
	ID                uuid.UUID  `json:"id"`
	FarmID            uuid.UUID  `json:"farm_id"`
	CoopID            *uuid.UUID `json:"coop_id,omitempty"` // Which coop this device belongs to
	DeviceID          string     `json:"device_id"`         // Hardware ID
	Name              string     `json:"name"`
	Type              string     `json:"type"` // "gpio", "relay", "pwm", "adc", "servo", "sensor"
	Model             *string    `json:"model,omitempty"`
	IsMainController  bool       `json:"is_main_controller"` // True if this is the Raspberry Pi for the coop
	FirmwareVersion   string     `json:"firmware_version"`
	HardwareID        string     `json:"hardware_id"` // Serial number
	Location          *string    `json:"location,omitempty"`
	IsActive          bool       `json:"is_active"`
	IsOnline          bool       `json:"is_online"`
	LastHeartbeat     *time.Time `json:"last_heartbeat,omitempty"`
	LastCommandStatus *string    `json:"last_command_status,omitempty"`
	LastCommandAt     *time.Time `json:"last_command_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// DeviceCommand represents a command sent to a device
type DeviceCommand struct {
	ID           uuid.UUID  `json:"id"`
	DeviceID     uuid.UUID  `json:"device_id"`
	FarmID       uuid.UUID  `json:"farm_id"`
	CoopID       *uuid.UUID `json:"coop_id,omitempty"` // Which coop this command is for
	IssuedBy     uuid.UUID  `json:"issued_by"`
	CommandType  string     `json:"command_type"` // "on", "off", "set_value"
	CommandValue *string    `json:"command_value,omitempty"`
	Status       string     `json:"status"` // "pending", "success", "failed", "timeout"
	Response     *string    `json:"response,omitempty"`
	IssuedAt     time.Time  `json:"issued_at"`
	ExecutedAt   *time.Time `json:"executed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Schedule represents an automated task
type Schedule struct {
	ID             uuid.UUID      `json:"id"`
	FarmID         uuid.UUID      `json:"farm_id"`
	CoopID         *uuid.UUID     `json:"coop_id,omitempty"` // Which coop this schedule applies to
	DeviceID       uuid.UUID      `json:"device_id"`
	Name           string         `json:"name"`
	ScheduleType   string         `json:"schedule_type"`             // "time_based", "duration_based", "condition_based"
	CronExpression *string        `json:"cron_expression,omitempty"` // For time_based: "0 6 * * *"
	OnDuration     *int           `json:"on_duration,omitempty"`     // For duration_based: seconds (how long to stay ON)
	OffDuration    *int           `json:"off_duration,omitempty"`    // For duration_based: seconds (how long to stay OFF)
	ConditionJSON  *string        `json:"condition_json,omitempty"`  // For condition_based: {"sensor":"temp","threshold":30}
	Action         string         `json:"action"`                    // "on", "off", "set_value"
	ActionValue    *string        `json:"action_value,omitempty"`    // Optional value (e.g., "75" for PWM)
	ActionDuration *int           `json:"action_duration,omitempty"` // For time_based: auto-turn-off after X seconds
	ActionSequence NullRawMessage `json:"action_sequence,omitempty"` // For time_based: multi-step pattern [{"action":"ON","duration":30},{"action":"OFF","duration":10}]
	Priority       int            `json:"priority"`                  // 0-10, higher = more important
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
	Status              string     `json:"status"` // "executed", "failed", "skipped"
	ExecutionDurationMs *int       `json:"execution_duration_ms,omitempty"`
	DeviceResponse      *string    `json:"device_response,omitempty"` // JSON response from device
	ErrorMessage        *string    `json:"error_message,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

// EventLog represents an audit trail entry
type EventLog struct {
	ID         uuid.UUID  `json:"id"`
	FarmID     uuid.UUID  `json:"farm_id"`
	UserID     uuid.UUID  `json:"user_id"`
	EventType  string     `json:"event_type"` // "login", "device_control", "schedule_update"
	ResourceID *uuid.UUID `json:"resource_id,omitempty"`
	OldValue   *string    `json:"old_value,omitempty"` // JSON
	NewValue   *string    `json:"new_value,omitempty"` // JSON
	IPAddress  string     `json:"ip_address"`
	CreatedAt  time.Time  `json:"created_at"`
}

// RegistrationKey for on-site account creation
type RegistrationKey struct {
	ID            uuid.UUID  `json:"id"`
	KeyCode       string     `json:"key_code"`
	FarmName      *string    `json:"farm_name,omitempty"`
	FarmLocation  *string    `json:"farm_location,omitempty"`
	CustomerName  *string    `json:"customer_name,omitempty"`
	CustomerPhone *string    `json:"customer_phone,omitempty"`
	IsUsed        bool       `json:"is_used"`
	UsedByUserID  *uuid.UUID `json:"used_by_user_id,omitempty"`
	UsedAt        *time.Time `json:"used_at,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedBy     string     `json:"created_by"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ===== REQUEST/RESPONSE MODELS =====

// SignupRequest for user registration
type SignupRequest struct {
	Email           *string `json:"email,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	Name            string  `json:"name"`
	Password        string  `json:"password"`
	Language        *string `json:"language,omitempty"`         // "km" or "en"
	RegistrationKey *string `json:"registration_key,omitempty"` // For on-site registration
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

// UserInfo for response objects
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    *string   `json:"email,omitempty"`
	Phone    *string   `json:"phone,omitempty"`
	Role     string    `json:"role"` // In context of a farm
	Language string    `json:"language"`
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

// JWTClaims for JWT token claims
type JWTClaims struct {
	UserID    uuid.UUID `json:"sub"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	FarmID    uuid.UUID `json:"farm_id"`
	Role      string    `json:"role"`
	IssuedAt  int64     `json:"iat"`
	ExpiresAt int64     `json:"exp"`
}

// Valid validates the token (required by jwt.Claims interface)
// Note: This must be a value receiver, not pointer, for jwt.NewWithClaims to work
func (j JWTClaims) Valid() error {
	if j.ExpiresAt < time.Now().Unix() {
		return jwt.NewValidationError("token expired", jwt.ValidationErrorExpired)
	}
	return nil
}
