from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime, date as date_type

class AssetCreate(BaseModel):
    type: str
    subtype: Optional[str] = None
    name: str
    amount: float = Field(..., ge=0)
    interest_rate: float = Field(default=0.0, ge=0)
    start_date: Optional[date_type] = None
    maturity_date: Optional[date_type] = None
    is_liability: bool = False
    generates_income: bool = False
    income_frequency: Optional[str] = None
    purchase_price: Optional[float] = None
    purchase_date: Optional[date_type] = None

class AssetOut(BaseModel):
    id: int
    user_id: int
    type: str
    subtype: Optional[str]
    name: str
    amount: float
    interest_rate: float
    start_date: Optional[date_type]
    maturity_date: Optional[date_type]
    is_liability: bool
    generates_income: bool
    income_frequency: Optional[str]
    purchase_price: Optional[float] = None
    purchase_date: Optional[date_type] = None
    created_at: datetime

class VehicleCreate(BaseModel):
    make_model: str
    purchase_cost: float = Field(..., ge=0)
    insurance_renewal_date: Optional[date_type] = None

class VehicleOut(BaseModel):
    id: int
    user_id: int
    make_model: str
    purchase_cost: float
    insurance_renewal_date: Optional[date_type]
    created_at: datetime
