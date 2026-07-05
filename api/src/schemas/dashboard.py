from pydantic import BaseModel
from typing import List, Optional
from datetime import datetime, date as date_type
from .goals import GoalOut
from .transactions import ExpenseOut, TrendPoint
from .bills import RecurringBillOut

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
