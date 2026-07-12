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
    years_of_deposit: Optional[int] = None
    is_liability: bool = False
    is_emergency: bool = False
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
    years_of_deposit: Optional[int] = None
    is_liability: bool
    is_emergency: bool = False
    generates_income: bool
    income_frequency: Optional[str]
    purchase_price: Optional[float] = None
    purchase_date: Optional[date_type] = None
    created_at: datetime

class VehicleCreate(BaseModel):
    make_model: str
    purchase_cost: float = Field(..., ge=0)
    insurance_renewal_date: Optional[date_type] = None
    model_year: Optional[int] = None
    purchase_year: Optional[int] = None
    fuel_type: Optional[str] = None
    mileage_kmpl: Optional[float] = None
    fuel_cost_total: Optional[float] = 0.0
    km_driven: Optional[float] = 0.0
    insurance_idv: Optional[float] = None
    insurance_renewal_amount: Optional[float] = None
    registration_number: Optional[str] = None

class VehicleOut(BaseModel):
    id: int
    user_id: int
    make_model: str
    purchase_cost: float
    insurance_renewal_date: Optional[date_type] = None
    model_year: Optional[int] = None
    purchase_year: Optional[int] = None
    fuel_type: Optional[str] = None
    mileage_kmpl: Optional[float] = None
    fuel_cost_total: float
    km_driven: float
    insurance_idv: Optional[float] = None
    insurance_renewal_amount: Optional[float] = None
    registration_number: Optional[str] = None
    created_at: datetime
    # Computes
    cost_per_km: float = 0.0
    suggested_idv: float = 0.0
    total_service_cost: float = 0.0
    next_service: Optional[dict] = None
