import logging
from datetime import datetime
from typing import List
from fastapi import APIRouter, Query, HTTPException, status
from ..services.market_data import (
    get_gold_rates,
    get_stock_quote,
    get_batch_stock_quotes
)
from ..schemas.market import (
    GoldRateOut,
    StockQuoteOut,
    BatchQuoteRequest,
    BatchQuoteOut
)

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/market", tags=["market"])


@router.get("/gold-rate", response_model=GoldRateOut)
def gold_rate():
    try:
        rates = get_gold_rates()
        return GoldRateOut(
            **rates,
            last_updated=datetime.now().isoformat(),
            source="goldprice.org"
        )
    except Exception as e:
        logger.error(f"Failed to fetch gold rate: {e}")
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Gold rate service unavailable"
        )


@router.get("/quote", response_model=StockQuoteOut)
def stock_quote(ticker: str = Query(..., description="Stock ticker (e.g. RELIANCE.NS)")):
    try:
        quote = get_stock_quote(ticker)
        if quote is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Could not fetch quote for {ticker}"
            )
        return StockQuoteOut(**quote)
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Failed to fetch stock quote for {ticker}: {e}")
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Stock quote service unavailable"
        )


@router.post("/batch-quotes", response_model=BatchQuoteOut)
def batch_quotes(request: BatchQuoteRequest):
    try:
        quotes_raw = get_batch_stock_quotes(request.tickers)
        result = {}
        for ticker, quote in quotes_raw.items():
            if quote:
                result[ticker] = StockQuoteOut(**quote)
            else:
                result[ticker] = None
        return BatchQuoteOut(quotes=result)
    except Exception as e:
        logger.error(f"Failed to fetch batch quotes: {e}")
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Batch quote service unavailable"
        )
