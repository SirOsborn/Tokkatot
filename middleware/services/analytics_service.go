package services

import (
	"middleware/database"
	"middleware/schemas"

	"github.com/google/uuid"
)

// AnalyticsService handles all reporting and dashboard logic
type AnalyticsService struct {
	farmService *FarmService
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		farmService: NewFarmService(),
	}
}

// GetFarmDashboard returns aggregated stats for a farm dashboard
func (s *AnalyticsService) GetFarmDashboard(userID, farmID uuid.UUID) (*schemas.DashboardResponse, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	var farmName string
	database.DB.QueryRow("SELECT name FROM farms WHERE id = $1", farmID).Scan(&farmName)

	var totalDevices, onlineDevices int64
	database.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE farm_id = $1 AND is_active = true", farmID).Scan(&totalDevices)
	database.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE farm_id = $1 AND is_active = true AND is_online = true", farmID).Scan(&onlineDevices)

	var activeAlerts, criticalAlerts, warningAlerts int64
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true", farmID).Scan(&activeAlerts)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true AND severity = 'critical'", farmID).Scan(&criticalAlerts)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true AND severity = 'warning'", farmID).Scan(&warningAlerts)

	rows, _ := database.DB.Query(`
		SELECT id, event_type, resource_id, ip_address, created_at
		FROM event_logs WHERE farm_id = $1
		ORDER BY created_at DESC LIMIT 10
	`, farmID)
	
	recentEvents := []schemas.EventEntry{}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var e schemas.EventEntry
			rows.Scan(&e.ID, &e.EventType, &e.ResourceID, &e.IPAddress, &e.Timestamp)
			recentEvents = append(recentEvents, e)
		}
	}

	resp := &schemas.DashboardResponse{}
	resp.Farm.ID = farmID
	resp.Farm.Name = farmName
	resp.DeviceStatus.Total = totalDevices
	resp.DeviceStatus.Online = onlineDevices
	resp.DeviceStatus.Offline = totalDevices - onlineDevices
	resp.Alerts.Active = activeAlerts
	resp.Alerts.Critical = criticalAlerts
	resp.Alerts.Warning = warningAlerts
	resp.RecentEvents = recentEvents

	return resp, nil
}
