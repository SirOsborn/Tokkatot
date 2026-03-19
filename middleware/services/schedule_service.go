package services

import (
	"database/sql"
	"middleware/database"
	"middleware/models"
	"time"

	"github.com/google/uuid"
)

// ScheduleService handles all business logic related to automated schedules
type ScheduleService struct {
	farmService *FarmService
}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{
		farmService: NewFarmService(),
	}
}

// ListSchedules returns all schedules for a farm with optional coop filter
func (s *ScheduleService) ListSchedules(userID, farmID uuid.UUID, coopID *uuid.UUID) ([]models.Schedule, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	query := `SELECT id, farm_id, coop_id, device_id, name, schedule_type, cron_expression, on_duration, off_duration, condition_json, action, action_value, action_duration, priority, is_active, created_at FROM schedules WHERE farm_id = $1`
	args := []interface{}{farmID}

	if coopID != nil {
		query += " AND coop_id = $2"
		args = append(args, *coopID)
	}

	query += " ORDER BY priority DESC, name ASC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var sc models.Schedule
		if err := rows.Scan(&sc.ID, &sc.FarmID, &sc.CoopID, &sc.DeviceID, &sc.Name, &sc.ScheduleType, &sc.CronExpression, 
			&sc.OnDuration, &sc.OffDuration, &sc.ConditionJSON, &sc.Action, &sc.ActionValue, &sc.ActionDuration, 
			&sc.Priority, &sc.IsActive, &sc.CreatedAt); err != nil {
			continue
		}
		schedules = append(schedules, sc)
	}

	return schedules, nil
}

// CreateSchedule creates a new schedule
func (s *ScheduleService) CreateSchedule(userID, farmID uuid.UUID, req models.Schedule) (*models.Schedule, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return nil, err
	}

	req.ID = uuid.New()
	req.FarmID = farmID
	req.CreatedAt = time.Now()

	_, err := database.DB.Exec(`
		INSERT INTO schedules (id, farm_id, coop_id, device_id, name, schedule_type, cron_expression, on_duration, off_duration, condition_json, action, action_value, action_duration, priority, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, req.ID, req.FarmID, req.CoopID, req.DeviceID, req.Name, req.ScheduleType, req.CronExpression, 
		req.OnDuration, req.OffDuration, req.ConditionJSON, req.Action, req.ActionValue, req.ActionDuration, 
		req.Priority, req.IsActive, req.CreatedAt)

	return &req, err
}

// UpdateSchedule updates an existing schedule
func (s *ScheduleService) UpdateSchedule(userID, farmID, scheduleID uuid.UUID, req models.Schedule) (*models.Schedule, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return nil, err
	}

	_, err := database.DB.Exec(`
		UPDATE schedules SET
			name = $1, schedule_type = $2, cron_expression = $3, on_duration = $4, 
			off_duration = $5, condition_json = $6, action = $7, action_value = $8, 
			action_duration = $9, priority = $10, is_active = $11, coop_id = $12, device_id = $13
		WHERE id = $14 AND farm_id = $15
	`, req.Name, req.ScheduleType, req.CronExpression, req.OnDuration, 
		req.OffDuration, req.ConditionJSON, req.Action, req.ActionValue, 
		req.ActionDuration, req.Priority, req.IsActive, req.CoopID, req.DeviceID, scheduleID, farmID)

	if err != nil {
		return nil, err
	}

	req.ID = scheduleID
	req.FarmID = farmID
	return &req, nil
}

// DeleteSchedule deletes a schedule
func (s *ScheduleService) DeleteSchedule(userID, farmID, scheduleID uuid.UUID) error {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return err
	}

	res, err := database.DB.Exec("DELETE FROM schedules WHERE id = $1 AND farm_id = $2", scheduleID, farmID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
