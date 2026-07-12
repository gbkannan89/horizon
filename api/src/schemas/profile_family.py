from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime, date as date_type


# ─── MEMBER SCHOOLING & PAYMENTS ──────────────────────────────────────────────
class SchoolingPaymentCreate(BaseModel):
    amount: float = Field(..., gt=0)
    paid_date: date_type
    receipt_ref: Optional[str] = None
    notes: Optional[str] = None

class SchoolingPaymentOut(BaseModel):
    id: int
    schooling_id: int
    amount: float
    paid_date: date_type
    receipt_ref: Optional[str] = None
    notes: Optional[str] = None
    created_at: datetime

class SchoolingCreate(BaseModel):
    institution_name: str
    fee_amount: float = Field(..., gt=0)
    fee_frequency: str = Field("monthly", pattern="^(monthly|quarterly|yearly)$")
    last_paid_date: Optional[date_type] = None
    next_due_date: Optional[date_type] = None
    notes: Optional[str] = None

class SchoolingOut(BaseModel):
    id: int
    member_id: int
    institution_name: str
    fee_amount: float
    fee_frequency: str
    last_paid_date: Optional[date_type] = None
    next_due_date: Optional[date_type] = None
    notes: Optional[str] = None
    created_at: datetime
    monthly_equivalent: float = 0.0


# ─── MEMBER HEALTH CHECKUPS ───────────────────────────────────────────────────
class CheckupCreate(BaseModel):
    checkup_type: str = Field("full_body", pattern="^(full_body|eye|dental|custom)$")
    frequency: str = Field("yearly", pattern="^(monthly|quarterly|yearly)$")
    recurring_cost: Optional[float] = 0.0
    last_checkup_date: Optional[date_type] = None
    next_due_date: Optional[date_type] = None
    notes: Optional[str] = None

class CheckupOut(BaseModel):
    id: int
    member_id: int
    checkup_type: str
    frequency: str
    recurring_cost: float
    last_checkup_date: Optional[date_type] = None
    next_due_date: Optional[date_type] = None
    notes: Optional[str] = None
    created_at: datetime


# ─── MEMBER MEDICINES ─────────────────────────────────────────────────────────
class MedicineCreate(BaseModel):
    medicine_name: str
    monthly_cost: float = Field(..., gt=0)
    purpose: Optional[str] = None
    is_regular: Optional[bool] = True
    prescribed_by: Optional[str] = None

class MedicineOut(BaseModel):
    id: int
    member_id: int
    medicine_name: str
    monthly_cost: float
    purpose: Optional[str] = None
    is_regular: bool
    prescribed_by: Optional[str] = None
    created_at: datetime


# ─── MEMBER VACCINATIONS ──────────────────────────────────────────────────────
class VaccinationCreate(BaseModel):
    vaccine_name: str
    frequency: str = Field("one_time", pattern="^(monthly|quarterly|yearly|one_time)$")
    recurring_cost: Optional[float] = 0.0
    last_vaccination_date: Optional[date_type] = None
    next_due_date: Optional[date_type] = None
    notes: Optional[str] = None

class VaccinationOut(BaseModel):
    id: int
    member_id: int
    vaccine_name: str
    frequency: str
    recurring_cost: float
    last_vaccination_date: Optional[date_type] = None
    next_due_date: Optional[date_type] = None
    notes: Optional[str] = None
    created_at: datetime


# ─── MEMBER EARNINGS ──────────────────────────────────────────────────────────
class EarningCreate(BaseModel):
    monthly_income: float = Field(..., ge=0)
    contribution_to_household: Optional[float] = 0.0
    occupation: Optional[str] = None

class EarningOut(BaseModel):
    id: int
    member_id: int
    monthly_income: float
    contribution_to_household: float
    occupation: Optional[str] = None
    created_at: datetime


# ─── FAMILY MEMBER ────────────────────────────────────────────────────────────
class FamilyMemberCreate(BaseModel):
    name: str
    dob: Optional[date_type] = None
    blood_group: Optional[str] = Field(None, pattern="^(A\+|A-|B\+|B-|AB\+|AB-|O\+|O-)$")
    relationship: str = Field("other", pattern="^(self|spouse|child|parent|other)$")
    avatar_color: Optional[str] = "#0D9488"
    is_self: Optional[bool] = False
    is_active: Optional[bool] = True

class FamilyMemberOut(BaseModel):
    id: int
    household_id: Optional[int] = None
    user_id: int
    name: str
    dob: Optional[date_type] = None
    blood_group: Optional[str] = None
    relationship: str
    avatar_color: str
    is_self: bool
    is_active: bool
    created_at: datetime
    # Nested fields
    schooling: List[SchoolingOut] = []
    checkups: List[CheckupOut] = []
    medicines: List[MedicineOut] = []
    vaccinations: List[VaccinationOut] = []
    earnings: Optional[EarningOut] = None
