from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime, date as date_type


class PFAssetCreate(BaseModel):
    start_date: date_type
    retirement_age: int = Field(default=58, ge=18, le=99)
    current_balance: float = Field(default=0, ge=0)
    monthly_contribution: float = Field(default=0, ge=0)
    employer_contribution: Optional[float] = None
    interest_rate: float = Field(default=8.25, ge=0, le=100)
    user_date_of_birth: Optional[date_type] = None
    notes: Optional[str] = None


class PFAssetUpdate(BaseModel):
    start_date: Optional[date_type] = None
    retirement_age: Optional[int] = Field(None, ge=18, le=99)
    current_balance: Optional[float] = Field(None, ge=0)
    monthly_contribution: Optional[float] = Field(None, ge=0)
    employer_contribution: Optional[float] = None
    interest_rate: Optional[float] = Field(None, ge=0, le=100)
    notes: Optional[str] = None


class PFAssetOut(BaseModel):
    id: int
    user_id: int
    start_date: date_type
    retirement_age: int
    current_balance: float
    monthly_contribution: float
    employer_contribution: Optional[float] = None
    interest_rate: float
    user_date_of_birth: Optional[date_type] = None
    projected_corpus: Optional[float] = None
    years_until_retirement: Optional[int] = None
    notes: Optional[str] = None
    created_at: datetime


class PFProjectionOut(BaseModel):
    current_balance: float
    monthly_contribution: float
    employer_contribution: Optional[float]
    interest_rate: float
    retirement_age: int
    years_until_retirement: int
    projected_corpus: float
    total_contributions: float
    total_interest_earned: float
