import math
import logging
from datetime import datetime, date, timedelta
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from ..database import get_db
from ..schemas import (
    SimulationInput, SimulationOut, NudgeOut,
    DebtOptimizerOut, SubscriptionInsightOut, UserOut
)
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/advisor", tags=["Advisor"])


# ─────────────────────────────────────────────────────────────────────────────
# SCENARIO SIMULATOR
# ─────────────────────────────────────────────────────────────────────────────
@router.post("/simulate", response_model=SimulationOut, status_code=201)
def run_simulation(
    sim: SimulationInput,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("SELECT SUM(amount) FROM incomes WHERE user_id = %s AND frequency = 'monthly'", (current_user.id,))
            total_income = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(amount) FROM assets WHERE user_id = %s", (current_user.id,))
            total_assets = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(outstanding) FROM liabilities WHERE user_id = %s", (current_user.id,))
            total_liabilities = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(monthly_equivalent) FROM recurring_bills WHERE user_id = %s", (current_user.id,))
            total_bills = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT amount FROM assets WHERE user_id = %s AND type IN ('bank', 'fd')", (current_user.id,))
            liquid_assets = sum(float(r[0]) for r in cur.fetchall())

        current_nw = total_assets - total_liabilities
        current_age = 30
        cur.execute("SELECT created_at FROM users WHERE id = %s", (current_user.id,))
        row = cur.fetchone()
        if row:
            years_since = (datetime.now() - row[0]).days / 365.0
            current_age = max(25, 30 + int(years_since))

        results = {}
        projected = []
        net_worth_impact = 0
        fin_score_impact = 0

        if sim.scenario_type == 'purchase' and sim.target_amount:
            down_payment = sim.target_amount * 0.2
            loan_amount = sim.target_amount - down_payment
            rate = sim.interest_rate or 9.0
            tenure = sim.tenure_months or 60
            monthly_rate = rate / 12 / 100
            emi = loan_amount * monthly_rate * (1 + monthly_rate) ** tenure / ((1 + monthly_rate) ** tenure - 1) if monthly_rate > 0 else loan_amount / tenure
            new_emi = emi
            new_debt = loan_amount

            net_worth_impact = -down_payment
            new_dti = ((total_bills + new_emi) / total_income * 100) if total_income > 0 else 100
            current_dti = ((total_bills) / total_income * 100) if total_income > 0 else 0
            dti_score_before = max(0, 30 - int((current_dti / 100 / 0.30) * 30))
            dti_score_after = max(0, 30 - int((new_dti / 100 / 0.30) * 30))
            fin_score_impact = dti_score_after - dti_score_before

            for year in range(1, 31):
                years_left = min(year, tenure // 12)
                remaining_principal = max(0, loan_amount - (new_emi * 12 * years_left))
                projected_nw = current_nw + net_worth_impact + (total_income * 0.2 * year) - remaining_principal
                projected.append({"year": current_age + year, "net_worth": round(projected_nw, 2)})

            results = {
                "type": "purchase",
                "purchase_amount": sim.target_amount,
                "down_payment": round(down_payment, 2),
                "loan_amount": round(loan_amount, 2),
                "monthly_emi": round(new_emi, 2),
                "net_worth_impact": round(net_worth_impact, 2),
                "dti_before": round(current_dti, 1),
                "dti_after": round(new_dti, 1),
                "fin_score_impact": fin_score_impact,
                "projection": projected,
            }

        elif sim.scenario_type == 'loan' and sim.loan_amount:
            loan_amount = sim.loan_amount
            rate = sim.interest_rate or 10.0
            tenure = sim.tenure_months or 60
            monthly_rate = rate / 12 / 100
            emi = loan_amount * monthly_rate * (1 + monthly_rate) ** tenure / ((1 + monthly_rate) ** tenure - 1) if monthly_rate > 0 else loan_amount / tenure
            total_payable = emi * tenure
            total_interest = total_payable - loan_amount

            net_worth_impact = loan_amount
            net_worth_impact -= down_payment if 'down_payment' in dir() else 0

            new_dti = ((total_bills + emi) / total_income * 100) if total_income > 0 else 100
            current_dti = ((total_bills) / total_income * 100) if total_income > 0 else 0
            dti_score_before = max(0, 30 - int((current_dti / 100 / 0.30) * 30))
            dti_score_after = max(0, 30 - int((new_dti / 100 / 0.30) * 30))
            fin_score_impact = dti_score_after - dti_score_before

            for year in range(1, 31):
                projected_nw = current_nw + net_worth_impact + (total_income * 0.2 * year) - max(0, loan_amount - (emi * 12 * min(year, tenure // 12)))
                projected.append({"year": current_age + year, "net_worth": round(projected_nw, 2)})

            results = {
                "type": "loan",
                "loan_amount": loan_amount,
                "interest_rate": rate,
                "tenure_months": tenure,
                "monthly_emi": round(emi, 2),
                "total_interest_payable": round(total_interest, 2),
                "total_payable": round(total_payable, 2),
                "dti_before": round(current_dti, 1),
                "dti_after": round(new_dti, 1),
                "fin_score_impact": fin_score_impact,
                "projection": projected,
            }

        elif sim.scenario_type == 'sip_increase' and sim.sip_amount:
            additional_sip = sim.sip_amount
            annual_return = 0.12
            monthly_return = annual_return / 12
            total_invested = 0
            for year in range(1, 31):
                months = year * 12
                future_value = additional_sip * (((1 + monthly_return) ** months - 1) / monthly_return) * (1 + monthly_return)
                total_invested = additional_sip * months
                wealth_gain = future_value - total_invested
                projected_nw = current_nw + future_value + (total_income * 0.15 * year)
                projected.append({
                    "year": current_age + year,
                    "net_worth": round(projected_nw, 2),
                    "invested": round(total_invested, 2),
                    "wealth": round(future_value, 2),
                })
            net_worth_impact = 0
            savings_rate_increase = (additional_sip / total_income * 100) if total_income > 0 else 0
            fin_score_impact = min(5, int(savings_rate_increase / 2))

            results = {
                "type": "sip_increase",
                "additional_sip": additional_sip,
                "annual_return": annual_return,
                "total_invested_30yrs": round(total_invested, 2),
                "wealth_created": round(wealth_gain if 'wealth_gain' in dir() else 0, 2),
                "savings_rate_increase": round(savings_rate_increase, 1),
                "fin_score_impact": fin_score_impact,
                "projection": projected,
            }

        elif sim.scenario_type == 'income_change' and sim.income_change:
            new_income = max(0, total_income + sim.income_change)
            delta = sim.income_change
            savings_boost = delta * 0.3

            for year in range(1, 31):
                accumulated = savings_boost * 12 * year * 1.08
                projected_nw = current_nw + (total_income * 0.2 * year) + accumulated
                projected.append({"year": current_age + year, "net_worth": round(projected_nw, 2)})

            results = {
                "type": "income_change",
                "income_delta": delta,
                "new_monthly_income": round(new_income, 2),
                "recommended_savings_boost": round(savings_boost, 2),
                "projection": projected,
            }

        elif sim.scenario_type == 'lumpsum' and sim.target_amount:
            lumpsum = sim.target_amount
            annual_return = 0.12
            for year in range(1, 31):
                future_value = lumpsum * (1 + annual_return) ** year
                projected_nw = current_nw + future_value + (total_income * 0.2 * year)
                projected.append({"year": current_age + year, "net_worth": round(projected_nw, 2), "investment_value": round(future_value, 2)})

            results = {
                "type": "lumpsum",
                "amount": lumpsum,
                "annual_return": annual_return,
                "projection": projected,
            }

        with conn.cursor() as cur:
            cur.execute(
                """INSERT INTO simulations (user_id, scenario_name, scenario_type, input_params, results)
                   VALUES (%s, %s, %s, %s, %s) RETURNING id, scenario_name, scenario_type, input_params, results, created_at""",
                (current_user.id, sim.scenario_name, sim.scenario_type,
                 sim.model_dump(), results)
            )
            row = cur.fetchone()
            conn.commit()

        return SimulationOut(
            id=row[0], scenario_name=row[1], scenario_type=row[2],
            input_params=row[3], results=row[4], created_at=row[5]
        )
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        logger.error(f"Simulation error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/simulations", response_model=List[SimulationOut])
def list_simulations(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute(
                "SELECT id, scenario_name, scenario_type, input_params, results, created_at FROM simulations WHERE user_id = %s ORDER BY created_at DESC LIMIT 20",
                (current_user.id,)
            )
            return [SimulationOut(
                id=r[0], scenario_name=r[1], scenario_type=r[2],
                input_params=r[3], results=r[4], created_at=r[5]
            ) for r in cur.fetchall()]
    except Exception as e:
        logger.error(f"List simulations error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ─────────────────────────────────────────────────────────────────────────────
# NUDGE ENGINE
# ─────────────────────────────────────────────────────────────────────────────
def generate_nudges(user_id: int, conn):
    try:
        today = date.today()
        with conn.cursor() as cur:
            cur.execute("SELECT SUM(amount) FROM incomes WHERE user_id = %s", (user_id,))
            total_income = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(monthly_equivalent) FROM recurring_bills WHERE user_id = %s AND bucket = 'Needs'", (user_id,))
            monthly_needs = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(monthly_equivalent) FROM recurring_bills WHERE user_id = %s", (user_id,))
            total_bills = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(amount) FROM assets WHERE user_id = %s AND type IN ('bank', 'fd')", (user_id,))
            liquid_assets = float(cur.fetchone()[0] or 0)

            cur.execute(
                "SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE user_id = %s AND EXTRACT(MONTH FROM date) = %s AND EXTRACT(YEAR FROM date) = %s AND bucket = 'Wants'",
                (user_id, today.month, today.year)
            )
            wants_spent = float(cur.fetchone()[0] or 0)

            cur.execute(
                "SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE user_id = %s AND EXTRACT(MONTH FROM date) = %s AND EXTRACT(YEAR FROM date) = %s AND bucket = 'Savings'",
                (user_id, today.month, today.year)
            )
            savings_spent = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT COALESCE(SUM(emi), 0) FROM liabilities WHERE user_id = %s", (user_id,))
            total_emis = float(cur.fetchone()[0] or 0)

            wants_budget = total_income * 0.30
            savings_budget = total_income * 0.20

            new_nudges = []

            if wants_spent > wants_budget * 0.80 and wants_budget > 0:
                new_nudges.append({
                    "category": "budget",
                    "severity": "warning",
                    "title": "Wants spending running high",
                    "message": f"You've spent ₹{wants_spent:,.0f} of your ₹{wants_budget:,.0f} wants budget ({(wants_spent/wants_budget*100):.0f}%). Consider cutting back on discretionary spending this month.",
                    "action_label": "View Budget",
                    "action_link": "budget",
                })

            if liquid_assets < monthly_needs * 3 and monthly_needs > 0:
                months_of_runway = liquid_assets / monthly_needs if monthly_needs > 0 else 0
                new_nudges.append({
                    "category": "emergency_fund",
                    "severity": "critical" if months_of_runway < 1 else "warning",
                    "title": "Emergency fund needs attention",
                    "message": f"Your liquid assets cover only {months_of_runway:.1f} months of needs. Target is 3-6 months (₹{monthly_needs*3:,.0f}).",
                    "action_label": "Build Fund",
                    "action_link": "portfolio",
                })

            if savings_spent < savings_budget * 0.5 and savings_budget > 0 and total_income > 0:
                new_nudges.append({
                    "category": "savings",
                    "severity": "warning",
                    "title": "Savings target behind",
                    "message": f"You've saved only ₹{savings_spent:,.0f} this month against your ₹{savings_budget:,.0f} goal. Try to allocate at least 20% of income to savings.",
                    "action_label": "Set Up SIP",
                    "action_link": "portfolio",
                })

            debt_income_ratio = (total_emis / total_income * 100) if total_income > 0 else 0
            if debt_income_ratio > 30:
                new_nudges.append({
                    "category": "debt",
                    "severity": "critical" if debt_income_ratio > 40 else "warning",
                    "title": "Debt-to-income ratio high",
                    "message": f"Your EMIs consume {debt_income_ratio:.0f}% of income (safe limit: 30%). Consider debt consolidation or accelerated repayment.",
                    "action_label": "Debt Strategy",
                    "action_link": "debt",
                })

            cur.execute("""
                SELECT b.id, b.name, b.amount, b.frequency
                FROM recurring_bills b
                LEFT JOIN subscription_insights si ON si.bill_id = b.id
                WHERE b.user_id = %s AND b.is_subscription = TRUE AND si.id IS NULL
                LIMIT 5
            """, (user_id,))
            unchecked_subs = cur.fetchall()

            for sub in unchecked_subs:
                monthly = sub[2] if sub[3] == 'monthly' else sub[2] / 12
                annual = monthly * 12
                if annual > 1000:
                    new_nudges.append({
                        "category": "subscription",
                        "severity": "info",
                        "title": f"Review: {sub[1]}",
                        "message": f"You're paying ₹{annual:,.0f}/year for {sub[1]}. Is this still needed?",
                        "action_label": "Review Subs",
                        "action_link": "subscriptions",
                    })

            for nudge in new_nudges:
                cur.execute(
                    """INSERT INTO nudges (user_id, category, severity, title, message, action_label, action_link)
                       SELECT %s, %s, %s, %s, %s, %s, %s
                       WHERE NOT EXISTS (
                           SELECT 1 FROM nudges
                           WHERE user_id = %s AND title = %s AND is_dismissed = FALSE AND created_at > CURRENT_DATE - INTERVAL '7 days'
                       )""",
                    (user_id, nudge["category"], nudge["severity"], nudge["title"], nudge["message"], nudge["action_label"], nudge["action_link"],
                     user_id, nudge["title"])
                )

            conn.commit()
    except Exception as e:
        conn.rollback()
        logger.error(f"Nudge generation error: {e}")


@router.get("/nudges", response_model=List[NudgeOut])
def get_nudges(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        generate_nudges(current_user.id, conn)
        with conn.cursor() as cur:
            cur.execute(
                "SELECT id, category, severity, title, message, action_label, action_link, is_dismissed, created_at FROM nudges WHERE user_id = %s AND is_dismissed = FALSE ORDER BY created_at DESC LIMIT 20",
                (current_user.id,)
            )
            return [NudgeOut(
                id=r[0], category=r[1], severity=r[2], title=r[3], message=r[4],
                action_label=r[5], action_link=r[6], is_dismissed=r[7], created_at=r[8]
            ) for r in cur.fetchall()]
    except Exception as e:
        logger.error(f"Get nudges error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.put("/nudges/{nudge_id}/dismiss")
def dismiss_nudge(
    nudge_id: int,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute(
                "UPDATE nudges SET is_dismissed = TRUE WHERE id = %s AND user_id = %s",
                (nudge_id, current_user.id)
            )
            conn.commit()
        return {"message": "Nudge dismissed"}
    except Exception as e:
        conn.rollback()
        logger.error(f"Dismiss nudge error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ─────────────────────────────────────────────────────────────────────────────
# DEBT OPTIMIZER
# ─────────────────────────────────────────────────────────────────────────────
@router.get("/debt-optimizer", response_model=DebtOptimizerOut)
def get_debt_optimizer(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute(
                "SELECT id, type, outstanding, emi, interest_rate, tenure_months FROM liabilities WHERE user_id = %s",
                (current_user.id,)
            )
            liabilities = cur.fetchall()

        if not liabilities:
            return DebtOptimizerOut(
                total_outstanding=0, total_monthly_emis=0, total_interest_paid=0,
                snowball={}, avalanche={}, recommended_strategy="Debt Free",
                estimated_freedom_months=0, total_interest_savable=0
            )

        debts = []
        for l in liabilities:
            debts.append({
                "id": l[0], "type": l[1], "outstanding": float(l[2]),
                "emi": float(l[3]), "rate": float(l[4]), "tenure": l[5] or 60
            })

        total_outstanding = sum(d["outstanding"] for d in debts)
        total_emis = sum(d["emi"] for d in debts)

        def calculate_snowball(debt_list):
            sorted_debts = sorted(debt_list, key=lambda x: x["outstanding"])
            remaining = [d.copy() for d in sorted_debts]
            total_interest = 0
            months = 0
            extra_payment = max(0, total_outstanding * 0.05)

            while any(d["outstanding"] > 0 for d in remaining):
                for d in remaining:
                    if d["outstanding"] <= 0:
                        continue
                    monthly_rate = d["rate"] / 12 / 100
                    interest = d["outstanding"] * monthly_rate
                    total_interest += interest
                    payment = min(d["emi"] + extra_payment, d["outstanding"] + interest)
                    d["outstanding"] -= (payment - interest)
                    if d["outstanding"] < 0:
                        d["outstanding"] = 0
                months += 1
                if months > 600:
                    break

            return {"months": months, "total_interest": round(total_interest, 2)}

        def calculate_avalanche(debt_list):
            sorted_debts = sorted(debt_list, key=lambda x: -x["rate"])
            remaining = [d.copy() for d in sorted_debts]
            total_interest = 0
            months = 0
            extra_payment = max(0, total_outstanding * 0.05)

            while any(d["outstanding"] > 0 for d in remaining):
                for d in remaining:
                    if d["outstanding"] <= 0:
                        continue
                    monthly_rate = d["rate"] / 12 / 100
                    interest = d["outstanding"] * monthly_rate
                    total_interest += interest
                    payment = min(d["emi"] + extra_payment, d["outstanding"] + interest)
                    d["outstanding"] -= (payment - interest)
                    if d["outstanding"] < 0:
                        d["outstanding"] = 0
                months += 1
                if months > 600:
                    break

            return {"months": months, "total_interest": round(total_interest, 2)}

        snowball = calculate_snowball(debts)
        avalanche = calculate_avalanche(debts)

        current_interest = sum(d["outstanding"] * d["rate"] / 100 for d in debts)
        best_strategy = "Avalanche" if avalanche["months"] <= snowball["months"] else "Snowball"
        interest_savable = max(0, current_interest - avalanche["total_interest"])

        return DebtOptimizerOut(
            total_outstanding=round(total_outstanding, 2),
            total_monthly_emis=round(total_emis, 2),
            total_interest_paid=round(current_interest, 2),
            snowball={"months_to_freedom": snowball["months"], "total_interest": snowball["total_interest"]},
            avalanche={"months_to_freedom": avalanche["months"], "total_interest": avalanche["total_interest"]},
            recommended_strategy=best_strategy,
            estimated_freedom_months=avalanche["months"],
            total_interest_savable=round(interest_savable, 2)
        )
    except Exception as e:
        logger.error(f"Debt optimizer error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ─────────────────────────────────────────────────────────────────────────────
# SUBSCRIPTION INSIGHTS
# ─────────────────────────────────────────────────────────────────────────────
@router.get("/subscription-insights", response_model=List[SubscriptionInsightOut])
def get_subscription_insights(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute(
                """SELECT b.id, b.name, b.amount, b.frequency, b.monthly_equivalent,
                          COALESCE(si.status, 'active'), COALESCE(si.savings_opportunity, 0),
                          si.notes, si.id
                   FROM recurring_bills b
                   LEFT JOIN subscription_insights si ON si.bill_id = b.id
                   WHERE b.user_id = %s AND b.is_subscription = TRUE
                   ORDER BY b.amount DESC""",
                (current_user.id,)
            )
            rows = cur.fetchall()
            results = []
            for r in rows:
                monthly = float(r[4]) if r[4] else (float(r[2]) if r[3] == 'monthly' else float(r[2]) / 12)
                annual = monthly * 12
                savings = float(r[6]) if r[6] else (annual if r[5] == 'unused' else 0)
                results.append(SubscriptionInsightOut(
                    id=r[8] or 0, bill_id=r[0], name=r[1],
                    amount=float(r[2]), monthly_cost=round(monthly, 2),
                    annual_cost=round(annual, 2), status=r[5],
                    savings_opportunity=round(savings, 2), notes=r[7]
                ))
            return results
    except Exception as e:
        logger.error(f"Subscription insights error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.put("/subscription-insights/{bill_id}")
def flag_subscription(
    bill_id: int, status: str,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        with conn.cursor() as cur:
            cur.execute("SELECT monthly_equivalent FROM recurring_bills WHERE id = %s AND user_id = %s", (bill_id, current_user.id))
            row = cur.fetchone()
            if not row:
                raise HTTPException(status_code=404, detail="Subscription not found")

            monthly = float(row[0])
            annual = monthly * 12
            savings = annual if status == 'unused' else 0

            cur.execute(
                """INSERT INTO subscription_insights (user_id, bill_id, status, monthly_cost, annual_cost, savings_opportunity)
                   VALUES (%s, %s, %s, %s, %s, %s)
                   ON CONFLICT ((SELECT id FROM subscription_insights WHERE bill_id = %s LIMIT 1)) DO UPDATE
                   SET status = EXCLUDED.status, savings_opportunity = EXCLUDED.savings_opportunity""",
                (current_user.id, bill_id, status, monthly, annual, savings, bill_id)
            )
            conn.commit()
        return {"message": f"Subscription flagged as '{status}'"}
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        logger.error(f"Flag subscription error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ─────────────────────────────────────────────────────────────────────────────
# LIFECYCLE SIMULATOR
# ─────────────────────────────────────────────────────────────────────────────
@router.post("/lifecycle-simulate")
def lifecycle_simulate(
    params: dict,
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    try:
        current_age = params.get("current_age", 30)
        retirement_age = params.get("retirement_age", 60)
        life_expectancy = params.get("life_expectancy", 100)
        inflation_rate = params.get("inflation_rate", 6.0) / 100
        expected_return = params.get("expected_return", 12.0) / 100
        monthly_expenses = params.get("monthly_expenses", 0)

        with conn.cursor() as cur:
            cur.execute("SELECT SUM(amount) FROM incomes WHERE user_id = %s", (current_user.id,))
            total_monthly_income = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(amount) FROM assets WHERE user_id = %s", (current_user.id,))
            total_assets = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(outstanding) FROM liabilities WHERE user_id = %s", (current_user.id,))
            total_liabilities = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT SUM(monthly_equivalent) FROM recurring_bills WHERE user_id = %s AND is_active = TRUE", (current_user.id,))
            total_bills = float(cur.fetchone()[0] or 0)

            cur.execute("SELECT amount FROM assets WHERE user_id = %s AND type IN ('bank', 'fd')", (current_user.id,))
            liquid_assets = sum(float(r[0]) for r in cur.fetchall())

        net_worth = total_assets - total_liabilities
        current_expenses = monthly_expenses if monthly_expenses > 0 else total_bills
        monthly_return = (1 + expected_return) ** (1 / 12) - 1
        monthly_inflation = (1 + inflation_rate) ** (1 / 12) - 1

        projection = []
        milestones = []
        yearly_net_worth = net_worth
        yearly_income = total_monthly_income * 12
        yearly_expenses = current_expenses * 12
        liquid_balance = liquid_assets

        for age in range(current_age, life_expectancy + 1):
            years_from_now = age - current_age
            is_retired = age >= retirement_age

            # Income stops at retirement
            if is_retired:
                annual_income = 0
                expense_income_gap = yearly_expenses
            else:
                annual_income = yearly_income * ((1 + inflation_rate) ** years_from_now)
                expense_income_gap = max(0, yearly_expenses - annual_income)

            # Expenses grow with inflation
            expenses_this_year = yearly_expenses * ((1 + inflation_rate) ** years_from_now)

            # Investment growth on net worth
            investment_return = yearly_net_worth * expected_return if not is_retired else yearly_net_worth * (expected_return * 0.7)

            if is_retired:
                net_worth_change = investment_return - expenses_this_year
            else:
                annual_savings = max(0, annual_income - expenses_this_year)
                net_worth_change = investment_return + annual_savings

            yearly_net_worth += net_worth_change
            liquid_balance += (annual_income - expenses_this_year) if not is_retired else (-expenses_this_year)

            if yearly_net_worth < 0:
                yearly_net_worth = 0

            is_valley = liquid_balance < 0 and years_from_now > 0
            if is_valley:
                milestones.append({
                    "age": age,
                    "type": "liquidity_valley",
                    "gap": round(abs(liquid_balance), 2),
                    "description": f"Cash flow gap of ₹{abs(liquid_balance):,.0f}/yr at age {age}",
                })

            if years_from_now % 5 == 0 or age == retirement_age or age == life_expectancy:
                projection.append({
                    "age": age,
                    "netWorth": round(yearly_net_worth, 2),
                    "income": round(annual_income, 2),
                    "expenses": round(expenses_this_year, 2),
                    "liquidBalance": round(liquid_balance, 2),
                    "isRetired": is_retired,
                })

            if age == retirement_age:
                milestones.append({
                    "age": age,
                    "type": "retirement",
                    "corpus": round(yearly_net_worth, 2),
                    "description": f"Retirement corpus of ₹{yearly_net_worth:,.0f} at age {age}",
                })

        milestones.append({
            "age": life_expectancy,
            "type": "end",
            "estate": round(max(0, yearly_net_worth), 2),
            "description": f"Estate value of ₹{max(0, yearly_net_worth):,.0f} at age {life_expectancy}",
        })

        return {
            "projection": projection,
            "milestones": milestones,
            "finalNetWorth": round(max(0, yearly_net_worth), 2),
            "yearsProjected": life_expectancy - current_age,
        }

    except Exception as e:
        logger.error(f"Lifecycle simulation error: {e}")
        raise HTTPException(status_code=500, detail=str(e))
