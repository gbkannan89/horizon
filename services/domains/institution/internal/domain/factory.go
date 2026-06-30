package domain

import (
	"time"
)

type InstitutionFactory struct{}

func NewInstitutionFactory() *InstitutionFactory { return &InstitutionFactory{} }

func (f *InstitutionFactory) Create(
	id, name, country string, instType InstitutionType,
	products []FinancialProduct, connectivity []ConnectivityCapability,
	website, phone, email, address, headquarters, regulator, regLicense, sourceOfTruth string,
	metadata map[string]string, tags []string, notes string,
) (*Institution, error) {
	if err := ValidateName(name); err != nil { return nil, err }
	if err := ValidateCountry(country); err != nil { return nil, err }
	cat := typeCategory[instType]
	if cat == "" { cat = "FinTech" }
	if products == nil { products = productMap[instType] }
	if products == nil { products = []FinancialProduct{} }
	if connectivity == nil { connectivity = []ConnectivityCapability{CCManual} }
	if metadata == nil { metadata = map[string]string{} }
	if tags == nil { tags = []string{} }
	if sourceOfTruth == "" { sourceOfTruth = "User" }

	now := time.Now().UTC()
	return &Institution{
		id: id, name: name, instType: instType, category: cat, status: StDraft,
		country: country, trustLevel: TLUnverified, health: IHStable, confidence: ICUserDefined,
		products: products, connectivity: connectivity, website: website,
		phone: phone, email: email, address: address, headquarters: headquarters,
		regulator: regulator, regLicense: regLicense, metadata: metadata, tags: tags,
		notes: notes, sourceOfTruth: sourceOfTruth, createdAt: now, updatedAt: now,
	}, nil
}
