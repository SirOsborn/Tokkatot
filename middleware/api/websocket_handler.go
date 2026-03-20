package api

import (
	"encoding/json"
	"log"
	"time"

	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// ===== WEBSOCKET HANDLERS =====

var WSHub = NewHub()

type WSClient struct {
	conn   *websocket.Conn
	userID uuid.UUID
	farmID uuid.UUID
	role   string
	send   chan []byte
}

type WSIncoming struct {
	Type   string          `json:"type"`
	FarmID string          `json:"farm_id,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// WebSocketUpgradeMiddleware ensures the request is a websocket upgrade.
func WebSocketUpgradeMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// WebSocketHandler handles the websocket connection lifecycle.
func WebSocketHandler(c *websocket.Conn) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"))
		_ = c.Close()
		return
	}

	farmID, _ := c.Locals("farm_id").(uuid.UUID)
	role, _ := c.Locals("role").(string)

	client := &WSClient{
		conn:   c,
		userID: userID,
		farmID: farmID,
		role:   role,
		send:   make(chan []byte, 32),
	}

	WSHub.register <- client
	defer func() { WSHub.unregister <- client }()

	go client.writePump()
	client.readPump()
}

func (c *WSClient) readPump() {
	defer func() {
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(64 * 1024)
	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Send welcome packet (direct)
	c.sendMessage("welcome", fiber.Map{
		"user_id": c.userID,
		"role":    c.role,
		"farm_id": c.farmID,
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var incoming WSIncoming
		if err := json.Unmarshal(msg, &incoming); err != nil {
			continue
		}

		switch incoming.Type {
		case "ping":
			c.sendMessage("pong", fiber.Map{
				"ts": time.Now().UTC().Format(time.RFC3339),
			})
		case "subscribe":
			// Security: only allow subscribing to the farm bound in JWT
			if incoming.FarmID != "" && incoming.FarmID != c.farmID.String() {
			c.sendMessage("error", fiber.Map{
				"message": "farm subscription not allowed",
			})
				continue
			}
			c.sendMessage("subscribed", fiber.Map{
				"farm_id": c.farmID.String(),
			})
		default:
			// Ignore unknown message types
		}
	}
}

func (c *WSClient) writePump() {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *WSClient) sendMessage(eventType string, data interface{}) {
	msg := WSMessage{
		Type:      eventType,
		FarmID:    c.farmID.String(),
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("WS send marshal error: %v", err)
		return
	}
	select {
	case c.send <- payload:
	default:
		// Drop if the client is backed up
	}
}

// GetWebSocketStatsHandler returns hub metrics
// @Summary WebSocket Stats
// @Description Returns internal hub status and client count
// @Tags WebSocket
// @Produce json
// @Success 200 {object} object
// @Router /v1/ws/stats [get]
func GetWebSocketStatsHandler(c *fiber.Ctx) error {
	stats := WSHub.Stats()
	return utils.SuccessResponse(c, fiber.StatusOK, stats, "WebSocket stats retrieved")
}
