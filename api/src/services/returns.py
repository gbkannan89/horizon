import logging
from typing import List, Tuple
from datetime import date

logger = logging.getLogger(__name__)

def calculate_cagr(begin_value: float, end_value: float, years: float) -> float:
    if begin_value <= 0 or years <= 0:
        return 0.0
    try:
        return (end_value / begin_value) ** (1.0 / years) - 1.0
    except Exception as e:
        logger.error(f"Error calculating CAGR: {e}")
        return 0.0

def calculate_xirr(cash_flows: List[Tuple[float, date]]) -> float:
    """
    Newton-Raphson solver for XIRR.
    cash_flows: list of (amount, date) tuples.
    Returns annualized XIRR as decimal (e.g. 0.12 = 12%).
    """
    if not cash_flows:
        return 0.0

    # Sort cash flows chronologically
    sorted_flows = sorted(cash_flows, key=lambda x: x[1])
    t0 = sorted_flows[0][1]

    def npv(r: float) -> float:
        val = 0.0
        for amt, dt in sorted_flows:
            d = (dt - t0).days
            val += amt / ((1.0 + r) ** (d / 365.25))
        return val

    def npv_prime(r: float) -> float:
        val = 0.0
        for amt, dt in sorted_flows:
            d = (dt - t0).days
            if d == 0:
                continue
            val -= (d / 365.25) * amt * ((1.0 + r) ** (-d / 365.25 - 1.0))
        return val

    # Try standard guesses
    r = 0.1
    for _ in range(100):
        try:
            n = npv(r)
            np_val = npv_prime(r)
            if abs(np_val) < 1e-12:
                break
            new_r = r - n / np_val
            if abs(new_r - r) < 1e-6:
                return new_r
            r = new_r
        except Exception:
            break
            
    return 0.0
