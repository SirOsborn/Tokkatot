package api

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"middleware/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// WebSocket client connection management
type Client struct {
	Conn     *websocket.Conn
	UserID   uuid.UUID
	FarmID   uuid.UUID
	CoopID   *uuid.UUID // Optional: filter updates by coop
	Send     chan []byte
	LastPing time.Time
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.RWMutex
}

var WSHub = &Hub{
	Clients:    make(map[*Client]bool),
	Broadcast:  make(chan []byte, 256),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

// WebSocket message types
type WSMessage struct {
	Type      string      `json:"type"` // "device_update", "command_update", "alert", "heartbeat"
	Timestamp time.Time   `json:"timestamp"`
	FarmID    uuid.UUID   `json:"farm_id"`
	CoopID    *uuid.UUID  `json:"coop_id,omitempty"`
	Data      interface{} `json:"data"`
}

// RunHub starts the WebSocket hub
func (h *Hub) RunHub() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client] = true
			h.Mutex.Unlock()
			log.Printf("✅ WebSocket client registered (User: %s, Farm: %s)", client.UserID, client.FarmID)

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.Mutex.Unlock()
			log.Printf("❌ WebSocket client unregistered (User: %s)", client.UserID)

		case message := <-h.Broadcast:
			// Parse message to check farm/coop filtering
			var wsMsg WSMessage
			if err := json.Unmarshal(message, &wsMsg); err != nil {
				log.Printf("⚠️  Failed to parse WebSocket message: %v", err)
				continue
			}

			h.Mutex.RLock()
			for client := range h.Clients {
				// Only send to clients subscribed to this farm
				if client.FarmID != wsMsg.FarmID {
					continue
				}

				// If message is coop-specific, filter by coop
				if wsMsg.CoopID != nil && client.CoopID != nil && *client.CoopID != *wsMsg.CoopID {
					continue
				}

				select {
				case client.Send <- message:
				default:
					// Client buffer full, disconnect
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.Mutex.RUnlock()
		}
	}
}

// WebSocketHandler handles WebSocket connections
// GET /ws?farm_id={farm_id}&coop_id={coop_id}
func WebSocketHandler(c *websocket.Conn) {
	// Extract farm_id from query params (validated in upgrade handler)
	farmIDStr := c.Query("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError, "Invalid farm_id"))
		c.Close()
		return
	}

	// Extract user_id from locals (set by AuthMiddleware)
	userID := c.Locals("user_id").(uuid.UUID)

	// Optional coop filter
	var coopID *uuid.UUID
	if coopIDStr := c.Query("coop_id"); coopIDStr != "" {
		parsed, err := uuid.Parse(coopIDStr)
		if err == nil {
			coopID = &parsed
		}
	}

	// Verify user has access to this farm
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Access denied to farm"))
		c.Close()
		return
	}

	// Create client
	client := &Client{
		Conn:     c,
		UserID:   userID,
		FarmID:   farmID,
		CoopID:   coopID,
		Send:     make(chan []byte, 256),
		LastPing: time.Now(),
	}

	// Register client
	WSHub.Register <- client

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// WebSocketUpgradeHandler validates auth before upgrading to WebSocket
// GET /v1/ws (protected route)
func WebSocketUpgradeHandler(c *fiber.Ctx) error {
	// Validate farm_id query parameter
	farmIDStr := c.Query("farm_id")
	if farmIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "farm_id query parameter is required",
		})
	}

	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid farm_id format",
		})
	}

	// Get user_id from auth middleware
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - user_id not found",
		})
	}

	// Verify access to farm
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this farm",
		})
	}

	// Upgrade to WebSocket
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(WebSocketHandler)(c)
	}

	return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{
		"error": "WebSocket upgrade required",
	})
}

// ReadPump reads messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		WSHub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.LastPing = time.Now()
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Handle client messages (ping, subscribe, etc.)
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// Handle different message types
		msgType, ok := msg["type"].(string)
		if !ok {
			continue
		}

		switch msgType {
		case "ping":
			c.Send <- []byte(`{"type":"pong","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`)

		case "subscribe_coop":
			// Update client's coop filter
			if coopIDStr, ok := msg["coop_id"].(string); ok {
				coopID, err := uuid.Parse(coopIDStr)
				if err == nil {
					c.CoopID = &coopID
					c.Send <- []byte(`{"type":"subscribed","coop_id":"` + coopIDStr + `"}`)
				}
			}

		case "unsubscribe_coop":
			c.CoopID = nil
			c.Send <- []byte(`{"type":"unsubscribed"}`)
		}
	}
}

// WritePump sends messages to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Channel closed
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ===== HELPER FUNCTIONS FOR BROADCASTING UPDATES =====

// BroadcastDeviceUpdate sends device status update to all connected clients
func BroadcastDeviceUpdate(device models.Device) {
	msg := WSMessage{
		Type:      "device_update",
		Timestamp: time.Now(),
		FarmID:    device.FarmID,
		CoopID:    device.CoopID,
		Data: map[string]interface{}{
			"device_id":      device.ID,
			"device_name":    device.Name,
			"is_online":      device.IsOnline,
			"last_heartbeat": device.LastHeartbeat,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal device update: %v", err)
		return
	}

	WSHub.Broadcast <- msgBytes
}

// BroadcastCommandUpdate sends command status update to all connected clients
func BroadcastCommandUpdate(command models.DeviceCommand, farmID uuid.UUID, coopID *uuid.UUID) {
	msg := WSMessage{
		Type:      "command_update",
		Timestamp: time.Now(),
		FarmID:    farmID,
		CoopID:    coopID,
		Data: map[string]interface{}{
			"command_id":   command.ID,
			"device_id":    command.DeviceID,
			"command_type": command.CommandType,
			"status":       command.Status,
			"executed_at":  command.ExecutedAt,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal command update: %v", err)
		return
	}

	WSHub.Broadcast <- msgBytes
}

// BroadcastAlert sends alert notification to all connected clients
func BroadcastAlert(farmID uuid.UUID, coopID *uuid.UUID, alertType, message string) {
	msg := WSMessage{
		Type:      "alert",
		Timestamp: time.Now(),
		FarmID:    farmID,
		CoopID:    coopID,
		Data: map[string]interface{}{
			"alert_type": alertType,
			"message":    message,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal alert: %v", err)
		return
	}

	WSHub.Broadcast <- msgBytes
}

// GetWebSocketStatsHandler returns WebSocket connection statistics
// GET /v1/ws/stats
func GetWebSocketStatsHandler(c *fiber.Ctx) error {
	// Only available to authenticated users
	userID := c.Locals("user_id").(uuid.UUID)

	WSHub.Mutex.RLock()
	totalClients := len(WSHub.Clients)

	// Count clients by farm
	farmCounts := make(map[uuid.UUID]int)
	for client := range WSHub.Clients {
		if client.UserID == userID {
			farmCounts[client.FarmID]++
		}
	}
	WSHub.Mutex.RUnlock()

	return c.JSON(fiber.Map{
		"total_connections": totalClients,
		"your_farms":        farmCounts,
		"timestamp":         time.Now(),
	})
}
