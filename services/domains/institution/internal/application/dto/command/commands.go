package command

type RegisterInstCommand struct {
	Name          string   `json:"name" validate:"required"`
	InstType      string   `json:"institution_type" validate:"required"`
	Country       string   `json:"country" validate:"required"`
	Website       string   `json:"website,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	Email         string   `json:"email,omitempty"`
	Address       string   `json:"address,omitempty"`
	HQ            string   `json:"headquarters,omitempty"`
	Regulator     string   `json:"regulator,omitempty"`
	RegLicense    string   `json:"regulatory_license,omitempty"`
	SourceOfTruth string   `json:"source_of_truth,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Notes         string   `json:"notes,omitempty"`
}
type InstIDCmd struct{ InstitutionID string `json:"institution_id" validate:"required"` }
type InstResult struct{ InstitutionID string `json:"institution_id"`; Status string `json:"status"`; Success bool `json:"success"` }
