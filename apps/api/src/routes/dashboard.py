from fastapi import APIRouter, Depends, HTTPException, status
from typing import List
import datetime
from ..database import get_db
from ..schemas import DashboardOverview, GoalOut, ExpenseOut, RecurringBillOut, ExpenseCreate, UserOut
from .auth import get_current_user

router = APIRouter(prefix="/api/v1/dashboard", tags=["Dashboard"])

@router.get("/overview", response_model=DashboardOverview)
def get_dashboard_overview(
    month: int = None,
    year: int = None,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            user_id = current_user.id  # FIX: was current_user['id'] — UserOut is a Pydantic model not a dict

            # Fetch Incomes
            cur.execute("SELECT amount, frequency FROM incomes WHERE user_id = %s", (user_id,))
            incomes = cur.fetchall()
            total_income = sum(float(i[0]) if i[1] == 'monthly' else float(i[0]) / 12 for i in incomes)

            # Fetch Household Income (Contributing Members)
            cur.execute("SELECT household_id FROM users WHERE id = %s", (user_id,))
            hh_row = cur.fetchone()
            if hh_row and hh_row[0]:
                cur.execute("SELECT contribution_to_household FROM contributing_members WHERE household_id = %s", (hh_row[0],))
                contrib_incomes = cur.fetchall()
                total_income += sum(float(c[0]) for c in contrib_incomes)

            # Fetch Assets
            cur.execute("SELECT amount FROM assets WHERE user_id = %s", (user_id,))
            total_assets = sum(float(a[0]) for a in cur.fetchall())

            # Fetch Liabilities
            cur.execute("SELECT outstanding FROM liabilities WHERE user_id = %s", (user_id,))
            total_liabilities = sum(float(l[0]) for l in cur.fetchall())

            net_worth = total_assets - total_liabilities

            # Save snapshot for today
            cur.execute("""
                INSERT INTO net_worth_snapshots (user_id, snapshot_date, net_worth) 
                VALUES (%s, CURRENT_DATE, %s) 
                ON CONFLICT (user_id, snapshot_date) 
                DO UPDATE SET net_worth = EXCLUDED.net_worth
            """, (user_id, net_worth))
            
            # Calculate net worth change
            cur.execute("""
                SELECT net_worth, snapshot_date FROM net_worth_snapshots 
                WHERE user_id = %s AND snapshot_date < CURRENT_DATE 
                ORDER BY snapshot_date ASC LIMIT 1
            """, (user_id,))
            past_nw_row = cur.fetchone()
            if past_nw_row:
                past_nw = float(past_nw_row[0])
                net_worth_change = net_worth - past_nw
                days_diff = (datetime.date.today() - past_nw_row[1]).days
                net_worth_change_period = f"in last {days_diff} days" if days_diff > 0 else "recently"
            else:
                net_worth_change = 0.0
                net_worth_change_period = "recently"

            # Fetch Expenses
            params = [user_id]
            date_filter = ""
            if month and year:
                date_filter = " AND EXTRACT(MONTH FROM date) = %s AND EXTRACT(YEAR FROM date) = %s"
                params.extend([month, year])
                
            cur.execute(
                f"SELECT id, user_id, name, amount, category, bucket, icon, date, created_at FROM expenses WHERE user_id = %s{date_filter} ORDER BY date DESC",
                tuple(params)
            )
            expense_rows = cur.fetchall()
            expenses = []
            needs_spent = wants_spent = savings_spent = 0.0

            for row in expense_rows:
                bucket = row[5]
                amt = float(row[3])
                if bucket == 'Needs':
                    needs_spent += amt
                elif bucket == 'Wants':
                    wants_spent += amt
                elif bucket == 'Savings':
                    savings_spent += amt

                expenses.append(ExpenseOut(
                    id=row[0], user_id=row[1], name=row[2], amount=amt, category=row[4],
                    bucket=bucket, icon=row[6], date=row[7], created_at=row[8]
                ))

            total_spent = needs_spent + wants_spent + savings_spent
            total_left = total_income - total_spent

            # Budget targets: 50/30/20 of total income
            needs_budget = total_income * 0.50 if total_income > 0 else 37500.0
            wants_budget = total_income * 0.30 if total_income > 0 else 22500.0
            savings_budget = total_income * 0.20 if total_income > 0 else 15000.0

            # Dynamic Fin Score
            savings_rate = (savings_spent / total_income) if total_income > 0 else 0
            savings_score = min(30, int((savings_rate / 0.20) * 30)) if total_income > 0 else 0

            cur.execute("SELECT emi FROM liabilities WHERE user_id = %s", (user_id,))
            total_emis = sum(float(l[0]) for l in cur.fetchall())
            dti_rate = (total_emis / total_income) if total_income > 0 else 0
            dti_score = max(0, 30 - int((dti_rate / 0.30) * 30)) if total_income > 0 else 30

            cur.execute("SELECT amount FROM assets WHERE user_id = %s AND type IN ('bank', 'fd')", (user_id,))
            liquid_assets = sum(float(a[0]) for a in cur.fetchall())
            ef_target = needs_budget * 3
            ef_score = min(40, int((liquid_assets / ef_target) * 40)) if ef_target > 0 else 0

            fin_score = savings_score + dti_score + ef_score
            fin_score = max(10, min(100, fin_score))  # clamp between 10 and 100

            # Fetch Goals
            cur.execute(
                "SELECT id, user_id, name, target_amount, current_amount, status, color, created_at FROM goals WHERE user_id = %s",
                (user_id,)
            )
            goals = [GoalOut(
                id=r[0], user_id=r[1], name=r[2], target_amount=float(r[3]), current_amount=float(r[4]),
                status=r[5], color=r[6], created_at=r[7]
            ) for r in cur.fetchall()]

            # Fetch Bills
            cur.execute(
                "SELECT id, user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription, is_emi, emi_total_months, emi_months_paid, start_date, created_at FROM recurring_bills WHERE user_id = %s",
                (user_id,)
            )
            bills = [RecurringBillOut(
                id=r[0], user_id=r[1], name=r[2], amount=float(r[3]), frequency=r[4], category=r[5],
                bucket=r[6], due_day=r[7], monthly_equivalent=float(r[8]), is_active=r[9], is_subscription=r[10],
                is_emi=r[11], emi_total_months=r[12], emi_months_paid=r[13], start_date=r[14], created_at=r[15]
            ) for r in cur.fetchall()]

            return DashboardOverview(
                net_worth=net_worth,
                net_worth_change=net_worth_change,
                net_worth_change_period=net_worth_change_period,
                fin_score=fin_score,
                total_income=total_income,
                total_spent=total_spent,
                total_left=total_left,
                needs_spent=needs_spent,
                wants_spent=wants_spent,
                savings_spent=savings_spent,
                needs_budget=needs_budget,
                wants_budget=wants_budget,
                savings_budget=savings_budget,
                goals=goals,
                recent_expenses=expenses,
                upcoming_bills=bills
            )
    except Exception as e:
        import traceback
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/expense", response_model=ExpenseOut, status_code=status.HTTP_201_CREATED)
def add_expense(
    expense: ExpenseCreate,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute(
                """
                INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                RETURNING id, user_id, name, amount, category, bucket, icon, date, created_at
                """,
                (
                    current_user.id,
                    expense.name,
                    expense.amount,
                    expense.category,
                    expense.bucket,
                    expense.icon or 'receipt',
                    expense.date
                )
            )
            row = cur.fetchone()
            conn.commit()

            return ExpenseOut(
                id=row[0], user_id=row[1], name=row[2], amount=float(row[3]),
                category=row[4], bucket=row[5], icon=row[6], date=row[7], created_at=row[8]
            )
    except Exception as e:
        conn.rollback()
        import traceback
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=str(e))
