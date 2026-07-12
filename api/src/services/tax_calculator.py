import logging

logger = logging.getLogger(__name__)

def estimate_tax(income_data: dict, regime: str) -> dict:
    """
    income_data: { gross_annual, basic, hra, pf_contrib, deductions_80c, deductions_80d }
    regime: 'old' or 'new'
    """
    gross_annual = float(income_data.get("gross_annual", 0.0))
    basic = float(income_data.get("basic", gross_annual * 0.4))
    hra = float(income_data.get("hra", gross_annual * 0.15))
    pf_contrib = float(income_data.get("pf_contrib", 0.0))
    deductions_80c = float(income_data.get("deductions_80c", 0.0))
    deductions_80d = float(income_data.get("deductions_80d", 0.0))

    if regime.lower() == "old":
        # Standard deduction: 50,000
        taxable = gross_annual - 50000.0
        
        # HRA Exemption: assume metro (50% basic) or normal, min HRA
        hra_exemption = min(hra, basic * 0.5)
        taxable -= hra_exemption
        
        # 80C: max 1.5L
        val_80c = min(150000.0, pf_contrib + deductions_80c)
        taxable -= val_80c
        
        # 80D: max 25,000
        val_80d = min(25000.0, deductions_80d)
        taxable -= val_80d
        
        taxable = max(0.0, taxable)
        
        # Slab calculation for old regime
        tax = 0.0
        if taxable <= 250000:
            tax = 0.0
        elif taxable <= 500000:
            tax = (taxable - 250000) * 0.05
        elif taxable <= 1000000:
            tax = 12500 + (taxable - 500000) * 0.20
        else:
            tax = 112500 + (taxable - 1000000) * 0.30
            
        # Rebate under section 87A (if taxable income <= 5L)
        if taxable <= 500000:
            tax = 0.0
            
    else:  # New regime (FY25-26/26-27 update)
        # Standard deduction: 75,000
        taxable = max(0.0, gross_annual - 75000.0)
        
        # New regime slabs:
        # Up to 3L: Nil
        # 3L to 7L: 5%
        # 7L to 10L: 10%
        # 10L to 12L: 15%
        # 12L to 15L: 20%
        # Above 15L: 30%
        tax = 0.0
        if taxable <= 300000:
            tax = 0.0
        elif taxable <= 700000:
            tax = (taxable - 300000) * 0.05
        elif taxable <= 1000000:
            tax = 20000 + (taxable - 700000) * 0.10
        elif taxable <= 1200000:
            tax = 50000 + (taxable - 1000000) * 0.15
        elif taxable <= 1500000:
            tax = 80000 + (taxable - 1200000) * 0.20
        else:
            tax = 140000 + (taxable - 1500000) * 0.30
            
        # Rebate under section 87A for new regime (if taxable income <= 7L, tax is zero)
        if taxable <= 700000:
            tax = 0.0

    cess = tax * 0.04
    total_tax = tax + cess
    effective_rate = (total_tax / gross_annual * 100) if gross_annual > 0 else 0.0

    return {
        "regime": regime,
        "taxable_income": taxable,
        "tax_before_cess": tax,
        "cess": cess,
        "total_tax": total_tax,
        "effective_rate": round(effective_rate, 2)
    }
