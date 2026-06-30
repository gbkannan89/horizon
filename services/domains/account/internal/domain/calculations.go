package domain

func CalculateCurrentBalance(events []struct{ Amount int64; State string }) int64 {
	var b int64
	for _, e := range events {
		if e.State == "Confirmed" || e.State == "Posted" { b += e.Amount }
	}
	return b
}

func CalculateAvailableBalance(current int64, uncleared int64, reserved int64, creditLimit *int64) int64 {
	avail := current + uncleared - reserved
	if creditLimit != nil && avail < 0 { avail = *creditLimit + avail }
	return avail
}

func CalculateSpendableBalance(available, allocated int64) int64 {
	s := available - allocated
	if s < 0 { return 0 }
	return s
}

func CalculateAccountHealth(balance int64, creditLimit *int64, status AccountStatus, daysSinceActivity int) AccountHealth {
	if status == StatusClosed { return AHClosed }
	if status == StatusDormant || daysSinceActivity > 90 { return AHDormant }
	if creditLimit != nil && *creditLimit > 0 {
		utilization := float64(-balance) / float64(*creditLimit) * 100
		if utilization > 90 { return AHCritical }
		if utilization > 70 { return AHWarning }
	}
	if status == StatusFrozen { return AHWarning }
	return AHHealthy
}

func CalculateLiquidityProfile(at AccountType, balance int64, status AccountStatus) LiquidityProfile {
	if status == StatusFrozen || status == StatusClosed { return LQLocked }
	return deriveLiquidity(at)
}

func CalculateBalanceFromEvents(events []struct{ Amount int64; State string }) int64 {
	return CalculateCurrentBalance(events)
}
