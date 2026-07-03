package domain

import (
	"time"
)

func (r *RecurringTransaction) CalculateNextOccurrence(from time.Time) (time.Time, error) {
	var next time.Time
	
	switch r.frequency {
	case FreqDaily:
		next = from.AddDate(0, 0, r.interval)
	case FreqWeekly:
		next = from.AddDate(0, 0, 7*r.interval)
	case FreqBiWeekly:
		next = from.AddDate(0, 0, 14*r.interval)
	case FreqMonthly:
		next = from.AddDate(0, r.interval, 0)
	case FreqQuarterly:
		next = from.AddDate(0, 3*r.interval, 0)
	case FreqSemiAnnual:
		next = from.AddDate(0, 6*r.interval, 0)
	case FreqAnnual:
		next = from.AddDate(r.interval, 0, 0)
	default: // FreqCustom
		next = from.AddDate(0, 0, r.interval)
	}

	if r.skipWeekends {
		for next.Weekday() == time.Saturday || next.Weekday() == time.Sunday {
			next = next.AddDate(0, 0, 1)
		}
	}
	
	return next, nil
}
