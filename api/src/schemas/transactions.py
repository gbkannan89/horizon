from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime, date as date_type

class UploadSummary(BaseModel):
    inserted: int
    skipped: int
    duplicates: int
    total_parsed: int

class ExpenseCreate(BaseModel):
    name: str
    amount: float = Field(..., ge=0)
    category: str
    bucket: str
    icon: Optional[str] = None
    date: date_type

class ExpenseOut(BaseModel):
    id: int
    user_id: int
    name: str
    amount: float
    category: str
    bucket: str
    icon: Optional[str]
    date: date_type
    created_at: datetime

class TrendPoint(BaseModel):
    month: str
    needs: float = 0
    wants: float = 0
    savings: float = 0
