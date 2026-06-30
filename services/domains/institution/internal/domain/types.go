package domain

type InstitutionStatus string
const (
	StDraft      InstitutionStatus = "Draft"
	StRegistered InstitutionStatus = "Registered"
	StVerified   InstitutionStatus = "Verified"
	StActive     InstitutionStatus = "Active"
	StSuspended  InstitutionStatus = "Suspended"
	StMerged     InstitutionStatus = "Merged"
	StClosed     InstitutionStatus = "Closed"
	StArchived   InstitutionStatus = "Archived"
)

type InstitutionType string
const (
	ITCommercialBank    InstitutionType = "CommercialBank"
	ITInvestmentBank    InstitutionType = "InvestmentBank"
	ITBroker            InstitutionType = "Broker"
	ITInsuranceCompany  InstitutionType = "InsuranceCompany"
	ITMutualFundHouse   InstitutionType = "MutualFundHouse"
	ITPensionProvider   InstitutionType = "PensionProvider"
	ITGovernment        InstitutionType = "Government"
	ITEmployer          InstitutionType = "Employer"
	ITPayrollProvider   InstitutionType = "PayrollProvider"
	ITCreditUnion       InstitutionType = "CreditUnion"
	ITCryptoExchange    InstitutionType = "CryptoExchange"
	ITLoanProvider      InstitutionType = "LoanProvider"
	ITDigitalWallet     InstitutionType = "DigitalWallet"
	ITFinTech           InstitutionType = "FinTech"
	ITOther             InstitutionType = "Other"
)

var typeCategory = map[InstitutionType]string{
	ITCommercialBank: "Banking", ITInvestmentBank: "Banking", ITBroker: "Investment",
	ITInsuranceCompany: "Insurance", ITMutualFundHouse: "Investment", ITPensionProvider: "Investment",
	ITGovernment: "Government", ITEmployer: "Employment", ITPayrollProvider: "Employment",
	ITCreditUnion: "Banking", ITCryptoExchange: "FinTech", ITLoanProvider: "FinTech",
	ITDigitalWallet: "FinTech", ITFinTech: "FinTech", ITOther: "FinTech",
}

type TrustLevel string
const (TLVerified TrustLevel = "Verified"; TLRegistered TrustLevel = "Registered"; TLUnverified TrustLevel = "Unverified"; TLUntrusted TrustLevel = "Untrusted")

type InstitutionHealth string
const (IHHealthy InstitutionHealth = "Healthy"; IHStable InstitutionHealth = "Stable"; IHWarning InstitutionHealth = "Warning"; IHRestricted InstitutionHealth = "Restricted"; IHSuspended InstitutionHealth = "Suspended"; IHHistorical InstitutionHealth = "Historical")

type InstitutionConfidence string
const (ICVerified InstitutionConfidence = "Verified"; ICTrusted InstitutionConfidence = "Trusted"; ICRegistered InstitutionConfidence = "Registered"; ICImported InstitutionConfidence = "Imported"; ICUserDefined InstitutionConfidence = "UserDefined"; ICUnknown InstitutionConfidence = "Unknown")

type ConnectivityCapability string
const (CCManual ConnectivityCapability = "Manual"; CCCSVImport ConnectivityCapability = "CSVImport"; CCOCRImport ConnectivityCapability = "OCRImport"; CCOpenBanking ConnectivityCapability = "OpenBanking"; CCBrokerAPI ConnectivityCapability = "BrokerAPI"; CCPayrollFeed ConnectivityCapability = "PayrollFeed"; CCGovernmentFeed ConnectivityCapability = "GovernmentFeed"; CCFileImport ConnectivityCapability = "FileImport"; CCFuture ConnectivityCapability = "FutureConnector")

type FinancialProduct string
const (
	FPChecking      FinancialProduct = "Checking"
	FPSavings       FinancialProduct = "Savings"
	FPCreditCard    FinancialProduct = "CreditCard"
	FPLoan          FinancialProduct = "Loan"
	FPMortgage      FinancialProduct = "Mortgage"
	FPBrokerage     FinancialProduct = "Brokerage"
	FPRetirement    FinancialProduct = "Retirement"
	FPInsurance     FinancialProduct = "Insurance"
	FPFixedDeposit  FinancialProduct = "FixedDeposit"
	FPForeignCurr   FinancialProduct = "ForeignCurrency"
	FPMutualFund    FinancialProduct = "MutualFund"
	FPExchange      FinancialProduct = "CryptoExchange"
	FPWallet        FinancialProduct = "DigitalWallet"
	FPPayroll       FinancialProduct = "Payroll"
	FPGovernment    FinancialProduct = "GovernmentDisbursement"
	FPOther         FinancialProduct = "Other"
)

var productMap = map[InstitutionType][]FinancialProduct{
	ITCommercialBank: {FPChecking, FPSavings, FPCreditCard, FPLoan, FPMortgage, FPFixedDeposit, FPForeignCurr},
	ITInvestmentBank: {FPBrokerage, FPRetirement, FPMutualFund},
	ITBroker: {FPBrokerage, FPRetirement, FPMutualFund},
	ITInsuranceCompany: {FPInsurance},
	ITMutualFundHouse: {FPMutualFund},
	ITPensionProvider: {FPRetirement},
	ITGovernment: {FPGovernment},
	ITEmployer: {FPPayroll},
	ITPayrollProvider: {FPPayroll},
	ITCreditUnion: {FPChecking, FPSavings, FPLoan, FPMortgage},
	ITCryptoExchange: {FPExchange},
	ITLoanProvider: {FPLoan},
	ITDigitalWallet: {FPWallet, FPForeignCurr},
	ITFinTech: {FPOther},
	ITOther: {FPOther},
}
