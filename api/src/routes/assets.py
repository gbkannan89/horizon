import logging
from datetime import date
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.assets import AssetCreate, AssetOut, VehicleCreate, VehicleOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/assets", tags=["assets"])

ASSET_SELECT = "id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, is_liability, generates_income, income_frequency, purchase_price, purchase_date, years_of_deposit, created_at"


def _row_to_asset(r):
    return AssetOut(
        id=r[0], user_id=r[1], type=r[2], subtype=r[3], name=r[4],
        amount=float(r[5]), interest_rate=float(r[6]),
        start_date=r[7], maturity_date=r[8],
        is_liability=r[9], generates_income=r[10], income_frequency=r[11],
        purchase_price=float(r[12]) if r[12] is not None else None,
        purchase_date=r[13],
        years_of_deposit=r[14],
        created_at=r[15]
    )


# POST Add Asset (Protected)
@router.post("", response_model=AssetOut, status_code=status.HTTP_201_CREATED)
def add_asset(asset_in: AssetCreate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        try:
            maturity_date = asset_in.maturity_date
            if asset_in.type == 'fd' and asset_in.start_date and asset_in.years_of_deposit and not maturity_date:
                maturity_date = date(
                    asset_in.start_date.year + asset_in.years_of_deposit,
                    asset_in.start_date.month,
                    asset_in.start_date.day
                )

            cur.execute(
                f"""INSERT INTO assets (user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date,
                is_liability, generates_income, income_frequency, purchase_price, purchase_date, years_of_deposit)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING {ASSET_SELECT}""",
                (
                    current_user.id, asset_in.type, asset_in.subtype, asset_in.name,
                    asset_in.amount, asset_in.interest_rate, asset_in.start_date, maturity_date,
                    asset_in.is_liability, asset_in.generates_income, asset_in.income_frequency,
                    asset_in.purchase_price, asset_in.purchase_date, asset_in.years_of_deposit
                )
            )
            row = cur.fetchone()

            if asset_in.generates_income and asset_in.interest_rate > 0:
                income_amount = asset_in.amount * (asset_in.interest_rate / 100)
                freq = asset_in.income_frequency or "yearly"
                if freq == "monthly":
                    income_amount = income_amount / 12
                elif freq == "quarterly":
                    income_amount = income_amount / 4
                cur.execute(
                    "INSERT INTO incomes (user_id, label, type, amount, frequency) VALUES (%s, %s, 'passive', %s, %s)",
                    (current_user.id, f"Income from {asset_in.name}", income_amount, freq)
                )

            db.commit()
            return _row_to_asset(row)
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add asset: {e}")
            raise HTTPException(status_code=500, detail="Could not add asset")


# GET List Assets (Protected)
@router.get("", response_model=List[AssetOut])
def list_assets(type: Optional[str] = None, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        if type:
            cur.execute(f"SELECT {ASSET_SELECT} FROM assets WHERE user_id = %s AND type = %s", (current_user.id, type))
        else:
            cur.execute(f"SELECT {ASSET_SELECT} FROM assets WHERE user_id = %s", (current_user.id,))
        return [_row_to_asset(r) for r in cur.fetchall()]

# Vehicles router has been extracted to a separate router in routes/vehicles.py

# PUT Update Asset (Protected)
@router.put("/{asset_id}", response_model=AssetOut)
def update_asset(asset_id: int, asset_in: AssetCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id FROM assets WHERE id = %s AND user_id = %s",
            (asset_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Asset not found"
            )
            
        try:
            maturity_date = asset_in.maturity_date
            if asset_in.type == 'fd' and asset_in.start_date and asset_in.years_of_deposit and not maturity_date:
                maturity_date = date(
                    asset_in.start_date.year + asset_in.years_of_deposit,
                    asset_in.start_date.month,
                    asset_in.start_date.day
                )

            cur.execute(
                f"""UPDATE assets SET type = %s, subtype = %s, name = %s, amount = %s, interest_rate = %s,
                start_date = %s, maturity_date = %s, is_liability = %s, generates_income = %s,
                income_frequency = %s, purchase_price = %s, purchase_date = %s, years_of_deposit = %s
                WHERE id = %s
                RETURNING {ASSET_SELECT}""",
                (
                    asset_in.type, asset_in.subtype, asset_in.name, asset_in.amount,
                    asset_in.interest_rate, asset_in.start_date, maturity_date,
                    asset_in.is_liability, asset_in.generates_income, asset_in.income_frequency,
                    asset_in.purchase_price, asset_in.purchase_date, asset_in.years_of_deposit, asset_id
                )
            )
            row = cur.fetchone()
            db.commit()
            return _row_to_asset(row)
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update asset: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not update asset"
            )

# DELETE Remove Asset (Protected)
@router.delete("/{asset_id}", status_code=status.HTTP_200_OK)
def delete_asset(asset_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id FROM assets WHERE id = %s AND user_id = %s",
            (asset_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Asset not found"
            )
            
        try:
            cur.execute("DELETE FROM assets WHERE id = %s", (asset_id,))
            db.commit()
            return {"message": "Asset deleted successfully"}
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to delete asset: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not delete asset"
            )
