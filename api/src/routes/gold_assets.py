import logging
from datetime import datetime, date
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.gold_assets import GoldAssetCreate, GoldAssetUpdate, GoldAssetOut
from ..services.market_data import get_gold_rate_for_carat
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/assets/gold", tags=["gold-assets"])

GOLD_SELECT = "id, user_id, carat, grams, purchase_price_per_gram, purchase_date, last_current_value, last_value_updated_at, notes, created_at"


def _row_to_out(r, current_value=None, total_cost=None, return_pct=None):
    return GoldAssetOut(
        id=r[0], user_id=r[1], carat=r[2], grams=float(r[3]),
        purchase_price_per_gram=float(r[4]), purchase_date=r[5],
        last_current_value=float(r[6]) if r[6] is not None else None,
        current_value=current_value,
        total_purchase_cost=total_cost,
        return_pct=return_pct,
        notes=r[8], created_at=r[9]
    )


def _compute_gold_value(r):
    carat = r[2]
    grams = float(r[3])
    total_cost = grams * float(r[4])
    try:
        rate = get_gold_rate_for_carat(carat)
        current_value = round(grams * rate, 2)
        return_pct = round((current_value - total_cost) / total_cost * 100, 2) if total_cost > 0 else 0.0
    except Exception:
        current_value = float(r[6]) if r[6] is not None else None
        return_pct = None
    return current_value, total_cost, return_pct


@router.get("", response_model=List[GoldAssetOut])
def list_gold_assets(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(f"SELECT {GOLD_SELECT} FROM gold_assets WHERE user_id = %s ORDER BY purchase_date DESC", (current_user.id,))
        rows = cur.fetchall()
        results = []
        for r in rows:
            cv, tc, rp = _compute_gold_value(r)
            results.append(_row_to_out(r, cv, tc, rp))
        return results


@router.post("", response_model=GoldAssetOut, status_code=status.HTTP_201_CREATED)
def add_gold_asset(asset_in: GoldAssetCreate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        try:
            cur.execute(
                f"""INSERT INTO gold_assets (user_id, carat, grams, purchase_price_per_gram, purchase_date, notes)
                VALUES (%s, %s, %s, %s, %s, %s) RETURNING {GOLD_SELECT}""",
                (current_user.id, asset_in.carat, asset_in.grams,
                 asset_in.purchase_price_per_gram, asset_in.purchase_date, asset_in.notes)
            )
            r = cur.fetchone()
            db.commit()
            cv, tc, rp = _compute_gold_value(r)
            return _row_to_out(r, cv, tc, rp)
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add gold asset: {e}")
            raise HTTPException(status_code=500, detail="Could not add gold asset")


@router.put("/{asset_id}", response_model=GoldAssetOut)
def update_gold_asset(asset_id: int, asset_in: GoldAssetUpdate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM gold_assets WHERE id = %s AND user_id = %s", (asset_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Gold asset not found")
        try:
            updates = {}
            for k, v in asset_in.model_dump(exclude_none=True).items():
                updates[k] = v
            if not updates:
                raise HTTPException(status_code=400, detail="No fields to update")
            set_clause = ", ".join(f"{k} = %s" for k in updates)
            values = list(updates.values()) + [asset_id]
            cur.execute(
                f"UPDATE gold_assets SET {set_clause} WHERE id = %s RETURNING {GOLD_SELECT}",
                values
            )
            r = cur.fetchone()
            db.commit()
            cv, tc, rp = _compute_gold_value(r)
            return _row_to_out(r, cv, tc, rp)
        except HTTPException:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update gold asset: {e}")
            raise HTTPException(status_code=500, detail="Could not update gold asset")


@router.delete("/{asset_id}", status_code=status.HTTP_200_OK)
def delete_gold_asset(asset_id: int, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM gold_assets WHERE id = %s AND user_id = %s", (asset_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Gold asset not found")
        cur.execute("DELETE FROM gold_assets WHERE id = %s", (asset_id,))
        db.commit()
        return {"message": "Gold asset deleted successfully"}


@router.post("/refresh", response_model=List[GoldAssetOut])
def refresh_gold_values(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(f"SELECT {GOLD_SELECT} FROM gold_assets WHERE user_id = %s", (current_user.id,))
        rows = cur.fetchall()
        results = []
        for r in rows:
            cv, tc, rp = _compute_gold_value(r)
            if cv is not None:
                cur.execute(
                    "UPDATE gold_assets SET last_current_value = %s, last_value_updated_at = %s WHERE id = %s",
                    (cv, datetime.now(), r[0])
                )
            results.append(_row_to_out(r, cv, tc, rp))
        db.commit()
        return results
