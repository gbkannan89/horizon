import logging
from typing import Dict, Any, Tuple
from ..database import get_db

logger = logging.getLogger(__name__)

def get_household_financial_data(household_id: int, db) -> Dict[str, Any]:
    """
    Gathers all financial data for a household (incomes, bank balances, EMIs, committed bills).
    """
    with db.cursor() as cur:
        # 1. Fetch user IDs in this household
        cur.execute("SELECT id FROM users WHERE household_id = %s", (household_id,))
        user_rows = cur.fetchall()
        user_ids = [r[0] for r in user_rows]
        
        if not user_ids:
            return {
                "user_incomes": 0.0,
                "contributing_income": 0.0,
                "total_net_income": 0.0,
                "liquid_savings": 0.0,
                "total_emis": 0.0,
                "committed_needs": 0.0,
                "committed_wants": 0.0,
                "committed_savings": 0.0,
                "total_committed_expenses": 0.0,
                "investment_classes_count": 0
            }

        # Format user_ids tuple for SQL IN query
        user_placeholders = ",".join(["%s"] * len(user_ids))
        
        # 2. Fetch incomes
        cur.execute(
            f"SELECT amount, frequency FROM incomes WHERE user_id IN ({user_placeholders})",
            tuple(user_ids)
        )
        income_rows = cur.fetchall()
        user_incomes = 0.0
        for amt, freq in income_rows:
            amt = float(amt)
            if freq.strip().lower() == "annual":
                user_incomes += amt / 12.0
            else:
                user_incomes += amt

        # 3. Fetch contributing member contributions
        cur.execute(
            "SELECT SUM(contribution_to_household) FROM contributing_members WHERE household_id = %s",
            (household_id,)
        )
        contrib_row = cur.fetchone()
        contributing_income = float(contrib_row[0]) if contrib_row and contrib_row[0] else 0.0
        
        total_net_income = user_incomes + contributing_income

        # 4. Fetch liquid savings (assets of type 'bank')
        cur.execute(
            f"SELECT SUM(amount) FROM assets WHERE user_id IN ({user_placeholders}) AND type = 'bank'",
            tuple(user_ids)
        )
        asset_row = cur.fetchone()
        liquid_savings = float(asset_row[0]) if asset_row and asset_row[0] else 0.0

        # 5. Fetch EMIs from liabilities
        cur.execute(
            f"SELECT SUM(emi) FROM liabilities WHERE user_id IN ({user_placeholders})",
            tuple(user_ids)
        )
        liab_row = cur.fetchone()
        total_emis = float(liab_row[0]) if liab_row and liab_row[0] else 0.0

        # 6. Fetch committed bills and group by 50/30/20 bucket
        cur.execute(
            f"SELECT monthly_equivalent, bucket FROM recurring_bills WHERE user_id IN ({user_placeholders}) AND is_active = TRUE",
            tuple(user_ids)
        )
        bill_rows = cur.fetchall()
        
        committed_needs = 0.0
        committed_wants = 0.0
        committed_savings = 0.0
        
        for amt, bucket in bill_rows:
            amt = float(amt)
            b = bucket.strip().lower()
            if b == "needs":
                committed_needs += amt
            elif b == "wants":
                committed_wants += amt
            elif b == "savings":
                committed_savings += amt

        # Total Committed Expenses = Rent/Bills + EMIs
        total_committed_expenses = committed_needs + committed_wants + committed_savings + total_emis

        # 7. Fetch investment diversity (unique asset types in the household: equity, debt, real estate, gold, insurance)
        # Note: In Phase 1, asset classes count is queried from the assets table types
        cur.execute(
            f"SELECT DISTINCT type FROM assets WHERE user_id IN ({user_placeholders})",
            tuple(user_ids)
        )
        type_rows = cur.fetchall()
        unique_types = {r[0].strip().lower() for r in type_rows if r[0]}
        
        # Also check if they have insurance
        cur.execute(
            f"SELECT COUNT(id) FROM recurring_bills WHERE user_id IN ({user_placeholders}) AND category = 'insurance'",
            tuple(user_ids)
        )
        has_insurance = cur.fetchone()[0] > 0
        if has_insurance:
            unique_types.add("insurance")
            
        investment_classes_count = len(unique_types)

        return {
            "user_incomes": user_incomes,
            "contributing_income": contributing_income,
            "total_net_income": total_net_income,
            "liquid_savings": liquid_savings,
            "total_emis": total_emis,
            "committed_needs": committed_needs,
            "committed_wants": committed_wants,
            "committed_savings": committed_savings,
            "total_committed_expenses": total_committed_expenses,
            "investment_classes_count": investment_classes_count
        }

def calculate_financial_health_score(household_id: int, db) -> Tuple[int, int, Dict[str, Any]]:
    """
    Calculates the Financial Health Score based on the 5-pillar algorithm.
    Returns: (score, max_possible_score, breakdown_details)
    """
    data = get_household_financial_data(household_id, db)
    income = data["total_net_income"]
    
    # Pillar 1: 50/30/20 Adherence (25 points)
    adherence_score = 0
    breakdown_50_30_20 = {"needs_ok": False, "wants_ok": False, "savings_ok": False}
    
    if income > 0:
        # Needs bucket: committed bills + EMIs <= 50%
        needs_ratio = (data["committed_needs"] + data["total_emis"]) / income
        # Wants bucket: subscriptions <= 30%
        wants_ratio = data["committed_wants"] / income
        # Savings bucket: committed savings >= 20%
        savings_ratio = data["committed_savings"] / income
        
        if needs_ratio <= 0.50:
            breakdown_50_30_20["needs_ok"] = True
        if wants_ratio <= 0.30:
            breakdown_50_30_20["wants_ok"] = True
        if savings_ratio >= 0.20:
            breakdown_50_30_20["savings_ok"] = True
            
        ok_count = sum(1 for v in breakdown_50_30_20.values() if v)
        if ok_count == 3:
            adherence_score = 25
        elif ok_count == 2:
            adherence_score = 15
        elif ok_count == 1:
            adherence_score = 8
        else:
            adherence_score = 0
    else:
        adherence_score = 0

    # Pillar 2: Emergency Fund (20 points)
    # Months of committed expenses covered by bank balance (liquid savings)
    # Committed Expenses = committed needs + EMIs
    committed_monthly_expense = data["committed_needs"] + data["total_emis"]
    emergency_score = 0
    emergency_months = 0.0
    
    if committed_monthly_expense > 0:
        emergency_months = data["liquid_savings"] / committed_monthly_expense
        if emergency_months >= 6.0:
            emergency_score = 20
        else:
            emergency_score = int((emergency_months / 6.0) * 20.0)
    elif data["liquid_savings"] > 0:
        emergency_score = 20
        emergency_months = 99.9  # arbitrary high value representing infinite coverage
    else:
        emergency_score = 0
        emergency_months = 0.0

    # Pillar 3: Debt-to-Income (20 points)
    debt_score = 0
    debt_ratio = 0.0
    if income > 0:
        debt_ratio = data["total_emis"] / income
        if debt_ratio < 0.30:
            debt_score = 20
        elif debt_ratio <= 0.50:
            # Proportional mapping from 30% (20 points) to 50% (0 points)
            debt_score = int(20 * (1 - (debt_ratio - 0.30) / 0.20))
        else:
            debt_score = 0
    else:
        debt_score = 0

    # Pillar 4: Goal Progress (20 points) - Locked in Phase 1
    goal_score = None
    
    # Pillar 5: Investment Diversity (15 points) - Locked in Phase 1
    diversity_score = None

    # Calculate final score (out of 65 points in Phase 1)
    calculated_score = adherence_score + emergency_score + debt_score
    max_possible_score = 65

    breakdown = {
        "income": income,
        "committedExpenses": committed_monthly_expense,
        "liquidSavings": data["liquid_savings"],
        "totalEMIs": data["total_emis"],
        "pillars": {
            "adherence": {
                "score": adherence_score,
                "max": 25,
                "details": breakdown_50_30_20
            },
            "emergency": {
                "score": emergency_score,
                "max": 20,
                "monthsCovered": round(emergency_months, 2)
            },
            "debtToIncome": {
                "score": debt_score,
                "max": 20,
                "ratio": round(debt_ratio * 100, 2)
            },
            "goalProgress": {
                "score": goal_score,
                "max": 20,
                "status": "locked_phase_1"
            },
            "investmentDiversity": {
                "score": diversity_score,
                "max": 15,
                "status": "locked_phase_1"
            }
        }
    }

    return calculated_score, max_possible_score, breakdown
