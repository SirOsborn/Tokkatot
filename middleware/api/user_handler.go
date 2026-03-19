package api

import (
	"middleware/schemas"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
)

// ===== USER PROFILE & SESSION HANDLERS =====

// GetCurrentUserHandler returns the profile of the authenticated user
// @Summary Get Current Profile
// @Description Returns the profile information of the currently authenticated user
// @Tags User
// @Produce json
// @Success 200 {object} models.User
// @Router /v1/user/profile [get]
func GetCurrentUserHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}

	user, err := authService.GetProfile(userID)
	if err != nil {
		return utils.InternalError(c, "Failed to retrieve profile")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"user": user,
	}, "User profile retrieved")
}

// UpdateProfileHandler updates the user's profile
// @Summary Update Profile
// @Description Updates the profile details of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param request body schemas.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} models.User
// @Router /v1/user/profile [put]
func UpdateProfileHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid user session")
	}

	var req schemas.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}

	user, err := authService.UpdateProfile(userID, req)
	if err != nil {
		return utils.InternalError(c, "Failed to update profile")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"user": user,
	}, "Profile updated successfully")
}

// ChangePasswordHandler changes the user's password
func ChangePasswordHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Password changed (mock)")
}

// GetUserSessionsHandler returns active sessions for the user
// @Summary List Active Sessions
// @Description Returns a list of all active sessions for the current user
// @Tags User, Sessions
// @Produce json
// @Success 200 {object} []models.UserSession
// @Router /v1/user/sessions [get]
func GetUserSessionsHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Sessions retrieved (mock)")
}

// RevokeUserSessionHandler revokes a specific session
func RevokeUserSessionHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Session revoked (mock)")
}

// GetUserActivityLogHandler returns recent activity for the user
func GetUserActivityLogHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Activity log retrieved (mock)")
}
