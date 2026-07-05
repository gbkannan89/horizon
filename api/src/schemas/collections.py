from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime, date as date_type

class CollectionMemberCreate(BaseModel):
    name: str
    expected_amount: float = Field(..., ge=0)
    paid_amount: float = 0
    paid_date: Optional[date_type] = None
    status: str = "pending"
    notes: Optional[str] = None

class CollectionCreate(BaseModel):
    label: str
    description: Optional[str] = None
    members: List[CollectionMemberCreate] = []

class CollectionMemberOut(BaseModel):
    id: int
    collection_id: int
    name: str
    expected_amount: float
    paid_amount: float
    paid_date: Optional[date_type]
    status: str
    notes: Optional[str]
    created_at: datetime

class CollectionOut(BaseModel):
    id: int
    user_id: int
    label: str
    description: Optional[str]
    total_expected: float
    total_collected: float
    status: str
    created_at: datetime
    member_count: int = 0
    paid_count: int = 0
    members: List[CollectionMemberOut] = []

class MemberPayment(BaseModel):
    paid_amount: float = Field(..., ge=0)
    paid_date: date_type
