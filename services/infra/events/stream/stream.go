package stream

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/horizon/core/packages/events"
	"github.com/horizon/core/services/internal/auth"
)

// StreamHandler handles HTTP SSE connections.
type StreamHandler struct {
	hub *Hub
}

func NewStreamHandler(hub *Hub) *StreamHandler {
	return &StreamHandler{hub: hub}
}

func (h *StreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	userID := auth.UserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// Optional: CORS headers if required
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := &Client{
		ID:     fmt.Sprintf("%s-%d", userID, r.Context().Value("request_id")), // simplistic ID
		UserID: userID,
		Send:   make(chan events.Envelope, 256),
	}

	h.hub.RegisterClient(client)
	defer h.hub.UnregisterClient(client)

	// Send an initial connected event
	connEvt := events.NewEnvelope("system.connected", 1, json.RawMessage(`{"status":"connected"}`))
	connEvt.UserID = userID
	client.Send <- connEvt

	for {
		select {
		case env, ok := <-client.Send:
			if !ok {
				// Hub closed the channel
				return
			}
			bytes, err := json.Marshal(env)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", bytes)
			flusher.Flush()
		case <-r.Context().Done():
			// Client disconnected
			return
		}
	}
}
