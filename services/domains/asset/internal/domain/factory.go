package domain

import (
	"time"
)

// AssetFactory creates new Asset entities.
type AssetFactory struct{}

// NewAssetFactory returns a new AssetFactory.
func NewAssetFactory() *AssetFactory { return &AssetFactory{} }

// Create creates a new asset in Planned state.
func (f *AssetFactory) Create(
	id, assetName string,
	classification AssetClassification,
	ownershipModel OwnershipModel,
	ownerID, currency string,
	valuationProfile ValuationProfile,
	valuationMethod ValuationMethod,
	ownershipPercentage float64,
	sourceOfTruth string,
	assetType, householdID, institutionID, accountID string,
	quantity *float64,
	externalIDs map[string]string,
	tags []string,
	notes string,
) (*Asset, error) {

	if err := ValidateAssetName(assetName); err != nil {
		return nil, err
	}
	if err := ValidateClassification(string(classification)); err != nil {
		return nil, err
	}
	if err := ValidateOwnershipPercentage(ownershipPercentage); err != nil {
		return nil, err
	}
	if err := ValidateQuantity(quantity); err != nil {
		return nil, err
	}
	if err := ValidateCurrency(currency); err != nil {
		return nil, err
	}
	if sourceOfTruth == "" {
		sourceOfTruth = "Manual"
	}
	if tags == nil {
		tags = []string{}
	}
	if externalIDs == nil {
		externalIDs = map[string]string{}
	}

	now := time.Now().UTC()
	return &Asset{
		id: id, assetName: assetName, classification: classification,
		ownershipModel: ownershipModel, ownerID: ownerID, currency: currency,
		valuationProfile: valuationProfile, valuationMethod: valuationMethod,
		ownershipPercentage: ownershipPercentage, sourceOfTruth: sourceOfTruth,
		assetHealth: AHStable, assetConfidence: ACUserVerified,
		liquidityProfile: deriveAssetLiquidity(classification),
		status: StatusPlanned, tags: tags, metadata: map[string]string{},
		externalIdentifiers: externalIDs, notes: notes, createdAt: now, updatedAt: now,
	}, nil
}

func deriveAssetLiquidity(c AssetClassification) LiquidityProfile {
	switch c {
	case ClsCash, ClsCashEquivalent:
		return LQImmediate
	case ClsEquity, ClsETF, ClsMutualFund, ClsGold, ClsSilver:
		return LQShortTerm
	case ClsBond, ClsCommodity, ClsCrypto, ClsDigitalAsset:
		return LQMediumTerm
	case ClsRealEstate, ClsBusinessOwnership:
		return LQIlliquid
	case ClsVehicle, ClsCollectible:
		return LQLongTerm
	case ClsRetirementAsset, ClsInsuranceCashValue:
		return LQLocked
	default:
		return LQMediumTerm
	}
}
