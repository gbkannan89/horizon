import logging
import time
import threading
from datetime import datetime, time as dt_time, date
from typing import Optional, Dict, List, Tuple

logger = logging.getLogger(__name__)

USD_INR_CACHE: Dict[str, any] = {"rate": None, "timestamp": 0}
GOLD_CACHE: Dict[str, any] = {"rate_24k": None, "rate_22k": None, "rate_18k": None, "timestamp": 0}
STOCK_CACHE: Dict[str, dict] = {}
STOCK_CACHE_TIMESTAMPS: Dict[str, float] = {}
_cache_lock = threading.Lock()

GOLD_CACHE_TTL = 6 * 3600
STOCK_MARKET_HOURS_TTL = 15 * 60
STOCK_NON_MARKET_TTL = 6 * 3600

MARKET_OPEN = dt_time(9, 15)
MARKET_CLOSE = dt_time(15, 30)


def _is_market_hours() -> bool:
    now = datetime.now()
    if now.weekday() >= 5:
        return False
    t = now.time()
    return MARKET_OPEN <= t <= MARKET_CLOSE


def _usd_to_inr() -> float:
    now = time.time()
    if USD_INR_CACHE["rate"] is not None and now - USD_INR_CACHE["timestamp"] < 3600:
        return USD_INR_CACHE["rate"]

    import httpx
    try:
        resp = httpx.get(
            "https://api.exchangerate-api.com/v4/latest/USD",
            timeout=5
        )
        if resp.status_code == 200:
            data = resp.json()
            rate = data["rates"]["INR"]
            with _cache_lock:
                USD_INR_CACHE["rate"] = rate
                USD_INR_CACHE["timestamp"] = now
            return rate
    except Exception as e:
        logger.warning(f"Failed to fetch USD/INR rate: {e}")

    if USD_INR_CACHE["rate"] is not None:
        return USD_INR_CACHE["rate"]
    return 83.50


def _get_gold_rate_from_api() -> Optional[float]:
    import httpx
    try:
        resp = httpx.get(
            "https://data-asg.goldprice.org/dbXRates/INR",
            timeout=10,
            headers={"User-Agent": "Mozilla/5.0"}
        )
        if resp.status_code == 200:
            data = resp.json()
            items = data.get("items", [])
            if items:
                item = items[0]
                ask = item.get("ask", 0)
                bid = item.get("bid", 0)
                xau_price_inr = (ask + bid) / 2
                oz_to_gram = 31.1035
                inr_per_gram_24k = xau_price_inr / oz_to_gram
                return round(inr_per_gram_24k, 2)
    except Exception as e:
        logger.warning(f"GoldPrice.org API failed: {e}")
    return None


def get_gold_rates() -> Dict[str, float]:
    now = time.time()
    with _cache_lock:
        if GOLD_CACHE["rate_24k"] is not None and now - GOLD_CACHE["timestamp"] < GOLD_CACHE_TTL:
            return {
                "rate_24k_per_gram": GOLD_CACHE["rate_24k"],
                "rate_22k_per_gram": GOLD_CACHE["rate_22k"],
                "rate_18k_per_gram": GOLD_CACHE["rate_18k"]
            }

    rate_24k = _get_gold_rate_from_api()
    if rate_24k is None:
        usd_inr = _usd_to_inr()
        rate_24k = round(83.50 * usd_inr / 31.1035, 2)

    rate_22k = round(rate_24k * 22 / 24, 2)
    rate_18k = round(rate_24k * 18 / 24, 2)

    with _cache_lock:
        GOLD_CACHE["rate_24k"] = rate_24k
        GOLD_CACHE["rate_22k"] = rate_22k
        GOLD_CACHE["rate_18k"] = rate_18k
        GOLD_CACHE["timestamp"] = now

    return {
        "rate_24k_per_gram": rate_24k,
        "rate_22k_per_gram": rate_22k,
        "rate_18k_per_gram": rate_18k
    }


def get_gold_rate_for_carat(carat: int) -> float:
    rates = get_gold_rates()
    key = f"rate_{carat}k_per_gram"
    if key in rates:
        return rates[key]
    base = rates["rate_24k_per_gram"]
    return round(base * carat / 24, 2)


def get_stock_quote(ticker: str) -> Optional[Dict[str, any]]:
    now = time.time()
    ttl = STOCK_MARKET_HOURS_TTL if _is_market_hours() else STOCK_NON_MARKET_TTL

    with _cache_lock:
        cached = STOCK_CACHE.get(ticker)
        cached_ts = STOCK_CACHE_TIMESTAMPS.get(ticker, 0)
        if cached is not None and now - cached_ts < ttl:
            return dict(cached)

    import httpx
    try:
        url = f"https://query1.finance.yahoo.com/v8/finance/chart/{ticker}"
        resp = httpx.get(
            url,
            timeout=10,
            headers={"User-Agent": "Mozilla/5.0"}
        )
        if resp.status_code == 200:
            data = resp.json()
            result = data.get("chart", {}).get("result", [])
            if result:
                meta = result[0].get("meta", {})
                quote = result[0].get("indicators", {}).get("quote", [{}])[0]
                current_price = meta.get("regularMarketPrice")
                if current_price is None:
                    prices = quote.get("close", [])
                    current_price = prices[-1] if prices else None

                prev_close = meta.get("chartPreviousClose") or meta.get("previousClose")
                change = round(current_price - prev_close, 2) if (current_price and prev_close) else None
                change_pct = round((change / prev_close) * 100, 2) if (change and prev_close) else None

                result_data = {
                    "ticker": ticker,
                    "name": meta.get("symbol", ticker),
                    "price": round(current_price, 2) if current_price else None,
                    "change": change,
                    "change_pct": change_pct,
                    "currency": meta.get("currency", "INR"),
                    "last_updated": datetime.now().isoformat(),
                    "source": "yahoo.finance"
                }

                with _cache_lock:
                    STOCK_CACHE[ticker] = result_data
                    STOCK_CACHE_TIMESTAMPS[ticker] = now

                return result_data
    except Exception as e:
        logger.warning(f"Yahoo Finance API failed for {ticker}: {e}")

    with _cache_lock:
        cached = STOCK_CACHE.get(ticker)
        if cached is not None:
            return dict(cached)
    return None


def get_batch_stock_quotes(tickers: List[str]) -> Dict[str, Optional[Dict[str, any]]]:
    results = {}
    for ticker in tickers:
        results[ticker] = get_stock_quote(ticker)
    return results
