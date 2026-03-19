package api

import (
	"log"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
)

// ===== WEBSOCKET HANDLERS =====

type Hub struct{}

func (h *Hub) RunHub() {
	log.Println("WebSocket Hub started (stub)")
}

var WSHub = &Hub{}

// WebSocketUpgradeHandler upgrades a request to a websocket connection
// @Summary WebSocket Upgrade
// @Description Real-time communication endpoint (requires websocket protocol)
// @Tags WebSocket
// @Success 101 {object} schemas.JSONResponse
// @Router /ws [get]
func WebSocketUpgradeHandler(c *fiber.Ctx) error {
	// Real implementation would use websocket library
	log.Println("WebSocket upgrade requested")
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "WebSocket upgrade not implemented in this refactor pass")
}

// GetWebSocketStatsHandler returns hub metrics
// @Summary WebSocket Stats
// @Description Returns internal hub status and client count
// @Tags WebSocket
// @Produce json
// @Success 200 {object} object
// @Router /api/ws/stats [get]
func GetWebSocketStatsHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "WebSocket stats retrieved (mock)")
}
