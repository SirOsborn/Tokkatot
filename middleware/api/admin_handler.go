package api

import (
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== ADMIN HANDLERS =====

func GetAdminStatsHandler(c *fiber.Ctx) error {
	stats, err := adminService.GetAdminStats()
	if err != nil {
		return utils.InternalError(c, "Failed to fetch admin stats")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, stats, "Admin stats retrieved")
}

func ListFarmersHandler(c *fiber.Ctx) error {
	farmers, err := adminService.ListAllFarmers()
	if err != nil {
		return utils.InternalError(c, "Failed to fetch farmers")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, farmers, "Farmers retrieved")
}

func RegisterFarmerHandler(c *fiber.Ctx) error {
	// Generates a one-time registration key for a new farmer
	var req struct {
		FarmName      string `json:"farm_name"`
		CustomerPhone string `json:"customer_phone"`
		NationalID    string `json:"national_id"`
		FullName      string `json:"full_name"`
		Sex           string `json:"sex"`
		Province      string `json:"province"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	key, err := adminService.GenerateRegistrationKey(req.FarmName, req.CustomerPhone, req.NationalID, req.FullName, req.Sex, req.Province, 90)
	if err != nil {
		return utils.InternalError(c, "Failed to generate registration key")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, key, "Registration key generated successfully")
}

func DeactivateFarmerHandler(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid user ID")
	}

	var req struct {
		Active bool `json:"active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	err = adminService.DeactivateUser(userID, req.Active)
	if err != nil {
		return utils.InternalError(c, "Failed to update user status")
	}

	status := "deactivated"
	if req.Active {
		status = "activated"
	}
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Farmer "+status+" successfully")
}

func GetFarmerProfileHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Farmer profile retrieved (mock)")
}

func ListViewersHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Viewers retrieved (mock)")
}

func ListRegKeysHandler(c *fiber.Ctx) error {
	keys, err := adminService.ListRegistrationKeys()
	if err != nil {
		return utils.InternalError(c, "Failed to fetch registration keys")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, keys, "Registration keys retrieved")
}

func UpdateAdminProfileHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Admin profile updated")
}

func ListGatewaysHandler(c *fiber.Ctx) error {
	gateways, err := adminService.ListGateways()
	if err != nil {
		return utils.InternalError(c, "Failed to fetch gateways")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, gateways, "Gateways retrieved")
}

func RevokeGatewayHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid gateway token ID")
	}

	err = adminService.RevokeGateway(id)
	if err != nil {
		return utils.InternalError(c, "Failed to revoke gateway")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Gateway access revoked")
}

func GetUnassignedGatewaysHandler(c *fiber.Ctx) error {
	gateways, err := deviceService.GetUnassignedGateways()
	if err != nil {
		return utils.InternalError(c, "Failed to fetch unassigned gateways")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, gateways, "Unassigned gateways retrieved")
}

func AssignGatewayHandler(c *fiber.Ctx) error {
	var req struct {
		HardwareID string    `json:"hardware_id"`
		FarmID     uuid.UUID `json:"farm_id"`
		CoopID     *uuid.UUID `json:"coop_id,omitempty"`
		Name       string    `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.HardwareID == "" || req.FarmID == uuid.Nil {
		return utils.BadRequest(c, "missing_fields", "Hardware ID and Farm ID are required")
	}

	err := deviceService.AssignGateway(req.HardwareID, req.FarmID, req.CoopID, req.Name)
	if err != nil {
		return utils.InternalError(c, "Failed to assign gateway: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Gateway assigned successfully")
}
