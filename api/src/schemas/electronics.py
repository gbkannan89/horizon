from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime, date as date_type


# ─── ELECTRONIC SERVICE LOG ──────────────────────────────────────────────────
class ElectronicServiceCreate(BaseModel):
    service_date: date_type
    service_type: Optional[str] = Field("repair", pattern="^(screen_replacement|battery|repair|annual_service)$")
    description: Optional[str] = None
    cost: float = Field(..., gt=0)
    service_center: Optional[str] = None

class ElectronicServiceOut(BaseModel):
    id: int
    electronic_id: int
    service_date: date_type
    service_type: str
    description: Optional[str] = None
    cost: float
    service_center: Optional[str] = None
    created_at: datetime


# ─── ELECTRONIC EMI ───────────────────────────────────────────────────────────
class ElectronicEmiCreate(BaseModel):
    bank_name: str
    emi_amount: float = Field(..., gt=0)
    interest_rate: Optional[float] = 0.0
    total_months: int = Field(..., gt=0)
    months_paid: Optional[int] = 0
    start_date: date_type
    start_immediately: Optional[bool] = True

class ElectronicEmiOut(BaseModel):
    id: int
    electronic_id: int
    bank_name: str
    emi_amount: float
    interest_rate: float
    total_months: int
    months_paid: int
    start_date: date_type
    start_immediately: bool
    created_at: datetime


# ─── ELECTRONICS ──────────────────────────────────────────────────────────────
class ElectronicCreate(BaseModel):
    name: str
    category: str = Field("mobile", pattern="^(mobile|laptop|tablet|tv|appliance|other)$")
    brand: Optional[str] = None
    model: Optional[str] = None
    purchase_date: date_type
    purchase_amount: float = Field(..., gt=0)
    warranty_years: Optional[int] = 1
    expected_life_years: Optional[int] = 3
    notes: Optional[str] = None

class ElectronicOut(BaseModel):
    id: int
    user_id: int
    name: str
    category: str
    brand: Optional[str] = None
    model: Optional[str] = None
    purchase_date: date_type
    purchase_amount: float
    warranty_years: int
    warranty_expiry_date: Optional[date_type] = None
    expected_life_years: int
    current_value: float
    notes: Optional[str] = None
    created_at: datetime
    # Computes
    warranty_status: str = "active" # active, expired, near_expiry
    warranty_days_remaining: int = 0
    total_service_cost: float = 0.0
    emi: Optional[ElectronicEmiOut] = None
