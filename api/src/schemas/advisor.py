from pydantic import BaseModel
from typing import Optional, List
from datetime import datetime

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
