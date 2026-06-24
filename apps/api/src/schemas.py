from pydantic import BaseModel, EmailStr, Field
from typing import Optional, List
from datetime import datetime, date as date_type

class UserRegister(BaseModel):
    email: EmailStr
    password: str = Field(..., min_length=6)
    name: Optional[str] = None
    phone: Optional[str] = None
    user_type: str = Field(default="salaried") # 'salaried' or 'self-employed'
    risk_profile: str = Field(default="moderate") # 'conservative', 'moderate', 'aggressive'

class UserLogin(BaseModel):
    email: EmailStr
    password: str

class DeleteAccountRequest(BaseModel):
    password: str

class UserOut(BaseModel):
    id: int
    email: EmailStr
    name: Optional[str] = None
    phone: Optional[str] = None
    user_type: str
    risk_profile: str
    household_id: Optional[int] = None

    class Config:
        from_attributes = True

class Token(BaseModel):
    access_token: str
    refresh_token: str
    token_type: str = "bearer"
    user: UserOut

class TokenRefresh(BaseModel):
    refresh_token: str

class TokenData(BaseModel):
    user_id: int
    email: str

class JoinHousehold(BaseModel):
    invite_code: str

class ContributingMemberCreate(BaseModel):
    name: str
    monthly_income: float = Field(..., ge=0)
    contribution_to_household: float = Field(..., ge=0)
    relationship: Optional[str] = None

class ContributingMemberUpdate(BaseModel):
    name: Optional[str] = None
    monthly_income: Optional[float] = Field(None, ge=0)
    contribution_to_household: Optional[float] = Field(None, ge=0)
    relationship: Optional[str] = None

# Income Schemas
class IncomeCreate(BaseModel):
    label: Optional[str] = None   # user-friendly name, e.g. 'My Salary', 'Partner Income'
    type: str  # 'salary', 'business', 'passive'
    amount: float = Field(..., ge=0)
    frequency: str = Field(default="monthly")  # 'monthly', 'annual'

class IncomeOut(BaseModel):
    id: int
    user_id: int
    label: Optional[str] = None
    type: str
    amount: float
    frequency: str
    created_at: datetime

# Asset Schemas
class AssetCreate(BaseModel):
    type: str  # 'bank', 'fd', 'rd', 'physical'
    subtype: Optional[str] = None  # 'property', 'gold', 'vehicle', etc.
    name: str
    amount: float = Field(..., ge=0)
    interest_rate: float = Field(default=0.0, ge=0)
    start_date: Optional[date_type] = None
    maturity_date: Optional[date_type] = None
    is_liability: bool = False
    generates_income: bool = False
    income_frequency: Optional[str] = None
    purchase_price: Optional[float] = None
    purchase_date: Optional[date_type] = None

class AssetOut(BaseModel):
    id: int
    user_id: int
    type: str
    subtype: Optional[str]
    name: str
    amount: float
    interest_rate: float
    start_date: Optional[date_type]
    maturity_date: Optional[date_type]
    is_liability: bool
    generates_income: bool
    income_frequency: Optional[str]
    purchase_price: Optional[float] = None
    purchase_date: Optional[date_type] = None
    created_at: datetime

# Vehicle Schemas
class VehicleCreate(BaseModel):
    make_model: str
    purchase_cost: float = Field(..., ge=0)
    insurance_renewal_date: Optional[date_type] = None

class VehicleOut(BaseModel):
    id: int
    user_id: int
    make_model: str
    purchase_cost: float
    insurance_renewal_date: Optional[date_type]
    created_at: datetime

# Liability Schemas
class LiabilityCreate(BaseModel):
    type: str  # 'home_loan', 'car_loan', 'personal_loan', 'credit_card'
    name: str
    outstanding: float = Field(..., ge=0)
    emi: float = Field(default=0.0, ge=0)
    interest_rate: float = Field(default=0.0, ge=0)
    tenure_months: Optional[int] = Field(default=None, ge=0)

class LiabilityOut(BaseModel):
    id: int
    user_id: int
    type: str
    name: str
    outstanding: float
    emi: float
    interest_rate: float
    tenure_months: Optional[int]
    created_at: datetime

# Recurring Bill Schemas
class RecurringBillCreate(BaseModel):
    name: str
    amount: float = Field(..., ge=0)
    frequency: str  # 'monthly', 'quarterly', 'half-yearly', 'yearly'
    category: str  # 'rent', 'utilities', 'insurance', 'education', etc.
    bucket: str  # 'Needs', 'Wants', 'Savings'
    due_day: int = Field(default=1, ge=1, le=31)
    is_active: bool = True
    is_subscription: bool = False
    is_emi: bool = False
    emi_total_months: Optional[int] = None
    emi_months_paid: int = 0
    start_date: Optional[date_type] = None

class RecurringBillOut(BaseModel):
    id: int
    user_id: int
    name: str
    amount: float
    frequency: str
    category: str
    bucket: str
    due_day: int
    monthly_equivalent: float
    is_active: bool
    is_subscription: bool
    is_emi: bool
    emi_total_months: Optional[int]
    emi_months_paid: Optional[int] = None
    start_date: Optional[date_type] = None
    created_at: datetime

# Goal Schemas
class GoalCreate(BaseModel):
    name: str
    target_amount: float = Field(..., ge=0)
    current_amount: float = Field(default=0.0, ge=0)
    status: str = Field(default="On Track")
    color: str = Field(default="#059669")

class GoalOut(BaseModel):
    id: int
    user_id: int
    name: str
    target_amount: float
    current_amount: float
    status: str
    color: str
    created_at: datetime
    projected_completion_date: Optional[str] = None
    monthly_saving_needed: Optional[float] = None

# Expense Schemas
class UploadSummary(BaseModel):
    inserted: int
    skipped: int
    duplicates: int
    total_parsed: int

class ExpenseCreate(BaseModel):
    name: str
    amount: float = Field(..., ge=0)
    category: str
    bucket: str
    icon: Optional[str] = None
    date: date_type

class ExpenseOut(BaseModel):
    id: int
    user_id: int
    name: str
    amount: float
    category: str
    bucket: str
    icon: Optional[str]
    date: date_type
    created_at: datetime

class TrendPoint(BaseModel):
    month: str
    needs: float = 0
    wants: float = 0
    savings: float = 0

# Dashboard Aggregated Schema
class DashboardOverview(BaseModel):
    net_worth: float
    net_worth_change: float
    net_worth_change_period: str
    fin_score: int
    total_income: float
    total_spent: float
    total_left: float
    needs_spent: float
    wants_spent: float
    savings_spent: float
    needs_budget: float
    wants_budget: float
    savings_budget: float
    goals: List[GoalOut]
    recent_expenses: List[ExpenseOut]
    upcoming_bills: List[RecurringBillOut]
    spending_trend: List[TrendPoint] = []

# Wishlist Schemas
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

# Discipline Aggregates
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
    status: str # 'Zero-Based', 'Surplus', 'Deficit'

class PayYourselfFirstOut(BaseModel):
    target_percentage: float = 20.0
    target_amount: float
    actual_savings: float
    status: str

# ── Advisor Schemas ───────────────────────────────────────────────────────────
class SimulationInput(BaseModel):
    scenario_name: str
    scenario_type: str
    target_amount: Optional[float] = None
    loan_amount: Optional[float] = None
    interest_rate: Optional[float] = None
    tenure_months: Optional[int] = None
    monthly_emi: Optional[float] = None
    sip_amount: Optional[float] = None
    income_change: Optional[float] = None

class SimulationOut(BaseModel):
    id: int
    scenario_name: str
    scenario_type: str
    input_params: dict
    results: dict
    created_at: datetime

class NudgeOut(BaseModel):
    id: int
    category: str
    severity: str
    title: str
    message: str
    action_label: Optional[str]
    action_link: Optional[str]
    is_dismissed: bool
    created_at: datetime

class DebtOptimizerOut(BaseModel):
    total_outstanding: float
    total_monthly_emis: float
    total_interest_paid: float
    snowball: dict
    avalanche: dict
    recommended_strategy: str
    estimated_freedom_months: int
    total_interest_savable: float

class SubscriptionInsightOut(BaseModel):
    id: int
    bill_id: int
    name: str
    amount: float
    monthly_cost: float
    annual_cost: float
    status: str
    savings_opportunity: float
    notes: Optional[str]

# Insurance Schemas
class InsuranceCreate(BaseModel):
    type: str
    provider: str
    policy_name: Optional[str] = None
    premium_amount: float = Field(..., ge=0)
    premium_frequency: str
    coverage_amount: Optional[float] = None
    renewal_date: Optional[date_type] = None

class InsuranceOut(BaseModel):
    id: int
    user_id: int
    type: str
    provider: str
    policy_name: Optional[str]
    premium_amount: float
    premium_frequency: str
    coverage_amount: Optional[float]
    renewal_date: Optional[date_type]
    created_at: datetime # 'On Track', 'Behind'
