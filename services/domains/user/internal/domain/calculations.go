package domain

func CalculateProfileCompleteness(u *User) float64 {
	weight := 0.0
	total := 0.0
	if u.DisplayName() != "" { weight += 15; total += 15 }
	if u.LegalName() != "" { weight += 15; total += 15 }
	if u.Country() != "" { weight += 15; total += 15 }
	if u.BaseCurrency() != "" { weight += 15; total += 15 }
	if u.Locale() != "" { weight += 10; total += 10 }
	if u.Timezone() != "" { weight += 10; total += 10 }
	if u.PreferredName() != "" { weight += 5; total += 5 }
	if u.DateOfBirth() != nil { weight += 5; total += 5 }
	if u.FinancialIdentityProfile() != "" && u.FinancialIdentityProfile() != "" { weight += 5; total += 5 }
	total += 5
	if len(u.Preferences()) > 0 { weight += 5 }

	return (weight / total) * 100
}

func CalculateOwnershipDiversity(hasAccounts, hasAssets, hasLiabilities, hasGoals, hasPortfolios bool) int {
	count := 0
	if hasAccounts { count++ }
	if hasAssets { count++ }
	if hasLiabilities { count++ }
	if hasGoals { count++ }
	if hasPortfolios { count++ }
	return count
}

func CalculateHouseholdParticipation(memberships []string) int {
	return len(memberships)
}

func CalculateFinancialActivityLevel(eventCount, accountCount, txVolume int, goalActivity bool) float64 {
	score := 0.0
	if eventCount > 100 { score += 0.3 } else if eventCount > 10 { score += 0.2 } else if eventCount > 0 { score += 0.1 }
	if accountCount > 5 { score += 0.2 } else if accountCount > 0 { score += 0.1 }
	if txVolume > 1000 { score += 0.3 } else if txVolume > 100 { score += 0.2 } else if txVolume > 0 { score += 0.1 }
	if goalActivity { score += 0.2 }
	if score > 1.0 { score = 1.0 }
	return score
}

func CalculatePreferenceCoverage(configured, total int) float64 {
	if total == 0 { return 0 }
	return float64(configured) / float64(total) * 100
}

func CalculateConsentCoverage(granted, total int) float64 {
	if total == 0 { return 0 }
	return float64(granted) / float64(total) * 100
}

func CalculateUserFinancialMaturity(accountAgeDays int, goalCount, portfolioCount, allocationChanges, engagementDurationDays int) float64 {
	score := 0.0
	if accountAgeDays > 365*3 { score += 0.3 } else if accountAgeDays > 365 { score += 0.2 } else if accountAgeDays > 30 { score += 0.1 }
	if goalCount > 3 { score += 0.2 } else if goalCount > 0 { score += 0.1 }
	if portfolioCount > 0 { score += 0.1 }
	if allocationChanges > 5 { score += 0.2 } else if allocationChanges > 0 { score += 0.1 }
	if engagementDurationDays > 365 { score += 0.2 } else if engagementDurationDays > 90 { score += 0.1 }
	if score > 1.0 { score = 1.0 }
	return score
}

func CalculateUserEngagementTrend(current, previous float64) float64 {
	return current - previous
}

func CalculateUserOwnershipDiversity(owned, total int) float64 {
	if total == 0 { return 0 }
	return float64(owned) / float64(total) * 100
}

func CalculateUserPlanningConsistency(goalStability, allocationStability, membershipStability float64, preferenceChangeFrequency int) float64 {
	score := (goalStability + allocationStability + membershipStability) / 3.0
	if preferenceChangeFrequency > 10 { score -= 0.2 } else if preferenceChangeFrequency > 5 { score -= 0.1 }
	if score < 0 { score = 0 }
	if score > 1 { score = 1 }
	_ = goalStability
	return score
}

func CalculateUserFinancialStability(health UserHealth, engagementTrend, planningConsistency, ownershipDiversity float64) float64 {
	healthScore := 0.5
	switch health {
	case HealthHealthy: healthScore = 1.0
	case HealthActive: healthScore = 0.8
	case HealthGrowing: healthScore = 0.6
	case HealthAtRisk: healthScore = 0.3
	case HealthInactive: healthScore = 0.1
	case HealthHistorical: healthScore = 0.0
	}
	return (healthScore + engagementTrend + planningConsistency + ownershipDiversity) / 4.0
}

func CalculateUserHouseholdParticipationScore(householdCount int, activeSharedGoals, sharedEntityCount int) float64 {
	score := 0.0
	if householdCount > 2 { score += 0.4 } else if householdCount > 0 { score += 0.2 }
	if activeSharedGoals > 3 { score += 0.3 } else if activeSharedGoals > 0 { score += 0.15 }
	if sharedEntityCount > 5 { score += 0.3 } else if sharedEntityCount > 0 { score += 0.15 }
	if score > 1.0 { score = 1.0 }
	return score
}

func CalculateFinancialEngagement(loginFrequency, featureUsage, goalUpdates, allocationChanges, householdActivity int) float64 {
	score := 0.0
	if loginFrequency > 20 { score += 0.3 } else if loginFrequency > 5 { score += 0.2 } else if loginFrequency > 0 { score += 0.1 }
	if featureUsage > 5 { score += 0.2 } else if featureUsage > 0 { score += 0.1 }
	if goalUpdates > 3 { score += 0.2 } else if goalUpdates > 0 { score += 0.1 }
	if allocationChanges > 3 { score += 0.15 } else if allocationChanges > 0 { score += 0.05 }
	if householdActivity > 3 { score += 0.15 } else if householdActivity > 0 { score += 0.05 }
	if score > 1.0 { score = 1.0 }
	return score
}
