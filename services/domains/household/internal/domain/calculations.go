package domain

func CalculateTotalNetWorth(totalAssets, totalLiabilities int64) int64 {
	return totalAssets - totalLiabilities
}

func DetermineHouseholdHealth(totalNetWorth int64, status HouseholdStatus, activeMemberRatio float64) HouseholdHealth {
	if status == HHStatusDissolved || status == HHStatusArchived {
		return HHCritical
	}
	if totalNetWorth < 0 {
		return HHCritical
	}
	if activeMemberRatio < 0.5 {
		return HHWarning
	}
	if totalNetWorth == 0 {
		return HHWarning
	}
	return HHHealthy
}

func CalculateActiveMemberRatio(totalMembers int, activeMembers int) float64 {
	if totalMembers == 0 { return 0 }
	return float64(activeMembers) / float64(totalMembers)
}

func DeriveDefaultCurrency(country string) string {
	switch country {
	case "IN", "India":
		return "INR"
	case "US", "USA", "United States":
		return "USD"
	case "GB", "UK", "United Kingdom":
		return "GBP"
	case "EU", "Eurozone":
		return "EUR"
	case "AE", "UAE":
		return "AED"
	case "SG", "Singapore":
		return "SGD"
	default:
		return "INR"
	}
}
