package api

import (
	"log"
	"middleware/services"
	"middleware/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== DEVICE MANAGEMENT HANDLERS =====

// ListDevicesHandler returns all devices in a farm or coop
// @Summary List Farm Devices
// @Description Returns all devices associated with a farm or specific coop
// @Tags Devices
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param coop_id query string false "Coop ID (UUID) to filter"
// @Param filter query string false "Status filter (all, online, offline)"
// @Param limit query int false "Max items" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} schemas.PaginatedResponse
// @Router /v1/farms/{farm_id}/devices [get]
func ListDevicesHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	statusFilter := c.Query("filter", "all")

	var coopID *uuid.UUID
	if coopIDStr := c.Query("coop_id"); coopIDStr != "" {
		parsedID, err := uuid.Parse(coopIDStr)
		if err == nil {
			coopID = &parsedID
		}
	}

	devices, total, err := deviceService.ListDevices(userID, farmID, coopID, statusFilter, limit, offset)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		log.Printf("List devices error: %v", err)
		return utils.InternalError(c, "Failed to fetch devices")
	}

	return utils.SuccessListResponse(c, devices, total, offset/limit+1, limit)
}

// AddDeviceHandler adds a new device
// @Summary Add Device
// @Description Registers a new device to a farm
// @Tags Devices
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param request body object true "Add Device Request"
// @Success 201 {object} models.Device
// @Router /v1/farms/{farm_id}/devices [post]
func AddDeviceHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusCreated, nil, "Device added (mock)")
}

// GetDeviceHandler returns device details
// @Summary Get Device Details
// @Description Returns details for a specific device
// @Tags Devices
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param device_id path string true "Device ID (UUID)"
// @Success 200 {object} models.Device
// @Router /v1/farms/{farm_id}/devices/{device_id} [get]
func GetDeviceHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	device, err := deviceService.GetDevice(userID, farmID, deviceID)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch device")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, device, "Device details retrieved")
}

func UpdateDeviceHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device updated (mock)")
}

func DeleteDeviceHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device deleted (mock)")
}

func GetDeviceHistoryHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Device history retrieved (mock)")
}

func GetDeviceStatusHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device status retrieved (mock)")
}

func GetDeviceConfigHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device config retrieved (mock)")
}

func UpdateDeviceConfigHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device config updated (mock)")
}

func CalibrateDeviceHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device calibrated (mock)")
}

// SendDeviceCommandHandler issues a control command
// @Summary Send Command
// @Description Sends a control command (ON/OFF/Value) to a device
// @Tags Devices, Commands
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param device_id path string true "Device ID (UUID)"
// @Param request body object true "Command Request"
// @Success 200 {object} models.DeviceCommand
// @Router /v1/farms/{farm_id}/devices/{device_id}/command [post]
func SendDeviceCommandHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	var req struct {
		CommandType string `json:"command_type"`
		Parameters  string `json:"parameters"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}

	cmd, err := deviceService.IssueCommand(userID, farmID, deviceID, req.CommandType, &req.Parameters)
	if err != nil {
		return utils.InternalError(c, "Failed to issue command")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cmd, "Command issued")
}

func GetDeviceCommandStatusHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Command status retrieved (mock)")
}

// ListDeviceCommandsHandler returns command history for a device
// @Summary List Device Commands
// @Description Returns recent commands issued to a device
// @Tags Devices, Commands
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param device_id path string true "Device ID (UUID)"
// @Success 200 {object} []models.DeviceCommand
// @Router /v1/farms/{farm_id}/devices/{device_id}/commands [get]
func ListDeviceCommandsHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	deviceID, err := uuid.Parse(c.Params("device_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid device ID")
	}

	cmds, err := deviceService.GetDeviceCommands(userID, farmID, deviceID, 20)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch commands")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cmds, "Commands retrieved")
}

func CancelCommandHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Command cancelled (mock)")
}

func GetFarmCommandHistoryHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Farm command history retrieved (mock)")
}

func EmergencyStopHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Emergency stop triggered (mock)")
}

// BatchDeviceCommandHandler issues commands to multiple devices
// @Summary Batch Commands
// @Description Sends a command to multiple devices in a farm/coop
// @Tags Devices, Commands
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param request body schemas.BatchCommandRequest true "Batch Command Request"
// @Success 200 {object} schemas.BatchResult
// @Router /v1/farms/{farm_id}/batch/command [post]
func BatchDeviceCommandHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Batch command processed (mock)")
}

func UpdateDeviceHeartbeatHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Heartbeat recorded (mock)")
}
