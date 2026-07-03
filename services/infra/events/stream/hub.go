package stream

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/horizon/core/packages/events"
)

// Client represents an active SSE connection.
type Client struct {
	ID     string
	UserID string
	Send   chan events.Envelope
}

// Hub maintains the set of active clients and broadcasts events to the clients.
type Hub struct {
	clients map[*Client]bool
	userMap map[string]map[*Client]bool
	
	register   chan *Client
	unregister chan *Client
	broadcast  chan events.Envelope
	
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		userMap:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan events.Envelope),
	}
}

func (h *Hub) Run(ctx context.Context) {
	// Send heartbeats every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if h.userMap[client.UserID] == nil {
				h.userMap[client.UserID] = make(map[*Client]bool)
			}
			h.userMap[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("SSE Client connected: %s (User: %s)", client.ID, client.UserID)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if h.userMap[client.UserID] != nil {
					delete(h.userMap[client.UserID], client)
					if len(h.userMap[client.UserID]) == 0 {
						delete(h.userMap, client.UserID)
					}
				}
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("SSE Client disconnected: %s", client.ID)
		case env := <-h.broadcast:
			h.mu.RLock()
			// If UserID is specified, broadcast only to that user
			if env.UserID != "" {
				if userClients, ok := h.userMap[env.UserID]; ok {
					for client := range userClients {
						select {
						case client.Send <- env:
						default:
							// If buffer full, drop and cleanup
							close(client.Send)
							delete(h.clients, client)
							delete(h.userMap[client.UserID], client)
						}
					}
				}
			} else {
				// Broadcast to all
				for client := range h.clients {
					select {
					case client.Send <- env:
					default:
						close(client.Send)
						delete(h.clients, client)
						if h.userMap[client.UserID] != nil {
							delete(h.userMap[client.UserID], client)
						}
					}
				}
			}
			h.mu.RUnlock()
		case <-ticker.C:
			// Heartbeat
			env := events.NewEnvelope("system.heartbeat", 1, json.RawMessage(`{}`))
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- env:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Publish implements the events.Publisher interface.
func (h *Hub) Publish(ctx context.Context, envelope events.Envelope) error {
	select {
	case h.broadcast <- envelope:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// RegisterClient adds a new client.
func (h *Hub) RegisterClient(c *Client) {
	h.register <- c
}

// UnregisterClient removes an active client.
func (h *Hub) UnregisterClient(c *Client) {
	h.unregister <- c
}
