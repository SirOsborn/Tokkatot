package schemas

import (
	"time"

	"github.com/google/uuid"
)

// WSMessage represents a standard WebSocket broadcast message
type WSMessage struct {
	Type      string      `json:"type" example:"device_update"` // "device_update", "command_update", "alert", "heartbeat"
	Timestamp time.Time   `json:"timestamp"`
	FarmID    uuid.UUID   `json:"farm_id"`
	CoopID    *uuid.UUID  `json:"coop_id,omitempty"`
	Data      interface{} `json:"data"`
}

// WSStatsResponse represents WebSocket connection statistics
type WSStatsResponse struct {
	TotalConnections int                       `json:"total_connections" example:"10"`
	YourFarms        map[uuid.UUID]int         `json:"your_farms"`
	Timestamp        time.Time                 `json:"timestamp"`
}
