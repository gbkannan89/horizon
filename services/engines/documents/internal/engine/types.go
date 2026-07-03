package engine

type Document struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category"`
	FileSize    int64    `json:"file_size"`
	MimeType    string   `json:"mime_type"`
	Tags        []string `json:"tags,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
