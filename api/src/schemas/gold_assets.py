from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime, date as date_type


class GoldAssetCreate(BaseModel):
    carat: int = Field(..., ge=18, le=24)
    grams: float = Field(..., gt=0)
    purchase_price_per_gram: float = Field(..., ge=0)
    purchase_date: date_type
    notes: Optional[str] = None


class GoldAssetUpdate(BaseModel):
    carat: Optional[int] = Field(None, ge=18, le=24)
    grams: Optional[float] = Field(None, gt=0)
    purchase_price_per_gram: Optional[float] = None
    purchase_date: Optional[date_type] = None
    notes: Optional[str] = None


class GoldAssetOut(BaseModel):
    id: int
    user_id: int
    carat: int
    grams: float
    purchase_price_per_gram: float
    purchase_date: date_type
    current_value: Optional[float] = None
    total_purchase_cost: Optional[float] = None
    return_pct: Optional[float] = None
    last_current_value: Optional[float] = None
    notes: Optional[str] = None
    created_at: datetime
