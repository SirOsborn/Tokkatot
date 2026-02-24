package api

import (
	"database/sql"
	"log"
	"strconv"

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

// GetUserSessionsHandler returns active sessions for the current user
// @GET /v1/users/sessions
func GetUserSessionsHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	rows, err := database.DB.Query(`
		SELECT id, device_name, ip_address, user_agent, last_activity, expires_at, created_at
		FROM user_sessions
		WHERE user_id = $1 AND expires_at > CURRENT_TIMESTAMP
		ORDER BY last_activity DESC
	`, userID)
	if err != nil {
		log.Printf("Get sessions error: %v", err)
		return utils.InternalError(c, "Failed to fetch sessions")
	}
	defer rows.Close()

	type SessionInfo struct {
		ID           string  `json:"id"`
		DeviceName   *string `json:"device_name"`
		IPAddress    *string `json:"ip_address"`
		UserAgent    *string `json:"user_agent"`
		LastActivity string  `json:"last_activity"`
		ExpiresAt    string  `json:"expires_at"`
	}

	sessions := []SessionInfo{}
	for rows.Next() {
		var s SessionInfo
		if err := rows.Scan(&s.ID, &s.DeviceName, &s.IPAddress, &s.UserAgent, &s.LastActivity, &s.ExpiresAt); err != nil {
			continue
		}
		sessions = append(sessions, s)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"sessions": sessions,
	}, "Sessions fetched successfully")
}

// RevokeUserSessionHandler revokes a specific user session
// @DELETE /v1/users/sessions/:session_id
func RevokeUserSessionHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	sessionID := c.Params("session_id")
	if sessionID == "" {
		return utils.BadRequest(c, "missing_id", "Session ID is required")
	}

	result, err := database.DB.Exec(
		"DELETE FROM user_sessions WHERE id = $1 AND user_id = $2",
		sessionID, userID,
	)
	if err != nil {
		log.Printf("Revoke session error: %v", err)
		return utils.InternalError(c, "Failed to revoke session")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NotFound(c, "Session not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Session revoked",
	}, "Session revoked")
}

// GetUserActivityLogHandler returns recent activity for the current user
// @GET /v1/users/activity-log?limit=50&offset=0&days=30
func GetUserActivityLogHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	days, _ := strconv.Atoi(c.Query("days", "30"))

	if limit > 200 {
		limit = 200
	}
	if days > 365 {
		days = 365
	}

	rows, err := database.DB.Query(`
		SELECT id, event_type, ip_address, created_at
		FROM event_logs
		WHERE user_id = $1 AND created_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, userID, days, limit, offset)
	if err != nil {
		log.Printf("Get activity log error: %v", err)
		return utils.InternalError(c, "Failed to fetch activity log")
	}
	defer rows.Close()

	type ActivityEntry struct {
		ID        string  `json:"id"`
		EventType string  `json:"event_type"`
		IPAddress *string `json:"ip_address"`
		Timestamp string  `json:"timestamp"`
	}

	activities := []ActivityEntry{}
	for rows.Next() {
		var a ActivityEntry
		if err := rows.Scan(&a.ID, &a.EventType, &a.IPAddress, &a.Timestamp); err != nil {
			continue
		}
		activities = append(activities, a)
	}

	var total int64
	database.DB.QueryRow(
		"SELECT COUNT(*) FROM event_logs WHERE user_id = $1", userID,
	).Scan(&total)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"activities": activities,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	}, "Activity log fetched successfully")
}
