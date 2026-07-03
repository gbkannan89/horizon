import logging
from datetime import date
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..database import get_db
from ..schemas import (
    UserOut, CollectionCreate, CollectionOut,
    CollectionMemberCreate, CollectionMemberOut, MemberPayment
)
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/collections", tags=["Collections"])


def _recalc_collection_totals(collection_id: int, conn):
    with conn.cursor() as cur:
        cur.execute("""
            SELECT COALESCE(SUM(expected_amount), 0),
                   COALESCE(SUM(paid_amount), 0),
                   COUNT(*),
                   COUNT(*) FILTER (WHERE status = 'paid')
            FROM collection_members WHERE collection_id = %s
        """, (collection_id,))
        total_exp, total_paid, member_count, paid_count = cur.fetchone()
        cur.execute("""
            UPDATE collections
            SET total_expected = %s, total_collected = %s
            WHERE id = %s
        """, (float(total_exp), float(total_paid), collection_id))
        conn.commit()
        return {
            "total_expected": float(total_exp),
            "total_collected": float(total_paid),
            "member_count": member_count,
            "paid_count": paid_count,
        }


@router.get("", response_model=List[CollectionOut])
def list_collections(
    status_filter: Optional[str] = None,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            if status_filter:
                cur.execute("""
                    SELECT id, user_id, label, description, total_expected, total_collected, status, created_at
                    FROM collections WHERE user_id = %s AND status = %s
                    ORDER BY created_at DESC
                """, (current_user.id, status_filter))
            else:
                cur.execute("""
                    SELECT id, user_id, label, description, total_expected, total_collected, status, created_at
                    FROM collections WHERE user_id = %s
                    ORDER BY created_at DESC
                """, (current_user.id,))

            rows = cur.fetchall()
            results = []
            for r in rows:
                cur.execute("""
                    SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'paid')
                    FROM collection_members WHERE collection_id = %s
                """, (r[0],))
                mc, pc = cur.fetchone()
                results.append(CollectionOut(
                    id=r[0], user_id=r[1], label=r[2], description=r[3],
                    total_expected=float(r[4]), total_collected=float(r[5]),
                    status=r[6], created_at=r[7],
                    member_count=mc, paid_count=pc,
                ))
            return results
    except Exception as e:
        logger.error(f"List collections error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("", response_model=CollectionOut, status_code=201)
def create_collection(
    data: CollectionCreate,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("""
                INSERT INTO collections (user_id, label, description, total_expected, total_collected, status)
                VALUES (%s, %s, %s, 0, 0, 'active')
                RETURNING id, user_id, label, description, total_expected, total_collected, status, created_at
            """, (current_user.id, data.label, data.description))
            row = cur.fetchone()

            collection_id = row[0]
            for member in data.members:
                cur.execute("""
                    INSERT INTO collection_members (collection_id, name, expected_amount, paid_amount, paid_date, status, notes)
                    VALUES (%s, %s, %s, %s, %s, %s, %s)
                """, (collection_id, member.name, member.expected_amount,
                      member.paid_amount, member.paid_date, member.status, member.notes))

            totals = _recalc_collection_totals(collection_id, conn)
            cur.execute("""
                SELECT id, collection_id, name, expected_amount, paid_amount, paid_date, status, notes, created_at
                FROM collection_members WHERE collection_id = %s ORDER BY id
            """, (collection_id,))
            members = [
                CollectionMemberOut(
                    id=m[0], collection_id=m[1], name=m[2],
                    expected_amount=float(m[3]), paid_amount=float(m[4]),
                    paid_date=m[5], status=m[6], notes=m[7], created_at=m[8]
                ) for m in cur.fetchall()
            ]

            return CollectionOut(
                id=row[0], user_id=row[1], label=row[2], description=row[3],
                total_expected=totals["total_expected"],
                total_collected=totals["total_collected"],
                status=row[6], created_at=row[7],
                member_count=totals["member_count"],
                paid_count=totals["paid_count"],
                members=members,
            )
    except Exception as e:
        conn.rollback()
        logger.error(f"Create collection error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/{collection_id}", response_model=CollectionOut)
def get_collection(
    collection_id: int,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("""
                SELECT id, user_id, label, description, total_expected, total_collected, status, created_at
                FROM collections WHERE id = %s AND user_id = %s
            """, (collection_id, current_user.id))
            row = cur.fetchone()
            if not row:
                raise HTTPException(status_code=404, detail="Collection not found")

            cur.execute("""
                SELECT id, collection_id, name, expected_amount, paid_amount, paid_date, status, notes, created_at
                FROM collection_members WHERE collection_id = %s ORDER BY status, id
            """, (collection_id,))
            members = [
                CollectionMemberOut(
                    id=m[0], collection_id=m[1], name=m[2],
                    expected_amount=float(m[3]), paid_amount=float(m[4]),
                    paid_date=m[5], status=m[6], notes=m[7], created_at=m[8]
                ) for m in cur.fetchall()
            ]

            total_exp = sum(m.expected_amount for m in members)
            total_paid = sum(m.paid_amount for m in members)
            paid_count = sum(1 for m in members if m.status == "paid")

            return CollectionOut(
                id=row[0], user_id=row[1], label=row[2], description=row[3],
                total_expected=total_exp, total_collected=total_paid,
                status=row[6], created_at=row[7],
                member_count=len(members), paid_count=paid_count,
                members=members,
            )
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Get collection error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.put("/{collection_id}", response_model=CollectionOut)
def update_collection(
    collection_id: int,
    label: Optional[str] = None,
    description: Optional[str] = None,
    status: Optional[str] = None,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("SELECT id FROM collections WHERE id = %s AND user_id = %s", (collection_id, current_user.id))
            if not cur.fetchone():
                raise HTTPException(status_code=404, detail="Collection not found")

            fields = []
            params = []
            if label is not None:
                fields.append("label = %s")
                params.append(label)
            if description is not None:
                fields.append("description = %s")
                params.append(description)
            if status is not None:
                fields.append("status = %s")
                params.append(status)

            if fields:
                params.append(collection_id)
                cur.execute(f"UPDATE collections SET {', '.join(fields)} WHERE id = %s", tuple(params))
                conn.commit()

            return get_collection(collection_id, current_user, conn)
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        logger.error(f"Update collection error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{collection_id}", status_code=204)
def delete_collection(
    collection_id: int,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("SELECT id FROM collections WHERE id = %s AND user_id = %s", (collection_id, current_user.id))
            if not cur.fetchone():
                raise HTTPException(status_code=404, detail="Collection not found")
            cur.execute("DELETE FROM collections WHERE id = %s", (collection_id,))
            conn.commit()
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        logger.error(f"Delete collection error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ── MEMBERS ────────────────────────────────────────────────────────────────

@router.post("/{collection_id}/members", response_model=List[CollectionMemberOut], status_code=201)
def add_members(
    collection_id: int,
    members: List[CollectionMemberCreate],
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("SELECT id FROM collections WHERE id = %s AND user_id = %s", (collection_id, current_user.id))
            if not cur.fetchone():
                raise HTTPException(status_code=404, detail="Collection not found")

            created = []
            for m in members:
                cur.execute("""
                    INSERT INTO collection_members (collection_id, name, expected_amount, paid_amount, paid_date, status, notes)
                    VALUES (%s, %s, %s, %s, %s, %s, %s)
                    RETURNING id, collection_id, name, expected_amount, paid_amount, paid_date, status, notes, created_at
                """, (collection_id, m.name, m.expected_amount, m.paid_amount, m.paid_date, m.status, m.notes))
                row = cur.fetchone()
                created.append(CollectionMemberOut(
                    id=row[0], collection_id=row[1], name=row[2],
                    expected_amount=float(row[3]), paid_amount=float(row[4]),
                    paid_date=row[5], status=row[6], notes=row[7], created_at=row[8]
                ))

            _recalc_collection_totals(collection_id, conn)
            return created
    except Exception as e:
        conn.rollback()
        logger.error(f"Add members error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.put("/members/{member_id}", response_model=CollectionMemberOut)
def record_payment(
    member_id: int,
    payment: MemberPayment,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("""
                SELECT cm.id, cm.collection_id, c.user_id
                FROM collection_members cm
                JOIN collections c ON c.id = cm.collection_id
                WHERE cm.id = %s AND c.user_id = %s
            """, (member_id, current_user.id))
            row = cur.fetchone()
            if not row:
                raise HTTPException(status_code=404, detail="Member not found")

            collection_id = row[1]
            cur.execute("SELECT expected_amount FROM collection_members WHERE id = %s", (member_id,))
            expected = float(cur.fetchone()[0])
            new_status = "paid" if payment.paid_amount >= expected else "partial"

            cur.execute("""
                UPDATE collection_members
                SET paid_amount = %s, paid_date = %s, status = %s
                WHERE id = %s
                RETURNING id, collection_id, name, expected_amount, paid_amount, paid_date, status, notes, created_at
            """, (payment.paid_amount, payment.paid_date, new_status, member_id))
            row = cur.fetchone()

            _recalc_collection_totals(collection_id, conn)

            return CollectionMemberOut(
                id=row[0], collection_id=row[1], name=row[2],
                expected_amount=float(row[3]), paid_amount=float(row[4]),
                paid_date=row[5], status=row[6], notes=row[7], created_at=row[8]
            )
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        logger.error(f"Record payment error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.delete("/members/{member_id}", status_code=204)
def remove_member(
    member_id: int,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("""
                SELECT cm.id, cm.collection_id
                FROM collection_members cm
                JOIN collections c ON c.id = cm.collection_id
                WHERE cm.id = %s AND c.user_id = %s
            """, (member_id, current_user.id))
            row = cur.fetchone()
            if not row:
                raise HTTPException(status_code=404, detail="Member not found")

            collection_id = row[1]
            cur.execute("DELETE FROM collection_members WHERE id = %s", (member_id,))
            _recalc_collection_totals(collection_id, conn)
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        logger.error(f"Remove member error: {e}")
        raise HTTPException(status_code=500, detail=str(e))
