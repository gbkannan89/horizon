package domain

import "time"

type PortfolioFactory struct{}

func NewPortfolioFactory() *PortfolioFactory { return &PortfolioFactory{} }

func (f *PortfolioFactory) Create(
	id, name string,
	pt PortfolioType,
	ownerID, baseCurrency string,
	membershipModel MembershipModel,
	riskProfile RiskProfile,
	benchmarkProfile BenchmarkProfile,
	benchmark, sourceOfTruth string,
	tags []string, notes string,
) (*Portfolio, error) {

	if err := ValidatePortfolioName(name); err != nil { return nil, err }
	if err := ValidateCurrency(baseCurrency); err != nil { return nil, err }
	if tags == nil { tags = []string{} }
	if sourceOfTruth == "" { sourceOfTruth = "User" }

	now := time.Now().UTC()
	return &Portfolio{
		id: id, name: name, portfolioType: pt, ownerID: ownerID,
		baseCurrency: baseCurrency, membershipModel: membershipModel,
		status: PSDraft, health: PHStable, confidence: PCUserVerifiedSt,
		riskProfile: riskProfile, benchmarkProfile: benchmarkProfile,
		benchmark: benchmark, sourceOfTruth: sourceOfTruth, tags: tags,
		metadata: map[string]string{}, notes: notes, members: []PortfolioMember{},
		createdAt: now, updatedAt: now,
	}, nil
}
