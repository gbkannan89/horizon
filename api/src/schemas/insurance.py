from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime, date as date_type

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
    created_at: datetime
