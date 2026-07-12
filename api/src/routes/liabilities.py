import logging
from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.liabilities import LiabilityCreate, LiabilityOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/liabilities", tags=["liabilities"])

LIABILITY_SELECT = "id, user_id, type, name, outstanding, emi, interest_rate, tenure_months, gold_grams, gold_carat, gold_items_count, loan_date, bank_name, created_at"


def _row_to_out(r):
    return LiabilityOut(
        id=r[0], user_id=r[1], type=r[2], name=r[3],
        outstanding=float(r[4]), emi=float(r[5]),
        interest_rate=float(r[6]), tenure_months=r[7],
        gold_grams=float(r[8]) if r[8] is not None else None,
        gold_carat=r[9], gold_items_count=r[10],
        loan_date=r[11], bank_name=r[12],
        created_at=r[13]
    )


@router.post("", response_model=LiabilityOut, status_code=status.HTTP_201_CREATED)
def add_liability(liability_in: LiabilityCreate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        try:
            cur.execute(
                f"""INSERT INTO liabilities (user_id, type, name, outstanding, emi, interest_rate, tenure_months,
                gold_grams, gold_carat, gold_items_count, loan_date, bank_name)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING {LIABILITY_SELECT}""",
                (
                    current_user.id,
                    liability_in.type,
                    liability_in.name,
                    liability_in.outstanding,
                    liability_in.emi,
                    liability_in.interest_rate,
                    liability_in.tenure_months,
                    liability_in.gold_grams,
                    liability_in.gold_carat,
                    liability_in.gold_items_count,
                    liability_in.loan_date,
                    liability_in.bank_name
                )
            )
            row = cur.fetchone()

            if liability_in.emi > 0:
                bill_name = f"{liability_in.name} EMI"
                cur.execute(
                    """INSERT INTO recurring_bills (user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription)
                    VALUES (%s, %s, %s, 'monthly', 'EMI', 'Needs', 1, %s, TRUE, FALSE)""",
                    (current_user.id, bill_name, liability_in.emi, liability_in.emi)
                )

            db.commit()
            return _row_to_out(row)
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add liability: {e}")
            raise HTTPException(status_code=500, detail="Could not add liability")


@router.get("", response_model=List[LiabilityOut])
def list_liabilities(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(f"SELECT {LIABILITY_SELECT} FROM liabilities WHERE user_id = %s", (current_user.id,))
        return [_row_to_out(r) for r in cur.fetchall()]


@router.put("/{liability_id}", response_model=LiabilityOut)
def update_liability(liability_id: int, liability_in: LiabilityCreate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM liabilities WHERE id = %s AND user_id = %s", (liability_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Liability not found")
        try:
            cur.execute(
                f"""UPDATE liabilities SET type = %s, name = %s, outstanding = %s, emi = %s, interest_rate = %s,
                tenure_months = %s, gold_grams = %s, gold_carat = %s, gold_items_count = %s, loan_date = %s, bank_name = %s
                WHERE id = %s RETURNING {LIABILITY_SELECT}""",
                (
                    liability_in.type, liability_in.name, liability_in.outstanding,
                    liability_in.emi, liability_in.interest_rate, liability_in.tenure_months,
                    liability_in.gold_grams, liability_in.gold_carat,
                    liability_in.gold_items_count, liability_in.loan_date, liability_in.bank_name,
                    liability_id
                )
            )
            row = cur.fetchone()
            db.commit()
            return _row_to_out(row)
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update liability: {e}")
            raise HTTPException(status_code=500, detail="Could not update liability")


@router.delete("/{liability_id}", status_code=status.HTTP_200_OK)
def delete_liability(liability_id: int, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM liabilities WHERE id = %s AND user_id = %s", (liability_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Liability not found")
        cur.execute("DELETE FROM liabilities WHERE id = %s", (liability_id,))
        db.commit()
        return {"message": "Liability deleted successfully"}
