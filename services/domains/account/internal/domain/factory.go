package domain

import (
	"time"
)

type AccountFactory struct{}

func NewAccountFactory() *AccountFactory { return &AccountFactory{} }

func (f *AccountFactory) Create(
	id, ownerID, name, currency string,
	acctType AccountType,
	cls AccountClassification,
	openedDate time.Time,
	visibility Visibility,
	ownership OwnershipModel,
	creditLimit *int64,
	interestRate *float64,
	institutionID, householdID, subType, country, sourceOfTruth, extRef string,
	tags []string,
	metadata map[string]string,
	notes string,
) (*Account, error) {

	if err := ValidateAccountName(name); err != nil { return nil, err }
	if err := ValidateAccountType(string(acctType)); err != nil { return nil, err }
	if err := ValidateCurrency(currency); err != nil { return nil, err }
	if err := ValidateClassificationConsistency(acctType, cls); err != nil { return nil, err }
	if err := ValidateCreditLimit(acctType, creditLimit); err != nil { return nil, err }
	if err := ValidateInterestRate(interestRate); err != nil { return nil, err }

	var instID, hhID, subTyp, cntry, sot, ref *string
	if institutionID != "" { instID = &institutionID }
	if householdID != "" { hhID = &householdID }
	if subType != "" { subTyp = &subType }
	if country != "" { cntry = &country }
	if sourceOfTruth != "" { sot = &sourceOfTruth }
	if extRef != "" { ref = &extRef }
	if tags == nil { tags = []string{} }
	if metadata == nil { metadata = map[string]string{} }

	now := time.Now().UTC()
	return &Account{
		id: id, institutionID: instID, ownerID: ownerID, householdID: hhID,
		classification: cls, accountType: acctType, accountSubType: subTyp,
		accountName: name, currency: currency, status: StatusDraft,
		openedDate: openedDate, liquidityProfile: deriveLiquidity(acctType),
		accountHealth: AHHealthy, creditLimit: creditLimit, interestRate: interestRate,
		country: cntry, sourceOfTruth: sot, visibility: visibility,
		externalReference: ref, tags: tags, metadata: metadata, notes: notes,
		createdAt: now, updatedAt: now,
	}, nil
}

func deriveLiquidity(t AccountType) LiquidityProfile {
	switch t {
	case ATCash, ATChecking, ATSavings, ATCashWallet:
		return LQImmediate
	case ATForeignCurrency, ATOverdraft:
		return LQShortTerm
	case ATFixedDeposit, ATGold, ATCryptoWallet:
		return LQMediumTerm
	case ATBrokerage, ATInvestmentAccount:
		return LQMediumTerm
	case ATRetirement, ATProvidentFund, ATEducationSavings, ATHealthSavings:
		return LQLongTerm
	case ATPropertyHolding, ATTrust:
		return LQIlliquid
	case ATVirtualEnvelope, ATEscrow, ATBusinessAccount:
		return LQRestricted
	default:
		return LQMediumTerm
	}
}
