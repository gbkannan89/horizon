import logging
from datetime import datetime
from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.stock_holdings import StockHoldingCreate, StockHoldingUpdate, StockHoldingOut
from ..services.market_data import get_stock_quote
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/assets/stocks", tags=["stock-holdings"])

STOCK_SELECT = "id, user_id, ticker, exchange, name, quantity, avg_purchase_price, purchase_date, total_invested, last_current_price, last_current_value, last_value_updated_at, created_at"


def _row_to_out(r, current_value=None, return_pct=None, last_price=None):
    return StockHoldingOut(
        id=r[0], user_id=r[1], ticker=r[2], exchange=r[3], name=r[4],
        quantity=r[5], avg_purchase_price=float(r[6]), purchase_date=r[7],
        total_invested=float(r[8]),
        last_current_price=float(last_price) if last_price else (float(r[9]) if r[9] else None),
        current_value=current_value,
        return_pct=return_pct,
        created_at=r[12]
    )


def _compute_stock_value(r):
    quantity = r[5]
    total_invested = float(r[8])
    ticker = r[2]
    try:
        quote = get_stock_quote(ticker)
        if quote and quote.get("price"):
            price = quote["price"]
            current_value = round(quantity * price, 2)
            return_pct = round((current_value - total_invested) / total_invested * 100, 2) if total_invested > 0 else 0.0
            return current_value, return_pct, price
    except Exception as e:
        logger.warning(f"Could not fetch price for {ticker}: {e}")

    last_val = float(r[10]) if r[10] else None
    last_price = float(r[9]) if r[9] else None
    if last_val is not None:
        return_pct = round((last_val - total_invested) / total_invested * 100, 2) if total_invested > 0 else 0.0
        return last_val, return_pct, last_price
    return None, None, None


@router.get("", response_model=List[StockHoldingOut])
def list_holdings(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(f"SELECT {STOCK_SELECT} FROM stock_holdings WHERE user_id = %s ORDER BY purchase_date DESC", (current_user.id,))
        rows = cur.fetchall()
        results = []
        for r in rows:
            cv, rp, lp = _compute_stock_value(r)
            results.append(_row_to_out(r, cv, rp, lp))
        return results


@router.post("", response_model=StockHoldingOut, status_code=status.HTTP_201_CREATED)
def add_holding(holding_in: StockHoldingCreate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        try:
            cur.execute(
                f"""INSERT INTO stock_holdings (user_id, ticker, exchange, name, quantity, avg_purchase_price, purchase_date, total_invested)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s) RETURNING {STOCK_SELECT}""",
                (current_user.id, holding_in.ticker, holding_in.exchange, "",
                 holding_in.quantity, holding_in.avg_purchase_price,
                 holding_in.purchase_date, holding_in.total_invested)
            )
            r = cur.fetchone()
            db.commit()
            cv, rp, lp = _compute_stock_value(r)
            return _row_to_out(r, cv, rp, lp)
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add stock holding: {e}")
            raise HTTPException(status_code=500, detail="Could not add stock holding")


@router.put("/{holding_id}", response_model=StockHoldingOut)
def update_holding(holding_id: int, holding_in: StockHoldingUpdate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM stock_holdings WHERE id = %s AND user_id = %s", (holding_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Stock holding not found")
        try:
            updates = {k: v for k, v in holding_in.model_dump(exclude_none=True).items()}
            if not updates:
                raise HTTPException(status_code=400, detail="No fields to update")
            set_clause = ", ".join(f"{k} = %s" for k in updates)
            values = list(updates.values()) + [holding_id]
            cur.execute(f"UPDATE stock_holdings SET {set_clause} WHERE id = %s RETURNING {STOCK_SELECT}", values)
            r = cur.fetchone()
            db.commit()
            cv, rp, lp = _compute_stock_value(r)
            return _row_to_out(r, cv, rp, lp)
        except HTTPException:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update stock holding: {e}")
            raise HTTPException(status_code=500, detail="Could not update stock holding")


@router.delete("/{holding_id}", status_code=status.HTTP_200_OK)
def delete_holding(holding_id: int, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM stock_holdings WHERE id = %s AND user_id = %s", (holding_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Stock holding not found")
        cur.execute("DELETE FROM stock_holdings WHERE id = %s", (holding_id,))
        db.commit()
        return {"message": "Stock holding deleted successfully"}


@router.post("/refresh", response_model=List[StockHoldingOut])
def refresh_stock_values(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(f"SELECT {STOCK_SELECT} FROM stock_holdings WHERE user_id = %s", (current_user.id,))
        rows = cur.fetchall()
        results = []
        for r in rows:
            cv, rp, lp = _compute_stock_value(r)
            if cv is not None and lp is not None:
                cur.execute(
                    "UPDATE stock_holdings SET last_current_price = %s, last_current_value = %s, last_value_updated_at = %s WHERE id = %s",
                    (lp, cv, datetime.now(), r[0])
                )
            results.append(_row_to_out(r, cv, rp, lp))
        db.commit()
        return results
