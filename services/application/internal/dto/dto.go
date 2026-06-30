package dto

// CommandRequest wraps every state-changing request.
type CommandRequest[T any] struct {
	IdempotencyKey string `json:"idempotency_key" validate:"required,uuid"`
	Data           T      `json:"data" validate:"required"`
}

// QueryRequest wraps every read request.
type QueryRequest struct {
	Cursor   string `json:"cursor,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

// CommandResponse is returned for all state-changing operations.
type CommandResponse struct {
	Success    bool   `json:"success"`
	Idempotent bool   `json:"idempotent,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`
	RequestID  string `json:"request_id"`
	Timestamp  string `json:"timestamp"`
}

// QueryResponse wraps all read responses.
type QueryResponse[T any] struct {
	Data     T        `json:"data"`
	Metadata Metadata `json:"metadata"`
}

// Metadata contains pagination and version info.
type Metadata struct {
	RequestID   string `json:"request_id"`
	Timestamp   string `json:"timestamp"`
	Cursor      string `json:"cursor,omitempty"`
	HasMore     bool   `json:"has_more"`
	TotalCount  int    `json:"total_count,omitempty"`
}
