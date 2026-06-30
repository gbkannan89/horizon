package command

type RegisterUserCommand struct {
	UserID                    string `json:"user_id"`
	DisplayName               string `json:"display_name" validate:"required"`
	LegalName                 string `json:"legal_name" validate:"required"`
	PreferredName             string `json:"preferred_name,omitempty"`
	DateOfBirth               string `json:"date_of_birth,omitempty"`
	Country                   string `json:"country" validate:"required"`
	BaseCurrency              string `json:"base_currency" validate:"required"`
	Locale                    string `json:"locale" validate:"required"`
	Timezone                  string `json:"timezone" validate:"required"`
	UserType                  string `json:"user_type"`
	FinancialIdentityProfile  string `json:"financial_identity_profile"`
	SourceOfTruth             string `json:"source_of_truth"`
	HouseholdID               string `json:"household_id,omitempty"`
}

type ActivateUserCommand struct {
	UserID string `json:"user_id" validate:"required"`
}

type SuspendUserCommand struct {
	UserID string `json:"user_id" validate:"required"`
	Reason string `json:"reason" validate:"required"`
}

type ReactivateUserCommand struct {
	UserID string `json:"user_id" validate:"required"`
}

type ArchiveUserCommand struct {
	UserID string `json:"user_id" validate:"required"`
	Reason string `json:"reason,omitempty"`
}

type UpdateProfileCommand struct {
	UserID         string `json:"user_id" validate:"required"`
	DisplayName    string `json:"display_name,omitempty"`
	LegalName      string `json:"legal_name,omitempty"`
	PreferredName  string `json:"preferred_name,omitempty"`
	DateOfBirth    string `json:"date_of_birth,omitempty"`
	Country        string `json:"country,omitempty"`
	BaseCurrency   string `json:"base_currency,omitempty"`
	Locale         string `json:"locale,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
}

type UpdatePreferencesCommand struct {
	UserID      string                 `json:"user_id" validate:"required"`
	Preferences map[string]interface{} `json:"preferences" validate:"required"`
}

type GrantConsentCommand struct {
	UserID      string `json:"user_id" validate:"required"`
	ConsentType string `json:"consent_type" validate:"required"`
	Scope       string `json:"scope" validate:"required"`
}

type RevokeConsentCommand struct {
	UserID      string `json:"user_id" validate:"required"`
	ConsentType string `json:"consent_type" validate:"required"`
}

type UpdatePrivacyCommand struct {
	UserID            string `json:"user_id" validate:"required"`
	DataSharing       string `json:"data_sharing,omitempty"`
	HouseholdVis      string `json:"household_visibility,omitempty"`
	ThirdPartyAccess  *bool  `json:"third_party_access,omitempty"`
}

type RegisterResult struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type CommandResult struct {
	Success bool   `json:"success"`
	UserID  string `json:"user_id,omitempty"`
	Status  string `json:"status,omitempty"`
}
