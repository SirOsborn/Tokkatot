package schemas

import (
	"middleware/models"
	"time"
)

// CoopWithDevices represents a coop with its associated device count
type CoopWithDevices struct {
	models.Coop
	DeviceCount int `json:"device_count"`
	Temperature *float64  `json:"temperature,omitempty"`
	Humidity    *float64  `json:"humidity,omitempty"`
	LastUpdated *time.Time `json:"last_updated,omitempty"`
}
