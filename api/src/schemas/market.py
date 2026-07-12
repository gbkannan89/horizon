from pydantic import BaseModel
from typing import Optional, Dict, List


class GoldRateOut(BaseModel):
    rate_24k_per_gram: float
    rate_22k_per_gram: float
    rate_18k_per_gram: float
    currency: str = "INR"
    last_updated: str
    source: str


class StockQuoteOut(BaseModel):
    ticker: str
    name: str
    price: Optional[float] = None
    change: Optional[float] = None
    change_pct: Optional[float] = None
    currency: str = "INR"
    last_updated: str
    source: str = "yahoo.finance"


class BatchQuoteRequest(BaseModel):
    tickers: List[str]


class BatchQuoteOut(BaseModel):
    quotes: Dict[str, Optional[StockQuoteOut]]
