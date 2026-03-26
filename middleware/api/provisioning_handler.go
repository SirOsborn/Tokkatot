package api

import (
	"middleware/schemas"
	"middleware/services"
	"middleware/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var provisioningService = services.NewProvisioningService()

// RequestProvisioningHandler starts the Zero-Config setup flow for a Pi
func RequestProvisioningHandler(c *fiber.Ctx) error {
	var req schemas.ProvisionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Hardware ID required")
	}

	code, expiresAt, err := provisioningService.RequestProvisioning(req.HardwareID)
	if err != nil {
		return utils.InternalError(c, "Provisioning failed")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"setup_code": code,
		"expires_at": expiresAt.Format(time.RFC3339),
	}, "Setup code generated. Waiting for pairing...")
}

// CheckProvisioningStatusHandler polls to see if a Pi has been claimed
func CheckProvisioningStatusHandler(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return utils.BadRequest(c, "missing_code", "Setup code required")
	}

	isClaimed, farmID, coopID, token, err := provisioningService.CheckProvisioningStatus(code)
	if err != nil {
		return utils.BadRequest(c, "check_failed", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"is_claimed": isClaimed,
		"farm_id":    farmID,
		"coop_id":    coopID,
		"token":      token,
	}, "Claim status retrieved")
}

// ClaimGatewayHandler associates a Pi with a Farmer's Farm/Coop
func ClaimGatewayHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_farm_id", "Farm ID required")
	}

	var req schemas.ClaimGatewayRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Setup Code and Coop ID required")
	}

	err = provisioningService.ClaimGateway(userID, farmID, req.CoopID, req.SetupCode)
	if err != nil {
		return utils.BadRequest(c, "claim_failed", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Gateway claimed successfully!")
}
