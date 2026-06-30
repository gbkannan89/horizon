package domain

import "math"

func CalculateAvailableCapacity(sourceBalance int64, existingReservations int64) int64 {
	c := sourceBalance - existingReservations; if c < 0 { return 0 }; return c
}

func CalculateReservedAmount(at AllocationType, available int64, weight *float64, fixedAmount *int64) int64 {
	switch at {
	case ATPercentage:
		if weight != nil { return int64(float64(available) * *weight / 100.0) }
	case ATFixedAmount:
		if fixedAmount != nil { return *fixedAmount }
	case ATRemainingBalance:
		return available
	case ATHybrid:
		pct := int64(0); fx := int64(0)
		if weight != nil { pct = int64(float64(available) * *weight / 100.0) }
		if fixedAmount != nil { fx = *fixedAmount }
		if pct > fx { return pct }; return fx
	}
	return 0
}

func CalculateAllocationHealth(reserved int64, available int64, allocated int64, target int64, status AllocationStatus) AllocationHealth {
	if status == StCompleted { return AHCompletedAlloc }
	if status == StPaused || status == StCancelled { return AHBlocked }
	if reserved > available { return AHOverallocated }
	if target > 0 && float64(allocated)/float64(target)*100 < 50 { return AHUnderfunded }
	if reserved <= 0 && status == StActive { return AHAtRisk }
	return AHHealthy
}

func CalculateAllocationConfidence(createdBy string, approvedBy *string, sourceOfTruth string) AllocationConfidence {
	if sourceOfTruth == "Migration" { return ACInferred }
	if createdBy == "Recommendation" && approvedBy == nil { return ACEstimated }
	if createdBy == "System" { return ACPlanned }
	return ACConfirmed
}

func CalculateWeightedDistribution(available int64, allocations []*Allocation) map[string]int64 {
	totalWeight := 0.0
	for _, a := range allocations {
		if a.status == StActive && a.weight != nil { totalWeight += *a.weight }
	}
	result := make(map[string]int64)
	if totalWeight == 0 { return result }
	remaining := available
	for _, a := range allocations {
		if a.status == StActive && a.weight != nil {
			amt := int64(float64(available) * *a.weight / totalWeight)
			if amt > remaining { amt = remaining }
			result[a.ID()] = amt
			remaining -= amt
		}
	}
	return result
}

func CalculatePriorityDistribution(available int64, allocations []*Allocation) map[string]int64 {
	result := make(map[string]int64)
	remaining := available
	for _, a := range allocations {
		if a.status == StActive && remaining > 0 {
			amt := int64(0)
			switch a.allocationType {
			case ATFixedAmount:
				if a.fixedAmount != nil { amt = *a.fixedAmount }
			case ATRemainingBalance:
				amt = remaining
			default:
				if a.fixedAmount != nil { amt = *a.fixedAmount } else { amt = remaining }
			}
			if amt > remaining { amt = remaining }
			result[a.ID()] = amt
			remaining -= amt
		}
	}
	return result
}

func CalculateCapacityUtilization(reserved, available int64) float64 {
	if available <= 0 { return 0 }
	return math.Round(float64(reserved)/float64(available)*100) / 100
}

func CalculateTotalAllocated(events []struct{ Amount int64 }) int64 {
	var t int64
	for _, e := range events { t += e.Amount }
	return t
}
