package domain

import (
	"fmt"
	"time"
)

// Portfolio represents an analytical aggregation of financial entities.
type Portfolio struct {
	id                 string
	name               string
	portfolioType      PortfolioType
	ownerID            string
	householdID        *string
	baseCurrency       string
	membershipModel    MembershipModel
	members            []PortfolioMember
	status             PortfolioStatus
	currentValue       int64
	costBasis          int64
	health             PortfolioHealth
	confidence         PortfolioConfidence
	riskProfile        RiskProfile
	benchmarkProfile   BenchmarkProfile
	benchmark          string
	liquidityProfile   string
	sourceOfTruth      string
	tags               []string
	metadata           map[string]string
	notes              string
	createdAt          time.Time
	updatedAt          time.Time
}

func (p *Portfolio) ID() string                     { return p.id }
func (p *Portfolio) Name() string                   { return p.name }
func (p *Portfolio) PortfolioType() PortfolioType    { return p.portfolioType }
func (p *Portfolio) OwnerID() string                { return p.ownerID }
func (p *Portfolio) BaseCurrency() string            { return p.baseCurrency }
func (p *Portfolio) MembershipModel() MembershipModel { return p.membershipModel }
func (p *Portfolio) Members() []PortfolioMember      { return p.members }
func (p *Portfolio) Status() PortfolioStatus         { return p.status }
func (p *Portfolio) CurrentValue() int64             { return p.currentValue }
func (p *Portfolio) CostBasis() int64               { return p.costBasis }
func (p *Portfolio) Health() PortfolioHealth         { return p.health }
func (p *Portfolio) Confidence() PortfolioConfidence { return p.confidence }
func (p *Portfolio) RiskProfile() RiskProfile        { return p.riskProfile }
func (p *Portfolio) BenchmarkProfile() BenchmarkProfile { return p.benchmarkProfile }
func (p *Portfolio) CreatedAt() time.Time            { return p.createdAt }
func (p *Portfolio) UpdatedAt() time.Time            { return p.updatedAt }

func (p *Portfolio) SetStatus(s PortfolioStatus)  { p.status = s; p.updatedAt = time.Now().UTC() }
func (p *Portfolio) SetCurrentValue(v int64)       { p.currentValue = v; p.updatedAt = time.Now().UTC() }
func (p *Portfolio) SetMembers(m []PortfolioMember) { p.members = m; p.updatedAt = time.Now().UTC() }

func (p *Portfolio) CanTransitionTo(target PortfolioStatus) error {
	m := map[PortfolioStatus]map[PortfolioStatus]bool{
		PSDraft:     {PSActive: true},
		PSActive:    {PSRebalanced: true, PSArchived: true},
		PSRebalanced: {PSActive: true},
		PSArchived:  {PSHistorical: true},
	}
	if al, ok := m[p.status]; ok && al[target] { return nil }
	if p.status == PSArchived { return fmt.Errorf("archived portfolios cannot be reactivated") }
	return fmt.Errorf("cannot transition from %s to %s", p.status, target)
}

func ReconstructFromDB(id, name string, pt PortfolioType, ownerID, baseCurrency string,
	mm MembershipModel, status PortfolioStatus, ph PortfolioHealth, pc PortfolioConfidence,
	rp RiskProfile, tags []string, createdAt time.Time) *Portfolio {
	return &Portfolio{
		id: id, name: name, portfolioType: pt, ownerID: ownerID, baseCurrency: baseCurrency,
		membershipModel: mm, status: status, health: ph, confidence: pc,
		riskProfile: rp, tags: tags, createdAt: createdAt,
	}
}