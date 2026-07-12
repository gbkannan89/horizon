import logging
from typing import List
from datetime import datetime, date, timedelta
from fastapi import APIRouter, Depends, HTTPException
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.analytics import (
    RecurringDetectionOut, RecurringConfirmIn,
    SpendingPatternOut, SubscriptionCandidateOut, LapsedSubscriptionOut,
    FullAnalysisOut
)
from ..services.analytics_engine import (
    detect_recurring_transactions, analyze_spending_patterns,
    detect_subscriptions, detect_lapsed_subscriptions,
    generate_spending_nudges, run_full_analysis, add_confirmed_recurring
)
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/analytics", tags=["Analytics"])

MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]


@router.get("/budget-breakdown")
def budget_breakdown(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    with conn.cursor() as cur:
        user_id = current_user.id
        today = date.today()
        first_current = today.replace(day=1)
        first_prev = (first_current - timedelta(days=1)).replace(day=1)

        # Current month spending by category
        cur.execute("""
            SELECT category, SUM(amount) as total
            FROM expenses
            WHERE user_id = %s AND date >= %s AND date <= %s
            GROUP BY category ORDER BY total DESC
        """, (user_id, first_current, today))
        category_rows = cur.fetchall()
        category_breakdown = [
            {"category": r[0], "amount": float(r[1])} for r in category_rows
        ]

        # Current month totals per bucket
        cur.execute("""
            SELECT bucket, SUM(amount)
            FROM expenses
            WHERE user_id = %s AND date >= %s AND date <= %s
            GROUP BY bucket
        """, (user_id, first_current, today))
        current_buckets = {r[0]: float(r[1]) for r in cur.fetchall()}

        # Previous month totals per bucket
        cur.execute("""
            SELECT bucket, SUM(amount)
            FROM expenses
            WHERE user_id = %s AND date >= %s AND date < %s
            GROUP BY bucket
        """, (user_id, first_prev, first_current))
        prev_buckets = {r[0]: float(r[1]) for r in cur.fetchall()}

        def pct_change(current, previous):
            if previous == 0:
                return 100.0 if current > 0 else 0.0
            return round((current - previous) / previous * 100, 1)

        bucket_comparison = []
        for bucket in ["Needs", "Wants", "Savings"]:
            curr_val = current_buckets.get(bucket, 0.0)
            prev_val = prev_buckets.get(bucket, 0.0)
            bucket_comparison.append({
                "bucket": bucket,
                "currentAmount": curr_val,
                "previousAmount": prev_val,
                "changePct": pct_change(curr_val, prev_val),
            })

        return {
            "categoryBreakdown": category_breakdown,
            "bucketComparison": bucket_comparison,
            "currentMonth": MONTHS[today.month - 1],
            "previousMonth": MONTHS[(first_prev).month - 1],
        }


@router.get("/portfolio-summary")
def portfolio_summary(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("""
                SELECT id, name, amount, type, purchase_price, purchase_date, created_at
                FROM assets WHERE user_id = %s ORDER BY amount DESC
            """, (current_user.id,))
            asset_rows = cur.fetchall()
    except Exception:
        # Fallback if purchase_price/purchase_date columns don't exist yet
        with conn.cursor() as cur:
            cur.execute("""
                SELECT id, name, amount, type, NULL, NULL, created_at
                FROM assets WHERE user_id = %s ORDER BY amount DESC
            """, (current_user.id,))
            asset_rows = cur.fetchall()

    if not asset_rows:
        return {"totalInvested": 0, "totalCurrent": 0, "totalReturn": 0, "estimatedCagr": 0, "allocation": [], "assets": []}

    type_labels = {"bank": "Savings", "fd": "Fixed Deposit", "equity": "Equity", "physical": "Physical", "rd": "Recurring Deposit"}
    type_colors = {"bank": "#3B82F6", "fd": "#8B5CF6", "equity": "#10B981", "physical": "#F59E0B", "rd": "#EC4899"}

    total_current = 0.0
    total_invested = 0.0
    total_return = 0.0
    type_totals = {}
    asset_list = []

    for r in asset_rows:
        current_val = float(r[2])
        total_current += current_val

        pp = float(r[4]) if r[4] is not None else None
        invested = pp if pp is not None and pp > 0 else current_val
        total_invested += invested

        diff = current_val - invested
        total_return += diff

        ret_pct = round((diff / invested * 100), 1) if invested > 0 else 0.0

        cagr = None
        if pp is not None and r[5] is not None and pp > 0:
            from datetime import date
            days = (date.today() - r[5]).days
            if days > 0:
                years = days / 365.25
                cagr = round(((current_val / pp) ** (1 / years) - 1) * 100, 1)

        asset_list.append({
            "id": r[0], "name": r[1], "type": r[3],
            "invested": round(invested, 2),
            "current": round(current_val, 2),
            "returnPct": ret_pct,
            "cagr": cagr
        })

        atype = r[3]
        if atype not in type_totals:
            type_totals[atype] = 0.0
        type_totals[atype] += current_val

    allocation = [
        {
            "type": t,
            "label": type_labels.get(t, t.capitalize()),
            "amount": round(v, 2),
            "percentage": round(v / total_current * 100, 1) if total_current > 0 else 0,
            "color": type_colors.get(t, "#6366F1"),
        }
        for t, v in sorted(type_totals.items(), key=lambda x: -x[1])
    ]

    estimated_cagr = 0.0
    cagr_values = [a["cagr"] for a in asset_list if a["cagr"] is not None]
    if cagr_values:
        estimated_cagr = round(sum(cagr_values) / len(cagr_values), 1)

    return {
        "totalInvested": round(total_invested, 2),
        "totalCurrent": round(total_current, 2),
        "totalReturn": round(total_return, 2),
        "estimatedCagr": estimated_cagr,
        "allocation": allocation,
        "assets": asset_list,
    }


# ── RECURRING DETECTION ────────────────────────────────────────────────────
@router.get("/recurring-detection", response_model=List[RecurringDetectionOut])
def get_recurring_detection(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        results = detect_recurring_transactions(current_user.id, conn)
        return results
    except Exception as e:
        logger.error(f"Recurring detection error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/recurring-detection/confirm")
def confirm_recurring(
    confirmation: RecurringConfirmIn,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        result = add_confirmed_recurring(current_user.id, confirmation.dict(), conn)
        if "error" in result:
            raise HTTPException(status_code=409, detail=result["error"])
        return result
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Confirm recurring error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ── SPENDING PATTERNS ──────────────────────────────────────────────────────
@router.get("/spending-patterns", response_model=SpendingPatternOut)
def get_spending_patterns(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        results = analyze_spending_patterns(current_user.id, conn)
        return results
    except Exception as e:
        logger.error(f"Spending patterns error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ── SUBSCRIPTION DETECTION ─────────────────────────────────────────────────
@router.get("/subscriptions/detect", response_model=List[SubscriptionCandidateOut])
def get_subscription_candidates(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        results = detect_subscriptions(current_user.id, conn)
        return results
    except Exception as e:
        logger.error(f"Subscription detection error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/subscriptions/lapsed", response_model=List[LapsedSubscriptionOut])
def get_lapsed_subscriptions(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        results = detect_lapsed_subscriptions(current_user.id, conn)
        return results
    except Exception as e:
        logger.error(f"Lapsed subscription detection error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ── FULL ANALYSIS ───────────────────────────────────────────────────────────
@router.post("/run-analysis", response_model=FullAnalysisOut)
def run_analysis(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        results = run_full_analysis(current_user.id, conn)
        return results
    except Exception as e:
        logger.error(f"Full analysis error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/returns")
def get_investment_returns(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    from ..services.returns import calculate_cagr
    with db.cursor() as cur:
        # Fetch gold, stock, PF, and FD assets
        cur.execute(
            """
            SELECT name, type, amount, purchase_price, purchase_date
            FROM assets WHERE user_id = %s
            """,
            (current_user.id,)
        )
        rows = cur.fetchall()
        
        results = []
        for name, atype, amt, pp, pdate in rows:
            amt = float(amt)
            pp = float(pp) if pp is not None else amt
            
            # Estimate years elapsed
            years = 1.0
            if pdate:
                days = (date.today() - pdate).days
                if days > 0:
                    years = days / 365.25
            
            cagr_val = calculate_cagr(pp, amt, years)
            results.append({
                "name": name,
                "type": atype,
                "invested": pp,
                "current": amt,
                "cagr": round(cagr_val * 100, 2),
                "return_type": "CAGR"
            })
            
        return results
