package domain

func CalculateInstitutionHealth(s InstitutionStatus, activeAccounts int, daysSinceActivity int, failedSyncs int) InstitutionHealth {
	if s == StClosed || s == StArchived { return IHHistorical }
	if s == StSuspended { return IHSuspended }
	if failedSyncs > 10 { return IHRestricted }
	if daysSinceActivity > 90 { return IHWarning }
	if activeAccounts == 0 { return IHStable }
	return IHHealthy
}

func CalculateInstitutionConfidence(s InstitutionStatus, trust TrustLevel, sourceOfTruth string) InstitutionConfidence {
	if sourceOfTruth == "Government" || sourceOfTruth == "Regulator" { return ICVerified }
	if trust == TLVerified { return ICTrusted }
	if s == StRegistered { return ICRegistered }
	if sourceOfTruth == "Import" { return ICImported }
	if sourceOfTruth == "User" { return ICUserDefined }
	return ICUnknown
}

func CalculateTrustLevel(s InstitutionStatus, verified bool, accountAgeDays int, syncSuccessRate float64) TrustLevel {
	if s == StClosed || s == StArchived { return TLUntrusted }
	if verified { return TLVerified }
	if s == StRegistered { return TLRegistered }
	if accountAgeDays > 365 && syncSuccessRate > 0.95 { return TLRegistered }
	return TLUnverified
}

func CalculateConnectivityScore(connectivity []ConnectivityCapability) int {
	s := 0
	for _, c := range connectivity {
		switch c {
		case CCOpenBanking: s += 5
		case CCBrokerAPI, CCPayrollFeed, CCGovernmentFeed: s += 4
		case CCCSVImport, CCOCRImport, CCFileImport: s += 2
		case CCManual: s += 1
		case CCFuture: s += 3
		}
	}
	return s
}

func CalculateProductCoverage(actual, expected []FinancialProduct) float64 {
	if len(expected) == 0 { return 1.0 }
	exp := map[FinancialProduct]bool{}
	for _, p := range expected { exp[p] = true }
	matched := 0
	for _, a := range actual { if exp[a] { matched++ } }
	return float64(matched) / float64(len(expected))
}

func ValidateProductForType(instType InstitutionType, product FinancialProduct) bool {
	allowed := productMap[instType]
	for _, p := range allowed { if p == product { return true } }
	return false
}
