package domain

import (
	"fmt"
	"time"
)

type Account struct {
	id                  string
	institutionID       *string
	ownerID             string
	householdID         *string
	classification      AccountClassification
	accountType         AccountType
	accountSubType      *string
	accountName         string
	currency            string
	status              AccountStatus
	openedDate          time.Time
	closedDate          *time.Time
	liquidityProfile    LiquidityProfile
	accountHealth       AccountHealth
	creditLimit         *int64
	interestRate        *float64
	country             *string
	sourceOfTruth       *string
	visibility          Visibility
	externalReference   *string
	tags                []string
	metadata            map[string]string
	notes               string
	createdAt           time.Time
	updatedAt           time.Time
}

func (a *Account) ID() string                    { return a.id }
func (a *Account) InstitutionID() *string        { return a.institutionID }
func (a *Account) OwnerID() string               { return a.ownerID }
func (a *Account) HouseholdID() *string          { return a.householdID }
func (a *Account) Classification() AccountClassification { return a.classification }
func (a *Account) AccountType() AccountType       { return a.accountType }
func (a *Account) AccountSubType() *string       { return a.accountSubType }
func (a *Account) AccountName() string            { return a.accountName }
func (a *Account) Currency() string               { return a.currency }
func (a *Account) Status() AccountStatus          { return a.status }
func (a *Account) OpenedDate() time.Time          { return a.openedDate }
func (a *Account) ClosedDate() *time.Time         { return a.closedDate }
func (a *Account) LiquidityProfile() LiquidityProfile { return a.liquidityProfile }
func (a *Account) AccountHealth() AccountHealth   { return a.accountHealth }
func (a *Account) CreditLimit() *int64            { return a.creditLimit }
func (a *Account) InterestRate() *float64         { return a.interestRate }
func (a *Account) Country() *string               { return a.country }
func (a *Account) SourceOfTruth() *string         { return a.sourceOfTruth }
func (a *Account) Visibility() Visibility         { return a.visibility }
func (a *Account) ExternalReference() *string     { return a.externalReference }
func (a *Account) Tags() []string                 { return a.tags }
func (a *Account) Metadata() map[string]string    { return a.metadata }
func (a *Account) Notes() string                  { return a.notes }
func (a *Account) CreatedAt() time.Time           { return a.createdAt }
func (a *Account) UpdatedAt() time.Time           { return a.updatedAt }

func (a *Account) SetStatus(s AccountStatus) {
	a.status = s; a.updatedAt = time.Now().UTC()
}

func (a *Account) CanTransitionTo(target AccountStatus) error {
	m := map[AccountStatus]map[AccountStatus]bool{
		StatusDraft:   {StatusActive: true},
		StatusActive:  {StatusDormant: true, StatusFrozen: true, StatusClosed: true},
		StatusDormant: {StatusActive: true, StatusFrozen: true},
		StatusFrozen:  {StatusActive: true, StatusClosed: true},
		StatusClosed:  {StatusArchived: true},
	}
	if al, ok := m[a.status]; ok && al[target] { return nil }
	if a.status == target { return fmt.Errorf("already in %s", target) }
	if a.status == StatusArchived { return fmt.Errorf("archived accounts cannot transition") }
	return fmt.Errorf("cannot transition from %s to %s", a.status, target)
}
