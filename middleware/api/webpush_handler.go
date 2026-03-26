package api

import (
	"middleware/schemas"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
)

// GetVapidPublicKeyHandler returns the public key required by the browser
// @Summary Get VAPID Public Key
// @Description Returns the VAPID public key for Web Push subscription
// @Tags WebPush
// @Produce json
// @Success 200 {object} schemas.VapidPublicKeyRes
// @Router /v1/users/push-key [get]
func GetVapidPublicKeyHandler(c *fiber.Ctx) error {
	pubKey := webPushService.GetPublicKey()
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"public_key": pubKey,
		},
	})
}

// SubscribePushHandler saves the browser push subscription
// @Summary Subscribe to Push Notifications
// @Description Saves a browser PushSubscription object to the database
// @Tags WebPush
// @Accept json
// @Produce json
// @Param request body schemas.WebPushSubscriptionReq true "Push Subscription"
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/users/push-subscribe [post]
func SubscribePushHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}

	var req schemas.WebPushSubscriptionReq
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid JSON body")
	}

	if req.Endpoint == "" || req.Keys.Auth == "" || req.Keys.P256dh == "" {
		return utils.BadRequest(c, "missing_fields", "Endpoint, auth, and p256dh are required")
	}

	userAgent := c.Get("User-Agent")
	if err := webPushService.Subscribe(userID, req, userAgent); err != nil {
		return utils.InternalError(c, "Failed to completely save push subscription")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Successfully subscribed to push notifications")
}

// UnsubscribePushHandler removes the browser push subscription
// @Summary Unsubscribe from Push Notifications
// @Description Removes a browser PushSubscription from the database
// @Tags WebPush
// @Accept json
// @Produce json
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/users/push-unsubscribe [post]
func UnsubscribePushHandler(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Invalid session")
	}

	// For simplicity, we accept the endpoint to delete
	type UnsubReq struct {
		Endpoint string `json:"endpoint"`
	}
	var req UnsubReq
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid JSON payload")
	}

	if req.Endpoint != "" {
		_ = webPushService.Unsubscribe(userID, req.Endpoint)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Successfully unsubscribed")
}
