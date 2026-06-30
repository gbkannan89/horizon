package domain

type AccountStatus string
const (
	StatusDraft   AccountStatus = "Draft"
	StatusActive  AccountStatus = "Active"
	StatusDormant AccountStatus = "Dormant"
	StatusFrozen  AccountStatus = "Frozen"
	StatusClosed  AccountStatus = "Closed"
	StatusArchived AccountStatus = "Archived"
)

type AccountClassification string
const (
	ClsAsset          AccountClassification = "Asset"
	ClsLiability      AccountClassification = "Liability"
	ClsEquity         AccountClassification = "Equity"
	ClsCashEquivalent AccountClassification = "CashEquivalent"
	ClsInvestment     AccountClassification = "Investment"
	ClsIncomeHolding  AccountClassification = "IncomeHolding"
	ClsExpenseHolding AccountClassification = "ExpenseHolding"
	ClsVirtual        AccountClassification = "Virtual"
	ClsOperational    AccountClassification = "Operational"
)

type AccountType string
const (
	ATCash             AccountType = "Cash"
	ATChecking         AccountType = "Checking"
	ATSavings          AccountType = "Savings"
	ATFixedDeposit     AccountType = "FixedDeposit"
	ATForeignCurrency  AccountType = "ForeignCurrency"
	ATCashWallet       AccountType = "CashWallet"
	ATCreditCard       AccountType = "CreditCard"
	ATLoan             AccountType = "Loan"
	ATMortgage         AccountType = "Mortgage"
	ATOverdraft        AccountType = "Overdraft"
	ATBrokerage        AccountType = "Brokerage"
	ATRetirement       AccountType = "Retirement"
	ATProvidentFund    AccountType = "ProvidentFund"
	ATEducationSavings AccountType = "EducationSavings"
	ATHealthSavings    AccountType = "HealthSavings"
	ATInvestmentAccount AccountType = "InvestmentAccount"
	ATGold             AccountType = "Gold"
	ATCryptoWallet     AccountType = "CryptoWallet"
	ATPropertyHolding  AccountType = "PropertyHolding"
	ATVirtualEnvelope  AccountType = "VirtualEnvelope"
	ATBusinessAccount  AccountType = "BusinessAccount"
	ATEscrow           AccountType = "Escrow"
	ATTrust            AccountType = "Trust"
)

type OwnershipModel string
const (
	OMPersonal  OwnershipModel = "Personal"
	OMJoint     OwnershipModel = "Joint"
	OMHousehold OwnershipModel = "Household"
	OMBusiness  OwnershipModel = "Business"
	OMMinor     OwnershipModel = "Minor"
	OMTrust     OwnershipModel = "Trust"
)

type Visibility string
const (
	VisPrivate   Visibility = "Private"
	VisHousehold Visibility = "Household"
	VisShared    Visibility = "Shared"
)

type LiquidityProfile string
const (
	LQImmediate   LiquidityProfile = "Immediate"
	LQShortTerm   LiquidityProfile = "ShortTerm"
	LQMediumTerm  LiquidityProfile = "MediumTerm"
	LQLongTerm    LiquidityProfile = "LongTerm"
	LQRestricted  LiquidityProfile = "Restricted"
	LQLocked      LiquidityProfile = "Locked"
	LQIlliquid    LiquidityProfile = "Illiquid"
)

type AccountHealth string
const (
	AHHealthy AccountHealth = "Healthy"
	AHWarning AccountHealth = "Warning"
	AHCritical AccountHealth = "Critical"
	AHDormant AccountHealth = "Dormant"
	AHClosed  AccountHealth = "Closed"
)

var typeClassifications = map[AccountType][]AccountClassification{
	ATCash: {ClsAsset, ClsCashEquivalent}, ATChecking: {ClsAsset, ClsCashEquivalent},
	ATSavings: {ClsAsset, ClsCashEquivalent}, ATFixedDeposit: {ClsAsset},
	ATForeignCurrency: {ClsAsset, ClsCashEquivalent}, ATCashWallet: {ClsCashEquivalent, ClsVirtual},
	ATCreditCard: {ClsLiability}, ATLoan: {ClsLiability}, ATMortgage: {ClsLiability},
	ATOverdraft: {ClsLiability}, ATBrokerage: {ClsInvestment, ClsAsset},
	ATRetirement: {ClsInvestment, ClsAsset}, ATProvidentFund: {ClsInvestment, ClsAsset},
	ATEducationSavings: {ClsInvestment, ClsAsset}, ATHealthSavings: {ClsInvestment, ClsAsset},
	ATInvestmentAccount: {ClsInvestment, ClsAsset}, ATGold: {ClsInvestment, ClsAsset},
	ATCryptoWallet: {ClsInvestment, ClsAsset}, ATPropertyHolding: {ClsAsset},
	ATVirtualEnvelope: {ClsVirtual}, ATBusinessAccount: {ClsEquity, ClsOperational},
	ATEscrow: {ClsOperational}, ATTrust: {ClsEquity},
}
