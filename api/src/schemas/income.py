from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime

class IncomeCreate(BaseModel):
    label: Optional[str] = None
    type: str
    amount: float = Field(..., ge=0)
    frequency: str = Field(default="monthly")

class IncomeOut(BaseModel):
    id: int
    user_id: int
    label: Optional[str] = None
    type: str
    amount: float
    frequency: str
    created_at: datetime
