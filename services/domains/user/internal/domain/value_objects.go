package domain

import "time"

type ConsentRecord struct {
	ConsentType ConsentType `json:"consent_type"`
	Granted     bool        `json:"granted"`
	GrantedAt   time.Time   `json:"granted_at"`
	RevokedAt   *time.Time  `json:"revoked_at,omitempty"`
	Scope       string      `json:"scope"`
}

type PrivacySetting struct {
	DataSharing        DataSharingLevel `json:"data_sharing"`
	HouseholdVisibility DataSharingLevel `json:"household_visibility"`
	ThirdPartyAccess   bool             `json:"third_party_access"`
}

func DefaultPrivacySetting() PrivacySetting {
	return PrivacySetting{
		DataSharing:        DSOptIn,
		HouseholdVisibility: DSRoleBased,
		ThirdPartyAccess:   false,
	}
}
