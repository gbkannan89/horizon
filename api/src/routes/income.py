import logging
from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.income import IncomeCreate, IncomeOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/incomes", tags=["incomes"])

# POST Add Income (Protected)
@router.post("", response_model=IncomeOut, status_code=status.HTTP_201_CREATED)
def add_income(income_in: IncomeCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    # Use label if provided, otherwise derive a default from type
    label = income_in.label or income_in.type.capitalize()
    with db.cursor() as cur:
        try:
            cur.execute(
                """
                INSERT INTO incomes (user_id, label, type, amount, frequency)
                VALUES (%s, %s, %s, %s, %s)
                RETURNING id, user_id, label, type, amount, frequency, created_at
                """,
                (current_user.id, label, income_in.type, income_in.amount, income_in.frequency)
            )
            row = cur.fetchone()
            db.commit()
            return IncomeOut(
                id=row[0],
                user_id=row[1],
                label=row[2],
                type=row[3],
                amount=float(row[4]),
                frequency=row[5],
                created_at=row[6]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add income: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not add income"
            )

# GET List Incomes (Protected)
@router.get("", response_model=List[IncomeOut])
def list_incomes(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id, user_id, label, type, amount, frequency, created_at FROM incomes WHERE user_id = %s ORDER BY created_at ASC",
            (current_user.id,)
        )
        rows = cur.fetchall()
        return [
            IncomeOut(
                id=r[0],
                user_id=r[1],
                label=r[2] or r[3].capitalize(),  # fallback for old rows without label
                type=r[3],
                amount=float(r[4]),
                frequency=r[5],
                created_at=r[6]
            )
            for r in rows
        ]

# PUT Update Income (Protected)
@router.put("/{income_id}", response_model=IncomeOut)
def update_income(income_id: int, income_in: IncomeCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    label = income_in.label or income_in.type.capitalize()
    with db.cursor() as cur:
        cur.execute(
            "SELECT id FROM incomes WHERE id = %s AND user_id = %s",
            (income_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Income source not found"
            )
            
        try:
            cur.execute(
                """
                UPDATE incomes SET label = %s, type = %s, amount = %s, frequency = %s
                WHERE id = %s
                RETURNING id, user_id, label, type, amount, frequency, created_at
                """,
                (label, income_in.type, income_in.amount, income_in.frequency, income_id)
            )
            row = cur.fetchone()
            db.commit()
            return IncomeOut(
                id=row[0],
                user_id=row[1],
                label=row[2],
                type=row[3],
                amount=float(row[4]),
                frequency=row[5],
                created_at=row[6]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update income: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not update income"
            )

# DELETE Remove Income (Protected)
@router.delete("/{income_id}", status_code=status.HTTP_200_OK)
def delete_income(income_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id FROM incomes WHERE id = %s AND user_id = %s",
            (income_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Income source not found"
            )
            
        try:
            cur.execute("DELETE FROM incomes WHERE id = %s", (income_id,))
            db.commit()
            return {"message": "Income source deleted successfully"}
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to delete income: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not delete income"
            )
