package domain

import "time"

type HouseholdFactory struct{}

func NewHouseholdFactory() *HouseholdFactory { return &HouseholdFactory{} }

func (f *HouseholdFactory) Create(
	id, name string,
	hhType HouseholdType,
	headID string,
	members []HouseholdMember,
	currency, country string,
	tags []string,
	notes string,
) (*Household, error) {

	if err := ValidateHouseholdName(name); err != nil { return nil, err }
	if err := ValidateHouseholdType(hhType); err != nil { return nil, err }
	if err := ValidateCurrency(currency); err != nil { return nil, err }
	if err := ValidateMemberCount(members); err != nil { return nil, err }
	if err := ValidateHeadOfHousehold(members, headID); err != nil { return nil, err }
	for _, m := range members {
		if err := ValidateMemberRole(m.Role); err != nil { return nil, err }
	}

	if tags == nil { tags = []string{} }

	now := time.Now().UTC()
	return &Household{
		id: id, name: name, householdType: hhType,
		headOfHouseholdID: headID, members: members,
		status: HHStatusDraft, currency: currency, country: country,
		totalAssets: 0, totalLiabilities: 0, totalNetWorth: 0,
		health: HHHealthy, metadata: map[string]string{},
		tags: tags, notes: notes, createdAt: now, updatedAt: now,
	}, nil
}
