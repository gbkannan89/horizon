package engine

func (t *TaxInput) totalIncome() int64 {
	return t.GrossSalary + t.BusinessIncome + t.CapitalGains + t.RentalIncome + t.OtherIncome
}

func (t *TaxInput) totalDeductions() int64 {
	return t.Section80C + t.Section80D + t.NPSContribution
}

func computeTax(income int64, slabs []TaxSlab) int64 {
	var tax int64
	remaining := income
	for _, slab := range slabs {
		if remaining <= 0 {
			break
		}
		var taxable int64
		if slab.MaxIncome == 0 {
			taxable = remaining
		} else {
			bandSize := slab.MaxIncome - slab.MinIncome
			taxable = bandSize
			if taxable > remaining {
				taxable = remaining
			}
		}
		// Use integer math: taxable * rate / 100
		tax += taxable * int64(slab.Rate) / 100
		remaining -= taxable
	}
	return tax
}

func min(a, b int64) int64 {
	if a < b { return a }
	return b
}

var oldRegimeSlabs = []TaxSlab{
	{MinIncome: 0, MaxIncome: 250000, Rate: 0},
	{MinIncome: 250000, MaxIncome: 500000, Rate: 5},
	{MinIncome: 500000, MaxIncome: 1000000, Rate: 20},
	{MinIncome: 1000000, MaxIncome: 0, Rate: 30},
}

var newRegimeSlabs = []TaxSlab{
	{MinIncome: 0, MaxIncome: 300000, Rate: 0},
	{MinIncome: 300000, MaxIncome: 600000, Rate: 5},
	{MinIncome: 600000, MaxIncome: 900000, Rate: 10},
	{MinIncome: 900000, MaxIncome: 1200000, Rate: 15},
	{MinIncome: 1200000, MaxIncome: 1500000, Rate: 20},
	{MinIncome: 1500000, MaxIncome: 0, Rate: 30},
}

func (t *TaxInput) ComputeOldRegime() TaxComputation {
	gross := t.totalIncome()
	deductions := t.totalDeductions()
	taxable := gross - deductions
	if taxable < 0 { taxable = 0 }

	tax := computeTax(taxable, oldRegimeSlabs)
	cess := int64(float64(tax) * 4.0 / 100.0)

	effective := 0.0
	if gross > 0 {
		effective = float64(tax+cess) / float64(gross) * 100.0
	}

	return TaxComputation{
		Regime: RegimeOld, GrossIncome: gross,
		Deductions: deductions, TaxableIncome: taxable,
		TaxAmount: tax, CessAmount: cess,
		TotalTax: tax + cess, EffectiveRate: effective,
	}
}

func (t *TaxInput) ComputeNewRegime() TaxComputation {
	gross := t.totalIncome()
	taxable := gross
	if taxable < 0 { taxable = 0 }

	tax := computeTax(taxable, newRegimeSlabs)
	cess := int64(float64(tax) * 4.0 / 100.0)

	effective := 0.0
	if gross > 0 {
		effective = float64(tax+cess) / float64(gross) * 100.0
	}

	return TaxComputation{
		Regime: RegimeNew, GrossIncome: gross,
		Deductions: 0, TaxableIncome: taxable,
		TaxAmount: tax, CessAmount: cess,
		TotalTax: tax + cess, EffectiveRate: effective,
	}
}

func (t *TaxInput) Analyze() TaxAnalysis {
	old := t.ComputeOldRegime()
	newCalc := t.ComputeNewRegime()

	analysis := TaxAnalysis{
		OldRegime: old,
		NewRegime: newCalc,
	}

	if old.TotalTax <= newCalc.TotalTax {
		analysis.Recommended = RegimeOld
		analysis.Savings = newCalc.TotalTax - old.TotalTax
	} else {
		analysis.Recommended = RegimeNew
		analysis.Savings = old.TotalTax - newCalc.TotalTax
	}
	if analysis.Savings < 0 { analysis.Savings = 0 }

	if t.Section80C < 150000 && old.TotalTax > 0 {
		analysis.Suggestions = append(analysis.Suggestions,
			"Maximise Section 80C investments up to ₹1,50,000 for additional tax savings")
	}
	if t.NPSContribution < 50000 && old.TotalTax > 0 {
		analysis.Suggestions = append(analysis.Suggestions,
			"Consider NPS contribution under Section 80CCD(1B) for extra ₹50,000 deduction")
	}
	if t.Section80D == 0 && old.TotalTax > 0 {
		analysis.Suggestions = append(analysis.Suggestions,
			"Health insurance premium under Section 80D can reduce taxable income")
	}

	return analysis
}
