package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"middleware/models"
	"middleware/schemas"
	"middleware/services"
	"middleware/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== SCHEDULE HANDLERS =====

// ListSchedulesHandler returns all schedules for a farm or coop
// @Summary List Farm Schedules
// @Description Returns all schedules associated with a farm or specific coop
// @Tags Schedules
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param coop_id query string false "Coop ID (UUID) to filter"
// @Success 200 {object} []models.Schedule
// @Router /v1/farms/{farm_id}/schedules [get]
func ListSchedulesHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	var coopID *uuid.UUID
	if coopIDStr := c.Query("coop_id"); coopIDStr != "" {
		parsedID, err := uuid.Parse(coopIDStr)
		if err == nil {
			coopID = &parsedID
		}
	}

	schedules, err := scheduleService.ListSchedules(userID, farmID, coopID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		log.Printf("List schedules error: %v", err)
		return utils.InternalError(c, "Failed to fetch schedules")
	}

	resp := make([]fiber.Map, 0, len(schedules))
	for _, sc := range schedules {
		resp = append(resp, scheduleToResponse(sc))
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"schedules": resp,
	}, "Schedules retrieved")
}

// CreateScheduleHandler creates a new schedule
// @Summary Create Schedule
// @Description Creates a new automated schedule for a device
// @Tags Schedules
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param request body object true "Create Schedule Request"
// @Success 201 {object} models.Schedule
// @Router /v1/farms/{farm_id}/schedules [post]
func CreateScheduleHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	var payload struct {
		schemas.CreateScheduleRequest
		IsEnabled *bool `json:"is_enabled,omitempty"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	req := payload.CreateScheduleRequest
	if payload.IsEnabled != nil {
		req.IsActive = payload.IsEnabled
	}

	// Infer action if not explicitly provided
	req.Action = normalizeScheduleAction(&req.Action, req.ActionValue)

	created, err := scheduleService.CreateSchedule(userID, farmID, req)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == sql.ErrNoRows {
		return utils.BadRequest(c, "invalid_device", "Device not found for this farm")
	}
	if err != nil {
		log.Printf("Create schedule error: %v", err)
		return utils.InternalError(c, "Failed to create schedule")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, scheduleToResponse(*created), "Schedule created")
}

// GetScheduleHandler returns details for a specific schedule
// @Summary Get Schedule Details
// @Description Returns details for a specific schedule
// @Tags Schedules
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param schedule_id path string true "Schedule ID (UUID)"
// @Success 200 {object} models.Schedule
// @Router /v1/farms/{farm_id}/schedules/{schedule_id} [get]
func GetScheduleHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	sc, err := scheduleService.GetSchedule(userID, farmID, scheduleID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.NotFound(c, "Schedule not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, scheduleToResponse(*sc), "Schedule retrieved")
}

func UpdateScheduleHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	var payload struct {
		schemas.UpdateScheduleRequest
		IsEnabled *bool `json:"is_enabled,omitempty"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	req := payload.UpdateScheduleRequest
	if payload.IsEnabled != nil {
		req.IsActive = payload.IsEnabled
	}
	if req.Action == nil || *req.Action == "" {
		v := normalizeScheduleAction(nil, req.ActionValue)
		req.Action = &v
	}

	updated, err := scheduleService.UpdateSchedule(userID, farmID, scheduleID, req)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		log.Printf("Update schedule error: %v", err)
		return utils.InternalError(c, "Failed to update schedule")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, scheduleToResponse(*updated), "Schedule updated")
}

// DeleteScheduleHandler deletes a schedule
// @Summary Delete Schedule
// @Description Deletes a specific automated schedule
// @Tags Schedules
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param schedule_id path string true "Schedule ID (UUID)"
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/farms/{farm_id}/schedules/{schedule_id} [delete]
func DeleteScheduleHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	if err := scheduleService.DeleteSchedule(userID, farmID, scheduleID); err != nil {
		if err == services.ErrFarmAccessDenied {
			return utils.Forbidden(c, "Access denied")
		}
		return utils.InternalError(c, "Failed to delete schedule")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Schedule deleted")
}

func GetScheduleExecutionHistoryHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	items, err := scheduleService.GetExecutionHistory(userID, farmID, scheduleID, 50)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch execution history")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, items, "Execution history retrieved")
}

func ExecuteScheduleNowHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid schedule ID")
	}

	cmd, err := scheduleService.ExecuteScheduleNow(userID, farmID, scheduleID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to execute schedule")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cmd, "Schedule execution queued")
}

func normalizeScheduleAction(action *string, actionValue *string) string {
	if action != nil && *action != "" {
		return *action
	}
	if actionValue == nil {
		return "on"
	}
	v := strings.ToLower(strings.TrimSpace(*actionValue))
	if v == "on" || v == "off" {
		return v
	}
	return "set_value"
}

func scheduleToResponse(sc models.Schedule) fiber.Map {
	var actionSequence interface{} = nil
	if len(sc.ActionSequence) > 0 {
		actionSequence = json.RawMessage(sc.ActionSequence)
	}

	actionValue := sc.ActionValue
	if actionValue == nil && (sc.Action == "on" || sc.Action == "off") {
		v := sc.Action
		actionValue = &v
	}

	return fiber.Map{
		"id":              sc.ID,
		"farm_id":         sc.FarmID,
		"coop_id":         sc.CoopID,
		"device_id":       sc.DeviceID,
		"name":            sc.Name,
		"schedule_type":   sc.ScheduleType,
		"cron_expression": sc.CronExpression,
		"on_duration":     sc.OnDuration,
		"off_duration":    sc.OffDuration,
		"condition_json":  sc.ConditionJSON,
		"action":          sc.Action,
		"action_value":    actionValue,
		"action_duration": sc.ActionDuration,
		"action_sequence": actionSequence,
		"priority":        sc.Priority,
		"is_enabled":      sc.IsActive,
		"next_execution":  sc.NextExecution,
		"last_execution":  sc.LastExecution,
		"execution_count": sc.ExecutionCount,
		"created_by":      sc.CreatedBy,
		"created_at":      sc.CreatedAt,
		"updated_at":      sc.UpdatedAt,
	}
}
