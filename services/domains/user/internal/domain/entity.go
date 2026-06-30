package domain

import (
	"fmt"
	"time"
)

type User struct {
	userID                    string
	displayName               string
	legalName                 string
	preferredName             string
	dateOfBirth               *time.Time
	country                   string
	baseCurrency              string
	locale                    string
	timezone                  string
	userHealth                UserHealth
	userConfidence            UserConfidence
	financialIdentityProfile  FinancialIdentityProfile
	status                    UserStatus
	userType                  UserType
	householdMemberships      []string
	privacyProfile            PrivacySetting
	consentProfile            []ConsentRecord
	financialBehaviourProfile FinBehaviourProfile
	preferences               map[string]interface{}
	profileCompleteness       float64
	sourceOfTruth             SOT
	tags                      []string
	metadata                  map[string]string
	notes                     string
	createdAt                 time.Time
	updatedAt                 time.Time
}

func (u *User) UserID() string                       { return u.userID }
func (u *User) DisplayName() string                  { return u.displayName }
func (u *User) LegalName() string                    { return u.legalName }
func (u *User) PreferredName() string                { return u.preferredName }
func (u *User) DateOfBirth() *time.Time              { return u.dateOfBirth }
func (u *User) Country() string                      { return u.country }
func (u *User) BaseCurrency() string                 { return u.baseCurrency }
func (u *User) Locale() string                       { return u.locale }
func (u *User) Timezone() string                     { return u.timezone }
func (u *User) UserHealth() UserHealth               { return u.userHealth }
func (u *User) UserConfidence() UserConfidence       { return u.userConfidence }
func (u *User) FinancialIdentityProfile() FinancialIdentityProfile { return u.financialIdentityProfile }
func (u *User) Status() UserStatus                   { return u.status }
func (u *User) UserType() UserType                   { return u.userType }
func (u *User) HouseholdMemberships() []string       { return u.householdMemberships }
func (u *User) PrivacyProfile() PrivacySetting       { return u.privacyProfile }
func (u *User) ConsentProfile() []ConsentRecord      { return u.consentProfile }
func (u *User) FinancialBehaviourProfile() FinBehaviourProfile { return u.financialBehaviourProfile }
func (u *User) Preferences() map[string]interface{}  { return u.preferences }
func (u *User) ProfileCompleteness() float64         { return u.profileCompleteness }
func (u *User) SourceOfTruth() SOT                   { return u.sourceOfTruth }
func (u *User) Tags() []string                       { return u.tags }
func (u *User) Metadata() map[string]string          { return u.metadata }
func (u *User) Notes() string                        { return u.notes }
func (u *User) CreatedAt() time.Time                 { return u.createdAt }
func (u *User) UpdatedAt() time.Time                 { return u.updatedAt }

func (u *User) SetStatus(s UserStatus)              { u.status = s; u.updatedAt = time.Now().UTC() }
func (u *User) SetUpdatedAt(t time.Time)             { u.updatedAt = t.UTC() }

func (u *User) CanTransitionTo(target UserStatus) error {
	transitions := map[UserStatus]map[UserStatus]bool{
		StatusRegistered: {StatusActive: true, StatusArchived: true},
		StatusActive:     {StatusInactive: true, StatusSuspended: true, StatusMerged: true, StatusArchived: true},
		StatusInactive:   {StatusActive: true, StatusArchived: true},
		StatusSuspended:  {StatusActive: true},
		StatusMerged:     {StatusArchived: true},
	}
	if u.status == StatusArchived && target == StatusHistorical {
		return nil
	}
	if allowed, ok := transitions[u.status]; ok {
		if allowed[target] {
			return nil
		}
	}
	if u.status == target {
		return fmt.Errorf("user is already in %s state", target)
	}
	if u.status == StatusArchived || u.status == StatusHistorical || u.status == StatusMerged {
		return fmt.Errorf("%s users cannot transition: %s", u.status, target)
	}
	return fmt.Errorf("cannot transition from %s to %s", u.status, target)
}

func (u *User) AddConsent(c ConsentRecord) {
	for i, existing := range u.consentProfile {
		if existing.ConsentType == c.ConsentType && existing.Granted {
			u.consentProfile[i].Granted = false
			now := time.Now().UTC()
			u.consentProfile[i].RevokedAt = &now
		}
	}
	u.consentProfile = append(u.consentProfile, c)
	u.updatedAt = time.Now().UTC()
}

func (u *User) RevokeConsent(ct ConsentType) bool {
	for i, c := range u.consentProfile {
		if c.ConsentType == ct && c.Granted {
			now := time.Now().UTC()
			u.consentProfile[i].Granted = false
			u.consentProfile[i].RevokedAt = &now
			u.updatedAt = time.Now().UTC()
			return true
		}
	}
	return false
}

func (u *User) GetActiveConsents() []ConsentRecord {
	var active []ConsentRecord
	for _, c := range u.consentProfile {
		if c.Granted {
			active = append(active, c)
		}
	}
	return active
}

// ReconstructFromDB creates a User from persistent storage without validation.
func ReconstructFromDB(id, displayName, legalName, preferredName, country, baseCurrency, locale, timezone string,
	userHealth UserHealth, userConfidence UserConfidence, fip FinancialIdentityProfile,
	status UserStatus, userType UserType, beh FinBehaviourProfile, profileCompleteness float64,
	sot SOT, tags []string, metadata map[string]string, notes string,
	dob *time.Time, createdAt, updatedAt time.Time) *User {
	return &User{
		userID: id, displayName: displayName, legalName: legalName, preferredName: preferredName,
		dateOfBirth: dob, country: country, baseCurrency: baseCurrency, locale: locale, timezone: timezone,
		userHealth: userHealth, userConfidence: userConfidence, financialIdentityProfile: fip,
		status: status, userType: userType, financialBehaviourProfile: beh,
		profileCompleteness: profileCompleteness, sourceOfTruth: sot,
		tags: tags, metadata: metadata, notes: notes, createdAt: createdAt, updatedAt: updatedAt,
	}
}