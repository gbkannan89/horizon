import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..database import get_db
from ..schemas import RecurringBillCreate, RecurringBillOut, UserOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/bills", tags=["bills"])

def calculate_monthly_equivalent(amount: float, frequency: str) -> float:
    freq = frequency.strip().lower()
    if freq == "monthly":
        return amount
    elif freq == "quarterly":
        return amount / 3.0
    elif freq == "half-yearly":
        return amount / 6.0
    elif freq == "yearly":
        return amount / 12.0
    else:
        return amount

# POST Add Recurring Bill / Subscription (Protected)
@router.post("", response_model=RecurringBillOut, status_code=status.HTTP_201_CREATED)
def add_recurring_bill(bill_in: RecurringBillCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    monthly_equivalent = calculate_monthly_equivalent(bill_in.amount, bill_in.frequency)
    
    with db.cursor() as cur:
        try:
            cur.execute(
                """
                INSERT INTO recurring_bills (user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription, created_at
                """,
                (
                    current_user.id,
                    bill_in.name,
                    bill_in.amount,
                    bill_in.frequency,
                    bill_in.category,
                    bill_in.bucket,
                    bill_in.due_day,
                    monthly_equivalent,
                    bill_in.is_active,
                    bill_in.is_subscription
                )
            )
            row = cur.fetchone()
            db.commit()
            return RecurringBillOut(
                id=row[0],
                user_id=row[1],
                name=row[2],
                amount=float(row[3]),
                frequency=row[4],
                category=row[5],
                bucket=row[6],
                due_day=row[7],
                monthly_equivalent=float(row[8]),
                is_active=row[9],
                is_subscription=row[10],
                created_at=row[11]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add recurring bill: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not add recurring bill"
            )

# GET List Bills (Protected)
@router.get("", response_model=List[RecurringBillOut])
def list_recurring_bills(is_subscription: Optional[bool] = None, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        if is_subscription is not None:
            cur.execute(
                """
                SELECT id, user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription, created_at
                FROM recurring_bills WHERE user_id = %s AND is_subscription = %s
                """,
                (current_user.id, is_subscription)
            )
        else:
            cur.execute(
                """
                SELECT id, user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription, created_at
                FROM recurring_bills WHERE user_id = %s
                """,
                (current_user.id,)
            )
        rows = cur.fetchall()
        return [
            RecurringBillOut(
                id=r[0],
                user_id=r[1],
                name=r[2],
                amount=float(r[3]),
                frequency=r[4],
                category=r[5],
                bucket=r[6],
                due_day=r[7],
                monthly_equivalent=float(r[8]),
                is_active=r[9],
                is_subscription=r[10],
                created_at=r[11]
            )
            for r in rows
        ]

# PUT Update Bill (Protected)
@router.put("/{bill_id}", response_model=RecurringBillOut)
def update_recurring_bill(bill_id: int, bill_in: RecurringBillCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    monthly_equivalent = calculate_monthly_equivalent(bill_in.amount, bill_in.frequency)
    
    with db.cursor() as cur:
        cur.execute(
            "SELECT id FROM recurring_bills WHERE id = %s AND user_id = %s",
            (bill_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Recurring bill not found"
            )
            
        try:
            cur.execute(
                """
                UPDATE recurring_bills SET name = %s, amount = %s, frequency = %s, category = %s, bucket = %s, due_day = %s, monthly_equivalent = %s, is_active = %s, is_subscription = %s
                WHERE id = %s
                RETURNING id, user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription, created_at
                """,
                (
                    bill_in.name,
                    bill_in.amount,
                    bill_in.frequency,
                    bill_in.category,
                    bill_in.bucket,
                    bill_in.due_day,
                    monthly_equivalent,
                    bill_in.is_active,
                    bill_in.is_subscription,
                    bill_id
                )
            )
            row = cur.fetchone()
            db.commit()
            return RecurringBillOut(
                id=row[0],
                user_id=row[1],
                name=row[2],
                amount=float(row[3]),
                frequency=row[4],
                category=row[5],
                bucket=row[6],
                due_day=row[7],
                monthly_equivalent=float(row[8]),
                is_active=row[9],
                is_subscription=row[10],
                created_at=row[11]
            )
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update recurring bill: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not update recurring bill"
            )

# DELETE Remove Bill (Protected)
@router.delete("/{bill_id}", status_code=status.HTTP_200_OK)
def delete_recurring_bill(bill_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id FROM recurring_bills WHERE id = %s AND user_id = %s",
            (bill_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Recurring bill not found"
            )
            
        try:
            cur.execute("DELETE FROM recurring_bills WHERE id = %s", (bill_id,))
            db.commit()
            return {"message": "Recurring bill deleted successfully"}
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to delete recurring bill: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not delete recurring bill"
            )
