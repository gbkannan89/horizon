from fastapi import APIRouter, Depends, HTTPException, status
from typing import List
from datetime import datetime, timedelta
import logging

from ..schemas import (
    WishlistItemCreate, WishlistItemOut, WishlistItemUpdate,
    FinancialGuardrailsOut, DebtRepaymentStrategyOut,
    ZeroBasedBudgetOut, PayYourselfFirstOut, UserOut
)
from ..database import get_db
from .auth import get_current_user

router = APIRouter(
    prefix="/api/discipline",
    tags=["Discipline"]
)

@router.get("/wishlist", response_model=List[WishlistItemOut])
def get_wishlist(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    try:
        with db.cursor() as cursor:
            cursor.execute(
                "SELECT id, user_id, name, amount, added_date, unlock_date, status, created_at FROM wishlist_items WHERE user_id = %s ORDER BY added_date DESC",
                (current_user.id,)
            )
            items = cursor.fetchall()
            return [
                WishlistItemOut(
                    id=row[0], user_id=row[1], name=row[2], amount=float(row[3]),
                    added_date=row[4], unlock_date=row[5], status=row[6], created_at=row[7]
                ) for row in items
            ]
    except Exception as e:
        logging.error(f"Error fetching wishlist: {e}")
        raise HTTPException(status_code=500, detail="Internal server error")

@router.post("/wishlist", response_model=WishlistItemOut, status_code=status.HTTP_201_CREATED)
def add_wishlist_item(item: WishlistItemCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    try:
        now = datetime.now()
        unlock_date = now + timedelta(days=item.lock_duration_days)
        with db.cursor() as cursor:
            cursor.execute(
                """
                INSERT INTO wishlist_items (user_id, name, amount, unlock_date, status, lock_duration_days)
                VALUES (%s, %s, %s, %s, 'locked', %s)
                RETURNING id, user_id, name, amount, added_date, unlock_date, status, created_at
                """,
                (current_user.id, item.name, item.amount, unlock_date, item.lock_duration_days)
            )
            row = cursor.fetchone()
            db.commit()
            return WishlistItemOut(
                id=row[0], user_id=row[1], name=row[2], amount=float(row[3]),
                added_date=row[4], unlock_date=row[5], status=row[6], created_at=row[7]
            )
    except Exception as e:
        db.rollback()
        logging.error(f"Error adding wishlist item: {e}")
        raise HTTPException(status_code=500, detail="Internal server error")

@router.put("/wishlist/{item_id}", response_model=WishlistItemOut)
def update_wishlist_item(item_id: int, update: WishlistItemUpdate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    try:
        with db.cursor() as cursor:
            cursor.execute("SELECT unlock_date FROM wishlist_items WHERE id = %s AND user_id = %s", (item_id, current_user.id))
            row = cursor.fetchone()
            if not row:
                raise HTTPException(status_code=404, detail="Item not found")
            
            unlock_date = row[0]
            now = datetime.now(unlock_date.tzinfo) if unlock_date.tzinfo else datetime.now()
            status_to_set = update.status
            
            if update.status == 'bought' and now < unlock_date:
                status_to_set = 'bought_early'
                
            cursor.execute(
                """
                UPDATE wishlist_items SET status = %s WHERE id = %s
                RETURNING id, user_id, name, amount, added_date, unlock_date, status, created_at
                """,
                (status_to_set, item_id)
            )
            row = cursor.fetchone()
            db.commit()
            return WishlistItemOut(
                id=row[0], user_id=row[1], name=row[2], amount=float(row[3]),
                added_date=row[4], unlock_date=row[5], status=row[6], created_at=row[7]
            )
    except HTTPException:
        raise
    except Exception as e:
        db.rollback()
        logging.error(f"Error updating wishlist item: {e}")
        raise HTTPException(status_code=500, detail="Internal server error")

@router.get("/guardrails", response_model=FinancialGuardrailsOut)
def get_guardrails(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    try:
        with db.cursor() as cursor:
            cursor.execute("SELECT SUM(amount) FROM assets WHERE user_id = %s AND type IN ('bank', 'fd')", (current_user.id,))
            liquid_assets = cursor.fetchone()[0] or 0.0
            
            cursor.execute("SELECT SUM(monthly_equivalent) FROM recurring_bills WHERE user_id = %s AND bucket = 'Needs'", (current_user.id,))
            monthly_needs = cursor.fetchone()[0] or 1.0 # prevent div by 0
            
            cursor.execute("SELECT SUM(amount) FROM incomes WHERE user_id = %s", (current_user.id,))
            total_income = cursor.fetchone()[0] or 1.0
            
            cursor.execute("SELECT SUM(monthly_equivalent) FROM recurring_bills WHERE user_id = %s AND category IN ('Rent', 'Mortgage', 'Housing')", (current_user.id,))
            housing_costs = cursor.fetchone()[0] or 0.0
            
            ef_ratio = float(liquid_assets) / float(monthly_needs)
            housing_ratio = (float(housing_costs) / float(total_income)) * 100
            
            runway_status = 'Healthy' if ef_ratio >= 3.0 else 'Warning' if ef_ratio >= 1.0 else 'Danger'
            housing_status = 'Healthy' if housing_ratio <= 30.0 else 'Warning' if housing_ratio <= 40.0 else 'Danger'
            
            return FinancialGuardrailsOut(
                emergency_fund_ratio=ef_ratio,
                emergency_target_amount=float(monthly_needs) * 6.0,
                emergency_current_amount=float(liquid_assets),
                housing_cost_ratio=housing_ratio,
                housing_status=housing_status,
                runway_status=runway_status
            )
    except Exception as e:
        logging.error(f"Error guardrails: {e}")
        raise HTTPException(status_code=500, detail="Internal server error")

@router.get("/debt-strategy", response_model=DebtRepaymentStrategyOut)
def get_debt_strategy(current_user: UserOut = Depends(get_current_user)):
    return DebtRepaymentStrategyOut(
        snowball_months_to_freedom=24,
        avalanche_months_to_freedom=22,
        snowball_total_interest=1200.0,
        avalanche_total_interest=950.0,
        recommended_strategy='Avalanche'
    )

@router.get("/zero-based-budget", response_model=ZeroBasedBudgetOut)
def get_zero_based_budget(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    try:
        with db.cursor() as cursor:
            cursor.execute("SELECT SUM(amount) FROM incomes WHERE user_id = %s AND frequency = 'monthly'", (current_user.id,))
            total_income = cursor.fetchone()[0] or 0.0
            
            cursor.execute("SELECT SUM(monthly_equivalent) FROM recurring_bills WHERE user_id = %s", (current_user.id,))
            total_allocated = cursor.fetchone()[0] or 0.0
            
            unallocated = float(total_income) - float(total_allocated)
            status = 'Zero-Based' if abs(unallocated) < 10.0 else 'Surplus' if unallocated > 0 else 'Deficit'
            
            return ZeroBasedBudgetOut(
                total_income=float(total_income),
                total_allocated=float(total_allocated),
                unallocated=unallocated,
                status=status
            )
    except Exception as e:
        raise HTTPException(status_code=500, detail="Internal server error")

@router.get("/pay-yourself-first", response_model=PayYourselfFirstOut)
def get_pay_yourself_first(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    try:
        with db.cursor() as cursor:
            cursor.execute("SELECT SUM(amount) FROM incomes WHERE user_id = %s AND frequency = 'monthly'", (current_user.id,))
            total_income = float(cursor.fetchone()[0] or 0.0)
            
            target_percentage = 20.0
            target_amount = total_income * (target_percentage / 100.0)
            
            cursor.execute("SELECT SUM(monthly_equivalent) FROM recurring_bills WHERE user_id = %s AND bucket = 'Savings'", (current_user.id,))
            actual_savings = float(cursor.fetchone()[0] or 0.0)
            
            status = 'On Track' if actual_savings >= target_amount else 'Behind'
            
            return PayYourselfFirstOut(
                target_percentage=target_percentage,
                target_amount=target_amount,
                actual_savings=actual_savings,
                status=status
            )
    except Exception as e:
        raise HTTPException(status_code=500, detail="Internal server error")
