package api

import "github.com/horizon/core/services/experiences/notifications/internal/engine"

type CenterResponse struct {
	Success  bool                      `json:"success"`
	Data     *engine.NotificationCenter `json:"data,omitempty"`
	Metadata *Metadata                 `json:"metadata,omitempty"`
}

type NotifListResponse struct {
	Success  bool               `json:"success"`
	Data     *NotifList         `json:"data,omitempty"`
	Metadata *Metadata          `json:"metadata,omitempty"`
}

type NotifList struct {
	Notifications []engine.Notification `json:"notifications"`
	Count         int                   `json:"count"`
}

type NotifResponse struct {
	Success  bool                `json:"success"`
	Data     *engine.Notification `json:"data,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
}

type PrefListResponse struct {
	Success  bool                 `json:"success"`
	Data     *PrefList            `json:"data,omitempty"`
	Metadata *Metadata            `json:"metadata,omitempty"`
}

type PrefList struct {
	Preferences []engine.Preference `json:"preferences"`
	Count       int                 `json:"count"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
