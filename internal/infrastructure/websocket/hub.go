package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/IlucielI/insurance-policy-core-api/internal/model"
)

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients (key: userID)
	clients map[string]map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Broadcast messages to specific user
	broadcast chan *BroadcastMessage

	mu sync.RWMutex
}

type BroadcastMessage struct {
	UserID       string
	Notification *model.Notification
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()
			log.Printf("📲 Client connected: user=%s", client.userID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("📴 Client disconnected: user=%s", client.userID)

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.UserID]
			h.mu.RUnlock()

			if clients != nil {
				data, err := json.Marshal(message.Notification)
				if err != nil {
					log.Printf("❌ Failed to marshal notification: %v", err)
					continue
				}

				for client := range clients {
					select {
					case client.send <- data:
					default:
						h.mu.Lock()
						close(client.send)
						delete(h.clients[message.UserID], client)
						if len(h.clients[message.UserID]) == 0 {
							delete(h.clients, message.UserID)
						}
						h.mu.Unlock()
					}
				}
				log.Printf("🔔 Notification sent to user %s: %s", message.UserID, message.Notification.Type)
			}
		}
	}
}

// SendNotification sends a notification to a specific user
func (h *Hub) SendNotification(userID string, notification *model.Notification) {
	h.broadcast <- &BroadcastMessage{
		UserID:       userID,
		Notification: notification,
	}
}

// GetConnectedUsers returns count of connected users
func (h *Hub) GetConnectedUsers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}