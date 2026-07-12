import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from datetime import datetime, date as date_type, timedelta
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.electronics import (
    ElectronicCreate, ElectronicOut,
    ElectronicServiceCreate, ElectronicServiceOut,
    ElectronicEmiCreate, ElectronicEmiOut
)
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/electronics", tags=["electronics"])


def compute_depreciation(purchase_amount: float, purchase_date: date_type, expected_life_years: int) -> float:
    today = date_type.today()
    age_days = (today - purchase_date).days
    age_years = max(0.0, age_days / 365.25)
    if age_years >= expected_life_years:
        return 0.0
    return max(0.0, purchase_amount * (1.0 - (age_years / expected_life_years)))


def compute_electronic_stats(device_id: int, purchase_amount: float, purchase_date: date_type, expected_life_years: int, warranty_expiry: Optional[date_type], cur) -> dict:
    curr_val = compute_depreciation(purchase_amount, purchase_date, expected_life_years)
    
    # Warranty calculations
    today = date_type.today()
    days_rem = 0
    w_status = "expired"
    
    if warranty_expiry:
        days_rem = (warranty_expiry - today).days
        if days_rem <= 0:
            w_status = "expired"
            days_rem = 0
        elif days_rem <= 30:
            w_status = "near_expiry"
        else:
            w_status = "active"
            
    # Total service cost
    cur.execute("SELECT SUM(cost) FROM electronics_service_log WHERE electronic_id = %s", (device_id,))
    total_service = cur.fetchone()[0] or 0.0
    
    # EMI details
    cur.execute(
        """
        SELECT id, electronic_id, bank_name, emi_amount, interest_rate, total_months, months_paid, start_date, start_immediately, created_at
        FROM electronics_emi WHERE electronic_id = %s
        """,
        (device_id,)
    )
    emi_row = cur.fetchone()
    emi_out = None
    if emi_row:
        emi_out = ElectronicEmiOut(
            id=emi_row[0], electronic_id=emi_row[1], bank_name=emi_row[2], emi_amount=float(emi_row[3]),
            interest_rate=float(emi_row[4]), total_months=emi_row[5], months_paid=emi_row[6],
            start_date=emi_row[7], start_immediately=emi_row[8], created_at=emi_row[9]
        )
        
    return {
        "current_value": float(curr_val),
        "warranty_status": w_status,
        "warranty_days_remaining": days_rem,
        "total_service_cost": float(total_service),
        "emi": emi_out
    }


# ─── PRIMARY ELECTRONICS ENDPOINTS ───────────────────────────────────────────

@router.post("", response_model=ElectronicOut, status_code=status.HTTP_201_CREATED)
def add_device(payload: ElectronicCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        try:
            # compute warranty expiry
            warranty_expiry = payload.purchase_date + timedelta(days=payload.warranty_years * 365)
            
            cur.execute(
                """
                INSERT INTO electronics (
                    user_id, name, category, brand, model, purchase_date, purchase_amount,
                    warranty_years, warranty_expiry_date, expected_life_years, notes
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, user_id, name, category, brand, model, purchase_date, purchase_amount,
                          warranty_years, warranty_expiry_date, expected_life_years, notes, created_at
                """,
                (
                    current_user.id, payload.name, payload.category, payload.brand, payload.model,
                    payload.purchase_date, payload.purchase_amount, payload.warranty_years,
                    warranty_expiry, payload.expected_life_years, payload.notes
                )
            )
            r = cur.fetchone()
            db.commit()
            
            stats = compute_electronic_stats(r[0], float(r[7]), r[6], r[10], r[9], cur)
            
            return ElectronicOut(
                id=r[0], user_id=r[1], name=r[2], category=r[3], brand=r[4], model=r[5],
                purchase_date=r[6], purchase_amount=float(r[7]), warranty_years=r[8],
                warranty_expiry_date=r[9], expected_life_years=r[10], notes=r[11], created_at=r[12],
                current_value=stats["current_value"], warranty_status=stats["warranty_status"],
                warranty_days_remaining=stats["warranty_days_remaining"],
                total_service_cost=stats["total_service_cost"], emi=stats["emi"]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.get("", response_model=List[ElectronicOut])
def list_devices(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT id, user_id, name, category, brand, model, purchase_date, purchase_amount,
                   warranty_years, warranty_expiry_date, expected_life_years, notes, created_at
            FROM electronics WHERE user_id = %s ORDER BY purchase_date DESC, id DESC
            """,
            (current_user.id,)
        )
        rows = cur.fetchall()
        result = []
        for r in rows:
            stats = compute_electronic_stats(r[0], float(r[7]), r[6], r[10], r[9], cur)
            
            # Update straight-line depreciation inside DB for persistent storage accuracy
            cur.execute(
                "UPDATE electronics SET current_value = %s WHERE id = %s",
                (stats["current_value"], r[0])
            )
            
            result.append(
                ElectronicOut(
                    id=r[0], user_id=r[1], name=r[2], category=r[3], brand=r[4], model=r[5],
                    purchase_date=r[6], purchase_amount=float(r[7]), warranty_years=r[8],
                    warranty_expiry_date=r[9], expected_life_years=r[10], notes=r[11], created_at=r[12],
                    current_value=stats["current_value"], warranty_status=stats["warranty_status"],
                    warranty_days_remaining=stats["warranty_days_remaining"],
                    total_service_cost=stats["total_service_cost"], emi=stats["emi"]
                )
            )
        db.commit()
        return result


@router.put("/{device_id}", response_model=ElectronicOut)
def update_device(device_id: int, payload: ElectronicCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Check ownership
        cur.execute("SELECT id FROM electronics WHERE id = %s AND user_id = %s", (device_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Device not found")
            
        try:
            warranty_expiry = payload.purchase_date + timedelta(days=payload.warranty_years * 365)
            
            cur.execute(
                """
                UPDATE electronics
                SET name = %s, category = %s, brand = %s, model = %s, purchase_date = %s, purchase_amount = %s,
                    warranty_years = %s, warranty_expiry_date = %s, expected_life_years = %s, notes = %s
                WHERE id = %s
                RETURNING id, user_id, name, category, brand, model, purchase_date, purchase_amount,
                          warranty_years, warranty_expiry_date, expected_life_years, notes, created_at
                """,
                (
                    payload.name, payload.category, payload.brand, payload.model, payload.purchase_date,
                    payload.purchase_amount, payload.warranty_years, warranty_expiry, payload.expected_life_years,
                    payload.notes, device_id
                )
            )
            r = cur.fetchone()
            db.commit()
            
            stats = compute_electronic_stats(r[0], float(r[7]), r[6], r[10], r[9], cur)
            
            return ElectronicOut(
                id=r[0], user_id=r[1], name=r[2], category=r[3], brand=r[4], model=r[5],
                purchase_date=r[6], purchase_amount=float(r[7]), warranty_years=r[8],
                warranty_expiry_date=r[9], expected_life_years=r[10], notes=r[11], created_at=r[12],
                current_value=stats["current_value"], warranty_status=stats["warranty_status"],
                warranty_days_remaining=stats["warranty_days_remaining"],
                total_service_cost=stats["total_service_cost"], emi=stats["emi"]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{device_id}", status_code=status.HTTP_200_OK)
def delete_device(device_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("DELETE FROM electronics WHERE id = %s AND user_id = %s", (device_id, current_user.id))
        db.commit()
        return {"message": "Device deleted successfully"}


# ─── ELECTRONIC SERVICE ENDPOINTS ────────────────────────────────────────────

@router.get("/{device_id}/service", response_model=List[ElectronicServiceOut])
def list_device_services(device_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM electronics WHERE id = %s AND user_id = %s", (device_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Device not found")
            
        cur.execute(
            """
            SELECT id, electronic_id, service_date, service_type, description, cost, service_center, created_at
            FROM electronics_service_log WHERE electronic_id = %s ORDER BY service_date DESC, id DESC
            """,
            (device_id,)
        )
        return [
            ElectronicServiceOut(
                id=r[0], electronic_id=r[1], service_date=r[2], service_type=r[3], description=r[4],
                cost=float(r[5]), service_center=r[6], created_at=r[7]
            )
            for r in cur.fetchall()
        ]


@router.post("/{device_id}/service", response_model=ElectronicServiceOut, status_code=status.HTTP_201_CREATED)
def add_device_service(device_id: int, payload: ElectronicServiceCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM electronics WHERE id = %s AND user_id = %s", (device_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Device not found")
            
        try:
            cur.execute(
                """
                INSERT INTO electronics_service_log (
                    electronic_id, service_date, service_type, description, cost, service_center
                ) VALUES (%s, %s, %s, %s, %s, %s)
                RETURNING id, electronic_id, service_date, service_type, description, cost, service_center, created_at
                """,
                (
                    device_id, payload.service_date, payload.service_type, payload.description, payload.cost,
                    payload.service_center
                )
            )
            r = cur.fetchone()
            db.commit()
            return ElectronicServiceOut(
                id=r[0], electronic_id=r[1], service_date=r[2], service_type=r[3], description=r[4],
                cost=float(r[5]), service_center=r[6], created_at=r[7]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{device_id}/service/{sid}", status_code=status.HTTP_200_OK)
def delete_device_service(device_id: int, sid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT es.id FROM electronics_service_log es
            JOIN electronics e ON es.electronic_id = e.id
            WHERE es.id = %s AND es.electronic_id = %s AND e.user_id = %s
            """,
            (sid, device_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Service log not found")
            
        cur.execute("DELETE FROM electronics_service_log WHERE id = %s", (sid,))
        db.commit()
        return {"message": "Service log deleted"}


# ─── ELECTRONIC EMI ENDPOINTS ────────────────────────────────────────────────

@router.get("/{device_id}/emi", response_model=Optional[ElectronicEmiOut])
def get_device_emi(device_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM electronics WHERE id = %s AND user_id = %s", (device_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Device not found")
            
        cur.execute(
            """
            SELECT id, electronic_id, bank_name, emi_amount, interest_rate, total_months, months_paid, start_date, start_immediately, created_at
            FROM electronics_emi WHERE electronic_id = %s
            """,
            (device_id,)
        )
        r = cur.fetchone()
        if not r:
            return None
        return ElectronicEmiOut(
            id=r[0], electronic_id=r[1], bank_name=r[2], emi_amount=float(r[3]), interest_rate=float(r[4]),
            total_months=r[5], months_paid=r[6], start_date=r[7], start_immediately=r[8], created_at=r[9]
        )


@router.post("/{device_id}/emi", response_model=ElectronicEmiOut, status_code=status.HTTP_201_CREATED)
def add_device_emi(device_id: int, payload: ElectronicEmiCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM electronics WHERE id = %s AND user_id = %s", (device_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Device not found")
            
        # Ensure only one EMI schedule per device
        cur.execute("SELECT id FROM electronics_emi WHERE electronic_id = %s", (device_id,))
        if cur.fetchone():
            raise HTTPException(status_code=400, detail="Device already has an EMI schedule")
            
        try:
            cur.execute(
                """
                INSERT INTO electronics_emi (
                    electronic_id, bank_name, emi_amount, interest_rate, total_months, months_paid, start_date, start_immediately
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, electronic_id, bank_name, emi_amount, interest_rate, total_months, months_paid, start_date, start_immediately, created_at
                """,
                (
                    device_id, payload.bank_name, payload.emi_amount, payload.interest_rate or 0.0,
                    payload.total_months, payload.months_paid or 0, payload.start_date, payload.start_immediately
                )
            )
            r = cur.fetchone()
            db.commit()
            return ElectronicEmiOut(
                id=r[0], electronic_id=r[1], bank_name=r[2], emi_amount=float(r[3]), interest_rate=float(r[4]),
                total_months=r[5], months_paid=r[6], start_date=r[7], start_immediately=r[8], created_at=r[9]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.put("/{device_id}/emi", response_model=ElectronicEmiOut)
def update_device_emi(device_id: int, payload: ElectronicEmiCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM electronics WHERE id = %s AND user_id = %s", (device_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Device not found")
            
        try:
            cur.execute(
                """
                UPDATE electronics_emi
                SET bank_name = %s, emi_amount = %s, interest_rate = %s, total_months = %s,
                    months_paid = %s, start_date = %s, start_immediately = %s
                WHERE electronic_id = %s
                RETURNING id, electronic_id, bank_name, emi_amount, interest_rate, total_months, months_paid, start_date, start_immediately, created_at
                """,
                (
                    payload.bank_name, payload.emi_amount, payload.interest_rate or 0.0, payload.total_months,
                    payload.months_paid or 0, payload.start_date, payload.start_immediately, device_id
                )
            )
            r = cur.fetchone()
            if not r:
                raise HTTPException(status_code=404, detail="No EMI schedule found to update")
            db.commit()
            return ElectronicEmiOut(
                id=r[0], electronic_id=r[1], bank_name=r[2], emi_amount=float(r[3]), interest_rate=float(r[4]),
                total_months=r[5], months_paid=r[6], start_date=r[7], start_immediately=r[8], created_at=r[9]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{device_id}/emi", status_code=status.HTTP_200_OK)
def delete_device_emi(device_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM electronics WHERE id = %s AND user_id = %s", (device_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Device not found")
            
        cur.execute("DELETE FROM electronics_emi WHERE electronic_id = %s", (device_id,))
        db.commit()
        return {"message": "EMI details deleted"}
