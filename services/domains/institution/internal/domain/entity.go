package domain

import (
	"fmt"
	"time"
)

type Institution struct {
	id              string
	name            string
	instType        InstitutionType
	category        string
	status          InstitutionStatus
	country         string
	trustLevel      TrustLevel
	health          InstitutionHealth
	confidence      InstitutionConfidence
	products        []FinancialProduct
	connectivity    []ConnectivityCapability
	website         string
	phone           string
	email           string
	address         string
	headquarters    string
	regulator       string
	regLicense      string
	metadata        map[string]string
	tags            []string
	notes           string
	sourceOfTruth   string
	createdAt       time.Time
	updatedAt       time.Time
}

func (i *Institution) ID() string { return i.id }
func (i *Institution) Name() string { return i.name }
func (i *Institution) Type() InstitutionType { return i.instType }
func (i *Institution) Category() string { return i.category }
func (i *Institution) Status() InstitutionStatus { return i.status }
func (i *Institution) Country() string { return i.country }
func (i *Institution) TrustLevel() TrustLevel { return i.trustLevel }
func (i *Institution) Health() InstitutionHealth { return i.health }
func (i *Institution) Confidence() InstitutionConfidence { return i.confidence }
func (i *Institution) Products() []FinancialProduct { return i.products }
func (i *Institution) Connectivity() []ConnectivityCapability { return i.connectivity }
func (i *Institution) CreatedAt() time.Time { return i.createdAt }
func (i *Institution) UpdatedAt() time.Time { return i.updatedAt }

func (i *Institution) SetStatus(s InstitutionStatus) { i.status = s; i.updatedAt = time.Now().UTC() }

func (i *Institution) CanTransitionTo(target InstitutionStatus) error {
	m := map[InstitutionStatus]map[InstitutionStatus]bool{
		StDraft: {StRegistered: true}, StRegistered: {StVerified: true, StMerged: true},
		StVerified: {StActive: true, StSuspended: true, StClosed: true},
		StActive: {StSuspended: true, StClosed: true},
		StSuspended: {StActive: true, StClosed: true},
		StMerged: {StArchived: true},
		StClosed: {StArchived: true},
	}
	if al, ok := m[i.status]; ok && al[target] { return nil }
	return fmt.Errorf("cannot transition from %s to %s", i.status, target)
}

func ReconstructFromDB(id, name string, instType InstitutionType, category, status, country string,
	tl TrustLevel, h InstitutionHealth, c InstitutionConfidence, tags []string, createdAt time.Time) *Institution {
	return &Institution{
		id: id, name: name, instType: instType, category: category, status: InstitutionStatus(status),
		country: country, trustLevel: tl, health: h, confidence: c, tags: tags, createdAt: createdAt,
	}
}