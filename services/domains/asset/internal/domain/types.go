package domain

// AssetStatus represents the lifecycle state of an asset.
type AssetStatus string

const (
	StatusPlanned         AssetStatus = "Planned"
	StatusAcquired        AssetStatus = "Acquired"
	StatusActive          AssetStatus = "Active"
	StatusPartiallyDisposed AssetStatus = "PartiallyDisposed"
	StatusFullyDisposed   AssetStatus = "FullyDisposed"
	StatusArchived        AssetStatus = "Archived"
	StatusHistorical      AssetStatus = "Historical"
)

// AssetClassification represents the economic category of an asset.
type AssetClassification string

const (
	ClsCash             AssetClassification = "Cash"
	ClsCashEquivalent   AssetClassification = "CashEquivalent"
	ClsEquity           AssetClassification = "Equity"
	ClsMutualFund       AssetClassification = "MutualFund"
	ClsETF              AssetClassification = "ETF"
	ClsBond             AssetClassification = "Bond"
	ClsRealEstate       AssetClassification = "RealEstate"
	ClsGold             AssetClassification = "Gold"
	ClsSilver           AssetClassification = "Silver"
	ClsCommodity        AssetClassification = "Commodity"
	ClsVehicle          AssetClassification = "Vehicle"
	ClsBusinessOwnership AssetClassification = "BusinessOwnership"
	ClsRetirementAsset  AssetClassification = "RetirementAsset"
	ClsInsuranceCashValue AssetClassification = "InsuranceCashValue"
	ClsCrypto           AssetClassification = "Crypto"
	ClsCollectible      AssetClassification = "Collectible"
	ClsDigitalAsset     AssetClassification = "DigitalAsset"
	ClsReceivable       AssetClassification = "Receivable"
	ClsOther            AssetClassification = "Other"
)

// ValuationProfile represents the approach category for valuation.
type ValuationProfile string

const (
	VPProfileMarketDriven      ValuationProfile = "MarketDriven"
	VPProfileRuleBased         ValuationProfile = "RuleBased"
	VPProfileExternalAppraisal ValuationProfile = "ExternalAppraisal"
	VPProfileGovernmentAssessment ValuationProfile = "GovernmentAssessment"
	VPProfileManualAssessment  ValuationProfile = "ManualAssessment"
	VPProfileIndexed           ValuationProfile = "Indexed"
	VPProfileEstimated         ValuationProfile = "Estimated"
	VPProfileHybrid            ValuationProfile = "Hybrid"
)

// ValuationMethod represents the specific calculation method for valuation.
type ValuationMethod string

const (
	VMMarketValue            ValuationMethod = "MarketValue"
	VMBookValue              ValuationMethod = "BookValue"
	VMFairValue              ValuationMethod = "FairValue"
	VMCostBasis              ValuationMethod = "CostBasis"
	VMReplacementValue       ValuationMethod = "ReplacementValue"
	VMGovernmentGuidanceValue ValuationMethod = "GovernmentGuidanceValue"
	VMManualValue            ValuationMethod = "ManualValue"
	VMEstimatedValue         ValuationMethod = "EstimatedValue"
	VMIndexedValue           ValuationMethod = "IndexedValue"
)

// OwnershipModel represents the ownership structure of an asset.
type OwnershipModel string

const (
	OMIndividual OwnershipModel = "Individual"
	OMJoint      OwnershipModel = "Joint"
	OMHousehold  OwnershipModel = "Household"
	OMBusiness   OwnershipModel = "Business"
	OMTrust      OwnershipModel = "Trust"
	OMMinor      OwnershipModel = "Minor"
	OMFractional OwnershipModel = "Fractional"
	OMNominee    OwnershipModel = "Nominee"
)

// AssetHealth represents the health assessment of an asset.
type AssetHealth string

const (
	AHHealthy     AssetHealth = "Healthy"
	AHAppreciating AssetHealth = "Appreciating"
	AHStable      AssetHealth = "Stable"
	AHDDepreciating AssetHealth = "Depreciating"
	AHIlliquid    AssetHealth = "Illiquid"
	AHImpaired    AssetHealth = "Impaired"
	AHHistorical  AssetHealth = "Historical"
)

// AssetConfidence represents the data reliability level.
type AssetConfidence string

const (
	ACVerified              AssetConfidence = "Verified"
	ACMarketVerified        AssetConfidence = "MarketVerified"
	ACInstitutionVerified   AssetConfidence = "InstitutionVerified"
	ACUserVerified          AssetConfidence = "UserVerified"
	ACImported              AssetConfidence = "Imported"
	ACEstimated             AssetConfidence = "Estimated"
	ACSimulated             AssetConfidence = "Simulated"
)

// LiquidityProfile represents how quickly the asset can be converted to cash.
type LiquidityProfile string

const (
	LQImmediate  LiquidityProfile = "Immediate"
	LQShortTerm  LiquidityProfile = "ShortTerm"
	LQMediumTerm LiquidityProfile = "MediumTerm"
	LQLongTerm   LiquidityProfile = "LongTerm"
	LQRestricted LiquidityProfile = "Restricted"
	LQLocked     LiquidityProfile = "Locked"
	LQIlliquid   LiquidityProfile = "Illiquid"
)

// ValuationRecord represents a single point-in-time valuation snapshot.
type ValuationRecord struct {
	AssetID        string          `json:"asset_id"`
	Value          int64           `json:"value"`
	Method         ValuationMethod `json:"method"`
	ValuationDate  string          `json:"valuation_date"`
	SourceOfTruth  string          `json:"source_of_truth"`
	RecordedAt     string          `json:"recorded_at"`
}
