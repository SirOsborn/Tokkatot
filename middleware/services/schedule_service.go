package services

import (
	"database/sql"
	"encoding/json"
	"middleware/database"
	"middleware/models"
	"middleware/schemas"
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

	query := `SELECT id, farm_id, coop_id, device_id, name, schedule_type, cron_expression, on_duration, off_duration, condition_json, action, action_value, action_duration, action_sequence, priority, is_active, next_execution, last_execution, execution_count, created_by, created_at, updated_at FROM schedules WHERE farm_id = $1`
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
			&sc.OnDuration, &sc.OffDuration, &sc.ConditionJSON, &sc.Action, &sc.ActionValue, &sc.ActionDuration, &sc.ActionSequence,
			&sc.Priority, &sc.IsActive, &sc.NextExecution, &sc.LastExecution, &sc.ExecutionCount, &sc.CreatedBy, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
			continue
		}
		schedules = append(schedules, sc)
	}

	return schedules, nil
}

// CreateSchedule creates a new schedule
func (s *ScheduleService) CreateSchedule(userID, farmID uuid.UUID, req schemas.CreateScheduleRequest) (*models.Schedule, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}

	if ok, err := scheduleDeviceBelongsToFarm(req.DeviceID, farmID); err != nil {
		return nil, err
	} else if !ok {
		return nil, sql.ErrNoRows
	}

	now := time.Now()
	priority := 0
	if req.Priority != nil {
		priority = *req.Priority
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var actionSequence interface{} = nil
	if len(req.ActionSequence) > 0 {
		actionSequence = models.NullRawMessage(req.ActionSequence)
	}

	schedule := models.Schedule{
		ID:             uuid.New(),
		FarmID:         farmID,
		CoopID:         req.CoopID,
		DeviceID:       req.DeviceID,
		Name:           req.Name,
		ScheduleType:   req.ScheduleType,
		CronExpression: req.CronExpression,
		OnDuration:     req.OnDuration,
		OffDuration:    req.OffDuration,
		ConditionJSON:  req.ConditionJSON,
		Action:         req.Action,
		ActionValue:    req.ActionValue,
		ActionDuration: req.ActionDuration,
		ActionSequence: models.NullRawMessage(nil),
		Priority:       priority,
		IsActive:       isActive,
		CreatedBy:      userID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if len(req.ActionSequence) > 0 {
		schedule.ActionSequence = models.NullRawMessage(req.ActionSequence)
	}

	_, err := database.DB.Exec(`
		INSERT INTO schedules (id, farm_id, coop_id, device_id, name, schedule_type, cron_expression, on_duration, off_duration, condition_json, action, action_value, action_duration, action_sequence, priority, is_active, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`, schedule.ID, schedule.FarmID, schedule.CoopID, schedule.DeviceID, schedule.Name, schedule.ScheduleType,
		schedule.CronExpression, schedule.OnDuration, schedule.OffDuration, schedule.ConditionJSON, schedule.Action,
		schedule.ActionValue, schedule.ActionDuration, actionSequence, schedule.Priority, schedule.IsActive,
		schedule.CreatedBy, schedule.CreatedAt, schedule.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

// UpdateSchedule updates an existing schedule
func (s *ScheduleService) UpdateSchedule(userID, farmID, scheduleID uuid.UUID, req schemas.UpdateScheduleRequest) (*models.Schedule, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}

	var actionSequence interface{} = nil
	if len(req.ActionSequence) > 0 {
		actionSequence = models.NullRawMessage(req.ActionSequence)
	}

	_, err := database.DB.Exec(`
		UPDATE schedules SET
			name = COALESCE($1, name),
			schedule_type = COALESCE($2, schedule_type),
			cron_expression = COALESCE($3, cron_expression),
			on_duration = COALESCE($4, on_duration),
			off_duration = COALESCE($5, off_duration),
			condition_json = COALESCE($6, condition_json),
			action = COALESCE($7, action),
			action_value = COALESCE($8, action_value),
			action_duration = COALESCE($9, action_duration),
			action_sequence = COALESCE($10, action_sequence),
			priority = COALESCE($11, priority),
			is_active = COALESCE($12, is_active),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $13 AND farm_id = $14
	`, req.Name, req.ScheduleType, req.CronExpression, req.OnDuration,
		req.OffDuration, req.ConditionJSON, req.Action, req.ActionValue,
		req.ActionDuration, actionSequence, req.Priority, req.IsActive, scheduleID, farmID)
	if err != nil {
		return nil, err
	}

	return s.GetSchedule(userID, farmID, scheduleID)
}

// DeleteSchedule deletes a schedule
func (s *ScheduleService) DeleteSchedule(userID, farmID, scheduleID uuid.UUID) error {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return err
	}

	res, err := database.DB.Exec("UPDATE schedules SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND farm_id = $2", scheduleID, farmID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetSchedule retrieves a schedule by ID
func (s *ScheduleService) GetSchedule(userID, farmID, scheduleID uuid.UUID) (*models.Schedule, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	var sc models.Schedule
	err := database.DB.QueryRow(`
		SELECT id, farm_id, coop_id, device_id, name, schedule_type, cron_expression, on_duration, off_duration,
		       condition_json, action, action_value, action_duration, action_sequence, priority, is_active,
		       next_execution, last_execution, execution_count, created_by, created_at, updated_at
		FROM schedules
		WHERE id = $1 AND farm_id = $2
	`, scheduleID, farmID).Scan(
		&sc.ID, &sc.FarmID, &sc.CoopID, &sc.DeviceID, &sc.Name, &sc.ScheduleType, &sc.CronExpression,
		&sc.OnDuration, &sc.OffDuration, &sc.ConditionJSON, &sc.Action, &sc.ActionValue, &sc.ActionDuration,
		&sc.ActionSequence, &sc.Priority, &sc.IsActive, &sc.NextExecution, &sc.LastExecution,
		&sc.ExecutionCount, &sc.CreatedBy, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

// GetExecutionHistory returns recent executions for a schedule
func (s *ScheduleService) GetExecutionHistory(userID, farmID, scheduleID uuid.UUID, limit int) ([]models.ScheduleExecution, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(`
		SELECT id, schedule_id, device_id, scheduled_time, actual_execution_time, status, execution_duration_ms, device_response, error_message, created_at
		FROM schedule_executions
		WHERE schedule_id = $1
		ORDER BY scheduled_time DESC
		LIMIT $2
	`, scheduleID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ScheduleExecution
	for rows.Next() {
		var it models.ScheduleExecution
		if err := rows.Scan(&it.ID, &it.ScheduleID, &it.DeviceID, &it.ScheduledTime, &it.ActualExecutionTime, &it.Status, &it.ExecutionDurationMs, &it.DeviceResponse, &it.ErrorMessage, &it.CreatedAt); err != nil {
			continue
		}
		items = append(items, it)
	}
	return items, nil
}

// ExecuteScheduleNow creates an immediate execution entry and device command
func (s *ScheduleService) ExecuteScheduleNow(userID, farmID, scheduleID uuid.UUID) (*models.DeviceCommand, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}

	sc, err := s.GetSchedule(userID, farmID, scheduleID)
	if err != nil {
		return nil, err
	}

	commandType := sc.Action
	commandValue := sc.ActionValue
	cmd, err := NewDeviceService().IssueCommand(userID, farmID, sc.DeviceID, commandType, commandValue)
	if err != nil {
		return nil, err
	}

	respJSON, _ := json.Marshal(map[string]interface{}{
		"command_id": cmd.ID.String(),
		"status":     cmd.Status,
	})

	_, _ = database.DB.Exec(`
		INSERT INTO schedule_executions (id, schedule_id, device_id, scheduled_time, actual_execution_time, status, device_response, created_at)
		VALUES ($1, $2, $3, $4, $5, 'executed', $6, $7)
	`, uuid.New(), scheduleID, sc.DeviceID, time.Now(), time.Now(), respJSON, time.Now())

	return cmd, nil
}

func scheduleDeviceBelongsToFarm(deviceID, farmID uuid.UUID) (bool, error) {
	var exists bool
	if err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM devices WHERE id = $1 AND farm_id = $2)", deviceID, farmID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
