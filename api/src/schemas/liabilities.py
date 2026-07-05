from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime

class LiabilityCreate(BaseModel):
    type: str
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
