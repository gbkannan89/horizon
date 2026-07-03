package domain

import (
	"testing"
	"time"
)

func TestBudgetFactoryCreate(t *testing.T) {
	f := NewBudgetFactory()
	start := time.Now()
	end := start.AddDate(0, 1, 0)
	cats := []BudgetCategory{
		{Category: "Food", BudgetedAmount: 10000},
		{Category: "Transport", BudgetedAmount: 5000},
	}

	b, err := f.Create("", "user-1", "", "Monthly Budget", PeriodMonthly, start, end, "INR", cats, nil)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if b.Name() != "Monthly Budget" {
		t.Errorf("expected 'Monthly Budget', got '%s'", b.Name())
	}
	if b.UserID() != "user-1" {
		t.Errorf("expected user-1, got %s", b.UserID())
	}
	if b.Status() != BStatusDraft {
		t.Errorf("expected Draft, got %s", b.Status())
	}
	if b.TotalBudgeted() != 15000 {
		t.Errorf("expected 15000, got %d", b.TotalBudgeted())
	}
	if b.TotalSpent() != 0 {
		t.Errorf("expected 0, got %d", b.TotalSpent())
	}
	if b.TotalRemaining() != 15000 {
		t.Errorf("expected 15000, got %d", b.TotalRemaining())
	}
	if b.Currency() != "INR" {
		t.Errorf("expected INR, got %s", b.Currency())
	}
}

func TestBudgetFactoryValidation(t *testing.T) {
	f := NewBudgetFactory()
	start := time.Now()

	tests := []struct {
		scenario string
		name     string
		period   BudgetPeriod
		end      time.Time
		cur      string
		err      string
	}{
		{"empty name", "", PeriodMonthly, start.AddDate(0, 1, 0), "INR", "budget name is required"},
		{"bad period", "Test", "Invalid", start.AddDate(0, 1, 0), "INR", "budget period is not recognized"},
		{"bad currency", "Test", PeriodMonthly, start.AddDate(0, 1, 0), "XYZ123", "currency must be a 3-letter ISO code"},
		{"end before start", "Test", PeriodMonthly, start.AddDate(0, 0, -1), "INR", "end date must be after start date"},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			_, err := f.Create("", "u1", "", tc.name, tc.period, start, tc.end, tc.cur, nil, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tc.err {
				t.Errorf("expected '%s', got '%s'", tc.err, err.Error())
			}
		})
	}
}

func TestBudgetCategoryValidation(t *testing.T) {
	f := NewBudgetFactory()
	start := time.Now()
	end := start.AddDate(0, 1, 0)

	emptyCat := []BudgetCategory{{Category: "", BudgetedAmount: 1000}}
	_, err := f.Create("", "u1", "", "Test", PeriodMonthly, start, end, "INR", emptyCat, nil)
	if err == nil {
		t.Fatal("expected error for empty category name")
	}

	negAmount := []BudgetCategory{{Category: "Food", BudgetedAmount: -100}}
	_, err = f.Create("", "u1", "", "Test", PeriodMonthly, start, end, "INR", negAmount, nil)
	if err == nil {
		t.Fatal("expected error for negative amount")
	}
}

func makeTestBudget(t *testing.T) *Budget {
	t.Helper()
	f := NewBudgetFactory()
	b, err := f.Create("budget-1", "user-1", "", "Test", PeriodMonthly, time.Now(), time.Now().AddDate(0, 1, 0), "INR", []BudgetCategory{
		{ID: "cat-1", Category: "Food", BudgetedAmount: 10000},
		{ID: "cat-2", Category: "Transport", BudgetedAmount: 5000},
	}, nil)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	return b
}

func TestActivateBudget(t *testing.T) {
	b := makeTestBudget(t)
	if err := ActivateBudget(b); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if b.Status() != BStatusActive {
		t.Errorf("expected Active, got %s", b.Status())
	}

	// Cannot activate again
	if err := ActivateBudget(b); err == nil {
		t.Error("expected error activating active budget")
	}
}

func TestPauseResumeCycle(t *testing.T) {
	b := makeTestBudget(t)
	ActivateBudget(b)

	if err := PauseBudget(b); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if b.Status() != BStatusPaused {
		t.Errorf("expected Paused, got %s", b.Status())
	}

	if err := ResumeBudget(b); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if b.Status() != BStatusActive {
		t.Errorf("expected Active, got %s", b.Status())
	}
}

func TestCompleteBudget(t *testing.T) {
	b := makeTestBudget(t)
	ActivateBudget(b)

	if err := CompleteBudget(b); err != nil {
		t.Fatalf("complete: %v", err)
	}
	if b.Status() != BStatusCompleted {
		t.Errorf("expected Completed, got %s", b.Status())
	}

	// Draft budgets cannot be completed
	b2 := makeTestBudget(t)
	if err := CompleteBudget(b2); err == nil {
		t.Error("expected error completing draft budget")
	}
}

func TestArchiveBudget(t *testing.T) {
	b := makeTestBudget(t)
	ActivateBudget(b)
	CompleteBudget(b)

	if err := ArchiveBudget(b); err != nil {
		t.Fatalf("archive: %v", err)
	}
	if b.Status() != BStatusArchived {
		t.Errorf("expected Archived, got %s", b.Status())
	}

	// Draft budgets cannot be archived
	b2 := makeTestBudget(t)
	if err := ArchiveBudget(b2); err == nil {
		t.Error("expected error archiving draft budget")
	}
}

func TestInvalidTransitions(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*Budget) error
	}{
		{"pause draft", func(b *Budget) error { return PauseBudget(b) }},
		{"resume draft", func(b *Budget) error { return ResumeBudget(b) }},
		{"resume active", func(b *Budget) error { ActivateBudget(b); return ResumeBudget(b) }},
		{"pause completed", func(b *Budget) error { ActivateBudget(b); CompleteBudget(b); return PauseBudget(b) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := makeTestBudget(t)
			if err := tc.fn(b); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestDetermineHealth(t *testing.T) {
	tests := []struct {
		name   string
		spent  int64
		budget int64
		health BudgetHealth
	}{
		{"on track", 3000, 10000, HealthOnTrack},
		{"overspent", 12000, 10000, HealthOverspent},
		{"critical", 9500, 10000, HealthCritical},
		{"under track", 2000, 10000, HealthUnderTrack},
		{"zero budget", 0, 0, HealthOnTrack},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &Budget{
				totalBudgeted: tc.budget,
				totalSpent:    tc.spent,
				status:        BStatusActive,
			}
			if h := b.determineHealth(); h != tc.health {
				t.Errorf("expected %s, got %s", tc.health, h)
			}
		})
	}
}

func TestDetermineHealthCompleted(t *testing.T) {
	b := &Budget{status: BStatusCompleted}
	if h := b.determineHealth(); h != HealthOnTrack {
		t.Errorf("expected OnTrack for completed, got %s", h)
	}
}

func TestCalculateCategoryRemaining(t *testing.T) {
	tests := []struct {
		name      string
		budgeted  int64
		spent     int64
		rollover  bool
		prev      int64
		expected  int64
	}{
		{"normal", 10000, 4000, false, 0, 6000},
		{"overspent no rollover", 10000, 12000, false, 0, 0},
		{"rollover", 10000, 4000, true, 2000, 8000},
		{"rollover overspent", 10000, 12000, true, 2000, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateCategoryRemaining(tc.budgeted, tc.spent, tc.rollover, tc.prev)
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}

func TestCalculateSpentPercentage(t *testing.T) {
	if p := CalculateSpentPercentage(2500, 10000); p != 25.0 {
		t.Errorf("expected 25, got %f", p)
	}
	if p := CalculateSpentPercentage(0, 0); p != 0 {
		t.Errorf("expected 0, got %f", p)
	}
	if p := CalculateSpentPercentage(10000, 10000); p != 100.0 {
		t.Errorf("expected 100, got %f", p)
	}
}

func TestCalculateTotalRemaining(t *testing.T) {
	cats := []BudgetCategory{
		{RemainingAmount: 6000},
		{RemainingAmount: 3000},
	}
	if total := CalculateTotalRemaining(cats); total != 9000 {
		t.Errorf("expected 9000, got %d", total)
	}
}

func TestReconstructFromDB(t *testing.T) {
	now := time.Now()
	b := ReconstructFromDB("b-1", "u-1", "", "Test", PeriodMonthly, now, now.AddDate(0, 1, 0),
		BStatusActive, 15000, 5000, 10000, "INR",
		[]BudgetCategory{{ID: "c-1", Category: "Food", BudgetedAmount: 10000, SpentAmount: 3000, RemainingAmount: 7000}},
		[]string{"important"}, now, now)

	if b.ID() != "b-1" { t.Errorf("expected b-1, got %s", b.ID()) }
	if b.Name() != "Test" { t.Errorf("expected Test, got %s", b.Name()) }
	if b.Status() != BStatusActive { t.Errorf("expected Active, got %s", b.Status()) }
	if len(b.Categories()) != 1 { t.Errorf("expected 1 category, got %d", len(b.Categories())) }
	if len(b.Tags()) != 1 { t.Errorf("expected 1 tag, got %d", len(b.Tags())) }
}
