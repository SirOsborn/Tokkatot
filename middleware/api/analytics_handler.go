package api

import (
	"log"
	"middleware/services"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== ANALYTICS & REPORTING HANDLERS =====

// GetFarmDashboardHandler returns overview stats for a farm
// @Summary Get Farm Dashboard
// @Description Returns an overview of farm status, sensor alerts, and device health
// @Tags Analytics
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Success 200 {object} schemas.DashboardResponse
// @Router /v1/farms/{farm_id}/dashboard [get]
func GetFarmDashboardHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	resp, err := analyticsService.GetFarmDashboard(userID, farmID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		log.Printf("Dashboard error: %v", err)
		return utils.InternalError(c, "Failed to fetch dashboard data")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, resp, "Dashboard data retrieved")
}

// GetCoopAnalyticsHandler returns detailed analytics for a coop
// @Summary Get Coop Analytics
// @Description Returns historical performance data for a coop
// @Tags Analytics
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param coop_id path string true "Coop ID (UUID)"
// @Success 200 {object} object
// @Router /v1/farms/{farm_id}/coops/{coop_id}/analytics [get]
func GetCoopAnalyticsHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Coop analytics retrieved (mock)")
}

func GetFarmSensorTrendsHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Sensor trends retrieved (mock)")
}

func GetDeviceHealthReportHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device health report retrieved (mock)")
}

func GetDeviceMetricsReportHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device metrics report retrieved (mock)")
}

func GetDeviceUsageReportHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device usage report retrieved (mock)")
}

func GetFarmPerformanceReportHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Performance report retrieved (mock)")
}

func ExportReportHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Report exported (mock)")
}

// GetFarmEventLogHandler returns recent system events
// @Summary Get Event logs
// @Description Returns recent activity and system events for the farm
// @Tags Analytics, Audit
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Success 200 {object} []models.EventLog
// @Router /v1/farms/{farm_id}/events [get]
func GetFarmEventLogHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Event logs retrieved (mock)")
}
