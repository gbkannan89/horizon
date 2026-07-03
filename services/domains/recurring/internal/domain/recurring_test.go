package domain

import (
	"testing"
	"time"
)

func fixedDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func ptr(t time.Time) *time.Time { return &t }

// ---------------------------------------------------------------------------
// Factory — NewRecurringTransaction
// ---------------------------------------------------------------------------

func TestNewRecurringTransaction_Valid(t *testing.T) {
	t.Parallel()

	start := fixedDate(2025, 1, 1)
	end := fixedDate(2025, 12, 31)
	tmpl := EventTemplate{Type: "transfer", Category: "bill", Source: "checking", Destination: "savings"}

	rt, err := NewRecurringTransaction(
		"tx-1", "u1", "h1", "Rent", "Monthly rent",
		150_000, "USD", FreqMonthly, 1,
		start, &end, tmpl, false, false,
		[]string{"rent", "housing"}, map[string]string{"key": "val"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rt.ID() != "tx-1" {
		t.Errorf("id = %q", rt.ID())
	}
	if rt.UserID() != "u1" {
		t.Errorf("userID = %q", rt.UserID())
	}
	if rt.HouseholdID() != "h1" {
		t.Errorf("householdID = %q", rt.HouseholdID())
	}
	if rt.Name() != "Rent" {
		t.Errorf("name = %q", rt.Name())
	}
	if rt.Description() != "Monthly rent" {
		t.Errorf("description = %q", rt.Description())
	}
	if rt.Amount() != 150_000 {
		t.Errorf("amount = %d", rt.Amount())
	}
	if rt.Currency() != "USD" {
		t.Errorf("currency = %q", rt.Currency())
	}
	if rt.Frequency() != FreqMonthly {
		t.Errorf("frequency = %q", rt.Frequency())
	}
	if rt.Interval() != 1 {
		t.Errorf("interval = %d", rt.Interval())
	}
	if !rt.StartDate().Equal(start) {
		t.Errorf("startDate = %v", rt.StartDate())
	}
	if rt.EndDate() == nil || !rt.EndDate().Equal(end) {
		t.Errorf("endDate = %v", rt.EndDate())
	}
	if rt.Status() != StatusActive {
		t.Errorf("status = %q", rt.Status())
	}
	if rt.EventTemplate() != tmpl {
		t.Errorf("eventTemplate = %v", rt.EventTemplate())
	}
	if rt.SkipHolidays() {
		t.Error("skipHolidays should be false")
	}
	if rt.SkipWeekends() {
		t.Error("skipWeekends should be false")
	}
	if len(rt.Tags()) != 2 || rt.Tags()[0] != "rent" {
		t.Errorf("tags = %v", rt.Tags())
	}
	if rt.Metadata()["key"] != "val" {
		t.Errorf("metadata = %v", rt.Metadata())
	}
	if rt.CreatedAt().IsZero() {
		t.Error("createdAt should not be zero")
	}
	if rt.UpdatedAt().IsZero() {
		t.Error("updatedAt should not be zero")
	}
	// next occurrence = start + 1 month = Feb 1
	expectedNext := fixedDate(2025, 2, 1)
	if rt.NextOccurrence() == nil || !rt.NextOccurrence().Equal(expectedNext) {
		t.Errorf("nextOccurrence = %v, want %v", rt.NextOccurrence(), expectedNext)
	}
	if rt.LastOccurrence() != nil {
		t.Errorf("lastOccurrence should be nil, got %v", rt.LastOccurrence())
	}
}

func TestNewRecurringTransaction_ValidationErrors(t *testing.T) {
	t.Parallel()

	start := fixedDate(2025, 1, 1)
	end := fixedDate(2025, 12, 31)
	tmpl := EventTemplate{Type: "transfer", Category: "bill", Source: "checking"}

	tests := []struct {
		name string
		edit func(rt *RecurringTransaction)
		want error
	}{
		{name: "empty name", edit: func(rt *RecurringTransaction) { rt.name = "" }, want: ErrInvalidName},
		{name: "zero amount", edit: func(rt *RecurringTransaction) { rt.amount = 0 }, want: ErrInvalidAmount},
		{name: "negative amount", edit: func(rt *RecurringTransaction) { rt.amount = -1 }, want: ErrInvalidAmount},
		{name: "zero interval", edit: func(rt *RecurringTransaction) { rt.interval = 0 }, want: ErrInvalidInterval},
		{name: "negative interval", edit: func(rt *RecurringTransaction) { rt.interval = -1 }, want: ErrInvalidInterval},
		{name: "end before start", edit: func(rt *RecurringTransaction) { rt.endDate = ptr(fixedDate(2024, 12, 31)) }, want: ErrInvalidDates},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We call the factory, then modify and validate manually to keep
			// the arrange section simple.
			rt := &RecurringTransaction{
				id: "id", userID: "u", householdID: "h",
				name: "Valid", description: "d",
				amount: 1000, currency: "USD",
				frequency: FreqMonthly, interval: 1,
				startDate: start, endDate: &end,
				status: StatusActive,
				eventTemplate: tmpl,
			}
			tt.edit(rt)
			err := rt.Validate()
			if err != tt.want {
				t.Errorf("got %v, want %v", err, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// State transitions
// ---------------------------------------------------------------------------

func makeActive(t *testing.T) *RecurringTransaction {
	t.Helper()
	start := fixedDate(2025, 1, 1)
	rt := &RecurringTransaction{
		id: "t", userID: "u", householdID: "h",
		name: "Test", amount: 1000, currency: "USD",
		frequency: FreqMonthly, interval: 1,
		startDate: start,
		status:    StatusActive,
	}
	rt.nextOccurrence = ptr(fixedDate(2025, 2, 1))
	return rt
}

func TestActivate(t *testing.T) {
	t.Parallel()

	t.Run("Paused_to_Active", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusPaused
		if err := rt.Activate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.Status() != StatusActive {
			t.Errorf("status = %q", rt.Status())
		}
	})

	t.Run("Active_returns_error", func(t *testing.T) {
		rt := makeActive(t)
		if err := rt.Activate(); err != ErrInvalidTransition {
			t.Errorf("got %v, want %v", err, ErrInvalidTransition)
		}
	})
}

func TestPause(t *testing.T) {
	t.Parallel()

	t.Run("Active_to_Paused", func(t *testing.T) {
		rt := makeActive(t)
		if err := rt.Pause(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.Status() != StatusPaused {
			t.Errorf("status = %q", rt.Status())
		}
	})

	t.Run("Paused_returns_error", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusPaused
		if err := rt.Pause(); err != ErrInvalidTransition {
			t.Errorf("got %v, want %v", err, ErrInvalidTransition)
		}
	})
}

func TestCancel(t *testing.T) {
	t.Parallel()

	t.Run("Active_to_Cancelled", func(t *testing.T) {
		rt := makeActive(t)
		if err := rt.Cancel(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.Status() != StatusCancelled {
			t.Errorf("status = %q", rt.Status())
		}
	})

	t.Run("Paused_to_Cancelled", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusPaused
		if err := rt.Cancel(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.Status() != StatusCancelled {
			t.Errorf("status = %q", rt.Status())
		}
	})

	t.Run("Completed_returns_error", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusCompleted
		if err := rt.Cancel(); err != ErrInvalidTransition {
			t.Errorf("got %v, want %v", err, ErrInvalidTransition)
		}
	})

	t.Run("Archived_returns_error", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusArchived
		if err := rt.Cancel(); err != ErrInvalidTransition {
			t.Errorf("got %v, want %v", err, ErrInvalidTransition)
		}
	})
}

func TestComplete(t *testing.T) {
	t.Parallel()

	t.Run("Active_to_Completed", func(t *testing.T) {
		rt := makeActive(t)
		if err := rt.Complete(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.Status() != StatusCompleted {
			t.Errorf("status = %q", rt.Status())
		}
	})

	t.Run("Paused_returns_error", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusPaused
		if err := rt.Complete(); err != ErrInvalidTransition {
			t.Errorf("got %v, want %v", err, ErrInvalidTransition)
		}
	})
}

func TestArchive(t *testing.T) {
	t.Parallel()

	t.Run("Completed_to_Archived", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusCompleted
		if err := rt.Archive(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.Status() != StatusArchived {
			t.Errorf("status = %q", rt.Status())
		}
	})

	t.Run("Cancelled_to_Archived", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusCancelled
		if err := rt.Archive(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.Status() != StatusArchived {
			t.Errorf("status = %q", rt.Status())
		}
	})

	t.Run("Active_returns_error", func(t *testing.T) {
		rt := makeActive(t)
		if err := rt.Archive(); err != ErrInvalidTransition {
			t.Errorf("got %v, want %v", err, ErrInvalidTransition)
		}
	})

	t.Run("Paused_returns_error", func(t *testing.T) {
		rt := makeActive(t)
		rt.status = StatusPaused
		if err := rt.Archive(); err != ErrInvalidTransition {
			t.Errorf("got %v, want %v", err, ErrInvalidTransition)
		}
	})
}

// ---------------------------------------------------------------------------
// AdvanceOccurrence
// ---------------------------------------------------------------------------

func TestAdvanceOccurrence(t *testing.T) {
	t.Parallel()

	t.Run("basic_advance", func(t *testing.T) {
		start := fixedDate(2025, 1, 1)
		rt := &RecurringTransaction{
			name: "Test", amount: 1000, currency: "USD",
			frequency: FreqMonthly, interval: 1,
			startDate:      start,
			status:         StatusActive,
			nextOccurrence: ptr(fixedDate(2025, 2, 1)),
		}

		if err := rt.AdvanceOccurrence(fixedDate(2025, 2, 1)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// last becomes old next (Feb 1)
		if rt.LastOccurrence() == nil || !rt.LastOccurrence().Equal(fixedDate(2025, 2, 1)) {
			t.Errorf("lastOccurrence = %v, want Feb 1", rt.LastOccurrence())
		}
		// next becomes Mar 1
		if rt.NextOccurrence() == nil || !rt.NextOccurrence().Equal(fixedDate(2025, 3, 1)) {
			t.Errorf("nextOccurrence = %v, want Mar 1", rt.NextOccurrence())
		}
		if rt.Status() != StatusActive {
			t.Errorf("status should remain Active, got %q", rt.Status())
		}
	})

	t.Run("advance_with_nil_next", func(t *testing.T) {
		rt := &RecurringTransaction{
			name: "Test", amount: 1000, currency: "USD",
			frequency: FreqDaily, interval: 1,
			startDate:      fixedDate(2025, 1, 1),
			status:         StatusActive,
			nextOccurrence: nil,
		}
		date := fixedDate(2025, 3, 15)
		if err := rt.AdvanceOccurrence(date); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.LastOccurrence() == nil || !rt.LastOccurrence().Equal(date) {
			t.Errorf("lastOccurrence = %v, want %v", rt.LastOccurrence(), date)
		}
	})

	t.Run("advance_past_endDate_completes", func(t *testing.T) {
		start := fixedDate(2025, 1, 1)
		end := fixedDate(2025, 2, 1)
		rt := &RecurringTransaction{
			name: "Test", amount: 1000, currency: "USD",
			frequency: FreqMonthly, interval: 1,
			startDate:      start,
			endDate:        &end,
			status:         StatusActive,
			nextOccurrence: ptr(fixedDate(2025, 2, 1)),
		}
		if err := rt.AdvanceOccurrence(fixedDate(2025, 2, 1)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Next occurrence (Mar 1) is after endDate (Feb 1) → complete
		if rt.NextOccurrence() != nil {
			t.Errorf("nextOccurrence should be nil, got %v", rt.NextOccurrence())
		}
		if rt.Status() != StatusCompleted {
			t.Errorf("status = %q, want Completed", rt.Status())
		}
	})

	t.Run("advance_on_endDate_stays_active", func(t *testing.T) {
		// If the next calculated occurrence *equals* endDate, it should not
		// complete — only strictly *after* triggers completion.
		start := fixedDate(2025, 1, 1)
		// Monthly: next from Feb 1 is Mar 1. Set endDate to Mar 1 exactly.
		end := fixedDate(2025, 3, 1)
		rt := &RecurringTransaction{
			name: "Test", amount: 1000, currency: "USD",
			frequency: FreqMonthly, interval: 1,
			startDate:      start,
			endDate:        &end,
			status:         StatusActive,
			nextOccurrence: ptr(fixedDate(2025, 2, 1)),
		}
		if err := rt.AdvanceOccurrence(fixedDate(2025, 2, 1)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt.NextOccurrence() == nil || !rt.NextOccurrence().Equal(fixedDate(2025, 3, 1)) {
			t.Errorf("nextOccurrence = %v, want Mar 1", rt.NextOccurrence())
		}
		if rt.Status() != StatusActive {
			t.Errorf("status should remain Active, got %q", rt.Status())
		}
	})
}

// ---------------------------------------------------------------------------
// ReconstructFromDB
// ---------------------------------------------------------------------------

func TestReconstructFromDB(t *testing.T) {
	t.Parallel()

	createdAt := fixedDate(2025, 1, 1)
	updatedAt := fixedDate(2025, 2, 1)
	start := fixedDate(2025, 1, 15)
	end := fixedDate(2025, 12, 31)
	next := fixedDate(2025, 2, 15)
	last := fixedDate(2025, 1, 15)

	t.Run("all_fields", func(t *testing.T) {
		rt := ReconstructFromDB(
			"tx-99", "u99", "h99", "Gym", "Monthly membership",
			5000, "EUR",
			"Monthly", 1,
			start, &end, &next, &last,
			"Active",
			"payment", "health", "checking", "",
			false, true,
			[]string{"fitness", "gym"}, map[string]string{"plan": "premium"},
			createdAt, updatedAt,
		)

		if rt.ID() != "tx-99" { t.Errorf("id = %q", rt.ID()) }
		if rt.Name() != "Gym" { t.Errorf("name = %q", rt.Name()) }
		if rt.Amount() != 5000 { t.Errorf("amount = %d", rt.Amount()) }
		if rt.Currency() != "EUR" { t.Errorf("currency = %q", rt.Currency()) }
		if rt.Frequency() != FreqMonthly { t.Errorf("frequency = %q", rt.Frequency()) }
		if rt.Interval() != 1 { t.Errorf("interval = %d", rt.Interval()) }
		if rt.Status() != StatusActive { t.Errorf("status = %q", rt.Status()) }
		if rt.SkipWeekends() != true { t.Errorf("skipWeekends = %v", rt.SkipWeekends()) }
		if !rt.StartDate().Equal(start) { t.Errorf("startDate mismatch") }
		if !rt.EndDate().Equal(end) { t.Errorf("endDate mismatch") }
		if !rt.NextOccurrence().Equal(next) { t.Errorf("next mismatch") }
		if !rt.LastOccurrence().Equal(last) { t.Errorf("last mismatch") }
		if rt.EventTemplate().Type != "payment" { t.Errorf("event type = %q", rt.EventTemplate().Type) }
		if len(rt.Tags()) != 2 { t.Errorf("tags len = %d", len(rt.Tags())) }
		if rt.Metadata()["plan"] != "premium" { t.Errorf("metadata = %v", rt.Metadata()) }
		if !rt.CreatedAt().Equal(createdAt) { t.Errorf("createdAt mismatch") }
		if !rt.UpdatedAt().Equal(updatedAt) { t.Errorf("updatedAt mismatch") }
	})

	t.Run("nil_tags_and_metadata", func(t *testing.T) {
		rt := ReconstructFromDB(
			"id", "u", "h", "N", "d",
			100, "USD",
			"Daily", 1,
			start, nil, nil, nil,
			"Active",
			"", "", "", "",
			false, false,
			nil, nil,
			createdAt, updatedAt,
		)
		if rt.Tags() == nil {
			t.Error("tags should be non-nil empty slice")
		}
		if len(rt.Tags()) != 0 {
			t.Errorf("tags = %v", rt.Tags())
		}
		if rt.Metadata() == nil {
			t.Error("metadata should be non-nil empty map")
		}
		if len(rt.Metadata()) != 0 {
			t.Errorf("metadata = %v", rt.Metadata())
		}
	})

	t.Run("nil_pointers", func(t *testing.T) {
		rt := ReconstructFromDB(
			"id", "u", "h", "N", "d",
			100, "USD",
			"Daily", 1,
			start, nil, nil, nil,
			"Paused",
			"", "", "", "",
			false, false,
			[]string{}, map[string]string{},
			createdAt, updatedAt,
		)
		if rt.EndDate() != nil {
			t.Error("endDate should be nil")
		}
		if rt.NextOccurrence() != nil {
			t.Error("nextOccurrence should be nil")
		}
		if rt.LastOccurrence() != nil {
			t.Error("lastOccurrence should be nil")
		}
		if rt.Status() != StatusPaused {
			t.Errorf("status = %q", rt.Status())
		}
	})
}

// ---------------------------------------------------------------------------
// CalculateNextOccurrence — all frequency types
// ---------------------------------------------------------------------------

func TestCalculateNextOccurrence(t *testing.T) {
	t.Parallel()

	base := fixedDate(2025, 1, 15)

	tests := []struct {
		name      string
		freq      Frequency
		interval  int
		from      time.Time
		want      time.Time
	}{
		{freq: FreqDaily,      interval: 1, from: base, want: fixedDate(2025, 1, 16)},
		{freq: FreqDaily,      interval: 3, from: base, want: fixedDate(2025, 1, 18)},
		{freq: FreqWeekly,     interval: 1, from: base, want: fixedDate(2025, 1, 22)},
		{freq: FreqWeekly,     interval: 2, from: base, want: fixedDate(2025, 1, 29)},
		{freq: FreqBiWeekly,   interval: 1, from: base, want: fixedDate(2025, 1, 29)},
		{freq: FreqBiWeekly,   interval: 2, from: base, want: fixedDate(2025, 2, 12)},
		{freq: FreqMonthly,    interval: 1, from: base, want: fixedDate(2025, 2, 15)},
		{freq: FreqMonthly,    interval: 3, from: base, want: fixedDate(2025, 4, 15)},
		{freq: FreqQuarterly,  interval: 1, from: base, want: fixedDate(2025, 4, 15)},
		{freq: FreqQuarterly,  interval: 2, from: base, want: fixedDate(2025, 7, 15)},
		{freq: FreqSemiAnnual, interval: 1, from: base, want: fixedDate(2025, 7, 15)},
		{freq: FreqSemiAnnual, interval: 2, from: base, want: fixedDate(2026, 1, 15)},
		{freq: FreqAnnual,     interval: 1, from: base, want: fixedDate(2026, 1, 15)},
		{freq: FreqAnnual,     interval: 2, from: base, want: fixedDate(2027, 1, 15)},
		{freq: FreqCustom,     interval: 5, from: base, want: fixedDate(2025, 1, 20)},
	}

	for _, tt := range tests {
		t.Run(string(tt.freq)+"_"+tt.name, func(t *testing.T) {
			rt := &RecurringTransaction{
				frequency: tt.freq,
				interval:  tt.interval,
			}
			got, err := rt.CalculateNextOccurrence(tt.from)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateNextOccurrence_WeekendSkipping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		from   time.Time
		skip   bool
		want   time.Time
	}{
		// Saturday = 6, Sunday = 0
		{name: "fri_skip_enabled",  from: fixedDate(2025, 1, 10), skip: true,  want: fixedDate(2025, 1, 13)}, // 10+1=11 Sat → 13 Mon
		{name: "fri_skip_disabled", from: fixedDate(2025, 1, 10), skip: false, want: fixedDate(2025, 1, 11)}, // 10+1=11 Sat
		{name: "sat_skip_enabled",  from: fixedDate(2025, 1, 10), skip: true,  want: fixedDate(2025, 1, 13)}, // 10+1=11 Sat→skip→12 Sun→skip→13 Mon
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RecurringTransaction{
				frequency:    FreqDaily,
				interval:     1,
				skipWeekends: tt.skip,
			}
			got, err := rt.CalculateNextOccurrence(tt.from)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Validate — edge cases
// ---------------------------------------------------------------------------

func TestValidate_EdgeCases(t *testing.T) {
	t.Parallel()

	start := fixedDate(2025, 1, 1)

	t.Run("same_start_and_end", func(t *testing.T) {
		rt := &RecurringTransaction{
			name: "X", amount: 100, interval: 1,
			startDate: start, endDate: ptr(start),
		}
		if err := rt.Validate(); err != nil {
			t.Errorf("same date should be valid: %v", err)
		}
	})

	t.Run("nil_endDate", func(t *testing.T) {
		rt := &RecurringTransaction{
			name: "X", amount: 100, interval: 1,
			startDate: start, endDate: nil,
		}
		if err := rt.Validate(); err != nil {
			t.Errorf("nil endDate should be valid: %v", err)
		}
	})
}
