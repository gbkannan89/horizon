from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import date as date_type

class RecurringDetectionOut(BaseModel):
    name: str
    amount: float
    frequency: str
    confidence: float
    occurrences: int
    first_date: str
    last_date: str
    category: str
    bucket: str
    is_subscription: bool
    amount_variation: float

class RecurringConfirmIn(BaseModel):
    name: str
    amount: float
    frequency: str
    category: str = "Misc"
    bucket: str = "Wants"
    due_day: int = 1
    is_subscription: bool = False

class CategoryTrendOut(BaseModel):
    category: str
    current_amount: float
    previous_amount: float
    change_pct: float

class AnomalyOut(BaseModel):
    id: int
    name: str
    amount: float
    category: str
    date: str
    z_score: float
    reason: str

class ForecastOut(BaseModel):
    current_spent: float
    projected_total: float
    average_monthly: float
    days_elapsed: int
    days_in_month: int

class SpendingPatternOut(BaseModel):
    category_breakdown: list = []
    category_trends: List[CategoryTrendOut]
    weekday_distribution: dict = {}
    weekend_boost_pct: float
    anomalies: List[AnomalyOut]
    forecast: ForecastOut
    top_merchants: list = []

class SubscriptionCandidateOut(BaseModel):
    name: str
    amount: float
    monthly_cost: float
    annual_cost: float
    frequency: str
    confidence: float
    is_new: bool
    savings_opportunity: float
    category: str
    bucket: str
    occurrences: int
    last_date: str

class LapsedSubscriptionOut(BaseModel):
    bill_id: int
    name: str
    monthly_cost: float
    annual_cost: float
    status: str
    savings_opportunity: float
    days_since_last_charge: int

class FullAnalysisOut(BaseModel):
    recurring_detections: List[RecurringDetectionOut]
    subscription_candidates: List[SubscriptionCandidateOut]
    lapsed_subscriptions: List[LapsedSubscriptionOut]
    spending_patterns: SpendingPatternOut
    nudges_generated: int

class UploadAnalysisOut(BaseModel):
    inserted: int
    skipped: int
    duplicates: int
    total_parsed: int
    analysis: Optional[FullAnalysisOut] = None
