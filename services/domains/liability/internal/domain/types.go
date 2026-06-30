package domain

// LiabilityStatus represents the lifecycle state of a liability.
type LiabilityStatus string
const (
	LSPlanned      LiabilityStatus = "Planned"
	LSApproved     LiabilityStatus = "Approved"
	LSActive       LiabilityStatus = "Active"
	LSGracePeriod  LiabilityStatus = "GracePeriod"
	LSDelinquent   LiabilityStatus = "Delinquent"
	LSRestructured LiabilityStatus = "Restructured"
	LSSettled      LiabilityStatus = "Settled"
	LSWrittenOff   LiabilityStatus = "WrittenOff"
	LSHistorical   LiabilityStatus = "Historical"
	LSArchived     LiabilityStatus = "Archived"
)

// LiabilityClassification represents the category of liability.
type LiabilityClassification string
const (
	LCMortgage       LiabilityClassification = "Mortgage"
	LCHomeLoan       LiabilityClassification = "HomeLoan"
	LCPersonalLoan   LiabilityClassification = "PersonalLoan"
	LCVehicleLoan    LiabilityClassification = "VehicleLoan"
	LCEducationLoan  LiabilityClassification = "EducationLoan"
	LCCreditCard     LiabilityClassification = "CreditCard"
	LCLineOfCredit   LiabilityClassification = "LineOfCredit"
	LCBusinessLoan   LiabilityClassification = "BusinessLoan"
	LCGoldLoan       LiabilityClassification = "GoldLoan"
	LCMarginLoan     LiabilityClassification = "MarginLoan"
	LCTaxLiability   LiabilityClassification = "TaxLiability"
	LCFamilyLoan     LiabilityClassification = "FamilyLoan"
	LCEmployerLoan   LiabilityClassification = "EmployerLoan"
	LCBuyNowPayLater LiabilityClassification = "BuyNowPayLater"
	LCOverdraft      LiabilityClassification = "Overdraft"
	LCOther          LiabilityClassification = "Other"
)

// RepaymentProfile represents the repayment structure.
type RepaymentProfile string
const (
	RPFixedInstallment    RepaymentProfile = "FixedInstallment"
	RPVariableInstallment RepaymentProfile = "VariableInstallment"
	RPInterestOnly        RepaymentProfile = "InterestOnly"
	RPBalloon             RepaymentProfile = "Balloon"
	RPRevolving           RepaymentProfile = "Revolving"
	RPFlexible            RepaymentProfile = "Flexible"
	RPBullet              RepaymentProfile = "Bullet"
	RPHybrid              RepaymentProfile = "Hybrid"
)

// RepaymentModel represents the calculation model for payments.
type RepaymentModel string
const (
	RMFixedEMI        RepaymentModel = "FixedEMI"
	RMReducingBalance  RepaymentModel = "ReducingBalance"
	RMInterestOnly     RepaymentModel = "InterestOnly"
	RMBalloonPayment   RepaymentModel = "BalloonPayment"
	RMBulletRepayment  RepaymentModel = "BulletRepayment"
	RMFlexible         RepaymentModel = "Flexible"
	RMRevolvingCredit  RepaymentModel = "RevolvingCredit"
	RMMinimumDue       RepaymentModel = "MinimumDue"
	RMCustom           RepaymentModel = "Custom"
)

// InterestModel represents the interest calculation approach.
type InterestModel string
const (
	IMSimple        InterestModel = "Simple"
	IMCompound      InterestModel = "Compound"
	IMReducingBal   InterestModel = "ReducingBalance"
	IMFixedRate     InterestModel = "FixedRate"
	IMFloatingRate  InterestModel = "FloatingRate"
)

// LiabilityHealth represents the health assessment.
type LiabilityHealth string
const (
	LHHealthy        LiabilityHealth = "Healthy"
	LHStable         LiabilityHealth = "Stable"
	LHWarning        LiabilityHealth = "Warning"
	LHCritical       LiabilityHealth = "Critical"
	LHDelinquentSt   LiabilityHealth = "Delinquent"
	LHSettled        LiabilityHealth = "Settled"
	LHHistorical     LiabilityHealth = "Historical"
)

// LiabilityConfidence represents data reliability.
type LiabilityConfidence string
const (
	LCVerifiedSt     LiabilityConfidence = "Verified"
	LCInstitutionSt  LiabilityConfidence = "InstitutionVerified"
	LCUserVerifiedSt LiabilityConfidence = "UserVerified"
	LCImportedSt     LiabilityConfidence = "Imported"
	LCEstimatedSt    LiabilityConfidence = "Estimated"
	LCSimulatedSt    LiabilityConfidence = "Simulated"
)
