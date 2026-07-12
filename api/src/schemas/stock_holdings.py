from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime, date as date_type


class StockHoldingCreate(BaseModel):
    ticker: str
    exchange: str = "NSE"
    quantity: int = Field(..., gt=0)
    avg_purchase_price: float = Field(..., ge=0)
    purchase_date: date_type
    total_invested: float


class StockHoldingUpdate(BaseModel):
    quantity: Optional[int] = Field(None, gt=0)
    avg_purchase_price: Optional[float] = None
    purchase_date: Optional[date_type] = None
    total_invested: Optional[float] = None


class StockHoldingOut(BaseModel):
    id: int
    user_id: int
    ticker: str
    exchange: str
    name: Optional[str] = None
    quantity: int
    avg_purchase_price: float
    purchase_date: date_type
    total_invested: float
    current_value: Optional[float] = None
    return_pct: Optional[float] = None
    last_current_price: Optional[float] = None
    created_at: datetime
