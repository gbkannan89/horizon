from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime, date as date_type

class RecurringBillCreate(BaseModel):
    name: str
    amount: float = Field(..., ge=0)
    frequency: str
    category: str
    bucket: str
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
