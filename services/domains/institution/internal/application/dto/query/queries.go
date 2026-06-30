package query

type GetByIDQuery struct{ InstitutionID string `json:"institution_id" validate:"required"` }
type ListByTypeQuery struct{ InstType string `json:"institution_type" validate:"required"`; Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"` }
type ListByCountryQuery struct{ Country string `json:"country" validate:"required"`; Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"` }
type ListByStatusQuery struct{ Status string `json:"status" validate:"required"`; Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"` }
type SearchQuery struct{ Query string `json:"query" validate:"required"`; Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"` }
type InstView struct {
	InstitutionID string `json:"institution_id"`; Name string `json:"name"`
	InstType      string `json:"institution_type"`; Category string `json:"category"`
	Status        string `json:"status"`; Country string `json:"country"`
	TrustLevel    string `json:"trust_level"`; Health string `json:"health"`
	Confidence    string `json:"confidence"`; CreatedAt string `json:"created_at"`
}
type PaginatedResult struct{ Institutions []InstView `json:"institutions"`; NextCursor string `json:"next_cursor,omitempty"`; HasMore bool `json:"has_more"` }
