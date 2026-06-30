package command

type CreateAssetCommand struct {
	OwnerID          string            `json:"owner_id" validate:"required"`
	AssetName        string            `json:"asset_name" validate:"required"`
	Classification   string            `json:"classification" validate:"required"`
	OwnershipModel   string            `json:"ownership_model" validate:"required"`
	Currency         string            `json:"currency" validate:"required"`
	ValuationProfile string            `json:"valuation_profile"`
	ValuationMethod  string            `json:"valuation_method"`
	OwnershipPct     float64           `json:"ownership_percentage"`
	SourceOfTruth    string            `json:"source_of_truth"`
	AssetType        string            `json:"asset_type,omitempty"`
	HouseholdID      string            `json:"household_id,omitempty"`
	InstitutionID    string            `json:"institution_id,omitempty"`
	AccountID        string            `json:"account_id,omitempty"`
	Quantity         *float64          `json:"quantity,omitempty"`
	ExternalIDs      map[string]string `json:"external_identifiers,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Notes            string            `json:"notes,omitempty"`
}

type AcquireAssetCommand struct {
	AssetID   string `json:"asset_id" validate:"required"`
	CostBasis int64  `json:"cost_basis" validate:"required"`
	AcqDate   string `json:"acquisition_date" validate:"required"`
}

type ActivateAssetCommand struct{ AssetID string `json:"asset_id" validate:"required"` }

type PartiallyDisposeCommand struct {
	AssetID      string  `json:"asset_id" validate:"required"`
	QtySold      float64 `json:"quantity_sold" validate:"required"`
	Proceeds     int64   `json:"proceeds" validate:"required"`
	RemainingQty float64 `json:"remaining_quantity" validate:"required"`
}

type FullyDisposeCommand struct {
	AssetID  string `json:"asset_id" validate:"required"`
	Proceeds int64  `json:"proceeds" validate:"required"`
	DispDate string `json:"disposition_date" validate:"required"`
}

type RevalueAssetCommand struct {
	AssetID string `json:"asset_id" validate:"required"`
	Value   int64  `json:"value" validate:"required"`
	Method  string `json:"method"`
}

type SplitAssetCommand struct {
	AssetID string  `json:"asset_id" validate:"required"`
	OldQty  float64 `json:"old_quantity" validate:"required"`
	NewQty  float64 `json:"new_quantity" validate:"required"`
}

type ArchiveAssetCommand struct{ AssetID string `json:"asset_id" validate:"required"` }

type AssetResult struct {
	AssetID string `json:"asset_id"`
	Status  string `json:"status"`
	Success bool   `json:"success"`
}
