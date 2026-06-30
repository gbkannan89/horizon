package domain

type EventState string

const (
	StateDraft     EventState = "DRAFT"
	StatePending   EventState = "PENDING"
	StateConfirmed EventState = "CONFIRMED"
	StatePosted    EventState = "POSTED"
	StateCancelled EventState = "CANCELLED"
	StateReversed  EventState = "REVERSED"
	StateArchived  EventState = "ARCHIVED"
)

func (s EventState) Valid() bool {
	switch s {
	case StateDraft, StatePending, StateConfirmed, StatePosted, StateCancelled, StateReversed, StateArchived:
		return true
	}
	return false
}

type EventType string

const (
	EventSalary               EventType = "Salary"
	EventFreelance            EventType = "Freelance"
	EventBusinessIncome       EventType = "BusinessIncome"
	EventPassiveIncome        EventType = "PassiveIncome"
	EventGift                 EventType = "Gift"
	EventRefund               EventType = "Refund"
	EventMiscIncome           EventType = "MiscellaneousIncome"
	EventPurchase             EventType = "Purchase"
	EventBill                 EventType = "Bill"
	EventSubscription         EventType = "Subscription"
	EventRent                 EventType = "Rent"
	EventHealthcare           EventType = "Healthcare"
	EventEducation            EventType = "Education"
	EventCharity              EventType = "Charity"
	EventMiscExpense          EventType = "MiscellaneousExpense"
	EventTransfer             EventType = "Transfer"
	EventInternalTransfer     EventType = "InternalTransfer"
	EventInvestmentPurchase   EventType = "InvestmentPurchase"
	EventInvestmentSale       EventType = "InvestmentSale"
	EventDividend             EventType = "Dividend"
	EventInterest             EventType = "Interest"
	EventCapitalGain          EventType = "CapitalGain"
	EventInvestmentFee        EventType = "InvestmentFee"
	EventLoanDisbursement     EventType = "LoanDisbursement"
	EventLoanPayment          EventType = "LoanPayment"
	EventLoanPrepayment       EventType = "LoanPrepayment"
	EventInterestCharge       EventType = "InterestCharge"
	EventLiabilityFee         EventType = "LiabilityFee"
	EventLiabilityClosed      EventType = "LiabilityClosed"
	EventAssetAcquisition     EventType = "AssetAcquisition"
	EventAssetDisposal        EventType = "AssetDisposal"
	EventAssetAppreciation    EventType = "AssetAppreciation"
	EventAssetDepreciation    EventType = "AssetDepreciation"
	EventTaxPayment           EventType = "TaxPayment"
	EventTaxRefund            EventType = "TaxRefund"
	EventTaxWithholding       EventType = "TaxWithholding"
	EventManualCorrection     EventType = "ManualCorrection"
	EventReconciliation       EventType = "Reconciliation"
	EventImportedEvent        EventType = "ImportedEvent"
	EventScheduledEvent       EventType = "ScheduledEvent"
	EventSystemAdjustment     EventType = "SystemAdjustment"
)

var AllEventTypes = map[EventType]bool{}

func init() {
	types := []EventType{
		EventSalary, EventFreelance, EventBusinessIncome,
		EventPassiveIncome, EventGift, EventRefund, EventMiscIncome,
		EventPurchase, EventBill, EventSubscription, EventRent,
		EventHealthcare, EventEducation, EventCharity, EventMiscExpense,
		EventTransfer, EventInternalTransfer,
		EventInvestmentPurchase, EventInvestmentSale, EventDividend,
		EventInterest, EventCapitalGain, EventInvestmentFee,
		EventLoanDisbursement, EventLoanPayment, EventLoanPrepayment,
		EventInterestCharge, EventLiabilityFee, EventLiabilityClosed,
		EventAssetAcquisition, EventAssetDisposal,
		EventAssetAppreciation, EventAssetDepreciation,
		EventTaxPayment, EventTaxRefund, EventTaxWithholding,
		EventManualCorrection, EventReconciliation,
		EventImportedEvent, EventScheduledEvent, EventSystemAdjustment,
	}
	for _, t := range types {
		AllEventTypes[t] = true
	}
}

func ValidEventType(s string) bool {
	_, ok := AllEventTypes[EventType(s)]
	return ok
}

type EventOrigin string

const (
	OriginUser       EventOrigin = "User"
	OriginBank       EventOrigin = "Bank"
	OriginBroker     EventOrigin = "Broker"
	OriginEmployer   EventOrigin = "Employer"
	OriginGovernment EventOrigin = "Government"
	OriginSystem     EventOrigin = "System"
	OriginImport     EventOrigin = "Import"
	OriginAPI        EventOrigin = "API"
	OriginSimulation EventOrigin = "Simulation"
)

var AllOrigins = map[EventOrigin]bool{
	OriginUser: true, OriginBank: true, OriginBroker: true,
	OriginEmployer: true, OriginGovernment: true, OriginSystem: true,
	OriginImport: true, OriginAPI: true, OriginSimulation: true,
}

type EventConfidence string

const (
	ConfidenceConfirmed  EventConfidence = "Confirmed"
	ConfidenceImported   EventConfidence = "Imported"
	ConfidenceEstimated  EventConfidence = "Estimated"
	ConfidenceSimulated  EventConfidence = "Simulated"
	ConfidenceInferred   EventConfidence = "Inferred"
)

var confidenceLevels = map[EventConfidence]int{
	ConfidenceConfirmed: 5, ConfidenceImported: 4, ConfidenceEstimated: 3,
	ConfidenceSimulated: 2, ConfidenceInferred: 1,
}

type SourceOfTruth string

const (
	SOTBankFeed         SourceOfTruth = "BankFeed"
	SOTBrokerStatement  SourceOfTruth = "BrokerStatement"
	SOTManualEntry      SourceOfTruth = "ManualEntry"
	SOTPayroll          SourceOfTruth = "Payroll"
	SOTOCR              SourceOfTruth = "OCR"
	SOTCSVImport        SourceOfTruth = "CSVImport"
	SOTAPI              SourceOfTruth = "API"
	SOTGovernmentRecord SourceOfTruth = "GovernmentRecord"
)

type CreatedBy string

const (
	CreatedByUser     CreatedBy = "User"
	CreatedBySystem   CreatedBy = "System"
	CreatedByImport   CreatedBy = "Import"
	CreatedByRecurring CreatedBy = "Recurring"
	CreatedByAPI      CreatedBy = "API"
)

var AllCreatedBy = map[CreatedBy]bool{
	CreatedByUser: true, CreatedBySystem: true,
	CreatedByImport: true, CreatedByRecurring: true, CreatedByAPI: true,
}

func (c EventConfidence) CanUpgradeTo(new EventConfidence) bool {
	return confidenceLevels[new] > confidenceLevels[c]
}

func (c EventConfidence) IsConfirmed() bool {
	return c == ConfidenceConfirmed
}
