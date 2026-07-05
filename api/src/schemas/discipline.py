from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime

class WishlistItemCreate(BaseModel):
    name: str
    amount: float = Field(..., ge=0)
    lock_duration_days: int = Field(default=30, ge=1)

class WishlistItemOut(BaseModel):
    id: int
    user_id: int
    name: str
    amount: float
    added_date: datetime
    unlock_date: datetime
    status: str
    created_at: datetime

class WishlistItemUpdate(BaseModel):
    status: str

class FinancialGuardrailsOut(BaseModel):
    emergency_fund_ratio: float
    emergency_target_amount: float
    emergency_current_amount: float
    housing_cost_ratio: float
    housing_status: str
    runway_status: str

class DebtRepaymentStrategyOut(BaseModel):
    snowball_months_to_freedom: int
    avalanche_months_to_freedom: int
    snowball_total_interest: float
    avalanche_total_interest: float
    recommended_strategy: str

class ZeroBasedBudgetOut(BaseModel):
    total_income: float
    total_allocated: float
    unallocated: float
    status: str

class PayYourselfFirstOut(BaseModel):
    target_percentage: float = 20.0
    target_amount: float
    actual_savings: float
    status: str
