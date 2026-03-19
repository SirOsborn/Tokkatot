package schemas

import (
	"middleware/models"
)

// CoopWithDevices represents a coop with its associated device count
type CoopWithDevices struct {
	models.Coop
	DeviceCount int `json:"device_count"`
}
