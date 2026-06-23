from fastapi import APIRouter, Depends, HTTPException, status
from typing import List
from ..database import get_db
from .auth import get_current_user
from ..schemas import UserOut, InsuranceCreate, InsuranceOut

router = APIRouter(prefix="/insurance", tags=["Insurance"])

@router.get("/", response_model=List[InsuranceOut])
def get_insurances(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute(
                """
                SELECT id, user_id, type, provider, policy_name, premium_amount,
                       premium_frequency, coverage_amount, renewal_date, created_at
                FROM insurances
                WHERE user_id = %s
                ORDER BY renewal_date ASC NULLS LAST
                """,
                (current_user.id,)
            )
            return [InsuranceOut(
                id=r[0], user_id=r[1], type=r[2], provider=r[3], policy_name=r[4],
                premium_amount=float(r[5]), premium_frequency=r[6], coverage_amount=float(r[7]) if r[7] else None,
                renewal_date=r[8], created_at=r[9]
            ) for r in cur.fetchall()]
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/", response_model=InsuranceOut, status_code=status.HTTP_201_CREATED)
def add_insurance(
    ins: InsuranceCreate,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute(
                """
                INSERT INTO insurances (user_id, type, provider, policy_name, premium_amount, premium_frequency, coverage_amount, renewal_date)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, user_id, type, provider, policy_name, premium_amount, premium_frequency, coverage_amount, renewal_date, created_at
                """,
                (current_user.id, ins.type, ins.provider, ins.policy_name, ins.premium_amount, ins.premium_frequency, ins.coverage_amount, ins.renewal_date)
            )
            r = cur.fetchone()
            conn.commit()
            return InsuranceOut(
                id=r[0], user_id=r[1], type=r[2], provider=r[3], policy_name=r[4],
                premium_amount=float(r[5]), premium_frequency=r[6], coverage_amount=float(r[7]) if r[7] else None,
                renewal_date=r[8], created_at=r[9]
            )
    except Exception as e:
        conn.rollback()
        raise HTTPException(status_code=500, detail=str(e))

@router.delete("/{ins_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_insurance(
    ins_id: int,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("DELETE FROM insurances WHERE id = %s AND user_id = %s RETURNING id", (ins_id, current_user.id))
            if not cur.fetchone():
                raise HTTPException(status_code=404, detail="Insurance not found")
            conn.commit()
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        raise HTTPException(status_code=500, detail=str(e))
