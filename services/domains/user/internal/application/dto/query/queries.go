package query

type GetUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
}

type ListUsersByStatusQuery struct {
	Status string `json:"status" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type SearchUsersQuery struct {
	Query  string `json:"query" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type UserResult struct {
	UserID                    string `json:"user_id"`
	DisplayName               string `json:"display_name"`
	LegalName                 string `json:"legal_name"`
	PreferredName             string `json:"preferred_name,omitempty"`
	Country                   string `json:"country"`
	BaseCurrency              string `json:"base_currency"`
	Locale                    string `json:"locale"`
	Timezone                  string `json:"timezone"`
	UserHealth                string `json:"user_health"`
	UserConfidence            string `json:"user_confidence"`
	FinancialIdentityProfile  string `json:"financial_identity_profile"`
	Status                    string `json:"status"`
	UserType                  string `json:"user_type"`
	ProfileCompleteness       float64 `json:"profile_completeness"`
	SourceOfTruth             string `json:"source_of_truth"`
	CreatedAt                 string `json:"created_at"`
	UpdatedAt                 string `json:"updated_at"`
}

type PaginatedResult struct {
	Users     []UserResult `json:"users"`
	NextCursor string       `json:"next_cursor,omitempty"`
	HasMore    bool         `json:"has_more"`
}

type ConsentResult struct {
	ConsentType string  `json:"consent_type"`
	Granted     bool    `json:"granted"`
	GrantedAt   string  `json:"granted_at"`
	RevokedAt   *string `json:"revoked_at,omitempty"`
	Scope       string  `json:"scope"`
}
