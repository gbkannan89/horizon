from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime, date as date_type


class BudgetItemCreate(BaseModel):
    category: str = Field(..., pattern="^(housing|food|transport|health|education|entertainment|savings|debt|other)$")
    label: str
    amount: float = Field(..., gt=0)
    frequency: Optional[str] = "monthly"
    bucket: str = Field(..., pattern="^(Needs|Wants|Savings)$")
    source: Optional[str] = "manual"
    source_id: Optional[int] = None
    source_label: Optional[str] = None
    is_committed: Optional[bool] = True
    sort_order: Optional[int] = 0


class BudgetItemUpdate(BaseModel):
    category: Optional[str] = Field(None, pattern="^(housing|food|transport|health|education|entertainment|savings|debt|other)$")
    label: Optional[str] = None
    amount: Optional[float] = Field(None, gt=0)
    frequency: Optional[str] = None
    bucket: Optional[str] = Field(None, pattern="^(Needs|Wants|Savings)$")
    is_committed: Optional[bool] = None
    sort_order: Optional[int] = None


class BudgetItemOut(BaseModel):
    id: int
    budget_plan_id: int
    category: str
    label: str
    amount: float
    frequency: str
    bucket: str
    source: str
    source_id: Optional[int] = None
    source_label: Optional[str] = None
    is_committed: bool
    sort_order: int
    created_at: datetime


class BudgetPlanOut(BaseModel):
    id: int
    user_id: int
    month: int
    year: int
    total_budgeted: float
    notes: Optional[str] = None
    created_at: datetime
    updated_at: datetime
    items: List[BudgetItemOut] = []


class BudgetSummary(BaseModel):
    total: float
    needs: float
    wants: float
    savings: float


class CompareItem(BaseModel):
    label: str
    budgeted: float
    actual: float
    remaining: float
    bucket: str


class BudgetComparisonOut(BaseModel):
    budgeted: BudgetSummary
    actual: BudgetSummary
    remaining: BudgetSummary
    items: List[CompareItem] = []
