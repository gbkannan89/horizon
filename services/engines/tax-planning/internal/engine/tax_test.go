package engine

import "testing"

func TestComputeTax(t *testing.T) {
	tests := []struct {
		name     string
		income   int64
		slabs    []TaxSlab
		expected int64
	}{
		{"below slab", 200000, oldRegimeSlabs, 0},
		{"basic slab", 300000, oldRegimeSlabs, 2500},
		{"mid slab", 600000, oldRegimeSlabs, 32500},
		{"high slab", 1500000, oldRegimeSlabs, 262500},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := computeTax(tc.income, tc.slabs); got != tc.expected {
				t.Errorf("computeTax(%d) = %d, want %d", tc.income, got, tc.expected)
			}
		})
	}
}

func TestTotalIncome(t *testing.T) {
	input := TaxInput{GrossSalary: 500000, BusinessIncome: 100000, RentalIncome: 60000}
	if got := input.totalIncome(); got != 660000 {
		t.Errorf("totalIncome() = %d, want 660000", got)
	}
}

func TestTotalDeductions(t *testing.T) {
	input := TaxInput{Section80C: 150000, Section80D: 25000, NPSContribution: 50000}
	if got := input.totalDeductions(); got != 225000 {
		t.Errorf("totalDeductions() = %d, want 225000", got)
	}
}

func TestOldRegime(t *testing.T) {
	input := TaxInput{GrossSalary: 1200000, Section80C: 150000, Section80D: 25000}
	result := input.ComputeOldRegime()

	if result.Regime != RegimeOld {
		t.Errorf("expected old regime, got %s", result.Regime)
	}
	if result.GrossIncome != 1200000 {
		t.Errorf("expected 1200000, got %d", result.GrossIncome)
	}
	if result.Deductions != 175000 {
		t.Errorf("expected 175000 deductions, got %d", result.Deductions)
	}
	if result.TaxableIncome != 1025000 {
		t.Errorf("expected 1025000, got %d", result.TaxableIncome)
	}
	if result.TaxAmount <= 0 {
		t.Error("expected positive tax amount")
	}
	if result.CessAmount <= 0 {
		t.Error("expected positive cess amount")
	}
	if result.EffectiveRate <= 0 {
		t.Error("expected positive effective rate")
	}
}

func TestNewRegime(t *testing.T) {
	input := TaxInput{GrossSalary: 1200000}
	result := input.ComputeNewRegime()

	if result.Regime != RegimeNew {
		t.Errorf("expected new regime, got %s", result.Regime)
	}
	if result.GrossIncome != 1200000 {
		t.Errorf("expected 1200000, got %d", result.GrossIncome)
	}
	if result.Deductions != 0 {
		t.Errorf("expected 0 deductions, got %d", result.Deductions)
	}
}

func TestAnalyze(t *testing.T) {
	input := TaxInput{GrossSalary: 1500000, Section80C: 150000, Section80D: 50000, NPSContribution: 50000}
	result := input.Analyze()

	if result.Recommended != RegimeOld && result.Recommended != RegimeNew {
		t.Errorf("expected old or new regime, got %s", result.Recommended)
	}
	if result.OldRegime.TotalTax <= 0 {
		t.Error("expected positive old regime tax")
	}
	if result.NewRegime.TotalTax <= 0 {
		t.Error("expected positive new regime tax")
	}
}

func TestAnalyze_LowIncome(t *testing.T) {
	input := TaxInput{GrossSalary: 300000}
	result := input.Analyze()

	if result.OldRegime.TotalTax > 5000 {
		t.Errorf("expected small tax for low income old regime, got %d", result.OldRegime.TotalTax)
	}
	if result.NewRegime.TotalTax != 0 {
		t.Errorf("expected 0 tax for low income new regime, got %d", result.NewRegime.TotalTax)
	}
}

func TestAnalyze_NoDeductions(t *testing.T) {
	input := TaxInput{GrossSalary: 1000000}
	result := input.Analyze()

	if result.Recommended != RegimeNew {
		t.Errorf("expected new regime recommended when no deductions, got %s", result.Recommended)
	}
}

func TestSuggestions(t *testing.T) {
	input := TaxInput{GrossSalary: 1500000}
	result := input.Analyze()

	if len(result.Suggestions) == 0 {
		t.Error("expected suggestions for high income with no deductions")
	}
}

func TestComputeTax_NewRegimeSlabs(t *testing.T) {
	tests := []struct {
		name     string
		income   int64
		minTax   int64
		maxTax   int64
	}{
		{"below threshold", 300000, 0, 0},
		{"first slab", 500000, 8000, 12000},
		{"mid slab", 900000, 40000, 50000},
		{"upper slab", 1200000, 80000, 100000},
		{"top slab", 2000000, 280000, 330000},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := TaxInput{GrossSalary: tc.income}
			result := input.ComputeNewRegime()
			if result.TotalTax < tc.minTax || result.TotalTax > tc.maxTax {
				t.Errorf("income %d: expected tax between %d-%d, got %d (tax:%d cess:%d)",
					tc.income, tc.minTax, tc.maxTax, result.TotalTax, result.TaxAmount, result.CessAmount)
			}
		})
	}
}
