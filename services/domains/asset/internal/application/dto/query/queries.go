package query

type GetAssetQuery struct{ AssetID string `json:"asset_id" validate:"required"` }
type ListByUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}
type ListByClassQuery struct {
	UserID         string `json:"user_id" validate:"required"`
	Classification string `json:"classification" validate:"required"`
	Cursor         string `json:"cursor,omitempty"`; Limit int `json:"limit,omitempty"`
}

type AssetView struct {
	AssetID          string   `json:"asset_id"`
	AssetName        string   `json:"asset_name"`
	Classification   string   `json:"classification"`
	Currency         string   `json:"currency"`
	Status           string   `json:"status"`
	CostBasis        int64    `json:"cost_basis"`
	CurrentValue     int64    `json:"current_value"`
	Quantity         *float64 `json:"quantity,omitempty"`
	UnitPrice        float64  `json:"unit_price"`
	OwnershipPct     float64  `json:"ownership_percentage"`
	ValuationMethod  string   `json:"valuation_method"`
	LiquidityProfile string   `json:"liquidity_profile"`
	AssetHealth      string   `json:"asset_health"`
	CreatedAt        string   `json:"created_at"`
}

type PaginatedResult struct {
	Assets     []AssetView `json:"assets"`
	NextCursor string      `json:"next_cursor,omitempty"`
	HasMore    bool        `json:"has_more"`
}
