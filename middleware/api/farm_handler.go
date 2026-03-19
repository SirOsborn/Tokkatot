package api

import (
	"log"
	"strconv"

	"middleware/schemas"
	"middleware/services"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== FARM MANAGEMENT HANDLERS =====

// ListFarmsHandler returns all farms the current user has access to
// @Summary List User Farms
// @Description Returns all farms associated with the authenticated user
// @Tags Farms
// @Produce json
// @Param limit query int false "Max items" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} schemas.PaginatedResponse
// @Router /v1/farms [get]
func ListFarmsHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	farms, total, err := farmService.ListFarms(userID, limit, offset)
	if err != nil {
		log.Printf("List farms error: %v", err)
		return utils.InternalError(c, "Failed to fetch farms")
	}

	return utils.SuccessListResponse(c, farms, total, offset/limit+1, limit)
}

// GetFarmHandler returns a single farm by ID
// @Summary Get Farm Details
// @Description Returns details for a specific farm
// @Tags Farms
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Success 200 {object} schemas.FarmWithRole
// @Router /v1/farms/{farm_id} [get]
func GetFarmHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	farm, err := farmService.GetFarm(userID, farmID)
	if err == services.ErrFarmNotFound {
		return utils.NotFound(c, "Farm not found")
	}
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch farm")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, farm, "Farm fetched successfully")
}

// CreateFarmHandler creates a new farm
// @Summary Create Farm
// @Description Creates a new farm and assigns the creator as owner
// @Tags Farms
// @Accept json
// @Produce json
// @Param request body schemas.CreateFarmRequest true "Create Farm Request"
// @Success 201 {object} models.Farm
// @Router /v1/farms [post]
func CreateFarmHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	var req schemas.CreateFarmRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	farm, err := farmService.CreateFarm(userID, req)
	if err != nil {
		return utils.InternalError(c, "Failed to create farm")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, farm, "Farm created successfully")
}

// UpdateFarmHandler updates farm details
// @Summary Update Farm
// @Description Updates details for an existing farm
// @Tags Farms
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param request body schemas.UpdateFarmRequest true "Update Farm Request"
// @Success 200 {object} models.Farm
// @Router /v1/farms/{farm_id} [put]
func UpdateFarmHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	var req schemas.UpdateFarmRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	farm, err := farmService.UpdateFarm(userID, farmID, req)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to update farm")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, farm, "Farm updated successfully")
}

// DeleteFarmHandler soft-deletes a farm
// @Summary Delete Farm
// @Description Soft-deletes a farm by setting is_active=false
// @Tags Farms
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/farms/{farm_id} [delete]
func DeleteFarmHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	err = farmService.DeleteFarm(userID, farmID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to delete farm")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Farm deleted successfully")
}
