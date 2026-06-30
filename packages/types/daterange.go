package types

import (
	"fmt"
	"time"
)

type DateRange struct {
	Start time.Time
	End   time.Time
}

func NewDateRange(start, end time.Time) (DateRange, error) {
	if start.After(end) {
		return DateRange{}, fmt.Errorf("start must be before or equal to end")
	}
	return DateRange{Start: start, End: end}, nil
}

func (d DateRange) Contains(t time.Time) bool {
	return !t.Before(d.Start) && !t.After(d.End)
}

func (d DateRange) Overlaps(other DateRange) bool {
	return d.Contains(other.Start) || d.Contains(other.End) || other.Contains(d.Start)
}

func (d DateRange) Duration() time.Duration {
	return d.End.Sub(d.Start)
}

func (d DateRange) IsEmpty() bool {
	return d.Start.IsZero() && d.End.IsZero()
}
