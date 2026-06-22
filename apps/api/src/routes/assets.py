import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..database import get_db
from ..schemas import AssetCreate, AssetOut, UserOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/assets", tags=["assets"])

# POST Add Asset (Protected)
@router.post("", response_model=AssetOut, status_code=status.HTTP_201_CREATED)
def add_asset(asset_in: AssetCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        try:
            cur.execute(
                """
                INSERT INTO assets (user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, created_at
                """,
                (
                    current_user.id,
                    asset_in.type,
                    asset_in.subtype,
                    asset_in.name,
                    asset_in.amount,
                    asset_in.interest_rate,
                    asset_in.start_date,
                    asset_in.maturity_date
                )
            )
            row = cur.fetchone()
            db.commit()
            return AssetOut(
                id=row[0],
                user_id=row[1],
                type=row[2],
                subtype=row[3],
                name=row[4],
                amount=float(row[5]),
                interest_rate=float(row[6]),
                start_date=row[7],
                maturity_date=row[8],
                created_at=row[9]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add asset: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not add asset"
            )

# GET List Assets (Protected)
@router.get("", response_model=List[AssetOut])
def list_assets(type: Optional[str] = None, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        if type:
            cur.execute(
                "SELECT id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, created_at FROM assets WHERE user_id = %s AND type = %s",
                (current_user.id, type)
            )
        else:
            cur.execute(
                "SELECT id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, created_at FROM assets WHERE user_id = %s",
                (current_user.id,)
            )
        rows = cur.fetchall()
        return [
            AssetOut(
                id=r[0],
                user_id=r[1],
                type=r[2],
                subtype=r[3],
                name=r[4],
                amount=float(r[5]),
                interest_rate=float(r[6]),
                start_date=r[7],
                maturity_date=r[8],
                created_at=r[9]
            )
            for r in rows
        ]

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
            cur.execute(
                """
                UPDATE assets SET type = %s, subtype = %s, name = %s, amount = %s, interest_rate = %s, start_date = %s, maturity_date = %s
                WHERE id = %s
                RETURNING id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, created_at
                """,
                (
                    asset_in.type,
                    asset_in.subtype,
                    asset_in.name,
                    asset_in.amount,
                    asset_in.interest_rate,
                    asset_in.start_date,
                    asset_in.maturity_date,
                    asset_id
                )
            )
            row = cur.fetchone()
            db.commit()
            return AssetOut(
                id=row[0],
                user_id=row[1],
                type=row[2],
                subtype=row[3],
                name=row[4],
                amount=float(row[5]),
                interest_rate=float(row[6]),
                start_date=row[7],
                maturity_date=row[8],
                created_at=row[9]
            )
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
