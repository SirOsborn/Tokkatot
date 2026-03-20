package api

import (
	"database/sql"
	"log"

	"middleware/database"
	"middleware/services"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== FARM MEMBER HANDLERS =====

// GetFarmMembersHandler returns all members of a farm
// @Summary List Farm Members
// @Description Returns all members (farmers and viewers) of a specific farm
// @Tags Farms
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Success 200 {object} schemas.PaginatedResponse
// @Router /v1/farms/{farm_id}/members [get]
func GetFarmMembersHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	members, total, err := farmService.GetFarmMembers(userID, farmID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "You do not have access to this farm")
	}
	if err != nil {
		log.Printf("Get farm members error: %v", err)
		return utils.InternalError(c, "Failed to fetch members")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"members": members,
		"total":   total,
	}, "Members fetched successfully")
}

// InviteFarmMemberHandler adds an existing user to a farm
// @Summary Invite Member
// @Description Adds an existing user to the farm with a specific role
// @Tags Farms
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param request body object true "Invite Member Request"
// @Success 201 {object} schemas.JSONResponse
// @Router /v1/farms/{farm_id}/members [post]
func InviteFarmMemberHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "farmer"); err != nil {
		return err
	}

	var req struct {
		Email *string `json:"email,omitempty"`
		Phone *string `json:"phone,omitempty"`
		Role  string  `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if req.Role != "viewer" && req.Role != "worker" {
		return utils.BadRequest(c, "invalid_role", "Role must be viewer or worker")
	}

	// Find target user
	var targetUserID uuid.UUID
	var query string
	var arg interface{}
	if req.Email != nil {
		query = "SELECT id FROM users WHERE email = $1 AND is_active = true"
		arg = req.Email
	} else {
		query = "SELECT id FROM users WHERE phone = $1 AND is_active = true"
		arg = req.Phone
	}

	err = database.DB.QueryRow(query, arg).Scan(&targetUserID)
	if err == sql.ErrNoRows {
		return utils.NotFound(c, "User not found with that contact")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to find user")
	}

	// Logic for adding member (moved from farm_handler.go)
	_, err = database.DB.Exec(`
		INSERT INTO farm_users (id, farm_id, user_id, role, invited_by, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (farm_id, user_id) DO UPDATE SET role = $4, is_active = true, updated_at = CURRENT_TIMESTAMP
	`, uuid.New(), farmID, targetUserID, req.Role, userID)
	if err != nil {
		return utils.InternalError(c, "Failed to add member")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, nil, "Member added successfully")
}

// UpdateFarmMemberRoleHandler changes a member's role
// @Summary Update Member Role
// @Description Changes the role of a farm member
// @Tags Farms
// @Accept json
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param user_id path string true "User ID (UUID)"
// @Param request body object true "Update Role Request"
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/farms/{farm_id}/members/{user_id} [put]
func UpdateFarmMemberRoleHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	targetUserID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid user ID")
	}

	if err := checkFarmAccess(userID, farmID, "farmer"); err != nil {
		return err
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if req.Role != "viewer" && req.Role != "worker" {
		return utils.BadRequest(c, "invalid_role", "Role must be viewer or worker")
	}

	_, err = database.DB.Exec(`
		UPDATE farm_users SET role = $1, updated_at = CURRENT_TIMESTAMP
		WHERE farm_id = $2 AND user_id = $3 AND is_active = true
	`, req.Role, farmID, targetUserID)
	if err != nil {
		return utils.InternalError(c, "Failed to update role")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Role updated successfully")
}

// RemoveFarmMemberHandler removes a member from a farm
// @Summary Remove Member
// @Description Soft-removes a member from the farm
// @Tags Farms
// @Produce json
// @Param farm_id path string true "Farm ID (UUID)"
// @Param user_id path string true "User ID (UUID)"
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/farms/{farm_id}/members/{user_id} [delete]
func RemoveFarmMemberHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	targetUserID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid user ID")
	}

	if err := checkFarmAccess(userID, farmID, "farmer"); err != nil {
		return err
	}

	_, err = database.DB.Exec(`
		UPDATE farm_users SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE farm_id = $1 AND user_id = $2 AND role != 'farmer'
	`, farmID, targetUserID)
	if err != nil {
		return utils.InternalError(c, "Failed to remove member")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Member removed successfully")
}
