package register

import (
	"context"
	"net/http"

	"github.com/horizon/core/services/infra/events/stream"
)

// InitHub initializes and starts the SSE Hub.
func InitHub(ctx context.Context) *stream.Hub {
	hub := stream.NewHub()
	go hub.Run(ctx)
	return hub
}

// RegisterRoutes registers the SSE stream endpoint.
func RegisterRoutes(mux *http.ServeMux, hub *stream.Hub) {
	handler := stream.NewStreamHandler(hub)
	mux.Handle("GET /api/v1/stream/events", handler)
}
