package http

import (
	"github.com/horizon/core/services/domains/financial-event/internal/application"
)

type EventHTTPHandler struct {
	svc *application.EventService
}

func NewEventHTTPHandler(svc *application.EventService) *EventHTTPHandler {
	return &EventHTTPHandler{svc: svc}
}
