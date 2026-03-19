package schemas

import (
	"github.com/google/uuid"
	"time"
)

// AlertHistoryEntry represents a historical alert with duration
type AlertHistoryEntry struct {
	ID              uuid.UUID  `json:"id"`
	AlertType       string     `json:"alert_type"`
	Severity        string     `json:"severity"`
	Message         string     `json:"message"`
	CreatedAt       time.Time  `json:"created_at"`
	AcknowledgedAt  *time.Time `json:"acknowledged_at,omitempty"`
	DurationSeconds int        `json:"duration_seconds"`
}

// AlertStats represents aggregate alert statistics
type AlertStats struct {
	Total         int64 `json:"total"`
	ActiveCount   int64 `json:"active_count"`
	CriticalCount int64 `json:"critical_count"`
}
