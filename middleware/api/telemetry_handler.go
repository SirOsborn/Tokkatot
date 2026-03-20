package api

import (
	"middleware/schemas"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ReportCoopDevicesHandler upserts devices and marks missing ones inactive
func ReportCoopDevicesHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	coopID, err := uuid.Parse(c.Params("coop_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	var req schemas.DeviceReportRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if req.HardwareID == "" || len(req.Devices) == 0 {
		return utils.BadRequest(c, "invalid_request", "hardware_id and devices are required")
	}

	devices, err := telemetryService.ReportCoopDevices(userID, farmID, coopID, req)
	if err != nil {
		return utils.InternalError(c, "Failed to report devices")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"devices": devices,
	}, "Devices reported")
}

// PostCoopTelemetryHandler ingests sensor telemetry from gateway
func PostCoopTelemetryHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	coopID, err := uuid.Parse(c.Params("coop_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	var req schemas.TelemetryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if err := telemetryService.IngestTelemetry(userID, farmID, coopID, req); err != nil {
		return utils.InternalError(c, "Failed to ingest telemetry")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Telemetry ingested")
}
