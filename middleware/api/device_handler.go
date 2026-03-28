package api

import (
	"log"
	"middleware/schemas"
	"middleware/services"
	"middleware/utils"
	"strconv"
	"strings"

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

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"devices": devices,
		"total":   total,
		"page":    offset/limit + 1,
		"limit":   limit,
	}, "Devices retrieved")
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
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	var req schemas.AddDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}
	if req.CoopID == nil {
		return utils.BadRequest(c, "invalid_coop", "coop_id is required")
	}

	device, err := deviceService.AddDevice(userID, farmID, req)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrCoopNotFound {
		return utils.BadRequest(c, "invalid_coop", "Coop not found for this farm")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to add device")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, device, "Device added")
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
	if err == services.ErrDeviceNotFound {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch device")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, device, "Device details retrieved")
}

func UpdateDeviceHandler(c *fiber.Ctx) error {
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

	var req schemas.UpdateDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}

	device, err := deviceService.UpdateDevice(userID, farmID, deviceID, req)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrDeviceNotFound {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to update device")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, device, "Device updated")
}

func DeleteDeviceHandler(c *fiber.Ctx) error {
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

	if err := deviceService.DeleteDevice(userID, farmID, deviceID); err != nil {
		if err == services.ErrFarmAccessDenied {
			return utils.Forbidden(c, "Access denied")
		}
		if err == services.ErrDeviceNotFound {
			return utils.NotFound(c, "Device not found")
		}
		return utils.InternalError(c, "Failed to delete device")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Device deleted")
}

func GetDeviceHistoryHandler(c *fiber.Ctx) error {
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

	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)
	sensorType := c.Query("sensor_type", "")

	history, total, err := deviceService.GetDeviceHistory(userID, farmID, deviceID, sensorType, limit, offset)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrDeviceNotFound {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch device history")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"readings": history,
		"total":    total,
	}, "Device history retrieved")
}

func GetDeviceStatusHandler(c *fiber.Ctx) error {
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

	status, err := deviceService.GetDeviceStatus(userID, farmID, deviceID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrDeviceNotFound {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch device status")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, status, "Device status retrieved")
}

func GetDeviceConfigHandler(c *fiber.Ctx) error {
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

	cfgs, err := deviceService.GetDeviceConfig(userID, farmID, deviceID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrDeviceNotFound {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch device config")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cfgs, "Device config retrieved")
}

func UpdateDeviceConfigHandler(c *fiber.Ctx) error {
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
		ParameterName  string   `json:"parameter_name"`
		ParameterValue string   `json:"parameter_value"`
		Unit           *string  `json:"unit,omitempty"`
		MinValue       *float64 `json:"min_value,omitempty"`
		MaxValue       *float64 `json:"max_value,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}
	if strings.TrimSpace(req.ParameterName) == "" {
		return utils.BadRequest(c, "invalid_param", "parameter_name is required")
	}

	cfg, err := deviceService.UpdateDeviceConfig(userID, farmID, deviceID, req.ParameterName, req.ParameterValue, req.Unit, req.MinValue, req.MaxValue)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrDeviceNotFound {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to update device config")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cfg, "Device config updated")
}

func CalibrateDeviceHandler(c *fiber.Ctx) error {
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
		ParameterName string `json:"parameter_name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}
	if strings.TrimSpace(req.ParameterName) == "" {
		return utils.BadRequest(c, "invalid_param", "parameter_name is required")
	}

	cfg, err := deviceService.CalibrateDevice(userID, farmID, deviceID, req.ParameterName)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrDeviceNotFound {
		return utils.NotFound(c, "Device not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to calibrate device")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cfg, "Device calibrated")
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
// @Router /v1/farms/{farm_id}/devices/{device_id}/commands [post]
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
		CommandType    string  `json:"command_type"`
		CommandValue   *string `json:"command_value,omitempty"`
		Parameters     *string `json:"parameters,omitempty"`
		ActionDuration *int    `json:"action_duration,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}

	value := req.CommandValue
	if value == nil && req.Parameters != nil {
		value = req.Parameters
	}

	cmd, err := deviceService.IssueCommand(userID, farmID, deviceID, req.CommandType, value, req.ActionDuration)
	if err != nil {
		return utils.InternalError(c, "Failed to issue command")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cmd, "Command issued")
}

func GetDeviceCommandStatusHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	commandID, err := uuid.Parse(c.Params("command_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid command ID")
	}

	cmd, err := deviceService.GetCommandStatus(userID, farmID, commandID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrCommandNotFound {
		return utils.NotFound(c, "Command not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch command")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, cmd, "Command status retrieved")
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
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	commandID, err := uuid.Parse(c.Params("command_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid command ID")
	}

	if err := deviceService.CancelCommand(userID, farmID, commandID); err != nil {
		if err == services.ErrFarmAccessDenied {
			return utils.Forbidden(c, "Access denied")
		}
		if err == services.ErrCommandNotFound {
			return utils.NotFound(c, "Command not found or already processed")
		}
		return utils.InternalError(c, "Failed to cancel command")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Command cancelled")
}

func GetFarmCommandHistoryHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	history, err := deviceService.GetFarmCommandHistory(userID, farmID, 50)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch farm command history")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, history, "Farm command history retrieved")
}

func EmergencyStopHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	count, err := deviceService.EmergencyStop(userID, farmID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to trigger emergency stop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{"commands_issued": count}, "Emergency stop triggered")
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
// @Router /v1/farms/{farm_id}/devices/batch-command [post]
func BatchDeviceCommandHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	var req schemas.BatchCommandRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}
	if len(req.DeviceIDs) == 0 || strings.TrimSpace(req.CommandType) == "" {
		return utils.BadRequest(c, "invalid_body", "device_ids and command_type are required")
	}

	results, err := deviceService.BatchCommands(userID, farmID, req)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to process batch command")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, results, "Batch command processed")
}

func UpdateDeviceHeartbeatHandler(c *fiber.Ctx) error {
	hardwareID := c.Params("hardware_id")
	if strings.TrimSpace(hardwareID) == "" {
		return utils.BadRequest(c, "invalid_id", "Invalid hardware ID")
	}

	var req struct {
		Status   *string `json:"status,omitempty"`
		Response *string `json:"response,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}

	status := ""
	if req.Status != nil {
		status = *req.Status
	}
	if err := deviceService.UpdateHeartbeat(hardwareID, status, req.Response); err != nil {
		if err == services.ErrDeviceNotFound {
			return utils.NotFound(c, "Device not found")
		}
		return utils.InternalError(c, "Failed to record heartbeat")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Heartbeat recorded")
}
// GetGatewayCommandsHandler returns pending commands for a specific hardware_id
func GetGatewayCommandsHandler(c *fiber.Ctx) error {
	hardwareID := c.Params("hardware_id")
	if strings.TrimSpace(hardwareID) == "" {
		return utils.BadRequest(c, "invalid_id", "Invalid hardware ID")
	}

	commands, err := deviceService.GetPendingCommands(hardwareID)
	if err != nil {
		log.Printf("Get gateway commands error: %v", err)
		return utils.InternalError(c, "Failed to fetch gateway commands")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, commands, "Gateway commands retrieved")
}

// UpdateGatewayCommandStatusHandler updates the status of a specific command
func UpdateGatewayCommandStatusHandler(c *fiber.Ctx) error {
	commandID, err := uuid.Parse(c.Params("command_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid command ID")
	}

	var req struct {
		Status   string `json:"status"`
		Response string `json:"response"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}

	if err := deviceService.UpdateCommandStatus(commandID, req.Status, req.Response); err != nil {
		if err == services.ErrCommandNotFound {
			return utils.NotFound(c, "Command not found")
		}
		return utils.InternalError(c, "Failed to update command status")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Command status updated")
}
