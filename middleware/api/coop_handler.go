package api

import (
	"log"
	"middleware/models"
	"middleware/services"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== COOP MANAGEMENT HANDLERS =====

// ListCoopsHandler returns all coops for a farm
// @Summary List Farm Coops
// @Description Returns all coops for a specific farm
// @Tags Coops
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Success 200 {object} []models.Coop
// @Router /v1/farms/{farm_id}/coops [get]
func ListCoopsHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	coops, err := coopService.ListCoops(userID, farmID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		log.Printf("List coops error: %v", err)
		return utils.InternalError(c, "Failed to fetch coops")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"coops": coops,
	}, "Coops retrieved")
}

// GetCoopHandler returns a single coop by ID
// @Summary Get Coop Details
// @Description Returns details for a specific coop
// @Tags Coops
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param coop_id path string true "Coop ID (UUID)"
// @Success 200 {object} models.Coop
// @Router /v1/farms/{farm_id}/coops/{coop_id} [get]
func GetCoopHandler(c *fiber.Ctx) error {
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

	coop, err := coopService.GetCoop(userID, farmID, coopID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrCoopNotFound {
		return utils.NotFound(c, "Coop not found")
	}
	if err != nil {
		log.Printf("Get coop error: %v", err)
		return utils.InternalError(c, "Failed to fetch coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, coop, "Coop retrieved")
}

// CreateCoopHandler creates a new coop
// @Summary Create Coop
// @Description Creates a new coop within a farm
// @Tags Coops
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param request body object true "Create Coop Request"
// @Success 201 {object} models.Coop
// @Router /v1/farms/{farm_id}/coops [post]
func CreateCoopHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	var req struct {
		Number       int     `json:"number"`
		Name         string  `json:"name"`
		Capacity     *int    `json:"capacity,omitempty"`
		CurrentCount *int    `json:"current_count,omitempty"`
		ChickenType  *string `json:"chicken_type,omitempty"`
		Description  *string `json:"description,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if req.Number < 1 || req.Name == "" {
		return utils.BadRequest(c, "invalid_request", "Coop number and name are required")
	}

	coop, err := coopService.CreateCoop(userID, farmID, models.Coop{
		Number:       req.Number,
		Name:         req.Name,
		Capacity:     req.Capacity,
		CurrentCount: req.CurrentCount,
		ChickenType:  req.ChickenType,
		Description:  req.Description,
	})
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		log.Printf("Create coop error: %v", err)
		return utils.InternalError(c, "Failed to create coop")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, coop, "Coop created")
}

// UpdateCoopHandler updates a coop
// @Summary Update Coop
// @Description Updates an existing coop
// @Tags Coops
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param coop_id path string true "Coop ID (UUID)"
// @Param request body object true "Update Coop Request"
// @Success 200 {object} models.Coop
// @Router /v1/farms/{farm_id}/coops/{coop_id} [put]
func UpdateCoopHandler(c *fiber.Ctx) error {
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

	var req struct {
		Number       *int    `json:"number,omitempty"`
		Name         *string `json:"name,omitempty"`
		Capacity     *int    `json:"capacity,omitempty"`
		CurrentCount *int    `json:"current_count,omitempty"`
		ChickenType  *string `json:"chicken_type,omitempty"`
		Description  *string `json:"description,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	coop, err := coopService.UpdateCoop(userID, farmID, coopID, req.Number, req.Name, req.Capacity, req.CurrentCount, req.ChickenType, req.Description)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err == services.ErrCoopNotFound {
		return utils.NotFound(c, "Coop not found")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to update coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, coop, "Coop updated")
}

// DeleteCoopHandler deletes a coop
// @Summary Delete Coop
// @Description Deletes a specific coop
// @Tags Coops
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param coop_id path string true "Coop ID (UUID)"
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/farms/{farm_id}/coops/{coop_id} [delete]
func DeleteCoopHandler(c *fiber.Ctx) error {
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

	if err := coopService.DeleteCoop(userID, farmID, coopID); err != nil {
		if err == services.ErrFarmAccessDenied {
			return utils.Forbidden(c, "Access denied")
		}
		if err == services.ErrCoopNotFound {
			return utils.NotFound(c, "Coop not found")
		}
		return utils.InternalError(c, "Failed to delete coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Coop deleted")
}

// TemperatureTimelineHandler returns temperature data points
// @Summary Get Temperature Timeline
// @Description Returns historical temperature data for a coop
// @Tags Coops
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param coop_id path string true "Coop ID (UUID)"
// @Success 200 {object} []interface{}
// @Router /v1/farms/{farm_id}/coops/{coop_id}/temperature-timeline [get]
func TemperatureTimelineHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Temperature timeline retrieved (mock)")
}
