import logging
from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from ..database import get_db
from ..schemas import LiabilityCreate, LiabilityOut, UserOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/liabilities", tags=["liabilities"])

# POST Add Liability (Protected)
@router.post("", response_model=LiabilityOut, status_code=status.HTTP_201_CREATED)
def add_liability(liability_in: LiabilityCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        try:
            cur.execute(
                """
                INSERT INTO liabilities (user_id, type, name, outstanding, emi, interest_rate, tenure_months)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                RETURNING id, user_id, type, name, outstanding, emi, interest_rate, tenure_months, created_at
                """,
                (
                    current_user.id,
                    liability_in.type,
                    liability_in.name,
                    liability_in.outstanding,
                    liability_in.emi,
                    liability_in.interest_rate,
                    liability_in.tenure_months
                )
            )
            row = cur.fetchone()
            db.commit()
            return LiabilityOut(
                id=row[0],
                user_id=row[1],
                type=row[2],
                name=row[3],
                outstanding=float(row[4]),
                emi=float(row[5]),
                interest_rate=float(row[6]),
                tenure_months=row[7],
                created_at=row[8]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add liability: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not add liability"
            )

# GET List Liabilities (Protected)
@router.get("", response_model=List[LiabilityOut])
def list_liabilities(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id, user_id, type, name, outstanding, emi, interest_rate, tenure_months, created_at FROM liabilities WHERE user_id = %s",
            (current_user.id,)
        )
        rows = cur.fetchall()
        return [
            LiabilityOut(
                id=r[0],
                user_id=r[1],
                type=r[2],
                name=r[3],
                outstanding=float(r[4]),
                emi=float(r[5]),
                interest_rate=float(r[6]),
                tenure_months=r[7],
                created_at=r[8]
            )
            for r in rows
        ]

# PUT Update Liability (Protected)
@router.put("/{liability_id}", response_model=LiabilityOut)
def update_liability(liability_id: int, liability_in: LiabilityCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id FROM liabilities WHERE id = %s AND user_id = %s",
            (liability_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Liability not found"
            )
            
        try:
            cur.execute(
                """
                UPDATE liabilities SET type = %s, name = %s, outstanding = %s, emi = %s, interest_rate = %s, tenure_months = %s
                WHERE id = %s
                RETURNING id, user_id, type, name, outstanding, emi, interest_rate, tenure_months, created_at
                """,
                (
                    liability_in.type,
                    liability_in.name,
                    liability_in.outstanding,
                    liability_in.emi,
                    liability_in.interest_rate,
                    liability_in.tenure_months,
                    liability_id
                )
            )
            row = cur.fetchone()
            db.commit()
            return LiabilityOut(
                id=row[0],
                user_id=row[1],
                type=row[2],
                name=row[3],
                outstanding=float(row[4]),
                emi=float(row[5]),
                interest_rate=float(row[6]),
                tenure_months=row[7],
                created_at=row[8]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update liability: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not update liability"
            )

# DELETE Remove Liability (Protected)
@router.delete("/{liability_id}", status_code=status.HTTP_200_OK)
def delete_liability(liability_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id FROM liabilities WHERE id = %s AND user_id = %s",
            (liability_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Liability not found"
            )
            
        try:
            cur.execute("DELETE FROM liabilities WHERE id = %s", (liability_id,))
            db.commit()
            return {"message": "Liability deleted successfully"}
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to delete liability: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not delete liability"
            )
