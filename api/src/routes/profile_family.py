import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from datetime import datetime, date as date_type, timedelta
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.profile_family import (
    FamilyMemberCreate, FamilyMemberOut,
    MemberInsuranceCreate, MemberInsuranceOut,
    SchoolingCreate, SchoolingOut, SchoolingPaymentCreate, SchoolingPaymentOut,
    CheckupCreate, CheckupOut, MedicineCreate, MedicineOut,
    VaccinationCreate, VaccinationOut, EarningCreate, EarningOut
)
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/profile/family", tags=["family"])


# Helper to compute monthly equivalent
def get_monthly_equiv(amt: float, freq: str) -> float:
    f = freq.lower()
    if f == "monthly":
        return amt
    elif f == "quarterly":
        return amt / 3.0
    elif f == "yearly":
        return amt / 12.0
    return 0.0


# ─── AGGREGATED RECURRING COSTS ───────────────────────────────────────────────

@router.get("/recurring-costs")
def get_recurring_costs(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Get household or user's own members
        cur.execute("SELECT id FROM family_members WHERE user_id = %s", (current_user.id,))
        member_ids = [r[0] for r in cur.fetchall()]
        if not member_ids:
            return {
                "total_monthly": 0.0,
                "schooling": 0.0,
                "medicines": 0.0,
                "checkups": 0.0,
                "vaccinations": 0.0
            }
            
        m_placeholders = ",".join(["%s"] * len(member_ids))
        
        # 1. Schooling
        cur.execute(
            f"SELECT fee_amount, fee_frequency FROM member_schooling WHERE member_id IN ({m_placeholders})",
            member_ids
        )
        school_total = sum(get_monthly_equiv(float(r[0]), r[1]) for r in cur.fetchall())
        
        # 2. Medicines
        cur.execute(
            f"SELECT monthly_cost FROM member_medicines WHERE member_id IN ({m_placeholders}) AND is_regular = TRUE",
            member_ids
        )
        meds_total = sum(float(r[0]) for r in cur.fetchall())
        
        # 3. Checkups
        cur.execute(
            f"SELECT recurring_cost, frequency FROM member_checkups WHERE member_id IN ({m_placeholders})",
            member_ids
        )
        check_total = sum(get_monthly_equiv(float(r[0]), r[1]) for r in cur.fetchall())
        
        # 4. Vaccinations
        cur.execute(
            f"SELECT recurring_cost, frequency FROM member_vaccinations WHERE member_id IN ({m_placeholders})",
            member_ids
        )
        vacc_total = sum(get_monthly_equiv(float(r[0]), r[1]) for r in cur.fetchall())
        
        # 5. Insurance
        cur.execute(
            f"SELECT premium_amount, premium_frequency FROM member_insurances WHERE member_id IN ({m_placeholders})",
            member_ids
        )
        insurance_total = sum(get_monthly_equiv(float(r[0]), r[1]) for r in cur.fetchall())
        
        total = school_total + meds_total + check_total + vacc_total + insurance_total
        
        return {
            "total_monthly": total,
            "schooling": school_total,
            "medicines": meds_total,
            "checkups": check_total,
            "vaccinations": vacc_total,
            "insurance": insurance_total
        }


# ─── FAMILY MEMBER CRUD ───────────────────────────────────────────────────────

@router.get("", response_model=List[FamilyMemberOut])
def list_family_members(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT id, household_id, user_id, name, dob, blood_group, relationship, avatar_color, is_self, is_active, earning_status, contribution_amount, phone, email, created_at
            FROM family_members WHERE user_id = %s ORDER BY is_self DESC, id ASC
            """,
            (current_user.id,)
        )
        rows = cur.fetchall()
        result = []
        for r in rows:
            mid = r[0]
            # Fetch nested items
            cur.execute("SELECT id, member_id, institution_name, fee_amount, fee_frequency, last_paid_date, next_due_date, notes, created_at FROM member_schooling WHERE member_id = %s", (mid,))
            school_rows = cur.fetchall()
            schooling_list = [
                SchoolingOut(
                    id=sr[0], member_id=sr[1], institution_name=sr[2], fee_amount=float(sr[3]),
                    fee_frequency=sr[4], last_paid_date=sr[5], next_due_date=sr[6], notes=sr[7], created_at=sr[8],
                    monthly_equivalent=get_monthly_equiv(float(sr[3]), sr[4])
                )
                for sr in school_rows
            ]
            
            cur.execute("SELECT id, member_id, checkup_type, frequency, recurring_cost, last_checkup_date, next_due_date, notes, created_at FROM member_checkups WHERE member_id = %s", (mid,))
            check_rows = cur.fetchall()
            checkups_list = [
                CheckupOut(
                    id=cr[0], member_id=cr[1], checkup_type=cr[2], frequency=cr[3], recurring_cost=float(cr[4]),
                    last_checkup_date=cr[5], next_due_date=cr[6], notes=cr[7], created_at=cr[8]
                )
                for cr in check_rows
            ]
            
            cur.execute("SELECT id, member_id, medicine_name, monthly_cost, purpose, is_regular, prescribed_by, created_at FROM member_medicines WHERE member_id = %s", (mid,))
            meds_rows = cur.fetchall()
            medicines_list = [
                MedicineOut(
                    id=mr[0], member_id=mr[1], medicine_name=mr[2], monthly_cost=float(mr[3]),
                    purpose=mr[4], is_regular=mr[5], prescribed_by=mr[6], created_at=mr[7]
                )
                for mr in meds_rows
            ]
            
            cur.execute("SELECT id, member_id, vaccine_name, frequency, recurring_cost, last_vaccination_date, next_due_date, notes, created_at FROM member_vaccinations WHERE member_id = %s", (mid,))
            vacc_rows = cur.fetchall()
            vaccinations_list = [
                VaccinationOut(
                    id=vr[0], member_id=vr[1], vaccine_name=vr[2], frequency=vr[3], recurring_cost=float(vr[4]),
                    last_vaccination_date=vr[5], next_due_date=vr[6], notes=vr[7], created_at=vr[8]
                )
                for vr in vacc_rows
            ]
            
            cur.execute("SELECT id, member_id, monthly_income, contribution_to_household, occupation, created_at FROM member_earnings WHERE member_id = %s", (mid,))
            earn_row = cur.fetchone()
            earnings_out = None
            if earn_row:
                earnings_out = EarningOut(
                    id=earn_row[0], member_id=earn_row[1], monthly_income=float(earn_row[2]),
                    contribution_to_household=float(earn_row[3]), occupation=earn_row[4], created_at=earn_row[5]
                )
                
            cur.execute("SELECT id, member_id, provider_name, policy_details, premium_frequency, premium_amount, coverage_amount, created_at FROM member_insurances WHERE member_id = %s", (mid,))
            ins_rows = cur.fetchall()
            insurances_list = [
                MemberInsuranceOut(
                    id=ir[0], member_id=ir[1], provider_name=ir[2], policy_details=ir[3],
                    premium_frequency=ir[4], premium_amount=float(ir[5]), coverage_amount=float(ir[6]), created_at=ir[7]
                )
                for ir in ins_rows
            ]
                
            result.append(
                FamilyMemberOut(
                    id=r[0], household_id=r[1], user_id=r[2], name=r[3], dob=r[4], blood_group=r[5],
                    relationship=r[6], avatar_color=r[7], is_self=r[8], is_active=r[9],
                    earning_status=r[10], contribution_amount=float(r[11]) if r[11] else 0.0,
                    phone=r[12], email=r[13], created_at=r[14],
                    schooling=schooling_list, checkups=checkups_list, medicines=medicines_list,
                    vaccinations=vaccinations_list, earnings=earnings_out, insurances=insurances_list
                )
            )
        return result


@router.post("", response_model=FamilyMemberOut, status_code=status.HTTP_201_CREATED)
def add_family_member(payload: FamilyMemberCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        try:
            # Check household id of user
            cur.execute("SELECT household_id FROM users WHERE id = %s", (current_user.id,))
            house_row = cur.fetchone()
            house_id = house_row[0] if house_row else None
            
            # Ensure unique is_self
            if payload.is_self:
                cur.execute("UPDATE family_members SET is_self = FALSE WHERE user_id = %s", (current_user.id,))
                
            cur.execute(
                """
                INSERT INTO family_members (household_id, user_id, name, dob, blood_group, relationship, avatar_color, is_self, is_active, earning_status, contribution_amount, phone, email)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, household_id, user_id, name, dob, blood_group, relationship, avatar_color, is_self, is_active, earning_status, contribution_amount, phone, email, created_at
                """,
                (
                    house_id, current_user.id, payload.name, payload.dob, payload.blood_group,
                    payload.relationship, payload.avatar_color or "#0D9488", payload.is_self, payload.is_active,
                    payload.earning_status or False, payload.contribution_amount or 0.0, payload.phone, payload.email
                )
            )
            r = cur.fetchone()
            db.commit()
            return FamilyMemberOut(
                id=r[0], household_id=r[1], user_id=r[2], name=r[3], dob=r[4], blood_group=r[5],
                relationship=r[6], avatar_color=r[7], is_self=r[8], is_active=r[9],
                earning_status=r[10], contribution_amount=float(r[11]) if r[11] else 0.0,
                phone=r[12], email=r[13], created_at=r[14]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.put("/{member_id}", response_model=FamilyMemberOut)
def update_family_member(member_id: int, payload: FamilyMemberCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
            
        try:
            if payload.is_self:
                cur.execute("UPDATE family_members SET is_self = FALSE WHERE user_id = %s AND id != %s", (current_user.id, member_id))
                
            cur.execute(
                """
                UPDATE family_members
                SET name = %s, dob = %s, blood_group = %s, relationship = %s, avatar_color = %s, is_self = %s, is_active = %s, phone = %s, email = %s
                WHERE id = %s
                RETURNING id, household_id, user_id, name, dob, blood_group, relationship, avatar_color, is_self, is_active, earning_status, contribution_amount, phone, email, created_at
                """,
                (
                    payload.name, payload.dob, payload.blood_group, payload.relationship,
                    payload.avatar_color or "#0D9488", payload.is_self, payload.is_active,
                    payload.phone, payload.email, member_id
                )
            )
            r = cur.fetchone()
            db.commit()
            return FamilyMemberOut(
                id=r[0], household_id=r[1], user_id=r[2], name=r[3], dob=r[4], blood_group=r[5],
                relationship=r[6], avatar_color=r[7], is_self=r[8], is_active=r[9],
                earning_status=r[10], contribution_amount=float(r[11]) if r[11] else 0.0,
                phone=r[12], email=r[13], created_at=r[14]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{member_id}", status_code=status.HTTP_200_OK)
def delete_family_member(member_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("DELETE FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        db.commit()
        return {"message": "Family member deleted"}


# ─── NESTED ITEM CREATE ENDPOINTS ─────────────────────────────────────────────

@router.post("/{member_id}/schooling", response_model=SchoolingOut)
def add_schooling(member_id: int, payload: SchoolingCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
        try:
            cur.execute(
                """
                INSERT INTO member_schooling (member_id, institution_name, fee_amount, fee_frequency, last_paid_date, next_due_date, notes)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                RETURNING id, member_id, institution_name, fee_amount, fee_frequency, last_paid_date, next_due_date, notes, created_at
                """,
                (
                    member_id, payload.institution_name, payload.fee_amount, payload.fee_frequency,
                    payload.last_paid_date, payload.next_due_date, payload.notes
                )
            )
            r = cur.fetchone()
            db.commit()
            return SchoolingOut(
                id=r[0], member_id=r[1], institution_name=r[2], fee_amount=float(r[3]),
                fee_frequency=r[4], last_paid_date=r[5], next_due_date=r[6], notes=r[7], created_at=r[8],
                monthly_equivalent=get_monthly_equiv(float(r[3]), r[4])
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{member_id}/schooling/{sid}", status_code=status.HTTP_200_OK)
def delete_schooling(member_id: int, sid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT s.id FROM member_schooling s
            JOIN family_members m ON s.member_id = m.id
            WHERE s.id = %s AND s.member_id = %s AND m.user_id = %s
            """,
            (sid, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Schooling record not found")
        cur.execute("DELETE FROM member_schooling WHERE id = %s", (sid,))
        db.commit()
        return {"message": "Schooling record removed"}


# ─── SCHOOLING PAYMENTS ────────────────────────────────────────────────────

PAYMENT_SELECT = "id, schooling_id, amount, paid_date, receipt_ref, notes, created_at"


@router.get("/{member_id}/schooling/{sid}/payments", response_model=List[SchoolingPaymentOut])
def list_schooling_payments(member_id: int, sid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT s.id FROM member_schooling s JOIN family_members m ON s.member_id = m.id WHERE s.id = %s AND s.member_id = %s AND m.user_id = %s",
            (sid, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Schooling record not found")
        cur.execute(
            f"SELECT {PAYMENT_SELECT} FROM member_schooling_payments WHERE schooling_id = %s ORDER BY paid_date DESC",
            (sid,)
        )
        return [
            SchoolingPaymentOut(
                id=r[0], schooling_id=r[1], amount=float(r[2]), paid_date=r[3],
                receipt_ref=r[4], notes=r[5], created_at=r[6]
            )
            for r in cur.fetchall()
        ]


@router.post("/{member_id}/schooling/{sid}/payments", response_model=SchoolingPaymentOut, status_code=status.HTTP_201_CREATED)
def add_schooling_payment(member_id: int, sid: int, payload: SchoolingPaymentCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT s.id FROM member_schooling s JOIN family_members m ON s.member_id = m.id WHERE s.id = %s AND s.member_id = %s AND m.user_id = %s",
            (sid, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Schooling record not found")
        try:
            cur.execute(
                """INSERT INTO member_schooling_payments (schooling_id, amount, paid_date, receipt_ref, notes)
                VALUES (%s, %s, %s, %s, %s)
                RETURNING id, schooling_id, amount, paid_date, receipt_ref, notes, created_at""",
                (sid, payload.amount, payload.paid_date, payload.receipt_ref, payload.notes)
            )
            r = cur.fetchone()
            db.commit()
            return SchoolingPaymentOut(
                id=r[0], schooling_id=r[1], amount=float(r[2]), paid_date=r[3],
                receipt_ref=r[4], notes=r[5], created_at=r[6]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{member_id}/schooling/{sid}/payments/{pid}", status_code=status.HTTP_200_OK)
def delete_schooling_payment(member_id: int, sid: int, pid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT sp.id FROM member_schooling_payments sp JOIN member_schooling s ON sp.schooling_id = s.id JOIN family_members m ON s.member_id = m.id WHERE sp.id = %s AND sp.schooling_id = %s AND s.member_id = %s AND m.user_id = %s",
            (pid, sid, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Payment record not found")
        cur.execute("DELETE FROM member_schooling_payments WHERE id = %s", (pid,))
        db.commit()
        return {"message": "Payment record deleted"}


@router.post("/{member_id}/checkups", response_model=CheckupOut)
def add_checkup(member_id: int, payload: CheckupCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
        try:
            # Auto-compute next due date if omitted
            next_due = payload.next_due_date
            if not next_due and payload.last_checkup_date:
                delta = 365 if payload.frequency == "yearly" else 90 if payload.frequency == "quarterly" else 30
                next_due = payload.last_checkup_date + timedelta(days=delta)
                
            cur.execute(
                """
                INSERT INTO member_checkups (member_id, checkup_type, frequency, recurring_cost, last_checkup_date, next_due_date, notes)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                RETURNING id, member_id, checkup_type, frequency, recurring_cost, last_checkup_date, next_due_date, notes, created_at
                """,
                (
                    member_id, payload.checkup_type, payload.frequency, payload.recurring_cost or 0.0,
                    payload.last_checkup_date, next_due, payload.notes
                )
            )
            r = cur.fetchone()
            db.commit()
            return CheckupOut(
                id=r[0], member_id=r[1], checkup_type=r[2], frequency=r[3], recurring_cost=float(r[4]),
                last_checkup_date=r[5], next_due_date=r[6], notes=r[7], created_at=r[8]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{member_id}/checkups/{cid}", status_code=status.HTTP_200_OK)
def delete_checkup(member_id: int, cid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT c.id FROM member_checkups c
            JOIN family_members m ON c.member_id = m.id
            WHERE c.id = %s AND c.member_id = %s AND m.user_id = %s
            """,
            (cid, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Checkup record not found")
        cur.execute("DELETE FROM member_checkups WHERE id = %s", (cid,))
        db.commit()
        return {"message": "Checkup removed"}


@router.post("/{member_id}/medicines", response_model=MedicineOut)
def add_medicine(member_id: int, payload: MedicineCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
        try:
            cur.execute(
                """
                INSERT INTO member_medicines (member_id, medicine_name, monthly_cost, purpose, is_regular, prescribed_by)
                VALUES (%s, %s, %s, %s, %s, %s)
                RETURNING id, member_id, medicine_name, monthly_cost, purpose, is_regular, prescribed_by, created_at
                """,
                (
                    member_id, payload.medicine_name, payload.monthly_cost, payload.purpose,
                    payload.is_regular, payload.prescribed_by
                )
            )
            r = cur.fetchone()
            db.commit()
            return MedicineOut(
                id=r[0], member_id=r[1], medicine_name=r[2], monthly_cost=float(r[3]),
                purpose=r[4], is_regular=r[5], prescribed_by=r[6], created_at=r[7]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{member_id}/medicines/{mid}", status_code=status.HTTP_200_OK)
def delete_medicine(member_id: int, mid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT med.id FROM member_medicines med
            JOIN family_members m ON med.member_id = m.id
            WHERE med.id = %s AND med.member_id = %s AND m.user_id = %s
            """,
            (mid, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Medicine record not found")
        cur.execute("DELETE FROM member_medicines WHERE id = %s", (mid,))
        db.commit()
        return {"message": "Medicine removed"}


@router.post("/{member_id}/vaccinations", response_model=VaccinationOut)
def add_vaccination(member_id: int, payload: VaccinationCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
        try:
            next_due = payload.next_due_date
            if not next_due and payload.last_vaccination_date and payload.frequency != "one_time":
                delta = 365 if payload.frequency == "yearly" else 90 if payload.frequency == "quarterly" else 30
                next_due = payload.last_vaccination_date + timedelta(days=delta)
                
            cur.execute(
                """
                INSERT INTO member_vaccinations (member_id, vaccine_name, frequency, recurring_cost, last_vaccination_date, next_due_date, notes)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                RETURNING id, member_id, vaccine_name, frequency, recurring_cost, last_vaccination_date, next_due_date, notes, created_at
                """,
                (
                    member_id, payload.vaccine_name, payload.frequency, payload.recurring_cost or 0.0,
                    payload.last_vaccination_date, next_due, payload.notes
                )
            )
            r = cur.fetchone()
            db.commit()
            return VaccinationOut(
                id=r[0], member_id=r[1], vaccine_name=r[2], frequency=r[3], recurring_cost=float(r[4]),
                last_vaccination_date=r[5], next_due_date=r[6], notes=r[7], created_at=r[8]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{member_id}/vaccinations/{vid}", status_code=status.HTTP_200_OK)
def delete_vaccination(member_id: int, vid: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT v.id FROM member_vaccinations v
            JOIN family_members m ON v.member_id = m.id
            WHERE v.id = %s AND v.member_id = %s AND m.user_id = %s
            """,
            (vid, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Vaccination record not found")
        cur.execute("DELETE FROM member_vaccinations WHERE id = %s", (vid,))
        db.commit()
        return {"message": "Vaccination removed"}


@router.post("/{member_id}/earnings", response_model=EarningOut)
def add_earning(member_id: int, payload: EarningCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
        # Ensure only one earning record per member
        cur.execute("SELECT id FROM member_earnings WHERE member_id = %s", (member_id,))
        if cur.fetchone():
            raise HTTPException(status_code=400, detail="Earning record already exists")
        try:
            cur.execute(
                """
                INSERT INTO member_earnings (member_id, monthly_income, contribution_to_household, occupation)
                VALUES (%s, %s, %s, %s)
                RETURNING id, member_id, monthly_income, contribution_to_household, occupation, created_at
                """,
                (member_id, payload.monthly_income, payload.contribution_to_household or 0.0, payload.occupation)
            )
            r = cur.fetchone()
            db.commit()
            return EarningOut(
                id=r[0], member_id=r[1], monthly_income=float(r[2]),
                contribution_to_household=float(r[3]), occupation=r[4], created_at=r[5]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{member_id}/earnings", status_code=status.HTTP_200_OK)
def delete_earning(member_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
        cur.execute("DELETE FROM member_earnings WHERE member_id = %s", (member_id,))
        db.commit()
        return {"message": "Earnings record deleted"}


# ─── MEMBER INSURANCE CRUD ─────────────────────────────────────────────────────

@router.get("/{member_id}/insurance", response_model=List[MemberInsuranceOut])
def list_member_insurance(member_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
        cur.execute(
            "SELECT id, member_id, provider_name, policy_details, premium_frequency, premium_amount, coverage_amount, created_at FROM member_insurances WHERE member_id = %s",
            (member_id,)
        )
        return [
            MemberInsuranceOut(
                id=r[0], member_id=r[1], provider_name=r[2], policy_details=r[3],
                premium_frequency=r[4], premium_amount=float(r[5]), coverage_amount=float(r[6]), created_at=r[7]
            )
            for r in cur.fetchall()
        ]


@router.post("/{member_id}/insurance", response_model=MemberInsuranceOut, status_code=status.HTTP_201_CREATED)
def add_member_insurance(member_id: int, payload: MemberInsuranceCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id FROM family_members WHERE id = %s AND user_id = %s", (member_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Member not found")
        try:
            cur.execute(
                """INSERT INTO member_insurances (member_id, provider_name, policy_details, premium_frequency, premium_amount, coverage_amount)
                VALUES (%s, %s, %s, %s, %s, %s)
                RETURNING id, member_id, provider_name, policy_details, premium_frequency, premium_amount, coverage_amount, created_at""",
                (member_id, payload.provider_name, payload.policy_details, payload.premium_frequency, payload.premium_amount, payload.coverage_amount)
            )
            r = cur.fetchone()
            db.commit()
            return MemberInsuranceOut(
                id=r[0], member_id=r[1], provider_name=r[2], policy_details=r[3],
                premium_frequency=r[4], premium_amount=float(r[5]), coverage_amount=float(r[6]), created_at=r[7]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.put("/{member_id}/insurance/{insurance_id}", response_model=MemberInsuranceOut)
def update_member_insurance(member_id: int, insurance_id: int, payload: MemberInsuranceCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT mi.id FROM member_insurances mi JOIN family_members fm ON mi.member_id = fm.id WHERE mi.id = %s AND mi.member_id = %s AND fm.user_id = %s",
            (insurance_id, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Insurance record not found")
        try:
            cur.execute(
                """UPDATE member_insurances SET provider_name = %s, policy_details = %s, premium_frequency = %s, premium_amount = %s, coverage_amount = %s
                WHERE id = %s
                RETURNING id, member_id, provider_name, policy_details, premium_frequency, premium_amount, coverage_amount, created_at""",
                (payload.provider_name, payload.policy_details, payload.premium_frequency, payload.premium_amount, payload.coverage_amount, insurance_id)
            )
            r = cur.fetchone()
            db.commit()
            return MemberInsuranceOut(
                id=r[0], member_id=r[1], provider_name=r[2], policy_details=r[3],
                premium_frequency=r[4], premium_amount=float(r[5]), coverage_amount=float(r[6]), created_at=r[7]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{member_id}/insurance/{insurance_id}", status_code=status.HTTP_200_OK)
def delete_member_insurance(member_id: int, insurance_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT mi.id FROM member_insurances mi JOIN family_members fm ON mi.member_id = fm.id WHERE mi.id = %s AND mi.member_id = %s AND fm.user_id = %s",
            (insurance_id, member_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Insurance record not found")
        cur.execute("DELETE FROM member_insurances WHERE id = %s", (insurance_id,))
        db.commit()
        return {"message": "Insurance record deleted"}
