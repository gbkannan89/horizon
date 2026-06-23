import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..database import get_db
from ..schemas import AssetCreate, AssetOut, VehicleCreate, VehicleOut, UserOut
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
                INSERT INTO assets (user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, is_liability, generates_income, income_frequency)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, is_liability, generates_income, income_frequency, created_at
                """,
                (
                    current_user.id,
                    asset_in.type,
                    asset_in.subtype,
                    asset_in.name,
                    asset_in.amount,
                    asset_in.interest_rate,
                    asset_in.start_date,
                    asset_in.maturity_date,
                    asset_in.is_liability,
                    asset_in.generates_income,
                    asset_in.income_frequency
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
                    """
                    INSERT INTO incomes (user_id, label, type, amount, frequency)
                    VALUES (%s, %s, 'passive', %s, %s)
                    """,
                    (current_user.id, f"Income from {asset_in.name}", income_amount, freq)
                )
                
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
                is_liability=row[9],
                generates_income=row[10],
                income_frequency=row[11],
                created_at=row[12]
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
                "SELECT id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, is_liability, generates_income, income_frequency, created_at FROM assets WHERE user_id = %s AND type = %s",
                (current_user.id, type)
            )
        else:
            cur.execute(
                "SELECT id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, is_liability, generates_income, income_frequency, created_at FROM assets WHERE user_id = %s",
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
                is_liability=r[9],
                generates_income=r[10],
                income_frequency=r[11],
                created_at=r[12]
            )
            for r in rows
        ]

# --- VEHICLES ---

@router.post("/vehicles", response_model=VehicleOut, status_code=status.HTTP_201_CREATED)
def add_vehicle(vehicle_in: VehicleCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        try:
            cur.execute(
                """
                INSERT INTO vehicles (user_id, make_model, purchase_cost, insurance_renewal_date)
                VALUES (%s, %s, %s, %s)
                RETURNING id, user_id, make_model, purchase_cost, insurance_renewal_date, created_at
                """,
                (
                    current_user.id,
                    vehicle_in.make_model,
                    vehicle_in.purchase_cost,
                    vehicle_in.insurance_renewal_date
                )
            )
            row = cur.fetchone()
            db.commit()
            return VehicleOut(
                id=row[0],
                user_id=row[1],
                make_model=row[2],
                purchase_cost=float(row[3]),
                insurance_renewal_date=row[4],
                created_at=row[5]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add vehicle: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not add vehicle"
            )

@router.get("/vehicles", response_model=List[VehicleOut])
def list_vehicles(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id, user_id, make_model, purchase_cost, insurance_renewal_date, created_at FROM vehicles WHERE user_id = %s",
            (current_user.id,)
        )
        rows = cur.fetchall()
        return [
            VehicleOut(
                id=r[0],
                user_id=r[1],
                make_model=r[2],
                purchase_cost=float(r[3]),
                insurance_renewal_date=r[4],
                created_at=r[5]
            )
            for r in rows
        ]

@router.delete("/vehicles/{vehicle_id}", status_code=status.HTTP_200_OK)
def delete_vehicle(vehicle_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("DELETE FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        db.commit()
        return {"message": "Vehicle deleted successfully"}

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
                UPDATE assets SET type = %s, subtype = %s, name = %s, amount = %s, interest_rate = %s, start_date = %s, maturity_date = %s, is_liability = %s, generates_income = %s, income_frequency = %s
                WHERE id = %s
                RETURNING id, user_id, type, subtype, name, amount, interest_rate, start_date, maturity_date, is_liability, generates_income, income_frequency, created_at
                """,
                (
                    asset_in.type,
                    asset_in.subtype,
                    asset_in.name,
                    asset_in.amount,
                    asset_in.interest_rate,
                    asset_in.start_date,
                    asset_in.maturity_date,
                    asset_in.is_liability,
                    asset_in.generates_income,
                    asset_in.income_frequency,
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
                is_liability=row[9],
                generates_income=row[10],
                income_frequency=row[11],
                created_at=row[12]
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
