from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime


class SalaryDetailCreate(BaseModel):
    company_name: Optional[str] = None
    from_year: int
    to_year: Optional[int] = None
    is_current: Optional[bool] = False
    fixed_pay: Optional[float] = None
    basic_pay: Optional[float] = None
    hra: Optional[float] = None
    lta: Optional[float] = None
    pf_employee: Optional[float] = None
    pf_employer: Optional[float] = None
    special_allowance: Optional[float] = None
    meal_card: Optional[float] = None
    variable_pay_percentage: Optional[float] = None
    variable_pay_amount: Optional[float] = None


class SalaryDetailUpdate(BaseModel):
    company_name: Optional[str] = None
    from_year: Optional[int] = None
    to_year: Optional[int] = None
    is_current: Optional[bool] = None
    fixed_pay: Optional[float] = None
    basic_pay: Optional[float] = None
    hra: Optional[float] = None
    lta: Optional[float] = None
    pf_employee: Optional[float] = None
    pf_employer: Optional[float] = None
    special_allowance: Optional[float] = None
    meal_card: Optional[float] = None
    variable_pay_percentage: Optional[float] = None
    variable_pay_amount: Optional[float] = None


class SalaryDetailOut(BaseModel):
    id: int
    income_id: int
    company_name: Optional[str] = None
    from_year: int
    to_year: Optional[int] = None
    is_current: bool
    fixed_pay: Optional[float] = None
    basic_pay: Optional[float] = None
    hra: Optional[float] = None
    lta: Optional[float] = None
    pf_employee: Optional[float] = None
    pf_employer: Optional[float] = None
    special_allowance: Optional[float] = None
    meal_card: Optional[float] = None
    variable_pay_percentage: Optional[float] = None
    variable_pay_amount: Optional[float] = None
    gross_annual: Optional[float] = None
    monthly_in_hand: Optional[float] = None
    created_at: datetime


class SalaryGrowthItem(BaseModel):
    year: int
    gross_annual: float
    company: Optional[str] = None
    growth_pct: Optional[float] = None


class SalaryGrowthOut(BaseModel):
    growth_data: List[SalaryGrowthItem] = []
    avg_growth_pct: float
    total_growth_pct: float
    time_period_years: int
