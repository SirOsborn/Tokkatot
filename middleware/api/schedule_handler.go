package api

import (
	"log"
	"middleware/services"
	"middleware/utils"

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

	return utils.SuccessResponse(c, fiber.StatusOK, schedules, "Schedules retrieved")
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
	return utils.SuccessResponse(c, fiber.StatusCreated, nil, "Schedule created (mock)")
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
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Schedule details retrieved (mock)")
}

func UpdateScheduleHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Schedule updated (mock)")
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
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Schedule deleted (mock)")
}

func GetScheduleExecutionHistoryHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Execution history retrieved (mock)")
}

func ExecuteScheduleNowHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Schedule execution triggered (mock)")
}
