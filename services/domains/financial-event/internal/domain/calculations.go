package domain

import (
	"fmt"
	"math"
	"time"
)

func ComputeRunningBalance(events []*FinancialEvent) int64 {
	var balance int64
	for _, e := range events {
		if e.State() == StateConfirmed || e.State() == StatePosted {
			balance += e.Amount()
		}
	}
	return balance
}

func ComputeNetWorthDelta(events []*FinancialEvent) int64 {
	var delta int64
	for _, e := range events {
		if e.State() == StatePosted {
			delta += e.Amount()
		}
	}
	return delta
}

func ComputeDuplicateScore(candidate *FinancialEvent, existing []*FinancialEvent) (maxScore float64, matchedID string) {
	maxScore = 0.0
	for _, ex := range existing {
		if ex.State() != StateConfirmed && ex.State() != StatePosted {
			continue
		}
		if ex.EventID() == candidate.EventID() {
			continue
		}
		score := computeSingleDuplicateScore(candidate, ex)
		if score > maxScore {
			maxScore = score
			matchedID = ex.EventID()
		}
	}
	return maxScore, matchedID
}

func computeSingleDuplicateScore(candidate, existing *FinancialEvent) float64 {
	score := 0.0
	totalWeight := 0.0

	if candidate.Amount() == existing.Amount() && candidate.Currency() == existing.Currency() {
		score += 0.4
	}
	totalWeight += 0.4

	srcMatch := candidate.Source() != "" && existing.Source() != "" && candidate.Source() == existing.Source()
	dstMatch := candidate.Destination() != "" && existing.Destination() != "" && candidate.Destination() == existing.Destination()
	if srcMatch || dstMatch {
		score += 0.3
	}
	totalWeight += 0.3

	dayDiff := math.Abs(float64(candidate.EffectiveDate().Sub(existing.EffectiveDate()).Hours()) / 24.0)
	dateScore := math.Max(0, 1.0-dayDiff/3.0)
	score += dateScore * 0.2
	totalWeight += 0.2

	if candidate.Reference() != "" && existing.Reference() != "" && candidate.Reference() == existing.Reference() {
		score += 0.1
	}
	totalWeight += 0.1

	return score / totalWeight
}

func ComputeOrderIndex(effectiveDate, eventDate time.Time, sequenceNumber int64) string {
	eff := effectiveDate.UnixMilli()
	evt := eventDate.UnixMilli()
	seq := float64(sequenceNumber) / 1e9
	val := float64(eff) + float64(evt)*1e-6 + seq
	return fmt.Sprintf("%.12f", val)
}

func ComputeReversalChainLength(eventID string, getOriginal func(string) (*FinancialEvent, error)) (int, error) {
	length := 0
	currentID := eventID
	for {
		event, err := getOriginal(currentID)
		if err != nil {
			return length, nil
		}
		if event == nil || event.ReversalOfEventID() == "" {
			break
		}
		length++
		currentID = event.ReversalOfEventID()
		if length > 100 {
			return length, fmt.Errorf("reversal chain exceeds maximum length of 100")
		}
	}
	return length, nil
}

func NewEventID() string {
	b := make([]byte, 16)
	now := time.Now().UnixMilli()
	b[0] = byte(now >> 40)
	b[1] = byte(now >> 32)
	b[2] = byte(now >> 24)
	b[3] = byte(now >> 16)
	b[4] = byte(now >> 8)
	b[5] = byte(now)
	b[6] = (b[6] & 0x0f) | 0x70
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
