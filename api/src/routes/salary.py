import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.salary import (
    SalaryDetailCreate, SalaryDetailUpdate, SalaryDetailOut,
    SalaryGrowthOut, SalaryGrowthItem
)
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/incomes", tags=["salary"])


def auto_compute_salary(data: dict) -> dict:
    fixed_pay = data.get("fixed_pay") or 0.0
    basic_pay = data.get("basic_pay")
    if fixed_pay > 0.0 and (basic_pay is None or basic_pay == 0.0):
        data["basic_pay"] = fixed_pay * 0.5
        
    var_pct = data.get("variable_pay_percentage") or 0.0
    var_amt = data.get("variable_pay_amount")
    if fixed_pay > 0.0 and var_pct > 0.0 and (var_amt is None or var_amt == 0.0):
        data["variable_pay_amount"] = fixed_pay * (var_pct / 100.0)
        
    final_var_amt = data.get("variable_pay_amount") or 0.0
    data["gross_annual"] = fixed_pay + final_var_amt
    
    pf_emp = data.get("pf_employee") or 0.0
    meal = data.get("meal_card") or 0.0
    monthly_gross = fixed_pay / 12.0
    data["monthly_in_hand"] = max(0.0, monthly_gross - pf_emp - meal)
    return data


@router.get("/{income_id}/salary", response_model=List[SalaryDetailOut])
def list_salary_details(income_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Verify ownership of the income source
        cur.execute("SELECT id FROM incomes WHERE id = %s AND user_id = %s", (income_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Income source not found")
            
        cur.execute(
            """
            SELECT id, income_id, company_name, from_year, to_year, is_current, fixed_pay, basic_pay, hra, lta,
                   pf_employee, pf_employer, special_allowance, meal_card, variable_pay_percentage,
                   variable_pay_amount, gross_annual, monthly_in_hand, created_at
            FROM salary_details WHERE income_id = %s ORDER BY from_year DESC
            """,
            (income_id,)
        )
        rows = cur.fetchall()
        return [
            SalaryDetailOut(
                id=r[0], income_id=r[1], company_name=r[2], from_year=r[3], to_year=r[4], is_current=r[5],
                fixed_pay=float(r[6]) if r[6] else None, basic_pay=float(r[7]) if r[7] else None,
                hra=float(r[8]) if r[8] else None, lta=float(r[9]) if r[9] else None,
                pf_employee=float(r[10]) if r[10] else None, pf_employer=float(r[11]) if r[11] else None,
                special_allowance=float(r[12]) if r[12] else None, meal_card=float(r[13]) if r[13] else None,
                variable_pay_percentage=float(r[14]) if r[14] else None, variable_pay_amount=float(r[15]) if r[15] else None,
                gross_annual=float(r[16]) if r[16] else None, monthly_in_hand=float(r[17]) if r[17] else None,
                created_at=r[18]
            )
            for r in rows
        ]


@router.post("/{income_id}/salary", response_model=SalaryDetailOut, status_code=status.HTTP_201_CREATED)
def create_salary_detail(income_id: int, payload: SalaryDetailCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Verify ownership of the income source
        cur.execute("SELECT id FROM incomes WHERE id = %s AND user_id = %s", (income_id, current_user.id))
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Income source not found")
            
        data = payload.model_dump()
        data = auto_compute_salary(data)
        
        try:
            if data["is_current"]:
                # Archive previous current salary details
                cur.execute(
                    "UPDATE salary_details SET is_current = FALSE, to_year = %s WHERE income_id = %s AND is_current = TRUE",
                    (data["from_year"], income_id)
                )
                
            cur.execute(
                """
                INSERT INTO salary_details (
                    income_id, company_name, from_year, to_year, is_current, fixed_pay, basic_pay, hra, lta,
                    pf_employee, pf_employer, special_allowance, meal_card, variable_pay_percentage,
                    variable_pay_amount, gross_annual, monthly_in_hand
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, income_id, company_name, from_year, to_year, is_current, fixed_pay, basic_pay, hra, lta,
                          pf_employee, pf_employer, special_allowance, meal_card, variable_pay_percentage,
                          variable_pay_amount, gross_annual, monthly_in_hand, created_at
                """,
                (
                    income_id, data["company_name"], data["from_year"], data["to_year"], data["is_current"],
                    data["fixed_pay"], data["basic_pay"], data["hra"], data["lta"], data["pf_employee"],
                    data["pf_employer"], data["special_allowance"], data["meal_card"], data["variable_pay_percentage"],
                    data["variable_pay_amount"], data["gross_annual"], data["monthly_in_hand"]
                )
            )
            r = cur.fetchone()
            
            # If this is the current salary, sync company_name & amount back to the primary income source
            if data["is_current"] and data["monthly_in_hand"]:
                cur.execute(
                    "UPDATE incomes SET company_name = %s, amount = %s WHERE id = %s",
                    (data["company_name"], data["monthly_in_hand"], income_id)
                )
                
            db.commit()
            
            return SalaryDetailOut(
                id=r[0], income_id=r[1], company_name=r[2], from_year=r[3], to_year=r[4], is_current=r[5],
                fixed_pay=float(r[6]) if r[6] else None, basic_pay=float(r[7]) if r[7] else None,
                hra=float(r[8]) if r[8] else None, lta=float(r[9]) if r[9] else None,
                pf_employee=float(r[10]) if r[10] else None, pf_employer=float(r[11]) if r[11] else None,
                special_allowance=float(r[12]) if r[12] else None, meal_card=float(r[13]) if r[13] else None,
                variable_pay_percentage=float(r[14]) if r[14] else None, variable_pay_amount=float(r[15]) if r[15] else None,
                gross_annual=float(r[16]) if r[16] else None, monthly_in_hand=float(r[17]) if r[17] else None,
                created_at=r[18]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.put("/{income_id}/salary/{id}", response_model=SalaryDetailOut)
def update_salary_detail(income_id: int, id: int, payload: SalaryDetailUpdate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Check ownership of salary_details via income source
        cur.execute(
            """
            SELECT sd.id, sd.fixed_pay, sd.basic_pay, sd.hra, sd.lta, sd.pf_employee, sd.pf_employer,
                   sd.special_allowance, sd.meal_card, sd.variable_pay_percentage, sd.variable_pay_amount
            FROM salary_details sd
            JOIN incomes i ON sd.income_id = i.id
            WHERE sd.id = %s AND sd.income_id = %s AND i.user_id = %s
            """,
            (id, income_id, current_user.id)
        )
        row = cur.fetchone()
        if not row:
            raise HTTPException(status_code=404, detail="Salary detail entry not found")
            
        update_data = payload.model_dump(exclude_unset=True)
        
        # Merge old and new values to recompute auto fields
        merged = {
            "fixed_pay": update_data.get("fixed_pay", float(row[1]) if row[1] else None),
            "basic_pay": update_data.get("basic_pay", float(row[2]) if row[2] else None),
            "hra": update_data.get("hra", float(row[3]) if row[3] else None),
            "lta": update_data.get("lta", float(row[4]) if row[4] else None),
            "pf_employee": update_data.get("pf_employee", float(row[5]) if row[5] else None),
            "pf_employer": update_data.get("pf_employer", float(row[6]) if row[6] else None),
            "special_allowance": update_data.get("special_allowance", float(row[7]) if row[7] else None),
            "meal_card": update_data.get("meal_card", float(row[8]) if row[8] else None),
            "variable_pay_percentage": update_data.get("variable_pay_percentage", float(row[9]) if row[9] else None),
            "variable_pay_amount": update_data.get("variable_pay_amount", float(row[10]) if row[10] else None),
        }
        merged = auto_compute_salary(merged)
        
        # Add computed values to update set
        update_data["basic_pay"] = merged["basic_pay"]
        update_data["variable_pay_amount"] = merged["variable_pay_amount"]
        update_data["gross_annual"] = merged["gross_annual"]
        update_data["monthly_in_hand"] = merged["monthly_in_hand"]
        
        updates = []
        params = []
        for k, v in update_data.items():
            updates.append(f"{k} = %s")
            params.append(v)
            
        params.append(id)
        query = f"UPDATE salary_details SET {', '.join(updates)} WHERE id = %s RETURNING id, income_id, company_name, from_year, to_year, is_current, fixed_pay, basic_pay, hra, lta, pf_employee, pf_employer, special_allowance, meal_card, variable_pay_percentage, variable_pay_amount, gross_annual, monthly_in_hand, created_at"
        
        try:
            cur.execute(query, tuple(params))
            r = cur.fetchone()
            
            # Sync to incomes if this updated entry is current
            if r[5] and r[17]:
                cur.execute(
                    "UPDATE incomes SET company_name = %s, amount = %s WHERE id = %s",
                    (r[2], r[17], income_id)
                )
                
            db.commit()
            return SalaryDetailOut(
                id=r[0], income_id=r[1], company_name=r[2], from_year=r[3], to_year=r[4], is_current=r[5],
                fixed_pay=float(r[6]) if r[6] else None, basic_pay=float(r[7]) if r[7] else None,
                hra=float(r[8]) if r[8] else None, lta=float(r[9]) if r[9] else None,
                pf_employee=float(r[10]) if r[10] else None, pf_employer=float(r[11]) if r[11] else None,
                special_allowance=float(r[12]) if r[12] else None, meal_card=float(r[13]) if r[13] else None,
                variable_pay_percentage=float(r[14]) if r[14] else None, variable_pay_amount=float(r[15]) if r[15] else None,
                gross_annual=float(r[16]) if r[16] else None, monthly_in_hand=float(r[17]) if r[17] else None,
                created_at=r[18]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{income_id}/salary/{id}", status_code=status.HTTP_200_OK)
def delete_salary_detail(income_id: int, id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT sd.id FROM salary_details sd
            JOIN incomes i ON sd.income_id = i.id
            WHERE sd.id = %s AND sd.income_id = %s AND i.user_id = %s
            """,
            (id, income_id, current_user.id)
        )
        if not cur.fetchone():
            raise HTTPException(status_code=404, detail="Salary detail entry not found")
            
        cur.execute("DELETE FROM salary_details WHERE id = %s", (id,))
        db.commit()
        return {"message": "Salary detail deleted successfully"}


@router.post("/{income_id}/salary/{id}/current", response_model=SalaryDetailOut)
def set_current_salary(income_id: int, id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Check ownership
        cur.execute(
            """
            SELECT sd.id, sd.company_name, sd.monthly_in_hand, sd.from_year FROM salary_details sd
            JOIN incomes i ON sd.income_id = i.id
            WHERE sd.id = %s AND sd.income_id = %s AND i.user_id = %s
            """,
            (id, income_id, current_user.id)
        )
        row = cur.fetchone()
        if not row:
            raise HTTPException(status_code=404, detail="Salary detail entry not found")
            
        comp_name, monthly_pay, start_yr = row[1], float(row[2]) if row[2] else 0.0, row[3]
        
        try:
            # Set all other salaries to not current
            cur.execute(
                "UPDATE salary_details SET is_current = FALSE, to_year = %s WHERE income_id = %s AND id != %s AND is_current = TRUE",
                (start_yr, income_id, id)
            )
            
            # Set target salary to current
            cur.execute(
                "UPDATE salary_details SET is_current = TRUE, to_year = NULL WHERE id = %s RETURNING id, income_id, company_name, from_year, to_year, is_current, fixed_pay, basic_pay, hra, lta, pf_employee, pf_employer, special_allowance, meal_card, variable_pay_percentage, variable_pay_amount, gross_annual, monthly_in_hand, created_at",
                (id,)
            )
            r = cur.fetchone()
            
            # Update incomes table
            cur.execute(
                "UPDATE incomes SET company_name = %s, amount = %s WHERE id = %s",
                (comp_name, monthly_pay, income_id)
            )
            
            db.commit()
            return SalaryDetailOut(
                id=r[0], income_id=r[1], company_name=r[2], from_year=r[3], to_year=r[4], is_current=r[5],
                fixed_pay=float(r[6]) if r[6] else None, basic_pay=float(r[7]) if r[7] else None,
                hra=float(r[8]) if r[8] else None, lta=float(r[9]) if r[9] else None,
                pf_employee=float(r[10]) if r[10] else None, pf_employer=float(r[11]) if r[11] else None,
                special_allowance=float(r[12]) if r[12] else None, meal_card=float(r[13]) if r[13] else None,
                variable_pay_percentage=float(r[14]) if r[14] else None, variable_pay_amount=float(r[15]) if r[15] else None,
                gross_annual=float(r[16]) if r[16] else None, monthly_in_hand=float(r[17]) if r[17] else None,
                created_at=r[18]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))


@router.get("/salary-growth", response_model=SalaryGrowthOut)
def get_salary_growth(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Get all salary detail records belonging to this user
        cur.execute(
            """
            SELECT sd.from_year, sd.gross_annual, sd.company_name
            FROM salary_details sd
            JOIN incomes i ON sd.income_id = i.id
            WHERE i.user_id = %s
            ORDER BY sd.from_year ASC
            """,
            (current_user.id,)
        )
        rows = cur.fetchall()
        
        growth_items = []
        prev_gross = 0.0
        
        for idx, r in enumerate(rows):
            yr, gross, comp = r[0], float(r[1]) if r[1] else 0.0, r[2]
            
            growth_pct = None
            if idx > 0 and prev_gross > 0.0:
                growth_pct = ((gross - prev_gross) / prev_gross) * 100.0
                
            growth_items.append(
                SalaryGrowthItem(year=yr, gross_annual=gross, company=comp, growth_pct=growth_pct)
            )
            prev_gross = gross
            
        # Calculation totals
        time_period = max(0, len(growth_items) - 1)
        
        valid_growth_rates = [item.growth_pct for item in growth_items if item.growth_pct is not None]
        avg_growth = sum(valid_growth_rates) / len(valid_growth_rates) if valid_growth_rates else 0.0
        
        total_growth = 0.0
        if len(growth_items) >= 2:
            first_gross = growth_items[0].gross_annual
            last_gross = growth_items[-1].gross_annual
            if first_gross > 0.0:
                total_growth = ((last_gross - first_gross) / first_gross) * 100.0
                
        return SalaryGrowthOut(
            growth_data=growth_items,
            avg_growth_pct=avg_growth,
            total_growth_pct=total_growth,
            time_period_years=time_period
        )
