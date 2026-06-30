package domain

// PortfolioStatus represents the lifecycle state.
type PortfolioStatus string
const (
	PSDraft     PortfolioStatus = "Draft"
	PSActive    PortfolioStatus = "Active"
	PSRebalanced PortfolioStatus = "Rebalanced"
	PSArchived  PortfolioStatus = "Archived"
	PSHistorical PortfolioStatus = "Historical"
)

// PortfolioType represents the kind of portfolio.
type PortfolioType string
const (
	PTNetWorth       PortfolioType = "NetWorth"
	PTInvestment     PortfolioType = "Investment"
	PTRetirement     PortfolioType = "Retirement"
	PTEmergencyFund  PortfolioType = "EmergencyFund"
	PTRealEstate     PortfolioType = "RealEstate"
	PTBusiness       PortfolioType = "Business"
	PTEducation      PortfolioType = "Education"
	PTDebt           PortfolioType = "Debt"
	PTInsurance      PortfolioType = "Insurance"
	PTHousehold      PortfolioType = "Household"
	PTCustom         PortfolioType = "Custom"
)

// MembershipModel represents how portfolio members are resolved.
type MembershipModel string
const (
	MMManual         MembershipModel = "Manual"
	MMRuleBased      MembershipModel = "RuleBased"
	MMDynamic        MembershipModel = "Dynamic"
	MMHousehold      MembershipModel = "Household"
	MMGoalBased      MembershipModel = "GoalBased"
	MMInstitutionBase MembershipModel = "InstitutionBased"
	MMAssetClassBase MembershipModel = "AssetClassBased"
	MMCustomMembership MembershipModel = "Custom"
)

// PortfolioHealth represents the health assessment.
type PortfolioHealth string
const (
	PHHealthy      PortfolioHealth = "Healthy"
	PHStable       PortfolioHealth = "Stable"
	PHWarning      PortfolioHealth = "Warning"
	PHCritical     PortfolioHealth = "Critical"
	PHDiversified  PortfolioHealth = "Diversified"
	PHConcentrated PortfolioHealth = "Concentrated"
	PHHistoricalSt PortfolioHealth = "Historical"
)

// PortfolioConfidence represents data reliability.
type PortfolioConfidence string
const (
	PCVerifiedSt    PortfolioConfidence = "Verified"
	PCMarketVerified PortfolioConfidence = "MarketVerified"
	PCUserVerifiedSt PortfolioConfidence = "UserVerified"
	PCImportedSt    PortfolioConfidence = "Imported"
	PCEstimatedSt   PortfolioConfidence = "Estimated"
	PCSimulatedSt   PortfolioConfidence = "Simulated"
)

// RiskProfile represents the risk level.
type RiskProfile string
const (
	RPConservative        RiskProfile = "Conservative"
	RPModeratelyConservative RiskProfile = "ModeratelyConservative"
	RPBalanced            RiskProfile = "Balanced"
	RPModeratelyAggressive RiskProfile = "ModeratelyAggressive"
	RPAggressive          RiskProfile = "Aggressive"
)

// BenchmarkProfile represents the benchmark type.
type BenchmarkProfile string
const (
	BPIndex    BenchmarkProfile = "Index"
	BPPeerGroup BenchmarkProfile = "PeerGroup"
	BPCustomBm  BenchmarkProfile = "Custom"
	BPInflation BenchmarkProfile = "Inflation"
	BPTarget    BenchmarkProfile = "Target"
	BPNone      BenchmarkProfile = "None"
)

// PortfolioMember represents a single entity in the portfolio.
type PortfolioMember struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"` // "asset", "liability", "account"
	Weight     float64 `json:"weight,omitempty"`
}
