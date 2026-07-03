package engine

type Integration struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	APIKey      string `json:"api_key,omitempty"`
	WebhookURL  string `json:"webhook_url,omitempty"`
	Enabled     bool   `json:"enabled"`
	LastSyncAt  string `json:"last_sync_at,omitempty"`
	CreatedAt   string `json:"created_at"`
}
