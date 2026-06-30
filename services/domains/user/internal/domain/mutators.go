package domain

import "time"

func (u *User) SetDisplayName(v string)          { u.displayName = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetLegalName(v string)             { u.legalName = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetPreferredName(v string)         { u.preferredName = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetCountry(v string)               { u.country = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetBaseCurrency(v string)          { u.baseCurrency = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetLocale(v string)                { u.locale = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetTimezone(v string)              { u.timezone = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetPreferences(v map[string]interface{}) { u.preferences = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetPrivacy(v PrivacySetting)       { u.privacyProfile = v; u.updatedAt = time.Now().UTC() }
func (u *User) SetProfileCompleteness(v float64)   { u.profileCompleteness = v }
func (u *User) SetFinancialBehaviourProfile(v FinBehaviourProfile) { u.financialBehaviourProfile = v }
func (u *User) AddHouseholdMembership(h string) {
	for _, m := range u.householdMemberships {
		if m == h { return }
	}
	u.householdMemberships = append(u.householdMemberships, h)
	u.updatedAt = time.Now().UTC()
}
func (u *User) RemoveHouseholdMembership(h string) {
	for i, m := range u.householdMemberships {
		if m == h {
			u.householdMemberships = append(u.householdMemberships[:i], u.householdMemberships[i+1:]...)
			u.updatedAt = time.Now().UTC()
			return
		}
	}
}
