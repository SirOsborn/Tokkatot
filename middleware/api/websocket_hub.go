package api

import (
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// ===== WEBSOCKET HUB =====

type wsBroadcast struct {
	farmID uuid.UUID
	data   []byte
}

type WSMessage struct {
	Type      string      `json:"type"`
	FarmID    string      `json:"farm_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

type Hub struct {
	register   chan *WSClient
	unregister chan *WSClient
	broadcast  chan wsBroadcast

	clients     map[*WSClient]struct{}
	farmClients map[uuid.UUID]map[*WSClient]struct{}
	clientCount atomic.Int64
}

func NewHub() *Hub {
	return &Hub{
		register:    make(chan *WSClient),
		unregister:  make(chan *WSClient),
		broadcast:   make(chan wsBroadcast, 256),
		clients:     make(map[*WSClient]struct{}),
		farmClients: make(map[uuid.UUID]map[*WSClient]struct{}),
	}
}

func (h *Hub) RunHub() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
			if client.farmID != uuid.Nil {
				if h.farmClients[client.farmID] == nil {
					h.farmClients[client.farmID] = make(map[*WSClient]struct{})
				}
				h.farmClients[client.farmID][client] = struct{}{}
			}
			h.clientCount.Add(1)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.farmID != uuid.Nil {
					if farmSet, ok := h.farmClients[client.farmID]; ok {
						delete(farmSet, client)
						if len(farmSet) == 0 {
							delete(h.farmClients, client.farmID)
						}
					}
				}
				h.clientCount.Add(-1)
				close(client.send)
			}

		case msg := <-h.broadcast:
			targets := h.farmClients[msg.farmID]
			for client := range targets {
				select {
				case client.send <- msg.data:
				default:
					// Slow consumer; drop connection to protect hub
					close(client.send)
					delete(targets, client)
					delete(h.clients, client)
					h.clientCount.Add(-1)
				}
			}
		}
	}
}

func (h *Hub) Stats() map[string]interface{} {
	return map[string]interface{}{
		"total_clients": h.clientCount.Load(),
		"farms_active":  len(h.farmClients),
	}
}

func (h *Hub) PublishFarmEvent(farmID uuid.UUID, eventType string, data interface{}) {
	if farmID == uuid.Nil {
		return
	}
	msg := WSMessage{
		Type:      eventType,
		FarmID:    farmID.String(),
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("WS publish marshal error: %v", err)
		return
	}
	h.broadcast <- wsBroadcast{farmID: farmID, data: payload}
}

