from pydantic import BaseModel, Field, EmailStr
from typing import Optional

class JoinHousehold(BaseModel):
    invite_code: str

class ContributingMemberCreate(BaseModel):
    name: str
    monthly_income: float = Field(..., ge=0)
    contribution_to_household: float = Field(..., ge=0)
    relationship: Optional[str] = None

class ContributingMemberUpdate(BaseModel):
    name: Optional[str] = None
    monthly_income: Optional[float] = Field(None, ge=0)
    contribution_to_household: Optional[float] = Field(None, ge=0)
    relationship: Optional[str] = None

class FamilyInviteRequest(BaseModel):
    email: EmailStr
    name: str
    relationship: Optional[str] = None
