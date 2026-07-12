import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from datetime import datetime, date as date_type
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.assets import VehicleCreate, VehicleOut
from ..schemas.vehicles import (
    VehicleServiceCreate, VehicleServiceOut,
    VehicleFuelCreate, VehicleFuelOut,
    VehicleLoanCreate, VehicleLoanOut
)
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/assets/vehicles", tags=["vehicles"])


def suggest_idv(purchase_cost: float, age_years: int) -> float:
    if age_years <= 0:
        return purchase_cost
    elif age_years == 1:
        return purchase_cost * 0.85
    elif age_years <= 3:
        return purchase_cost * (0.85 - (age_years - 1) * 0.10)
    elif age_years <= 5:
        return purchase_cost * (0.65 - (age_years - 3) * 0.10)
    else:
        return purchase_cost * 0.50


def compute_vehicle_stats(vehicle_id: int, purchase_cost: float, purchase_year: Optional[int], model_year: Optional[int], km_driven: float, fuel_cost_total: float, cur) -> dict:
    cost_per_km = fuel_cost_total / km_driven if km_driven > 0 else 0.0
    
    # Age calculation
    current_year = datetime.now().year
    age = current_year - (purchase_year or model_year or current_year)
    suggested = suggest_idv(purchase_cost, max(0, age))
    
    # Total service cost
    cur.execute("SELECT SUM(cost) FROM vehicle_service_log WHERE vehicle_id = %s", (vehicle_id,))
    total_service = cur.fetchone()[0] or 0.0
    
    # Next service due
    cur.execute("SELECT next_service_km FROM vehicle_service_log WHERE vehicle_id = %s ORDER BY service_date DESC, id DESC LIMIT 1", (vehicle_id,))
    last_service_row = cur.fetchone()
    next_service_due = None
    if last_service_row and last_service_row[0]:
        next_service_km = float(last_service_row[0])
        next_service_due = {
            "type": "km",
            "at_km": next_service_km,
            "remaining_km": max(0.0, next_service_km - km_driven)
        }
        
    return {
        "cost_per_km": cost_per_km,
        "suggested_idv": suggested,
        "total_service_cost": float(total_service),
        "next_service": next_service_due
    }


# ─── PRIMARY VEHICLE ENDPOINTS ───────────────────────────────────────────────

@router.post("", response_model=VehicleOut, status_code=status.HTTP_201_CREATED)
def add_vehicle(payload: VehicleCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        try:
            cur.execute(
                """
                INSERT INTO vehicles (
                    user_id, make_model, purchase_cost, insurance_renewal_date,
                    model_year, purchase_year, fuel_type, mileage_kmpl,
                    fuel_cost_total, km_driven, insurance_idv, insurance_renewal_amount, registration_number
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, user_id, make_model, purchase_cost, insurance_renewal_date,
                          model_year, purchase_year, fuel_type, mileage_kmpl,
                          fuel_cost_total, km_driven, insurance_idv, insurance_renewal_amount, registration_number, created_at
                """,
                (
                    current_user.id, payload.make_model, payload.purchase_cost, payload.insurance_renewal_date,
                    payload.model_year, payload.purchase_year, payload.fuel_type, payload.mileage_kmpl,
                    payload.fuel_cost_total or 0.0, payload.km_driven or 0.0, payload.insurance_idv,
                    payload.insurance_renewal_amount, payload.registration_number
                )
            )
            r = cur.fetchone()
            db.commit()
            
            stats = compute_vehicle_stats(r[0], float(r[3]), r[6], r[5], float(r[10]), float(r[9]), cur)
            
            return VehicleOut(
                id=r[0], user_id=r[1], make_model=r[2], purchase_cost=float(r[3]), insurance_renewal_date=r[4],
                model_year=r[5], purchase_year=r[6], fuel_type=r[7], mileage_kmpl=float(r[8]) if r[8] else None,
                fuel_cost_total=float(r[9]), km_driven=float(r[10]), insurance_idv=float(r[11]) if r[11] else None,
                insurance_renewal_amount=float(r[12]) if r[12] else None, registration_number=r[13], created_at=r[14],
                **stats
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.get("", response_model=List[VehicleOut])
def list_vehicles(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT id, user_id, make_model, purchase_cost, insurance_renewal_date,
                   model_year, purchase_year, fuel_type, mileage_kmpl,
                   fuel_cost_total, km_driven, insurance_idv, insurance_renewal_amount, registration_number, created_at
            FROM vehicles WHERE user_id = %s ORDER BY id DESC
            """,
            (current_user.id,)
        )
        rows = cur.fetchall()
        result = []
        for r in rows:
            stats = compute_vehicle_stats(r[0], float(r[3]), r[6], r[5], float(r[10]), float(r[9]), cur)
            result.append(
                VehicleOut(
                    id=r[0], user_id=r[1], make_model=r[2], purchase_cost=float(r[3]), insurance_renewal_date=r[4],
                    model_year=r[5], purchase_year=r[6], fuel_type=r[7], mileage_kmpl=float(r[8]) if r[8] else None,
                    fuel_cost_total=float(r[9]), km_driven=float(r[10]), insurance_idv=float(r[11]) if r[11] else None,
                    insurance_renewal_amount=float(r[12]) if r[12] else None, registration_number=r[13], created_at=r[14],
                    **stats
                )
            )
        return result


@router.delete("/{vehicle_id}", status_code=status.HTTP_200_OK)
def delete_vehicle(vehicle_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("DELETE FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        db.commit()
        return {"message": "Vehicle deleted successfully"}


# ─── VEHICLE SERVICE ENDPOINTS ───────────────────────────────────────────────

@router.get("/{vehicle_id}/service", response_model=List[VehicleServiceOut])
def list_service_records(vehicle_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Check ownership
        cur.execute("SELECT id FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vehicle not found")
            
        cur.execute(
            """
            SELECT id, vehicle_id, service_date, service_type, description, cost, km_at_service, service_center, next_service_km, created_at
            FROM vehicle_service_log WHERE vehicle_id = %s ORDER BY service_date DESC, id DESC
            """,
            (vehicle_id,)
        )
        return [
            VehicleServiceOut(
                id=r[0], vehicle_id=r[1], service_date=r[2], service_type=r[3], description=r[4],
                cost=float(r[5]), km_at_service=float(r[6]) if r[6] else None, service_center=r[7],
                next_service_km=float(r[8]) if r[8] else None, created_at=r[9]
            )
            for r in cur.fetchall()
        ]


@router.post("/{vehicle_id}/service", response_model=VehicleServiceOut, status_code=status.HTTP_201_CREATED)
def add_service_record(vehicle_id: int, payload: VehicleServiceCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vehicle not found")
            
        try:
            cur.execute(
                """
                INSERT INTO vehicle_service_log (
                    vehicle_id, service_date, service_type, description, cost, km_at_service, service_center, next_service_km
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, vehicle_id, service_date, service_type, description, cost, km_at_service, service_center, next_service_km, created_at
                """,
                (
                    vehicle_id, payload.service_date, payload.service_type, payload.description, payload.cost,
                    payload.km_at_service, payload.service_center, payload.next_service_km
                )
            )
            r = cur.fetchone()
            
            # If km_at_service is provided, update vehicle's current km_driven
            if payload.km_at_service:
                cur.execute(
                    "UPDATE vehicles SET km_driven = GREATEST(km_driven, %s) WHERE id = %s",
                    (payload.km_at_service, vehicle_id)
                )
                
            db.commit()
            return VehicleServiceOut(
                id=r[0], vehicle_id=r[1], service_date=r[2], service_type=r[3], description=r[4],
                cost=float(r[5]), km_at_service=float(r[6]) if r[6] else None, service_center=r[7],
                next_service_km=float(r[8]) if r[8] else None, created_at=r[9]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{vehicle_id}/service/{sid}", status_code=status.HTTP_200_OK)
def delete_service_record(vehicle_id: int, sid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT vs.id FROM vehicle_service_log vs
            JOIN vehicles v ON vs.vehicle_id = v.id
            WHERE vs.id = %s AND vs.vehicle_id = %s AND v.user_id = %s
            """,
            (sid, vehicle_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Service record not found")
            
        cur.execute("DELETE FROM vehicle_service_log WHERE id = %s", (sid,))
        db.commit()
        return {"message": "Service record deleted"}


# ─── VEHICLE FUEL ENDPOINTS ──────────────────────────────────────────────────

@router.get("/{vehicle_id}/fuel", response_model=List[VehicleFuelOut])
def list_fuel_records(vehicle_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vehicle not found")
            
        cur.execute(
            """
            SELECT id, vehicle_id, fill_date, amount, liters, km_at_fill, price_per_liter, is_full_tank, created_at
            FROM vehicle_fuel_log WHERE vehicle_id = %s ORDER BY fill_date DESC, id DESC
            """,
            (vehicle_id,)
        )
        return [
            VehicleFuelOut(
                id=r[0], vehicle_id=r[1], fill_date=r[2], amount=float(r[3]), liters=float(r[4]),
                km_at_fill=float(r[5]) if r[5] else None, price_per_liter=float(r[6]) if r[6] else None,
                is_full_tank=r[7], created_at=r[8]
            )
            for r in cur.fetchall()
        ]


@router.post("/{vehicle_id}/fuel", response_model=VehicleFuelOut, status_code=status.HTTP_201_CREATED)
def add_fuel_record(vehicle_id: int, payload: VehicleFuelCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vehicle not found")
            
        try:
            price_pl = payload.price_per_liter or (payload.amount / payload.liters if payload.liters > 0 else 0.0)
            cur.execute(
                """
                INSERT INTO vehicle_fuel_log (
                    vehicle_id, fill_date, amount, liters, km_at_fill, price_per_liter, is_full_tank
                ) VALUES (%s, %s, %s, %s, %s, %s, %s)
                RETURNING id, vehicle_id, fill_date, amount, liters, km_at_fill, price_per_liter, is_full_tank, created_at
                """,
                (
                    vehicle_id, payload.fill_date, payload.amount, payload.liters, payload.km_at_fill,
                    price_pl, payload.is_full_tank
                )
            )
            r = cur.fetchone()
            
            # Update aggregate fuel cost and check if km_at_fill pushes km_driven higher
            cur.execute(
                """
                UPDATE vehicles
                SET fuel_cost_total = fuel_cost_total + %s,
                    km_driven = GREATEST(km_driven, COALESCE(%s, km_driven))
                WHERE id = %s
                """,
                (payload.amount, payload.km_at_fill, vehicle_id)
            )
            
            db.commit()
            return VehicleFuelOut(
                id=r[0], vehicle_id=r[1], fill_date=r[2], amount=float(r[3]), liters=float(r[4]),
                km_at_fill=float(r[5]) if r[5] else None, price_per_liter=float(r[6]) if r[6] else None,
                is_full_tank=r[7], created_at=r[8]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{vehicle_id}/fuel/{fid}", status_code=status.HTTP_200_OK)
def delete_fuel_record(vehicle_id: int, fid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Check ownership and fetch amount to deduct
        cur.execute(
            """
            SELECT vf.id, vf.amount FROM vehicle_fuel_log vf
            JOIN vehicles v ON vf.vehicle_id = v.id
            WHERE vf.id = %s AND vf.vehicle_id = %s AND v.user_id = %s
            """,
            (fid, vehicle_id, current_user.id)
        )
        row = cur.fetchone()
        if not row:
            raise HTTPException(status_code=404, detail="Fuel record not found")
            
        amt = float(row[1])
        cur.execute("DELETE FROM vehicle_fuel_log WHERE id = %s", (fid,))
        cur.execute("UPDATE vehicles SET fuel_cost_total = GREATEST(0.00, fuel_cost_total - %s) WHERE id = %s", (amt, vehicle_id))
        
        db.commit()
        return {"message": "Fuel record deleted"}


# ─── VEHICLE LOAN ENDPOINTS ──────────────────────────────────────────────────

@router.get("/{vehicle_id}/loan", response_model=Optional[VehicleLoanOut])
def get_vehicle_loan(vehicle_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vehicle not found")
            
        cur.execute(
            """
            SELECT id, vehicle_id, bank_name, loan_amount, interest_rate, tenure_months, emi, start_date, emi_paid, outstanding, created_at
            FROM vehicle_loans WHERE vehicle_id = %s
            """,
            (vehicle_id,)
        )
        r = cur.fetchone()
        if not r:
            return None
            
        return VehicleLoanOut(
            id=r[0], vehicle_id=r[1], bank_name=r[2], loan_amount=float(r[3]), interest_rate=float(r[4]),
            tenure_months=r[5], emi=float(r[6]), start_date=r[7], emi_paid=r[8], outstanding=float(r[9]), created_at=r[10]
        )


@router.post("/{vehicle_id}/loan", response_model=VehicleLoanOut, status_code=status.HTTP_201_CREATED)
def add_vehicle_loan(vehicle_id: int, payload: VehicleLoanCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vehicle not found")
            
        # Ensure only one loan per vehicle
        cur.execute("SELECT id FROM vehicle_loans WHERE vehicle_id = %s", (vehicle_id,))
        if cur.fetchone():
            raise HTTPException(status_code=400, detail="Vehicle already has a loan attached")
            
        try:
            cur.execute(
                """
                INSERT INTO vehicle_loans (
                    vehicle_id, bank_name, loan_amount, interest_rate, tenure_months, emi, start_date, emi_paid, outstanding
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, vehicle_id, bank_name, loan_amount, interest_rate, tenure_months, emi, start_date, emi_paid, outstanding, created_at
                """,
                (
                    vehicle_id, payload.bank_name, payload.loan_amount, payload.interest_rate, payload.tenure_months,
                    payload.emi, payload.start_date, payload.emi_paid or 0, payload.outstanding
                )
            )
            r = cur.fetchone()
            db.commit()
            return VehicleLoanOut(
                id=r[0], vehicle_id=r[1], bank_name=r[2], loan_amount=float(r[3]), interest_rate=float(r[4]),
                tenure_months=r[5], emi=float(r[6]), start_date=r[7], emi_paid=r[8], outstanding=float(r[9]), created_at=r[10]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.put("/{vehicle_id}/loan", response_model=VehicleLoanOut)
def update_vehicle_loan(vehicle_id: int, payload: VehicleLoanCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vehicle not found")
            
        try:
            cur.execute(
                """
                UPDATE vehicle_loans
                SET bank_name = %s, loan_amount = %s, interest_rate = %s, tenure_months = %s,
                    emi = %s, start_date = %s, emi_paid = %s, outstanding = %s
                WHERE vehicle_id = %s
                RETURNING id, vehicle_id, bank_name, loan_amount, interest_rate, tenure_months, emi, start_date, emi_paid, outstanding, created_at
                """,
                (
                    payload.bank_name, payload.loan_amount, payload.interest_rate, payload.tenure_months,
                    payload.emi, payload.start_date, payload.emi_paid or 0, payload.outstanding, vehicle_id
                )
            )
            r = cur.fetchone()
            if not r:
                raise HTTPException(status_code=404, detail="No loan found to update")
            db.commit()
            return VehicleLoanOut(
                id=r[0], vehicle_id=r[1], bank_name=r[2], loan_amount=float(r[3]), interest_rate=float(r[4]),
                tenure_months=r[5], emi=float(r[6]), start_date=r[7], emi_paid=r[8], outstanding=float(r[9]), created_at=r[10]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{vehicle_id}/loan", status_code=status.HTTP_200_OK)
def delete_vehicle_loan(vehicle_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM vehicles WHERE id = %s AND user_id = %s", (vehicle_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vehicle not found")
            
        cur.execute("DELETE FROM vehicle_loans WHERE vehicle_id = %s", (vehicle_id,))
        db.commit()
        return {"message": "Vehicle loan deleted"}
