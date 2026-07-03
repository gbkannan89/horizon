package domain

func (b *Budget) determineHealth() BudgetHealth {
	if b.status == BStatusArchived || b.status == BStatusCompleted {
		return HealthOnTrack
	}
	if b.totalBudgeted == 0 {
		return HealthOnTrack
	}
	ratio := float64(b.totalSpent) / float64(b.totalBudgeted)
	switch {
	case ratio > 1.0:
		return HealthOverspent
	case ratio > 0.9:
		return HealthCritical
	case ratio < 0.3 && b.totalSpent > 0:
		return HealthUnderTrack
	default:
		return HealthOnTrack
	}
}

func CalculateCategoryRemaining(budgeted, spent int64, rollover bool, prevRemaining int64) int64 {
	if rollover {
		return budgeted - spent + prevRemaining
	}
	remaining := budgeted - spent
	if remaining < 0 {
		return 0
	}
	return remaining
}

func CalculateTotalRemaining(categories []BudgetCategory) int64 {
	var total int64
	for _, c := range categories {
		total += c.RemainingAmount
	}
	return total
}

func CalculateSpentPercentage(spent, budgeted int64) float64 {
	if budgeted == 0 { return 0 }
	return float64(spent) / float64(budgeted) * 100
}
