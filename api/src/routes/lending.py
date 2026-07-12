import logging
from datetime import date
from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.lending import LendingCreate, LendingUpdate, LendingOut, LendingOverview
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/lending", tags=["lending"])

LENDING_SELECT = "id, user_id, direction, person_name, amount, date_given, promised_return_date, actual_return_date, returned_amount, status, notes, created_at"


def _row_to_out(r):
    remaining = float(r[4]) - float(r[8])
    return LendingOut(
        id=r[0], user_id=r[1], direction=r[2], person_name=r[3],
        amount=float(r[4]), date_given=r[5],
        promised_return_date=r[6], actual_return_date=r[7],
        returned_amount=float(r[8]), remaining=round(max(remaining, 0), 2),
        status=r[9], notes=r[10], created_at=r[11]
    )


def _update_status(cur, record_id: int):
    cur.execute(
        "SELECT amount, returned_amount, promised_return_date FROM lending_records WHERE id = %s",
        (record_id,)
    )
    r = cur.fetchone()
    if not r:
        return
    total = float(r[0])
    returned = float(r[1])
    promised = r[2]

    if returned >= total:
        new_status = "closed"
    elif returned > 0:
        new_status = "partial"
    elif promised and date.today() > promised:
        new_status = "overdue"
    else:
        new_status = "open"

    cur.execute("UPDATE lending_records SET status = %s WHERE id = %s", (new_status, record_id))


@router.get("", response_model=List[LendingOut])
def list_lending(direction: str = None, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        if direction:
            cur.execute(
                f"SELECT {LENDING_SELECT} FROM lending_records WHERE user_id = %s AND direction = %s ORDER BY date_given DESC",
                (current_user.id, direction)
            )
        else:
            cur.execute(
                f"SELECT {LENDING_SELECT} FROM lending_records WHERE user_id = %s ORDER BY date_given DESC",
                (current_user.id,)
            )
        rows = cur.fetchall()
        for r in rows:
            _update_status(cur, r[0])
        db.commit()
        cur.execute(
            f"SELECT {LENDING_SELECT} FROM lending_records WHERE user_id = %s ORDER BY date_given DESC",
            (current_user.id,)
        )
        return [_row_to_out(r) for r in cur.fetchall()]


@router.post("", response_model=LendingOut, status_code=status.HTTP_201_CREATED)
def add_lending(record_in: LendingCreate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        try:
            status_val = "open"
            if record_in.promised_return_date and date.today() > record_in.promised_return_date:
                status_val = "overdue"

            cur.execute(
                f"""INSERT INTO lending_records (user_id, direction, person_name, amount, date_given, promised_return_date, status, notes)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s) RETURNING {LENDING_SELECT}""",
                (current_user.id, record_in.direction, record_in.person_name,
                 record_in.amount, record_in.date_given, record_in.promised_return_date,
                 status_val, record_in.notes)
            )
            r = cur.fetchone()
            db.commit()
            return _row_to_out(r)
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to add lending record: {e}")
            raise HTTPException(status_code=500, detail="Could not add lending record")


@router.put("/{record_id}", response_model=LendingOut)
def update_lending(record_id: int, record_in: LendingUpdate, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM lending_records WHERE id = %s AND user_id = %s", (record_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Lending record not found")
        try:
            updates = {k: v for k, v in record_in.model_dump(exclude_none=True).items()}
            if not updates:
                raise HTTPException(status_code=400, detail="No fields to update")

            if "returned_amount" in updates:
                cur.execute("SELECT amount FROM lending_records WHERE id = %s", (record_id,))
                total = float(cur.fetchone()[0])
                if updates["returned_amount"] > total:
                    raise HTTPException(status_code=400, detail="Returned amount cannot exceed total amount")

            set_clause = ", ".join(f"{k} = %s" for k in updates)
            values = list(updates.values()) + [record_id]
            cur.execute(f"UPDATE lending_records SET {set_clause} WHERE id = %s RETURNING {LENDING_SELECT}", values)
            r = cur.fetchone()
            _update_status(cur, record_id)
            db.commit()
            cur.execute(f"SELECT {LENDING_SELECT} FROM lending_records WHERE id = %s", (record_id,))
            r = cur.fetchone()
            return _row_to_out(r)
        except HTTPException:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update lending record: {e}")
            raise HTTPException(status_code=500, detail="Could not update lending record")


@router.delete("/{record_id}", status_code=status.HTTP_200_OK)
def delete_lending(record_id: int, current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM lending_records WHERE id = %s AND user_id = %s", (record_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Lending record not found")
        cur.execute("DELETE FROM lending_records WHERE id = %s", (record_id,))
        db.commit()
        return {"message": "Lending record deleted successfully"}


@router.get("/overview", response_model=LendingOverview)
def get_lending_overview(current_user: UserOut = Depends(get_current_user), db=Depends(get_db)):
    with db.cursor() as cur:
        for r in cur.execute(
            f"SELECT {LENDING_SELECT} FROM lending_records WHERE user_id = %s", (current_user.id,)
        ).fetchall():
            _update_status(cur, r[0])
        db.commit()

        cur.execute("""
            SELECT
                COALESCE(SUM(CASE WHEN direction = 'lent' THEN amount - returned_amount ELSE 0 END), 0),
                COALESCE(SUM(CASE WHEN direction = 'borrowed' THEN amount - returned_amount ELSE 0 END), 0),
                COALESCE(SUM(CASE WHEN status = 'overdue' THEN 1 ELSE 0 END), 0),
                COALESCE(SUM(CASE WHEN status IN ('open', 'partial') THEN 1 ELSE 0 END), 0),
                COALESCE(SUM(CASE WHEN status = 'closed' THEN 1 ELSE 0 END), 0)
            FROM lending_records WHERE user_id = %s
        """, (current_user.id,))
        r = cur.fetchone()
        total_lent = float(r[0])
        total_borrowed = float(r[1])
        return LendingOverview(
            total_lent=total_lent,
            total_borrowed=total_borrowed,
            net_receivable=round(total_lent - total_borrowed, 2),
            overdue_count=r[2],
            open_count=r[3],
            closed_count=r[4]
        )
