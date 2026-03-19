package schemas

import "github.com/google/uuid"

// EventEntry represents a single entry in the event audit log
type EventEntry struct {
	ID         string  `json:"id"`
	EventType  string  `json:"event_type" example:"device_command"`
	UserID     string  `json:"user_id,omitempty"`
	ResourceID *string `json:"resource_id,omitempty"`
	IPAddress  *string `json:"ip_address,omitempty"`
	Timestamp  string  `json:"timestamp"`
}

// DataPoint represents a single aggregated sensor data point
type DataPoint struct {
	Timestamp  string  `json:"timestamp"`
	SensorType string  `json:"sensor_type" example:"temperature"`
	Avg        float64 `json:"avg" example:"25.5"`
	Min        float64 `json:"min" example:"22.0"`
	Max        float64 `json:"max" example:"31.2"`
	Unit       string  `json:"unit" example:"C"`
}

// DashboardResponse represents the overview statistics for a farm dashboard
type DashboardResponse struct {
	Farm struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	} `json:"farm"`
	DeviceStatus struct {
		Total   int64 `json:"total"`
		Online  int64 `json:"online"`
		Offline int64 `json:"offline"`
	} `json:"device_status"`
	Alerts struct {
		Active   int64 `json:"active"`
		Critical int64 `json:"critical"`
		Warning  int64 `json:"warning"`
	} `json:"alerts"`
	QuickStats struct {
		Last24hCommands int64 `json:"last_24h_commands"`
		Last24hAlerts   int64 `json:"last_24h_alerts"`
	} `json:"quick_stats"`
	RecentEvents []EventEntry `json:"recent_events"`
}
