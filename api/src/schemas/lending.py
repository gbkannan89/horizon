from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime, date as date_type


class LendingCreate(BaseModel):
    direction: str = Field(..., pattern="^(lent|borrowed)$")
    person_name: str
    amount: float = Field(..., gt=0)
    date_given: date_type
    promised_return_date: Optional[date_type] = None
    notes: Optional[str] = None


class LendingUpdate(BaseModel):
    person_name: Optional[str] = None
    amount: Optional[float] = Field(None, gt=0)
    date_given: Optional[date_type] = None
    promised_return_date: Optional[date_type] = None
    returned_amount: Optional[float] = Field(None, ge=0)
    status: Optional[str] = Field(None, pattern="^(open|partial|closed|overdue)$")
    notes: Optional[str] = None


class LendingOut(BaseModel):
    id: int
    user_id: int
    direction: str
    person_name: str
    amount: float
    date_given: date_type
    promised_return_date: Optional[date_type] = None
    actual_return_date: Optional[date_type] = None
    returned_amount: float
    remaining: float
    status: str
    notes: Optional[str] = None
    created_at: datetime


class LendingOverview(BaseModel):
    total_lent: float
    total_borrowed: float
    net_receivable: float
    overdue_count: int
    open_count: int
    closed_count: int
