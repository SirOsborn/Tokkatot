package schemas

// UpdateProfileRequest represents the request to update user profile
type UpdateProfileRequest struct {
	Name             *string `json:"name,omitempty" example:"John Doe"`
	Email            *string `json:"email,omitempty" example:"john@example.com"`
	Phone            *string `json:"phone,omitempty" example:"+85512345678"`
	NationalIDNumber *string `json:"national_id_number,omitempty" example:"123456789"`
	Sex              *string `json:"sex,omitempty" example:"male"`
	Province         *string `json:"province,omitempty" example:"Kandal"`
}

// ChangePasswordRequest represents the request to change user password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" example:"oldpassword123"`
	NewPassword     string `json:"new_password" example:"newpassword123"`
}

// SessionInfo represents an active user session
type SessionInfo struct {
	ID           string  `json:"id"`
	DeviceName   *string `json:"device_name"`
	IPAddress    *string `json:"ip_address"`
	UserAgent    *string `json:"user_agent"`
	LastActivity string  `json:"last_activity"`
	ExpiresAt    string  `json:"expires_at"`
}

// ActivityEntry represents a single entry in the user activity log
type ActivityEntry struct {
	ID        string  `json:"id"`
	EventType string  `json:"event_type"`
	IPAddress *string `json:"ip_address"`
	Timestamp string  `json:"timestamp"`
}
