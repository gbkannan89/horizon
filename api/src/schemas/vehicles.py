from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime, date as date_type


# ─── VEHICLE SERVICE LOG ──────────────────────────────────────────────────────
class VehicleServiceCreate(BaseModel):
    service_date: date_type
    service_type: Optional[str] = Field("regular", pattern="^(regular|major|repair|annual)$")
    description: Optional[str] = None
    cost: float = Field(..., gt=0)
    km_at_service: Optional[float] = None
    service_center: Optional[str] = None
    next_service_km: Optional[float] = None

class VehicleServiceOut(BaseModel):
    id: int
    vehicle_id: int
    service_date: date_type
    service_type: str
    description: Optional[str] = None
    cost: float
    km_at_service: Optional[float] = None
    service_center: Optional[str] = None
    next_service_km: Optional[float] = None
    created_at: datetime


# ─── VEHICLE FUEL LOG ─────────────────────────────────────────────────────────
class VehicleFuelCreate(BaseModel):
    fill_date: date_type
    amount: float = Field(..., gt=0)
    liters: float = Field(..., gt=0)
    km_at_fill: Optional[float] = None
    price_per_liter: Optional[float] = None
    is_full_tank: Optional[bool] = True

class VehicleFuelOut(BaseModel):
    id: int
    vehicle_id: int
    fill_date: date_type
    amount: float
    liters: float
    km_at_fill: Optional[float] = None
    price_per_liter: Optional[float] = None
    is_full_tank: bool
    created_at: datetime


# ─── VEHICLE LOAN ─────────────────────────────────────────────────────────────
class VehicleLoanCreate(BaseModel):
    bank_name: Optional[str] = None
    loan_amount: float = Field(..., gt=0)
    interest_rate: float = Field(..., gt=0)
    tenure_months: int = Field(..., gt=0)
    emi: float = Field(..., gt=0)
    start_date: date_type
    emi_paid: Optional[int] = 0
    outstanding: float = Field(..., gte=0)

class VehicleLoanOut(BaseModel):
    id: int
    vehicle_id: int
    bank_name: Optional[str] = None
    loan_amount: float
    interest_rate: float
    tenure_months: int
    emi: float
    start_date: date_type
    emi_paid: int
    outstanding: float
    created_at: datetime
