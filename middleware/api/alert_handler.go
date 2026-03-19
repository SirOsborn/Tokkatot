package api

import (
	"log"
	"middleware/services"
	"middleware/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== ALERT HANDLERS =====

// ===== ALERT HANDLERS =====

// GetFarmAlertsHandler returns filtered alerts for a farm
// @Summary List Farm Alerts
// @Description Returns recent alerts for a specific farm
// @Tags Alerts
// @Security ApiKeyAuth
// @Param farm_id path string true "Farm ID"
// @Param limit query int false "Limit (default 50)"
// @Param offset query int false "Offset"
// @Param is_active query bool false "Filter by active status"
// @Param severity query string false "Filter by severity (info, warning, critical)"
// @Success 200 {object} schemas.ErrorResponse
// @Router /api/farms/{farm_id}/alerts [get]
func GetFarmAlertsHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	isActiveFilter := c.Query("is_active", "true") == "true"
	severity := c.Query("severity", "all")

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	alerts, total, activeCount, criticalCount, err := alertService.GetFarmAlerts(userID, farmID, limit, offset, isActiveFilter, severity)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		log.Printf("Get farm alerts error: %v", err)
		return utils.InternalError(c, "Failed to fetch alerts")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"alerts":         alerts,
		"total":          total,
		"active_count":   activeCount,
		"critical_count": criticalCount,
	}, "Alerts retrieved")
}

// GetAlertHistoryHandler returns historical alerts
// @Summary Get Alert History
// @Description Returns a paginated history of alerts for a farm
// @Tags Alerts
// @Security ApiKeyAuth
// @Param farm_id path string true "Farm ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} schemas.PaginatedResponse{data=[]models.Alert}
// @Router /api/farms/{farm_id}/alerts/history [get]
func GetAlertHistoryHandler(c *fiber.Ctx) error {
	// Reusing GetFarmAlerts with active_only=false for history
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	alerts, total, _, _, err := alertService.GetFarmAlerts(userID, farmID, limit, offset, false, "all")
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Access denied")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to fetch alert history")
	}

	return utils.SuccessListResponse(c, alerts, total, offset/limit+1, limit)
}

// AcknowledgeAlertHandler marks an alert as read
// @Summary Acknowledge Alert
// @Description Marks a specific alert as acknowledged
// @Tags Alerts
// @Security ApiKeyAuth
// @Param farm_id path string true "Farm ID"
// @Param alert_id path string true "Alert ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/farms/{farm_id}/alerts/{alert_id}/acknowledge [post]
func AcknowledgeAlertHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}
	farmID, err := uuid.Parse(c.Params("farm_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}
	alertID, err := uuid.Parse(c.Params("alert_id"))
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid alert ID")
	}

	err = alertService.AcknowledgeAlert(userID, farmID, alertID)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "Permission denied")
	}
	if err == services.ErrAlertNotFound {
		return utils.NotFound(c, "Alert not found or already acknowledged")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to acknowledge alert")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Alert acknowledged")
}

// GetAlertHandler returns a single alert by ID
func GetAlertHandler(c *fiber.Ctx) error {
	// Simple implementation using list logic if needed
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Not implemented yet - use list with filter")
}

// ===== ALERT SUBSCRIPTION HANDLERS =====

// CreateAlertSubscriptionHandler creates a new subscription
func CreateAlertSubscriptionHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusCreated, nil, "Subscription created (mock)")
}

// GetAlertSubscriptionsHandler returns all subscriptions for a user
func GetAlertSubscriptionsHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, []interface{}{}, "Subscriptions retrieved (mock)")
}

// UpdateAlertSubscriptionHandler updates a subscription
func UpdateAlertSubscriptionHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Subscription updated (mock)")
}

// DeleteAlertSubscriptionHandler deletes a subscription
func DeleteAlertSubscriptionHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Subscription deleted (mock)")
}
