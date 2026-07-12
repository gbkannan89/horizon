import logging
import math
from datetime import date, datetime
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.pf_assets import PFAssetCreate, PFAssetUpdate, PFAssetOut, PFProjectionOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/assets/pf", tags=["pf-assets"])

PF_SELECT = "id, user_id, start_date, retirement_age, current_balance, monthly_contribution, employer_contribution, interest_rate, user_date_of_birth, notes, created_at"


def _project_corpus(balance: float, monthly: float, emp_contrib: Optional[float], rate: float, years: int) -> dict:
    total_monthly = monthly + (emp_contrib or 0)
    monthly_rate = (rate / 100) / 12
    total_months = years * 12

    if monthly_rate > 0:
        fv_current = balance * ((1 + monthly_rate) ** total_months)
        fv_contribs = total_monthly * (((1 + monthly_rate) ** total_months - 1) / monthly_rate) * (1 + monthly_rate)
    else:
        fv_current = balance
        fv_contribs = total_monthly * total_months

    projected = round(fv_current + fv_contribs, 2)
    total_contributions = round(balance + total_monthly * total_months, 2)
    total_interest = round(projected - total_contributions, 2)
    return {"projected": projected, "total_contributions": total_contributions, "total_interest": total_interest}


@router.get("", response_model=List[PFAssetOut])
def list_pf_assets(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(f"SELECT {PF_SELECT} FROM pf_assets WHERE user_id = %s ORDER BY start_date DESC", (current_user.id,))
        rows = cur.fetchall()
        results = []
        for r in rows:
            today = date.today()
            start_date = r[2]
            years_until_retirement = max(0, r[3] - (today.year - start_date.year))
            proj = _project_corpus(float(r[4]), float(r[5]), float(r[6]) if r[6] else None, float(r[7]), years_until_retirement)
            results.append(PFAssetOut(
                id=r[0], user_id=r[1], start_date=start_date, retirement_age=r[3],
                current_balance=float(r[4]), monthly_contribution=float(r[5]),
                employer_contribution=float(r[6]) if r[6] else None,
                interest_rate=float(r[7]),
                user_date_of_birth=r[8],
                projected_corpus=proj["projected"],
                years_until_retirement=years_until_retirement,
                notes=r[9], created_at=r[10]
            ))
        return results


@router.post("", response_model=PFAssetOut, status_code=status.HTTP_201_CREATED)
def add_pf_asset(asset_in: PFAssetCreate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        try:
            cur.execute(
                f"""INSERT INTO pf_assets (user_id, start_date, retirement_age, current_balance, monthly_contribution,
                employer_contribution, interest_rate, user_date_of_birth, notes)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s) RETURNING {PF_SELECT}""",
                (current_user.id, asset_in.start_date, asset_in.retirement_age,
                 asset_in.current_balance, asset_in.monthly_contribution,
                 asset_in.employer_contribution, asset_in.interest_rate,
                 asset_in.user_date_of_birth, asset_in.notes)
            )
            r = cur.fetchone()
            db.commit()
            today = date.today()
            years_until = max(0, r[3] - (today.year - r[2].year))
            proj = _project_corpus(float(r[4]), float(r[5]), float(r[6]) if r[6] else None, float(r[7]), years_until)
            return PFAssetOut(
                id=r[0], user_id=r[1], start_date=r[2], retirement_age=r[3],
                current_balance=float(r[4]), monthly_contribution=float(r[5]),
                employer_contribution=float(r[6]) if r[6] else None,
                interest_rate=float(r[7]),
                user_date_of_birth=r[8],
                projected_corpus=proj["projected"],
                years_until_retirement=years_until,
                notes=r[9], created_at=r[10]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add PF asset: {e}")
            raise HTTPException(status_code=500, detail="Could not add PF asset")


@router.put("", response_model=PFAssetOut)
def update_pf_asset(asset_in: PFAssetUpdate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(f"SELECT {PF_SELECT} FROM pf_assets WHERE user_id = %s", (current_user.id,))
        existing = cur.fetchone()
        if not existing:
            raise HTTPException(status_code=404, detail="PF asset not found")
        try:
            updates = {k: v for k, v in asset_in.model_dump(exclude_none=True).items()}
            if not updates:
                raise HTTPException(status_code=400, detail="No fields to update")
            set_clause = ", ".join(f"{k} = %s" for k in updates)
            values = list(updates.values()) + [existing[0]]
            cur.execute(f"UPDATE pf_assets SET {set_clause} WHERE id = %s RETURNING {PF_SELECT}", values)
            r = cur.fetchone()
            db.commit()
            today = date.today()
            years_until = max(0, r[3] - (today.year - r[2].year))
            proj = _project_corpus(float(r[4]), float(r[5]), float(r[6]) if r[6] else None, float(r[7]), years_until)
            return PFAssetOut(
                id=r[0], user_id=r[1], start_date=r[2], retirement_age=r[3],
                current_balance=float(r[4]), monthly_contribution=float(r[5]),
                employer_contribution=float(r[6]) if r[6] else None,
                interest_rate=float(r[7]),
                user_date_of_birth=r[8],
                projected_corpus=proj["projected"],
                years_until_retirement=years_until,
                notes=r[9], created_at=r[10]
            )
        except HTTPException:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update PF asset: {e}")
            raise HTTPException(status_code=500, detail="Could not update PF asset")


@router.delete("", status_code=status.HTTP_200_OK)
def delete_pf_asset(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM pf_assets WHERE user_id = %s", (current_user.id,))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="PF asset not found")
        cur.execute("DELETE FROM pf_assets WHERE user_id = %s", (current_user.id,))
        db.commit()
        return {"message": "PF asset deleted successfully"}


@router.get("/projection", response_model=PFProjectionOut)
def get_pf_projection(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(f"SELECT {PF_SELECT} FROM pf_assets WHERE user_id = %s", (current_user.id,))
        r = cur.fetchone()
        if not r:
            raise HTTPException(status_code=404, detail="No PF asset found")

        today = date.today()
        years_until = max(0, r[3] - (today.year - r[2].year))
        balance = float(r[4])
        monthly = float(r[5])
        emp_contrib = float(r[6]) if r[6] else None
        rate = float(r[7])
        total_monthly = monthly + (emp_contrib or 0)

        proj = _project_corpus(balance, monthly, emp_contrib, rate, years_until)
        return PFProjectionOut(
            current_balance=balance,
            monthly_contribution=monthly,
            employer_contribution=emp_contrib,
            interest_rate=rate,
            retirement_age=r[3],
            years_until_retirement=years_until,
            projected_corpus=proj["projected"],
            total_contributions=proj["total_contributions"],
            total_interest_earned=proj["total_interest"]
        )
