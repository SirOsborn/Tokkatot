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
)

// ===== ALERT HANDLERS =====

// GetFarmAlertsHandler returns alerts for a farm
// @GET /v1/farms/:farm_id/alerts?limit=50&is_active=true&severity=all
func GetFarmAlertsHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	isActiveFilter := c.Query("is_active", "true")
	severity := c.Query("severity", "all")
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT a.id, a.farm_id, a.device_id, a.alert_type, a.severity, a.message,
		       a.threshold_value, a.actual_value, a.is_active, a.triggered_at,
		       a.acknowledged_by, a.acknowledged_at, a.resolved_at, a.created_at
		FROM alerts a
		WHERE a.farm_id = $1
	`
	args := []interface{}{farmID}

	if isActiveFilter == "true" {
		query += " AND a.is_active = true"
	} else if isActiveFilter == "false" {
		query += " AND a.is_active = false"
	}

	if severity != "all" {
		query += " AND a.severity = $" + strconv.Itoa(len(args)+1)
		args = append(args, severity)
	}

	query += " ORDER BY a.triggered_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Get farm alerts error: %v", err)
		return utils.InternalError(c, "Failed to fetch alerts")
	}
	defer rows.Close()

	alerts := []models.Alert{}
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(
			&a.ID, &a.FarmID, &a.DeviceID, &a.AlertType, &a.Severity, &a.Message,
			&a.ThresholdValue, &a.ActualValue, &a.IsActive, &a.TriggeredAt,
			&a.AcknowledgedBy, &a.AcknowledgedAt, &a.ResolvedAt, &a.CreatedAt,
		); err != nil {
			continue
		}
		alerts = append(alerts, a)
	}

	var total, activeCount, criticalCount int64
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1", farmID).Scan(&total)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true", farmID).Scan(&activeCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND is_active = true AND severity = 'critical'", farmID).Scan(&criticalCount)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"alerts":         alerts,
		"total":          total,
		"active_count":   activeCount,
		"critical_count": criticalCount,
	}, "Alerts fetched successfully")
}

// GetAlertHandler returns a single alert by ID
// @GET /v1/farms/:farm_id/alerts/:alert_id
func GetAlertHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	alertID, err := uuid.Parse(c.Params("alert_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid alert ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	var a models.Alert
	err = database.DB.QueryRow(`
		SELECT id, farm_id, device_id, alert_type, severity, message,
		       threshold_value, actual_value, is_active, triggered_at,
		       acknowledged_by, acknowledged_at, resolved_at, created_at
		FROM alerts WHERE id = $1 AND farm_id = $2
	`, alertID, farmID).Scan(
		&a.ID, &a.FarmID, &a.DeviceID, &a.AlertType, &a.Severity, &a.Message,
		&a.ThresholdValue, &a.ActualValue, &a.IsActive, &a.TriggeredAt,
		&a.AcknowledgedBy, &a.AcknowledgedAt, &a.ResolvedAt, &a.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Alert not found")
	}
	if err != nil {
		log.Printf("Get alert error: %v", err)
		return utils.InternalError(c, "Failed to fetch alert")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, a, "Alert fetched successfully")
}

// AcknowledgeAlertHandler marks an alert as acknowledged
// @PUT /v1/farms/:farm_id/alerts/:alert_id/acknowledge
func AcknowledgeAlertHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	alertID, err := uuid.Parse(c.Params("alert_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid alert ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	result, err := database.DB.Exec(`
		UPDATE alerts
		SET is_active = false, acknowledged_by = $1, acknowledged_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND farm_id = $3
	`, userID, alertID, farmID)
	if err != nil {
		log.Printf("Acknowledge alert error: %v", err)
		return utils.InternalError(c, "Failed to acknowledge alert")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NotFound(c, "Alert not found")
	}

	var a models.Alert
	database.DB.QueryRow(`
		SELECT id, is_active, acknowledged_by, acknowledged_at FROM alerts WHERE id = $1
	`, alertID).Scan(&a.ID, &a.IsActive, &a.AcknowledgedBy, &a.AcknowledgedAt)

	return utils.SuccessResponse(c, fiber.StatusOK, a, "Alert acknowledged")
}

// GetAlertHistoryHandler returns historical alerts for a farm (time range)
// @GET /v1/farms/:farm_id/alerts/history?days=30&limit=200
func GetAlertHistoryHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	if err := checkFarmAccess(userID, farmID, "viewer"); err != nil {
		return err
	}

	days, _ := strconv.Atoi(c.Query("days", "30"))
	limit, _ := strconv.Atoi(c.Query("limit", "200"))
	if days > 365 {
		days = 365
	}
	if limit > 500 {
		limit = 500
	}

	rows, err := database.DB.Query(`
		SELECT id, alert_type, severity, triggered_at, resolved_at,
		       EXTRACT(EPOCH FROM (COALESCE(resolved_at, CURRENT_TIMESTAMP) - triggered_at))/60 AS duration_minutes
		FROM alerts
		WHERE farm_id = $1 AND triggered_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')
		ORDER BY triggered_at DESC
		LIMIT $3
	`, farmID, days, limit)
	if err != nil {
		log.Printf("Get alert history error: %v", err)
		return utils.InternalError(c, "Failed to fetch alert history")
	}
	defer rows.Close()

	type AlertHistoryEntry struct {
		ID              uuid.UUID `json:"id"`
		AlertType       string    `json:"alert_type"`
		Severity        string    `json:"severity"`
		TriggeredAt     string    `json:"triggered_at"`
		ResolvedAt      *string   `json:"resolved_at,omitempty"`
		DurationMinutes float64   `json:"duration_minutes"`
	}

	history := []AlertHistoryEntry{}
	for rows.Next() {
		var e AlertHistoryEntry
		if err := rows.Scan(&e.ID, &e.AlertType, &e.Severity, &e.TriggeredAt, &e.ResolvedAt, &e.DurationMinutes); err != nil {
			continue
		}
		history = append(history, e)
	}

	var total, criticalCount, warningCount, infoCount int64
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND triggered_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')", farmID, days).Scan(&total)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND severity = 'critical' AND triggered_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')", farmID, days).Scan(&criticalCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND severity = 'warning' AND triggered_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')", farmID, days).Scan(&warningCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM alerts WHERE farm_id = $1 AND severity = 'info' AND triggered_at > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 day')", farmID, days).Scan(&infoCount)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"alerts":         history,
		"total":          total,
		"critical_count": criticalCount,
		"warning_count":  warningCount,
		"info_count":     infoCount,
	}, "Alert history fetched successfully")
}

// ===== ALERT SUBSCRIPTION HANDLERS =====

// CreateAlertSubscriptionHandler creates a notification subscription for the user
// @POST /v1/users/alert-subscriptions
func CreateAlertSubscriptionHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	var req struct {
		AlertType       string  `json:"alert_type"`
		Channel         string  `json:"channel"`
		IsEnabled       bool    `json:"is_enabled"`
		QuietHoursStart *string `json:"quiet_hours_start,omitempty"`
		QuietHoursEnd   *string `json:"quiet_hours_end,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	if req.AlertType == "" {
		return utils.BadRequest(c, "missing_fields", "alert_type is required")
	}
	if req.Channel == "" {
		req.Channel = "push"
	}

	id := uuid.New()
	var sub models.AlertSubscription
	err := database.DB.QueryRow(`
		INSERT INTO alert_subscriptions (id, user_id, alert_type, channel, is_enabled, quiet_hours_start, quiet_hours_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, alert_type, channel) DO UPDATE
			SET is_enabled = $5, quiet_hours_start = $6, quiet_hours_end = $7, updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, alert_type, channel, is_enabled, quiet_hours_start, quiet_hours_end, created_at, updated_at
	`, id, userID, req.AlertType, req.Channel, req.IsEnabled, req.QuietHoursStart, req.QuietHoursEnd,
	).Scan(
		&sub.ID, &sub.UserID, &sub.AlertType, &sub.Channel, &sub.IsEnabled,
		&sub.QuietHoursStart, &sub.QuietHoursEnd, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		log.Printf("Create alert subscription error: %v", err)
		return utils.InternalError(c, "Failed to create subscription")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, sub, "Subscription created")
}

// GetAlertSubscriptionsHandler returns alert subscriptions for the current user
// @GET /v1/users/alert-subscriptions
func GetAlertSubscriptionsHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	rows, err := database.DB.Query(`
		SELECT id, user_id, alert_type, channel, is_enabled, quiet_hours_start, quiet_hours_end, created_at, updated_at
		FROM alert_subscriptions WHERE user_id = $1 ORDER BY created_at ASC
	`, userID)
	if err != nil {
		log.Printf("Get alert subscriptions error: %v", err)
		return utils.InternalError(c, "Failed to fetch subscriptions")
	}
	defer rows.Close()

	subs := []models.AlertSubscription{}
	for rows.Next() {
		var s models.AlertSubscription
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.AlertType, &s.Channel, &s.IsEnabled,
			&s.QuietHoursStart, &s.QuietHoursEnd, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			continue
		}
		subs = append(subs, s)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"subscriptions": subs,
	}, "Subscriptions fetched successfully")
}

// UpdateAlertSubscriptionHandler updates an alert subscription
// @PUT /v1/users/alert-subscriptions/:subscription_id
func UpdateAlertSubscriptionHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	subID, err := uuid.Parse(c.Params("subscription_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid subscription ID")
	}

	var req struct {
		IsEnabled       *bool   `json:"is_enabled,omitempty"`
		QuietHoursStart *string `json:"quiet_hours_start,omitempty"`
		QuietHoursEnd   *string `json:"quiet_hours_end,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	result, err := database.DB.Exec(`
		UPDATE alert_subscriptions
		SET is_enabled = COALESCE($1, is_enabled),
		    quiet_hours_start = COALESCE($2, quiet_hours_start),
		    quiet_hours_end = COALESCE($3, quiet_hours_end),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5
	`, req.IsEnabled, req.QuietHoursStart, req.QuietHoursEnd, subID, userID)
	if err != nil {
		log.Printf("Update subscription error: %v", err)
		return utils.InternalError(c, "Failed to update subscription")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NotFound(c, "Subscription not found")
	}

	var sub models.AlertSubscription
	database.DB.QueryRow(`
		SELECT id, user_id, alert_type, channel, is_enabled, quiet_hours_start, quiet_hours_end, created_at, updated_at
		FROM alert_subscriptions WHERE id = $1
	`, subID).Scan(
		&sub.ID, &sub.UserID, &sub.AlertType, &sub.Channel, &sub.IsEnabled,
		&sub.QuietHoursStart, &sub.QuietHoursEnd, &sub.CreatedAt, &sub.UpdatedAt,
	)

	return utils.SuccessResponse(c, fiber.StatusOK, sub, "Subscription updated")
}

// DeleteAlertSubscriptionHandler deletes an alert subscription
// @DELETE /v1/users/alert-subscriptions/:subscription_id
func DeleteAlertSubscriptionHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	subID, err := uuid.Parse(c.Params("subscription_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid subscription ID")
	}

	result, err := database.DB.Exec(
		"DELETE FROM alert_subscriptions WHERE id = $1 AND user_id = $2",
		subID, userID,
	)
	if err != nil {
		log.Printf("Delete subscription error: %v", err)
		return utils.InternalError(c, "Failed to delete subscription")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NotFound(c, "Subscription not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Subscription deleted",
	}, "Subscription deleted")
}
