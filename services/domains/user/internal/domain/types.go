package domain

type UserStatus string

const (
	StatusRegistered UserStatus = "Registered"
	StatusActive     UserStatus = "Active"
	StatusInactive   UserStatus = "Inactive"
	StatusSuspended  UserStatus = "Suspended"
	StatusMerged     UserStatus = "Merged"
	StatusArchived   UserStatus = "Archived"
	StatusHistorical UserStatus = "Historical"
)

var AllStatuses = map[UserStatus]bool{
	StatusRegistered: true, StatusActive: true, StatusInactive: true,
	StatusSuspended: true, StatusMerged: true, StatusArchived: true, StatusHistorical: true,
}

type UserType string

const (
	TypeIndividual      UserType = "Individual"
	TypeMinor           UserType = "Minor"
	TypeGuardian        UserType = "Guardian"
	TypeAdvisor         UserType = "Advisor"
	TypeBusinessOwner   UserType = "BusinessOwner"
	TypeRetiree         UserType = "Retiree"
	TypeStudent         UserType = "Student"
	TypeDependent       UserType = "Dependent"
	TypeCustom          UserType = "Custom"
)

var AllUserTypes = map[UserType]bool{
	TypeIndividual: true, TypeMinor: true, TypeGuardian: true,
	TypeAdvisor: true, TypeBusinessOwner: true, TypeRetiree: true,
	TypeStudent: true, TypeDependent: true, TypeCustom: true,
}

type FinancialIdentityProfile string

const (
	FIPPersonal   FinancialIdentityProfile = "Personal"
	FIPHousehold  FinancialIdentityProfile = "Household"
	FIPBusiness   FinancialIdentityProfile = "Business"
	FIPInvestor   FinancialIdentityProfile = "Investor"
	FIPRetiree    FinancialIdentityProfile = "Retiree"
	FIPStudent    FinancialIdentityProfile = "Student"
	FIPGuardian   FinancialIdentityProfile = "Guardian"
	FIPAdvisor    FinancialIdentityProfile = "Advisor"
	FIPHybrid     FinancialIdentityProfile = "Hybrid"
)

var AllFIPs = map[FinancialIdentityProfile]bool{
	FIPPersonal: true, FIPHousehold: true, FIPBusiness: true,
	FIPInvestor: true, FIPRetiree: true, FIPStudent: true,
	FIPGuardian: true, FIPAdvisor: true, FIPHybrid: true,
}

type UserHealth string

const (
	HealthHealthy  UserHealth = "Healthy"
	HealthActive   UserHealth = "Active"
	HealthGrowing  UserHealth = "Growing"
	HealthAtRisk   UserHealth = "AtRisk"
	HealthInactive UserHealth = "Inactive"
	HealthHistorical UserHealth = "Historical"
)

type UserConfidence string

const (
	ConfVerified      UserConfidence = "Verified"
	ConfKYCVerified   UserConfidence = "KYCVerified"
	ConfUserVerified  UserConfidence = "UserVerified"
	ConfImported      UserConfidence = "Imported"
	ConfEstimated     UserConfidence = "Estimated"
	ConfSimulated     UserConfidence = "Simulated"
)

var AllConfidences = map[UserConfidence]bool{
	ConfVerified: true, ConfKYCVerified: true, ConfUserVerified: true,
	ConfImported: true, ConfEstimated: true, ConfSimulated: true,
}

type FinBehaviourProfile string

const (
	BehavSaver    FinBehaviourProfile = "Saver"
	BehavInvestor FinBehaviourProfile = "Investor"
	BehavPlanner  FinBehaviourProfile = "Planner"
	BehavTrader   FinBehaviourProfile = "Trader"
	BehavSpender  FinBehaviourProfile = "Spender"
	BehavBalanced FinBehaviourProfile = "Balanced"
)

type ConsentType string

const (
	ConsentDataAccess   ConsentType = "DataAccess"
	ConsentDataSharing  ConsentType = "DataSharing"
	ConsentAIProcessing ConsentType = "AIProcessing"
	ConsentAnalytics    ConsentType = "Analytics"
	ConsentMarketing    ConsentType = "Marketing"
	ConsentResearch     ConsentType = "Research"
)

var AllConsentTypes = map[ConsentType]bool{
	ConsentDataAccess: true, ConsentDataSharing: true, ConsentAIProcessing: true,
	ConsentAnalytics: true, ConsentMarketing: true, ConsentResearch: true,
}

type DataSharingLevel string

const (
	DSAllShared  DataSharingLevel = "AllShared"
	DSOptIn      DataSharingLevel = "OptIn"
	DSEntity     DataSharingLevel = "PerEntity"
	DSRoleBased  DataSharingLevel = "RoleBased"
	DSCustom     DataSharingLevel = "Custom"
)

type SOT string

const (
	SOTSystem   SOT = "System"
	SOTManual   SOT = "Manual"
	SOTImport   SOT = "Import"
	SOTMigration SOT = "Migration"
)
