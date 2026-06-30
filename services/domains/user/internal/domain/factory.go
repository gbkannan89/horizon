package domain

import (
	"errors"
	"fmt"
	"time"
)

type UserFactory struct{}

func NewUserFactory() *UserFactory { return &UserFactory{} }

func (f *UserFactory) Register(
	userID, displayName, legalName, country, baseCurrency, locale, timezone string,
	userType UserType, preferredName string, dateOfBirth *time.Time,
	financialIdentityProfile FinancialIdentityProfile,
	sourceOfTruth SOT, householdID string,
) (*User, error) {

	if err := ValidateDisplayName(displayName); err != nil {
		return nil, err
	}
	if err := ValidateLegalName(legalName); err != nil {
		return nil, err
	}
	if err := ValidateCountry(country); err != nil {
		return nil, err
	}
	if err := ValidateBaseCurrency(baseCurrency); err != nil {
		return nil, err
	}
	if err := ValidateLocale(locale); err != nil {
		return nil, err
	}
	if err := ValidateTimezone(timezone); err != nil {
		return nil, err
	}
	if err := ValidateDateOfBirth(dateOfBirth); err != nil {
		return nil, err
	}
	if err := ValidateMinorUser(userType, dateOfBirth); err != nil {
		return nil, err
	}
	if !AllUserTypes[userType] {
		return nil, fmt.Errorf("invalid user type: %s", userType)
	}
	if !AllFIPs[financialIdentityProfile] {
		return nil, fmt.Errorf("invalid financial identity profile: %s", financialIdentityProfile)
	}
	if err := ValidateSOT(sourceOfTruth); err != nil {
		return nil, err
	}
	if preferredName != "" && len([]rune(preferredName)) > 100 {
		return nil, errors.New("preferred name must be 100 characters or fewer")
	}

	now := time.Now().UTC()
	memberships := []string{}
	if householdID != "" {
		memberships = append(memberships, householdID)
	}

	return &User{
		userID:                    userID,
		displayName:               displayName,
		legalName:                 legalName,
		preferredName:             preferredName,
		dateOfBirth:               dateOfBirth,
		country:                   country,
		baseCurrency:              baseCurrency,
		locale:                    locale,
		timezone:                  timezone,
		userHealth:                HealthHealthy,
		userConfidence:            ConfUserVerified,
		financialIdentityProfile:  financialIdentityProfile,
		status:                    StatusRegistered,
		userType:                  userType,
		householdMemberships:      memberships,
		privacyProfile:            DefaultPrivacySetting(),
		consentProfile:            []ConsentRecord{},
		financialBehaviourProfile: BehavBalanced,
		preferences:               make(map[string]interface{}),
		profileCompleteness:       0,
		sourceOfTruth:             sourceOfTruth,
		tags:                      []string{},
		metadata:                  make(map[string]string),
		createdAt:                 now,
		updatedAt:                 now,
	}, nil
}
