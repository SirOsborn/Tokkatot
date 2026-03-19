package services

import (
	"middleware/database"
	"middleware/models"
	"strconv"
	"errors"

	"github.com/google/uuid"
)

// AlertService handles all business logic related to alerts
type AlertService struct {
	farmService *FarmService
}

func NewAlertService() *AlertService {
	return &AlertService{
		farmService: NewFarmService(),
	}
}

var (
	ErrAlertNotFound = errors.New("alert not found")
)

// GetFarmAlerts returns filtered alerts for a farm
func (s *AlertService) GetFarmAlerts(userID, farmID uuid.UUID, limit, offset int, isActiveOnly bool, severity string) ([]models.Alert, int64, int64, int64, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, 0, 0, 0, err
	}

	query := `SELECT id, farm_id, coop_id, alert_type, message, severity, is_active, is_acknowledged, triggered_at, resolved_at FROM alerts WHERE farm_id = $1`
	args := []interface{}{farmID}
	nextArg := 2

	if isActiveOnly {
		query += " AND is_active = true"
	}
	if severity != "all" && severity != "" {
		query += " AND severity = $" + strconv.Itoa(nextArg)
		args = append(args, severity)
		nextArg++
	}

	query += " ORDER BY triggered_at DESC LIMIT $" + strconv.Itoa(nextArg) + " OFFSET $" + strconv.Itoa(nextArg+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.FarmID, &a.CoopID, &a.AlertType, &a.Message, &a.Severity, &a.IsActive, &a.IsAcknowledged, &a.TriggeredAt, &a.ResolvedAt); err != nil {
			continue
		}
		alerts = append(alerts, a)
	}

	var total, activeCount, criticalCount int64
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1", farmID).Scan(&total)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true", farmID).Scan(&activeCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true AND severity = 'critical'", farmID).Scan(&criticalCount)

	return alerts, total, activeCount, criticalCount, nil
}

// AcknowledgeAlert marks an alert as acknowledged
func (s *AlertService) AcknowledgeAlert(userID, farmID, alertID uuid.UUID) error {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return err
	}

	res, err := database.DB.Exec("UPDATE alerts SET is_acknowledged = true, is_active = false, resolved_at = CURRENT_TIMESTAMP WHERE id = $1 AND farm_id = $2 AND is_acknowledged = false", alertID, farmID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrAlertNotFound
	}

	return nil
}
