package api

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"middleware/database"
	"middleware/models"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

// ===== SCHEDULE MANAGEMENT ENDPOINTS =====

// CreateScheduleHandler creates a new automated schedule
// POST /v1/farms/:farm_id/schedules
func CreateScheduleHandler(c *fiber.Ctx) error {
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check farm access (manager+ required)
	userID := c.Locals("user_id").(uuid.UUID)
	err = checkFarmAccess(userID, farmID, "manager")
	if err != nil {
		return err
	}

	// Parse request body
	var req struct {
		DeviceID       uuid.UUID       `json:"device_id"`
		CoopID         *uuid.UUID      `json:"coop_id,omitempty"`
		Name           string          `json:"name"`
		ScheduleType   string          `json:"schedule_type"` // "time_based", "duration_based", "condition_based"
		CronExpression *string         `json:"cron_expression,omitempty"`
		OnDuration     *int            `json:"on_duration,omitempty"`     // For duration_based: seconds ON
		OffDuration    *int            `json:"off_duration,omitempty"`    // For duration_based: seconds OFF
		ConditionJSON  *string         `json:"condition_json,omitempty"`  // For condition_based: sensor rules
		Action         string          `json:"action"`                    // "on", "off", "set_value"
		ActionValue    *string         `json:"action_value,omitempty"`    // PWM value, etc.
		ActionDuration *int            `json:"action_duration,omitempty"` // For time_based: auto-off after X seconds
		ActionSequence json.RawMessage `json:"action_sequence,omitempty"` // For time_based: multi-step [{"action":"ON","duration":30}]
		Priority       *int            `json:"priority,omitempty"`
		IsActive       *bool           `json:"is_active,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	// Validate required fields
	if req.Name == "" {
		return utils.BadRequest(c, "missing_name", "Schedule name is required")
	}
	if req.ScheduleType == "" {
		return utils.BadRequest(c, "missing_type", "Schedule type is required (time_based, duration_based, condition_based)")
	}
	if req.Action == "" {
		return utils.BadRequest(c, "missing_action", "Action is required (on, off, set_value)")
	}

	// Validate schedule type and required fields
	switch req.ScheduleType {
	case "time_based":
		if req.CronExpression == nil || *req.CronExpression == "" {
			return utils.BadRequest(c, "missing_cron", "Cron expression is required for time_based schedules")
		}
		// Validate cron expression
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(*req.CronExpression)
		if err != nil {
			return utils.BadRequest(c, "invalid_cron", "Invalid cron expression: "+err.Error())
		}
		// Calculate next execution time
		nextExecution := schedule.Next(time.Now())
		_ = nextExecution // We'll use this later

	case "duration_based":
		if req.OnDuration == nil || req.OffDuration == nil {
			return utils.BadRequest(c, "missing_duration", "on_duration and off_duration are required for duration_based schedules")
		}
		if *req.OnDuration <= 0 || *req.OffDuration <= 0 {
			return utils.BadRequest(c, "invalid_duration", "on_duration and off_duration must be positive")
		}

	case "condition_based":
		if req.ConditionJSON == nil || *req.ConditionJSON == "" {
			return utils.BadRequest(c, "missing_condition", "condition_json is required for condition_based schedules")
		}
		// TODO: Validate JSON structure

	default:
		return utils.BadRequest(c, "invalid_type", "Invalid schedule_type. Must be: time_based, duration_based, or condition_based")
	}

	// Verify device exists and belongs to this farm
	var deviceFarmID uuid.UUID
	err = database.DB.QueryRow(`
		SELECT farm_id FROM devices WHERE id = $1 AND is_active = true
	`, req.DeviceID).Scan(&deviceFarmID)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error: "+err.Error())
	}
	if deviceFarmID != farmID {
		return fiber.NewError(fiber.StatusForbidden, "Device does not belong to this farm")
	}

	// Set defaults
	priority := 0
	if req.Priority != nil {
		priority = *req.Priority
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Calculate next execution time for time_based schedules
	var nextExecution *time.Time
	if req.ScheduleType == "time_based" && req.CronExpression != nil {
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, _ := parser.Parse(*req.CronExpression)
		next := schedule.Next(time.Now())
		nextExecution = &next
	}

	// Create schedule
	scheduleID := uuid.New()
	now := time.Now()

	_, err = database.DB.Exec(`
		INSERT INTO schedules (
			id, farm_id, coop_id, device_id, name, schedule_type,
			cron_expression, on_duration, off_duration, condition_json,
			action, action_value, action_duration, action_sequence, priority, is_active,
			next_execution, execution_count, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`, scheduleID, farmID, req.CoopID, req.DeviceID, req.Name, req.ScheduleType,
		req.CronExpression, req.OnDuration, req.OffDuration, req.ConditionJSON,
		req.Action, req.ActionValue, req.ActionDuration, req.ActionSequence, priority, isActive,
		nextExecution, 0, userID, now, now)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return utils.Conflict(c, "duplicate_schedule", "A schedule with this name already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create schedule: "+err.Error())
	}

	// Fetch created schedule
	schedule := models.Schedule{}
	err = database.DB.QueryRow(`
		SELECT id, farm_id, coop_id, device_id, name, schedule_type,
			   cron_expression, on_duration, off_duration, condition_json,
			   action, action_value, action_duration, action_sequence, priority, is_active,
			   next_execution, last_execution, execution_count,
			   created_by, created_at, updated_at
		FROM schedules WHERE id = $1
	`, scheduleID).Scan(
		&schedule.ID, &schedule.FarmID, &schedule.CoopID, &schedule.DeviceID,
		&schedule.Name, &schedule.ScheduleType, &schedule.CronExpression,
		&schedule.OnDuration, &schedule.OffDuration, &schedule.ConditionJSON,
		&schedule.Action, &schedule.ActionValue, &schedule.ActionDuration, &schedule.ActionSequence, &schedule.Priority, &schedule.IsActive,
		&schedule.NextExecution, &schedule.LastExecution, &schedule.ExecutionCount,
		&schedule.CreatedBy, &schedule.CreatedAt, &schedule.UpdatedAt,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch created schedule")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, schedule, "Schedule created successfully")
}

// ListSchedulesHandler lists all schedules for a farm
// GET /v1/farms/:farm_id/schedules
func ListSchedulesHandler(c *fiber.Ctx) error {
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check farm access
	userID := c.Locals("user_id").(uuid.UUID)
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	// Parse query parameters
	deviceIDParam := c.Query("device_id")
	limit := c.QueryInt("limit", 50)
	if limit > 200 {
		limit = 200
	}

	// Build query
	query := `
		SELECT s.id, s.farm_id, s.coop_id, s.device_id, s.name, s.schedule_type,
			   s.cron_expression, s.on_duration, s.off_duration, s.condition_json,
			   s.action, s.action_value, s.action_duration, s.action_sequence, s.priority, s.is_active,
			   s.next_execution, s.last_execution, s.execution_count,
			   s.created_by, s.created_at, s.updated_at,
			   d.name as device_name
		FROM schedules s
		LEFT JOIN devices d ON s.device_id = d.id
		WHERE s.farm_id = $1 AND s.is_active = true
	`
	args := []interface{}{farmID}

	// Filter by device if provided
	if deviceIDParam != "" {
		deviceID, err := uuid.Parse(deviceIDParam)
		if err != nil {
			return utils.BadRequest(c, "invalid_id", "Invalid device_id")
		}
		query += ` AND s.device_id = $2`
		args = append(args, deviceID)
	}

	query += ` ORDER BY s.priority DESC, s.created_at DESC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit)

	// Execute query
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch schedules: "+err.Error())
	}
	defer rows.Close()

	type ScheduleWithDevice struct {
		models.Schedule
		DeviceName *string `json:"device_name,omitempty"`
	}

	schedules := []ScheduleWithDevice{}
	for rows.Next() {
		var s ScheduleWithDevice
		err := rows.Scan(
			&s.ID, &s.FarmID, &s.CoopID, &s.DeviceID, &s.Name, &s.ScheduleType,
			&s.CronExpression, &s.OnDuration, &s.OffDuration, &s.ConditionJSON,
			&s.Action, &s.ActionValue, &s.ActionDuration, &s.ActionSequence, &s.Priority, &s.IsActive,
			&s.NextExecution, &s.LastExecution, &s.ExecutionCount,
			&s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
			&s.DeviceName,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan schedule: "+err.Error())
		}
		schedules = append(schedules, s)
	}

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM schedules WHERE farm_id = $1 AND is_active = true`
	countArgs := []interface{}{farmID}
	if deviceIDParam != "" {
		deviceID, _ := uuid.Parse(deviceIDParam)
		countQuery += ` AND device_id = $2`
		countArgs = append(countArgs, deviceID)
	}
	database.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	return utils.SuccessListResponse(c, schedules, total, 1, limit)
}

// GetScheduleHandler gets details of a specific schedule
// GET /v1/farms/:farm_id/schedules/:schedule_id
func GetScheduleHandler(c *fiber.Ctx) error {
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	// Check farm access
	userID := c.Locals("user_id").(uuid.UUID)
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	// Fetch schedule with device name
	type ScheduleDetails struct {
		models.Schedule
		DeviceName *string `json:"device_name,omitempty"`
	}

	var s ScheduleDetails
	err = database.DB.QueryRow(`
		SELECT s.id, s.farm_id, s.coop_id, s.device_id, s.name, s.schedule_type,
			   s.cron_expression, s.on_duration, s.off_duration, s.condition_json,
			   s.action, s.action_value, s.action_duration, s.action_sequence, s.priority, s.is_active,
			   s.next_execution, s.last_execution, s.execution_count,
			   s.created_by, s.created_at, s.updated_at,
			   d.name as device_name
		FROM schedules s
		LEFT JOIN devices d ON s.device_id = d.id
		WHERE s.id = $1 AND s.farm_id = $2
	`, scheduleID, farmID).Scan(
		&s.ID, &s.FarmID, &s.CoopID, &s.DeviceID, &s.Name, &s.ScheduleType,
		&s.CronExpression, &s.OnDuration, &s.OffDuration, &s.ConditionJSON,
		&s.Action, &s.ActionValue, &s.ActionDuration, &s.ActionSequence, &s.Priority, &s.IsActive,
		&s.NextExecution, &s.LastExecution, &s.ExecutionCount,
		&s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
		&s.DeviceName,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Schedule not found")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch schedule: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, s, "Schedule fetched successfully")
}

// UpdateScheduleHandler updates a schedule
// PUT /v1/farms/:farm_id/schedules/:schedule_id
func UpdateScheduleHandler(c *fiber.Ctx) error {
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	// Check farm access (manager+ required)
	userID := c.Locals("user_id").(uuid.UUID)
	err = checkFarmAccess(userID, farmID, "manager")
	if err != nil {
		return err
	}

	// Parse request
	var req struct {
		Name           *string         `json:"name,omitempty"`
		ScheduleType   *string         `json:"schedule_type,omitempty"`
		CronExpression *string         `json:"cron_expression,omitempty"`
		OnDuration     *int            `json:"on_duration,omitempty"`
		OffDuration    *int            `json:"off_duration,omitempty"`
		ConditionJSON  *string         `json:"condition_json,omitempty"`
		Action         *string         `json:"action,omitempty"`
		ActionValue    *string         `json:"action_value,omitempty"`
		ActionDuration *int            `json:"action_duration,omitempty"`
		ActionSequence json.RawMessage `json:"action_sequence,omitempty"`
		Priority       *int            `json:"priority,omitempty"`
		IsActive       *bool           `json:"is_active,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	// Verify schedule exists and belongs to farm
	var existingScheduleType string
	err = database.DB.QueryRow(`
		SELECT schedule_type FROM schedules WHERE id = $1 AND farm_id = $2
	`, scheduleID, farmID).Scan(&existingScheduleType)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Schedule not found")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error: "+err.Error())
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		updates = append(updates, "name = $"+strconv.Itoa(argCount))
		args = append(args, *req.Name)
		argCount++
	}
	if req.ScheduleType != nil {
		updates = append(updates, "schedule_type = $"+strconv.Itoa(argCount))
		args = append(args, *req.ScheduleType)
		argCount++
	}
	if req.CronExpression != nil {
		// Validate cron expression
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(*req.CronExpression)
		if err != nil {
			return utils.BadRequest(c, "invalid_cron", "Invalid cron expression: "+err.Error())
		}
		// Update next execution time
		nextExecution := schedule.Next(time.Now())
		updates = append(updates, "cron_expression = $"+strconv.Itoa(argCount))
		args = append(args, *req.CronExpression)
		argCount++
		updates = append(updates, "next_execution = $"+strconv.Itoa(argCount))
		args = append(args, nextExecution)
		argCount++
	}
	if req.OnDuration != nil {
		updates = append(updates, "on_duration = $"+strconv.Itoa(argCount))
		args = append(args, *req.OnDuration)
		argCount++
	}
	if req.OffDuration != nil {
		updates = append(updates, "off_duration = $"+strconv.Itoa(argCount))
		args = append(args, *req.OffDuration)
		argCount++
	}
	if req.ConditionJSON != nil {
		updates = append(updates, "condition_json = $"+strconv.Itoa(argCount))
		args = append(args, *req.ConditionJSON)
		argCount++
	}
	if req.Action != nil {
		updates = append(updates, "action = $"+strconv.Itoa(argCount))
		args = append(args, *req.Action)
		argCount++
	}
	if req.ActionValue != nil {
		updates = append(updates, "action_value = $"+strconv.Itoa(argCount))
		args = append(args, *req.ActionValue)
		argCount++
	}
	if req.ActionDuration != nil {
		updates = append(updates, "action_duration = $"+strconv.Itoa(argCount))
		args = append(args, *req.ActionDuration)
		argCount++
	}
	if len(req.ActionSequence) > 0 {
		updates = append(updates, "action_sequence = $"+strconv.Itoa(argCount))
		args = append(args, req.ActionSequence)
		argCount++
	}
	if req.Priority != nil {
		updates = append(updates, "priority = $"+strconv.Itoa(argCount))
		args = append(args, *req.Priority)
		argCount++
	}
	if req.IsActive != nil {
		updates = append(updates, "is_active = $"+strconv.Itoa(argCount))
		args = append(args, *req.IsActive)
		argCount++
	}

	if len(updates) == 0 {
		return utils.BadRequest(c, "no_updates", "No fields to update")
	}

	// Add updated_at
	updates = append(updates, "updated_at = $"+strconv.Itoa(argCount))
	args = append(args, time.Now())
	argCount++

	// Add WHERE clause parameters
	args = append(args, scheduleID, farmID)

	query := `UPDATE schedules SET ` + strings.Join(updates, ", ") +
		` WHERE id = $` + strconv.Itoa(argCount) + ` AND farm_id = $` + strconv.Itoa(argCount+1)

	_, err = database.DB.Exec(query, args...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update schedule: "+err.Error())
	}

	// Fetch updated schedule
	var schedule models.Schedule
	err = database.DB.QueryRow(`
		SELECT id, farm_id, coop_id, device_id, name, schedule_type,
			   cron_expression, on_duration, off_duration, condition_json,
			   action, action_value, action_duration, action_sequence, priority, is_active,
			   next_execution, last_execution, execution_count,
			   created_by, created_at, updated_at
		FROM schedules WHERE id = $1
	`, scheduleID).Scan(
		&schedule.ID, &schedule.FarmID, &schedule.CoopID, &schedule.DeviceID,
		&schedule.Name, &schedule.ScheduleType, &schedule.CronExpression,
		&schedule.OnDuration, &schedule.OffDuration, &schedule.ConditionJSON,
		&schedule.Action, &schedule.ActionValue, &schedule.ActionDuration, &schedule.ActionSequence, &schedule.Priority, &schedule.IsActive,
		&schedule.NextExecution, &schedule.LastExecution, &schedule.ExecutionCount,
		&schedule.CreatedBy, &schedule.CreatedAt, &schedule.UpdatedAt,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch updated schedule")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, schedule, "Schedule updated successfully")
}

// DeleteScheduleHandler soft deletes a schedule
// DELETE /v1/farms/:farm_id/schedules/:schedule_id
func DeleteScheduleHandler(c *fiber.Ctx) error {
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	// Check farm access (manager+ required)
	userID := c.Locals("user_id").(uuid.UUID)
	err = checkFarmAccess(userID, farmID, "manager")
	if err != nil {
		return err
	}

	// Verify schedule exists
	var exists bool
	err = database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM schedules WHERE id = $1 AND farm_id = $2)
	`, scheduleID, farmID).Scan(&exists)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error: "+err.Error())
	}
	if !exists {
		return utils.NotFound(c, "Schedule not found")
	}

	// Soft delete (set is_active = false)
	_, err = database.DB.Exec(`
		UPDATE schedules SET is_active = false, updated_at = $1
		WHERE id = $2 AND farm_id = $3
	`, time.Now(), scheduleID, farmID)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete schedule: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Schedule deleted successfully",
	})
}

// GetScheduleExecutionHistoryHandler gets execution history for a schedule
// GET /v1/farms/:farm_id/schedules/:schedule_id/executions
func GetScheduleExecutionHistoryHandler(c *fiber.Ctx) error {
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	// Check farm access
	userID := c.Locals("user_id").(uuid.UUID)
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	// Verify schedule exists
	var exists bool
	err = database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM schedules WHERE id = $1 AND farm_id = $2)
	`, scheduleID, farmID).Scan(&exists)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error: "+err.Error())
	}
	if !exists {
		return utils.NotFound(c, "Schedule not found")
	}

	// Parse query parameters
	limit := c.QueryInt("limit", 100)
	if limit > 500 {
		limit = 500
	}
	days := c.QueryInt("days", 30)

	// Fetch executions
	rows, err := database.DB.Query(`
		SELECT id, schedule_id, device_id, scheduled_time, actual_execution_time,
			   status, execution_duration_ms, device_response, error_message, created_at
		FROM schedule_executions
		WHERE schedule_id = $1
		  AND scheduled_time >= $2
		ORDER BY scheduled_time DESC
		LIMIT $3
	`, scheduleID, time.Now().AddDate(0, 0, -days), limit)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch executions: "+err.Error())
	}
	defer rows.Close()

	executions := []models.ScheduleExecution{}
	successCount := 0
	for rows.Next() {
		var exec models.ScheduleExecution
		err := rows.Scan(
			&exec.ID, &exec.ScheduleID, &exec.DeviceID, &exec.ScheduledTime,
			&exec.ActualExecutionTime, &exec.Status, &exec.ExecutionDurationMs,
			&exec.DeviceResponse, &exec.ErrorMessage, &exec.CreatedAt,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan execution: "+err.Error())
		}
		executions = append(executions, exec)
		if exec.Status == "executed" {
			successCount++
		}
	}

	// Calculate success rate
	successRate := 0.0
	if len(executions) > 0 {
		successRate = (float64(successCount) / float64(len(executions))) * 100
	}

	return c.JSON(fiber.Map{
		"executions":   executions,
		"total":        len(executions),
		"success_rate": successRate,
	})
}

// ExecuteScheduleNowHandler manually executes a schedule immediately
// POST /v1/farms/:farm_id/schedules/:schedule_id/execute-now
func ExecuteScheduleNowHandler(c *fiber.Ctx) error {
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	// Check farm access (manager+ required)
	userID := c.Locals("user_id").(uuid.UUID)
	err = checkFarmAccess(userID, farmID, "manager")
	if err != nil {
		return err
	}

	// Fetch schedule details	// Fetch schedule details
	var schedule models.Schedule
	err = database.DB.QueryRow(`
		SELECT id, device_id, action, action_value, is_active
		FROM schedules
		WHERE id = $1 AND farm_id = $2
	`, scheduleID, farmID).Scan(
		&schedule.ID, &schedule.DeviceID, &schedule.Action,
		&schedule.ActionValue, &schedule.IsActive,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Schedule not found")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch schedule: "+err.Error())
	}

	// Create a device command
	commandID := uuid.New()
	now := time.Now()

	_, err = database.DB.Exec(`
		INSERT INTO device_commands (
			id, device_id, farm_id, issued_by, command_type,
			command_value, status, issued_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, commandID, schedule.DeviceID, farmID, userID, schedule.Action,
		schedule.ActionValue, "pending", now, now)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create command: "+err.Error())
	}

	// Create execution log
	executionID := uuid.New()
	_, err = database.DB.Exec(`
		INSERT INTO schedule_executions (
			id, schedule_id, device_id, scheduled_time, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`, executionID, scheduleID, schedule.DeviceID, now, "executed", now)

	if err != nil {
		// Non-critical error, log but don't fail
		// TODO: Add proper logging
	}

	// TODO: Publish MQTT command to device

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"execution_id": executionID,
		"command_id":   commandID,
		"schedule_id":  scheduleID,
		"status":       "queued",
		"message":      "Schedule queued for immediate execution",
	})
}
