package api

import (
	"database/sql"
	"log"

	"middleware/database"
	"middleware/models"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ===== USER PROFILE HANDLERS =====

// GetCurrentUserHandler returns the current authenticated user's profile
// @GET /v1/users/me
func GetCurrentUserHandler(c *fiber.Ctx) error {
	// Get user ID from JWT context (set by AuthMiddleware)
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	var user models.User
	query := `
	SELECT id, email, phone, name, language, timezone, avatar_url, is_active, contact_verified, last_login, created_at
	FROM users
	WHERE id = $1
	`

	err := database.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.Name,
		&user.Language,
		&user.Timezone,
		&user.AvatarURL,
		&user.IsActive,
		&user.ContactVerified,
		&user.LastLogin,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "User not found")
	}

	if err != nil {
		log.Printf("Get user error: %v", err)
		return utils.InternalError(c, "Failed to fetch user")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user, "User fetched successfully")
}

// UpdateProfileHandler updates the current user's profile
// @PUT /v1/users/me
func UpdateProfileHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	var req struct {
		Name     *string `json:"name,omitempty"`
		Language *string `json:"language,omitempty"`
		Timezone *string `json:"timezone,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	// Build dynamic UPDATE query
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Name != nil && *req.Name != "" {
		if err := utils.ValidateName(*req.Name); err != nil {
			return utils.BadRequest(c, "invalid_name", err.Error())
		}
		updates = append(updates, "name = $"+string(rune('0'+argCount)))
		args = append(args, *req.Name)
		argCount++
	}

	if req.Language != nil && *req.Language != "" {
		if err := utils.ValidateLanguage(*req.Language); err != nil {
			return utils.BadRequest(c, "invalid_language", err.Error())
		}
		updates = append(updates, "language = $"+string(rune('0'+argCount)))
		args = append(args, *req.Language)
		argCount++
	}

	if req.Timezone != nil && *req.Timezone != "" {
		updates = append(updates, "timezone = $"+string(rune('0'+argCount)))
		args = append(args, *req.Timezone)
		argCount++
	}

	if len(updates) == 0 {
		return utils.BadRequest(c, "no_updates", "No fields to update")
	}

	// Add updated_at and user_id
	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, userID)

	query := "UPDATE users SET " + updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = $" + string(rune('0'+len(args))) + " RETURNING id, name, language, timezone"

	var user struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Language string    `json:"language"`
		Timezone string    `json:"timezone"`
	}

	err := database.DB.QueryRow(query, args...).Scan(&user.ID, &user.Name, &user.Language, &user.Timezone)
	if err != nil {
		log.Printf("Update profile error: %v", err)
		return utils.InternalError(c, "Failed to update profile")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user, "Profile updated successfully")
}

// ChangePasswordHandler changes the current user's password
// @POST /v1/users/me/change-password
func ChangePasswordHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return utils.BadRequest(c, "missing_fields", "Current and new passwords are required")
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return utils.BadRequest(c, "weak_password", err.Error())
	}

	// Get current password hash
	var currentHash string
	err := database.DB.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&currentHash)
	if err != nil {
		log.Printf("Get password error: %v", err)
		return utils.InternalError(c, "Failed to verify password")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword)); err != nil {
		return utils.Unauthorized(c, "Current password is incorrect")
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hash error: %v", err)
		return utils.InternalError(c, "Failed to update password")
	}

	// Update password
	_, err = database.DB.Exec("UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", string(newHash), userID)
	if err != nil {
		log.Printf("Update password error: %v", err)
		return utils.InternalError(c, "Failed to update password")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Password changed successfully",
	}, "Password updated")
}
