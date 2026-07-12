import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from datetime import datetime, date as date_type
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.budget import (
    BudgetItemCreate, BudgetItemUpdate, BudgetItemOut,
    BudgetPlanOut, BudgetComparisonOut, BudgetSummary, CompareItem
)
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/budget", tags=["budget"])

# Helper to map categories to default buckets
CATEGORY_BUCKET_MAP = {
    "housing": "Needs",
    "food": "Needs",
    "transport": "Needs",
    "health": "Needs",
    "education": "Needs",
    "entertainment": "Wants",
    "savings": "Savings",
    "debt": "Needs",
    "other": "Needs"
}

def get_or_create_plan(cur, user_id: int, month: int, year: int) -> int:
    cur.execute(
        "SELECT id FROM budget_plans WHERE user_id = %s AND month = %s AND year = %s",
        (user_id, month, year)
    )
    row = cur.fetchone()
    if row:
        return row[0]
        
    cur.execute(
        """
        INSERT INTO budget_plans (user_id, month, year, total_budgeted)
        VALUES (%s, %s, %s, 0.00)
        RETURNING id
        """,
        (user_id, month, year)
    )
    return cur.fetchone()[0]

def update_plan_total(cur, plan_id: int):
    cur.execute("SELECT SUM(amount) FROM budget_items WHERE budget_plan_id = %s", (plan_id,))
    total = cur.fetchone()[0] or 0.0
    cur.execute(
        "UPDATE budget_plans SET total_budgeted = %s, updated_at = CURRENT_TIMESTAMP WHERE id = %s",
        (total, plan_id)
    )

def calc_monthly_equivalent(amount: float, frequency: str) -> float:
    freq = frequency.lower()
    if freq == 'monthly':
        return amount
    elif freq == 'quarterly':
        return amount / 3.0
    elif freq == 'yearly':
        return amount / 12.0
    return 0.0

@router.get("/plan", response_model=BudgetPlanOut)
def get_budget_plan(month: Optional[int] = None, year: Optional[int] = None, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    m = month or datetime.now().month
    y = year or datetime.now().year
    
    with db.cursor() as cur:
        plan_id = get_or_create_plan(cur, current_user.id, m, y)
        db.commit()
        
        cur.execute("SELECT id, user_id, month, year, total_budgeted, notes, created_at, updated_at FROM budget_plans WHERE id = %s", (plan_id,))
        p = cur.fetchone()
        
        cur.execute(
            """
            SELECT id, budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed, sort_order, created_at
            FROM budget_items WHERE budget_plan_id = %s ORDER BY sort_order ASC, id ASC
            """,
            (plan_id,)
        )
        items = [
            BudgetItemOut(
                id=r[0], budget_plan_id=r[1], category=r[2], label=r[3], amount=float(r[4]),
                frequency=r[5], bucket=r[6], source=r[7], source_id=r[8], source_label=r[9],
                is_committed=r[10], sort_order=r[11], created_at=r[12]
            )
            for r in cur.fetchall()
        ]
        
        return BudgetPlanOut(
            id=p[0], user_id=p[1], month=p[2], year=p[3], total_budgeted=float(p[4]),
            notes=p[5], created_at=p[6], updated_at=p[7], items=items
        )

@router.post("/plan/auto-populate")
def auto_populate_budget(month: Optional[int] = None, year: Optional[int] = None, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    m = month or datetime.now().month
    y = year or datetime.now().year
    
    with db.cursor() as cur:
        plan_id = get_or_create_plan(cur, current_user.id, m, y)
        
        items_added = 0
        items_skipped = 0
        
        # 1. Fetch active recurring bills (non-EMI)
        cur.execute(
            "SELECT id, name, amount, category, bucket, frequency FROM recurring_bills WHERE user_id = %s AND is_active = TRUE AND is_emi = FALSE",
            (current_user.id,)
        )
        for bill in cur.fetchall():
            b_id, b_name, b_amt, b_cat, b_bucket, b_freq = bill
            # Category fallback check
            cat = b_cat if b_cat in CATEGORY_BUCKET_MAP else 'other'
            # Idempotence check
            cur.execute(
                "SELECT id FROM budget_items WHERE budget_plan_id = %s AND source = 'auto_bill' AND source_id = %s",
                (plan_id, b_id)
            )
            if cur.fetchone():
                items_skipped += 1
                continue
                
            cur.execute(
                """
                INSERT INTO budget_items (budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed)
                VALUES (%s, %s, %s, %s, %s, %s, 'auto_bill', %s, %s, TRUE)
                """,
                (plan_id, cat, b_name, float(b_amt), b_freq, b_bucket, b_id, f"Bill: {b_name}")
            )
            items_added += 1

        # 2. Fetch EMIs
        cur.execute(
            "SELECT id, name, amount, category, bucket, frequency FROM recurring_bills WHERE user_id = %s AND is_active = TRUE AND is_emi = TRUE",
            (current_user.id,)
        )
        for emi in cur.fetchall():
            e_id, e_name, e_amt, e_cat, e_bucket, e_freq = emi
            cat = e_cat if e_cat in CATEGORY_BUCKET_MAP else 'debt'
            cur.execute(
                "SELECT id FROM budget_items WHERE budget_plan_id = %s AND source = 'auto_emi' AND source_id = %s",
                (plan_id, e_id)
            )
            if cur.fetchone():
                items_skipped += 1
                continue
                
            cur.execute(
                """
                INSERT INTO budget_items (budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed)
                VALUES (%s, %s, %s, %s, %s, %s, 'auto_emi', %s, %s, TRUE)
                """,
                (plan_id, cat, e_name, float(e_amt), e_freq, e_bucket, e_id, f"EMI: {e_name}")
            )
            items_added += 1

        # 3. Fetch family member recurring costs (Phase 5)
        if current_user.household_id:
            cur.execute(
                "SELECT id, name FROM family_members WHERE household_id = %s AND is_active = TRUE",
                (current_user.household_id,)
            )
            members = cur.fetchall()
            for member in members:
                m_id, m_name = member
                
                # Schooling
                cur.execute("SELECT id, institution_name, fee_amount, fee_frequency FROM member_schooling WHERE member_id = %s", (m_id,))
                for s in cur.fetchall():
                    s_id, inst, amt, freq = s
                    cur.execute(
                        "SELECT id FROM budget_items WHERE budget_plan_id = %s AND source = 'auto_profile' AND source_label = %s",
                        (plan_id, f"Schooling: {inst} ({m_name})")
                    )
                    if cur.fetchone():
                        items_skipped += 1
                        continue
                    monthly_val = calc_monthly_equivalent(float(amt), freq)
                    cur.execute(
                        """
                        INSERT INTO budget_items (budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed)
                        VALUES (%s, 'education', %s, %s, 'monthly', 'Needs', 'auto_profile', %s, %s, TRUE)
                        """,
                        (plan_id, f"{m_name} school fee", monthly_val, s_id, f"Schooling: {inst} ({m_name})")
                    )
                    items_added += 1
                    
                # Medicines
                cur.execute("SELECT id, medicine_name, monthly_cost FROM member_medicines WHERE member_id = %s", (m_id,))
                for med in cur.fetchall():
                    med_id, med_name, amt = med
                    cur.execute(
                        "SELECT id FROM budget_items WHERE budget_plan_id = %s AND source = 'auto_profile' AND source_label = %s",
                        (plan_id, f"Medicine: {med_name} ({m_name})")
                    )
                    if cur.fetchone():
                        items_skipped += 1
                        continue
                    cur.execute(
                        """
                        INSERT INTO budget_items (budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed)
                        VALUES (%s, 'health', %s, %s, 'monthly', 'Needs', 'auto_profile', %s, %s, TRUE)
                        """,
                        (plan_id, f"{m_name} medicines", float(amt), med_id, f"Medicine: {med_name} ({m_name})")
                    )
                    items_added += 1

                # Checkups
                cur.execute("SELECT id, checkup_type, recurring_cost, frequency FROM member_checkups WHERE member_id = %s", (m_id,))
                for check in cur.fetchall():
                    c_id, c_type, amt, freq = check
                    cur.execute(
                        "SELECT id FROM budget_items WHERE budget_plan_id = %s AND source = 'auto_profile' AND source_label = %s",
                        (plan_id, f"Checkup: {c_type} ({m_name})")
                    )
                    if cur.fetchone():
                        items_skipped += 1
                        continue
                    monthly_val = calc_monthly_equivalent(float(amt), freq)
                    cur.execute(
                        """
                        INSERT INTO budget_items (budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed)
                        VALUES (%s, 'health', %s, %s, 'monthly', 'Needs', 'auto_profile', %s, %s, TRUE)
                        """,
                        (plan_id, f"{m_name} checkup", monthly_val, c_id, f"Checkup: {c_type} ({m_name})")
                    )
                    items_added += 1

                # Vaccinations
                cur.execute("SELECT id, vaccine_name, recurring_cost, frequency FROM member_vaccinations WHERE member_id = %s AND frequency != 'one_time'", (m_id,))
                for vac in cur.fetchall():
                    v_id, v_name, amt, freq = vac
                    cur.execute(
                        "SELECT id FROM budget_items WHERE budget_plan_id = %s AND source = 'auto_profile' AND source_label = %s",
                        (plan_id, f"Vaccination: {v_name} ({m_name})")
                    )
                    if cur.fetchone():
                        items_skipped += 1
                        continue
                    monthly_val = calc_monthly_equivalent(float(amt), freq)
                    cur.execute(
                        """
                        INSERT INTO budget_items (budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed)
                        VALUES (%s, 'health', %s, %s, 'monthly', 'Needs', 'auto_profile', %s, %s, TRUE)
                        """,
                        (plan_id, f"{m_name} vaccine", monthly_val, v_id, f"Vaccination: {v_name} ({m_name})")
                    )
                    items_added += 1

        update_plan_total(cur, plan_id)
        db.commit()
        
        cur.execute("SELECT total_budgeted FROM budget_plans WHERE id = %s", (plan_id,))
        total = cur.fetchone()[0] or 0.0
        
        return {
            "items_added": items_added,
            "items_skipped": items_skipped,
            "total_budgeted": float(total)
        }

@router.post("/items", response_model=BudgetItemOut, status_code=status.HTTP_201_CREATED)
def add_budget_item(item: BudgetItemCreate, month: Optional[int] = None, year: Optional[int] = None, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    m = month or datetime.now().month
    y = year or datetime.now().year
    
    with db.cursor() as cur:
        plan_id = get_or_create_plan(cur, current_user.id, m, y)
        
        try:
            cur.execute(
                """
                INSERT INTO budget_items (budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed, sort_order)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed, sort_order, created_at
                """,
                (plan_id, item.category, item.label, item.amount, item.frequency, item.bucket,
                 item.source, item.source_id, item.source_label, item.is_committed, item.sort_order)
            )
            row = cur.fetchone()
            update_plan_total(cur, plan_id)
            db.commit()
            return BudgetItemOut(
                id=row[0], budget_plan_id=row[1], category=row[2], label=row[3], amount=float(row[4]),
                frequency=row[5], bucket=row[6], source=row[7], source_id=row[8], source_label=row[9],
                is_committed=row[10], sort_order=row[11], created_at=row[12]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))

@router.put("/items/{id}", response_model=BudgetItemOut)
def update_budget_item(id: int, item: BudgetItemUpdate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        # Check ownership of budget item's parent plan
        cur.execute(
            """
            SELECT bi.budget_plan_id FROM budget_items bi
            JOIN budget_plans bp ON bi.budget_plan_id = bp.id
            WHERE bi.id = %s AND bp.user_id = %s
            """,
            (id, current_user.id)
        )
        row = cur.fetchone()
        if not row:
            raise HTTPException(status_code=404, detail="Budget item not found")
        plan_id = row[0]
        
        updates = []
        params = []
        if item.category is not None:
            updates.append("category = %s")
            params.append(item.category)
        if item.label is not None:
            updates.append("label = %s")
            params.append(item.label)
        if item.amount is not None:
            updates.append("amount = %s")
            params.append(item.amount)
        if item.frequency is not None:
            updates.append("frequency = %s")
            params.append(item.frequency)
        if item.bucket is not None:
            updates.append("bucket = %s")
            params.append(item.bucket)
        if item.is_committed is not None:
            updates.append("is_committed = %s")
            params.append(item.is_committed)
        if item.sort_order is not None:
            updates.append("sort_order = %s")
            params.append(item.sort_order)
            
        if not updates:
            raise HTTPException(status_code=400, detail="No fields to update")
            
        params.append(id)
        query = f"UPDATE budget_items SET {', '.join(updates)} WHERE id = %s RETURNING id, budget_plan_id, category, label, amount, frequency, bucket, source, source_id, source_label, is_committed, sort_order, created_at"
        
        try:
            cur.execute(query, tuple(params))
            row = cur.fetchone()
            update_plan_total(cur, plan_id)
            db.commit()
            return BudgetItemOut(
                id=row[0], budget_plan_id=row[1], category=row[2], label=row[3], amount=float(row[4]),
                frequency=row[5], bucket=row[6], source=row[7], source_id=row[8], source_label=row[9],
                is_committed=row[10], sort_order=row[11], created_at=row[12]
            )
        except Exception as e:
            db.rollback()
            raise HTTPException(status_code=500, detail=str(e))

@router.delete("/items/{id}", status_code=status.HTTP_200_OK)
def delete_budget_item(id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT bi.budget_plan_id FROM budget_items bi
            JOIN budget_plans bp ON bi.budget_plan_id = bp.id
            WHERE bi.id = %s AND bp.user_id = %s
            """,
            (id, current_user.id)
        )
        row = cur.fetchone()
        if not row:
            raise HTTPException(status_code=404, detail="Budget item not found")
        plan_id = row[0]
        
        cur.execute("DELETE FROM budget_items WHERE id = %s", (id,))
        update_plan_total(cur, plan_id)
        db.commit()
        return {"message": "Budget item deleted successfully"}

@router.get("/compare", response_model=BudgetComparisonOut)
def compare_budget_vs_actual(month: Optional[int] = None, year: Optional[int] = None, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    m = month or datetime.now().month
    y = year or datetime.now().year
    
    with db.cursor() as cur:
        plan_id = get_or_create_plan(cur, current_user.id, m, y)
        
        # 1. Fetch budget plan items
        cur.execute(
            "SELECT category, label, amount, bucket FROM budget_items WHERE budget_plan_id = %s",
            (plan_id,)
        )
        budget_rows = cur.fetchall()
        
        budget_total = 0.0
        budget_needs = 0.0
        budget_wants = 0.0
        budget_savings = 0.0
        
        # We index budget targets by category and case-insensitive label
        budget_by_category = {}
        for r in budget_rows:
            cat, label, amt, bucket = r
            amt = float(amt)
            budget_total += amt
            if bucket == 'Needs':
                budget_needs += amt
            elif bucket == 'Wants':
                budget_wants += amt
            elif bucket == 'Savings':
                budget_savings += amt
                
            budget_by_category[cat.lower()] = budget_by_category.get(cat.lower(), 0.0) + amt
            budget_by_category[label.lower()] = budget_by_category.get(label.lower(), 0.0) + amt
            
        # 2. Fetch actual expenses for the month
        cur.execute(
            """
            SELECT category, bucket, SUM(amount) FROM expenses
            WHERE user_id = %s AND EXTRACT(MONTH FROM date) = %s AND EXTRACT(YEAR FROM date) = %s
            GROUP BY category, bucket
            """,
            (current_user.id, m, y)
        )
        expense_rows = cur.fetchall()
        
        act_total = 0.0
        act_needs = 0.0
        act_wants = 0.0
        act_savings = 0.0
        
        actual_by_category = {}
        for r in expense_rows:
            cat, bucket, amt = r
            amt = float(amt)
            act_total += amt
            if bucket == 'Needs':
                act_needs += amt
            elif bucket == 'Wants':
                act_wants += amt
            elif bucket == 'Savings':
                act_savings += amt
                
            actual_by_category[cat.lower()] = actual_by_category.get(cat.lower(), 0.0) + amt
            
        # Compile comparison items
        compare_items = []
        for r in budget_rows:
            cat, label, amt, bucket = r
            amt = float(amt)
            
            # Find actual spending: try to match by exact label first (e.g. Groceries), then by category
            actual_amt = actual_by_category.get(label.lower(), 0.0)
            if actual_amt == 0.0:
                actual_amt = actual_by_category.get(cat.lower(), 0.0)
                # clear category so we don't double count if multiple items map to same category
                # (approximate heuristic for visual representation)
                actual_by_category[cat.lower()] = 0.0
                
            compare_items.append(
                CompareItem(
                    label=label,
                    budgeted=amt,
                    actual=actual_amt,
                    remaining=max(0.0, amt - actual_amt),
                    bucket=bucket
                )
            )
            
        return BudgetComparisonOut(
            budgeted=BudgetSummary(total=budget_total, needs=budget_needs, wants=budget_wants, savings=budget_savings),
            actual=BudgetSummary(total=act_total, needs=act_needs, wants=act_wants, savings=act_savings),
            remaining=BudgetSummary(
                total=max(0.0, budget_total - act_total),
                needs=max(0.0, budget_needs - act_needs),
                wants=max(0.0, budget_wants - act_wants),
                savings=max(0.0, budget_savings - act_savings)
            ),
            items=compare_items
        )
