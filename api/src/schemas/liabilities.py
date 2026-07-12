from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime, date as date_type

class LiabilityCreate(BaseModel):
    type: str
    name: str
    outstanding: float = Field(..., ge=0)
    emi: float = Field(default=0.0, ge=0)
    interest_rate: float = Field(default=0.0, ge=0)
    tenure_months: Optional[int] = Field(default=None, ge=0)
    gold_grams: Optional[float] = None
    gold_carat: Optional[int] = None
    gold_items_count: Optional[int] = None
    loan_date: Optional[date_type] = None
    bank_name: Optional[str] = None

class LiabilityOut(BaseModel):
    id: int
    user_id: int
    type: str
    name: str
    outstanding: float
    emi: float
    interest_rate: float
    tenure_months: Optional[int]
    gold_grams: Optional[float] = None
    gold_carat: Optional[int] = None
    gold_items_count: Optional[int] = None
    loan_date: Optional[date_type] = None
    bank_name: Optional[str] = None
    created_at: datetime
