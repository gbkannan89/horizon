import logging
import statistics
from datetime import date, timedelta, datetime
from typing import List, Dict, Any, Optional, Tuple
from collections import defaultdict

logger = logging.getLogger(__name__)

AMOUNT_TOLERANCE = 0.10  # 10% variance allowed for same merchant amounts

FREQUENCY_RANGES = {
    "weekly": (6, 8),
    "monthly": (25, 35),
    "quarterly": (85, 95),
    "half_yearly": (170, 190),
    "yearly": (360, 370),
}


def _detect_frequency(intervals: List[float]) -> Tuple[Optional[str], float]:
    if len(intervals) < 2:
        return None, 0.0
    mean_interval = statistics.mean(intervals)
    if len(intervals) >= 2:
        std_dev = statistics.stdev(intervals) if len(intervals) > 1 else 0
    else:
        std_dev = 0

    cv = std_dev / mean_interval if mean_interval > 0 else 1.0

    best_freq = None
    best_dist = float("inf")
    for freq, (lo, hi) in FREQUENCY_RANGES.items():
        if lo <= mean_interval <= hi:
            dist = abs(mean_interval - (lo + hi) / 2)
            if dist < best_dist:
                best_dist = dist
                best_freq = freq

    confidence = 0.0
    if best_freq and cv < 0.30:
        freq_boost = max(0, 1.0 - cv / 0.30)
        count_boost = min(1.0, (len(intervals) + 1) / 12)
        confidence = min(1.0, 0.3 + freq_boost * 0.4 + count_boost * 0.3)

    return best_freq, round(confidence, 2)


def detect_recurring_transactions(user_id: int, conn) -> List[Dict[str, Any]]:
    with conn.cursor() as cur:
        cur.execute("""
            SELECT name, amount, date
            FROM expenses
            WHERE user_id = %s AND amount > 0
            ORDER BY name, amount, date
        """, (user_id,))
        rows = cur.fetchall()

    if not rows:
        return []

    groups: Dict[str, List[dict]] = defaultdict(list)
    for name, amount, txn_date in rows:
        key = name.strip().lower()
        for existing_key in list(groups.keys()):
            if existing_key == key:
                break
        else:
            key = name.strip().lower()
        groups[key].append({"name": name, "amount": float(amount), "date": txn_date})

    suggestions = []

    for merchant, txns in groups.items():
        display_name = txns[0]["name"]
        amounts = [t["amount"] for t in txns]
        dates = sorted([t["date"] for t in txns])

        if len(dates) < 3:
            continue

        mean_amt = statistics.mean(amounts)
        if mean_amt > 0:
            amount_variation = max(abs(a - mean_amt) / mean_amt for a in amounts)
        else:
            amount_variation = 1.0

        intervals = []
        for i in range(1, len(dates)):
            diff = (dates[i] - dates[i - 1]).days
            if diff > 0:
                intervals.append(float(diff))

        frequency, confidence = _detect_frequency(intervals)

        if frequency and confidence >= 0.4:
            if amount_variation <= AMOUNT_TOLERANCE:
                confidence = min(1.0, confidence + 0.15)

            category, bucket = _suggest_category(display_name)
            is_sub = confidence >= 0.6 and mean_amt <= 2000 and frequency == "monthly"

            suggestions.append({
                "name": display_name,
                "amount": round(mean_amt, 2),
                "frequency": frequency,
                "confidence": confidence,
                "occurrences": len(dates),
                "first_date": dates[0].isoformat(),
                "last_date": dates[-1].isoformat(),
                "category": category,
                "bucket": bucket,
                "is_subscription": is_sub,
                "amount_variation": round(amount_variation, 3),
            })

    suggestions.sort(key=lambda s: -s["confidence"])
    return suggestions


def _suggest_category(description: str) -> Tuple[str, str]:
    desc = description.upper()
    keywords = {
        "ENTERTAINMENT": ("Entertainment", "Wants"),
        "NETFLIX": ("Entertainment", "Wants"),
        "AMAZON PRIME": ("Entertainment", "Wants"),
        "SPOTIFY": ("Entertainment", "Wants"),
        "YOUTUBE": ("Entertainment", "Wants"),
        "HOTSTAR": ("Entertainment", "Wants"),
        "ZOOM": ("Entertainment", "Wants"),
        "RENT": ("Housing", "Needs"),
        "ELECTRICITY": ("Utilities", "Needs"),
        "WATER": ("Utilities", "Needs"),
        "INTERNET": ("Utilities", "Needs"),
        "BESCOM": ("Utilities", "Needs"),
        "AIRTEL": ("Utilities", "Needs"),
        "JIO": ("Utilities", "Needs"),
        "GYM": ("Health", "Wants"),
        "CULT": ("Health", "Wants"),
        "INSURANCE": ("Insurance", "Needs"),
        "LIC": ("Savings", "Savings"),
        "SIP": ("Investment", "Savings"),
        "MUTUAL FUND": ("Investment", "Savings"),
        "PPF": ("Savings", "Savings"),
        "GROWW": ("Investment", "Savings"),
        "ZERODHA": ("Investment", "Savings"),
        "LOAN": ("Debt", "Needs"),
        "EMI": ("Debt", "Needs"),
        "DMART": ("Groceries", "Needs"),
        "BIGBASKET": ("Groceries", "Needs"),
        "ZEITO": ("Groceries", "Needs"),
        "SWIGGY": ("Food", "Wants"),
        "ZOMATO": ("Food", "Wants"),
        "UBER": ("Transport", "Needs"),
        "OLA": ("Transport", "Needs"),
        "RAPIDO": ("Transport", "Needs"),
        "PETROL": ("Transport", "Needs"),
    }
    for keyword, (cat, bucket) in keywords.items():
        if keyword in desc:
            return cat, bucket
    return ("Misc", "Wants")


def analyze_spending_patterns(user_id: int, conn) -> Dict[str, Any]:
    with conn.cursor() as cur:
        cur.execute("SELECT MAX(date) FROM expenses WHERE user_id = %s", (user_id,))
        max_date_row = cur.fetchone()
        today_date = max_date_row[0] if max_date_row and max_date_row[0] else date.today()

    first_current = today_date.replace(day=1)
    first_prev = (first_current - timedelta(days=1)).replace(day=1)
    first_prev_prev = (first_prev - timedelta(days=1)).replace(day=1)

    with conn.cursor() as cur:
        cur.execute("""
            SELECT category, SUM(amount) as total
            FROM expenses
            WHERE user_id = %s AND date >= %s AND date <= %s
            GROUP BY category ORDER BY total DESC
        """, (user_id, first_current, today_date))
        current_category = {r[0]: float(r[1]) for r in cur.fetchall()}

        cur.execute("""
            SELECT category, SUM(amount) as total
            FROM expenses
            WHERE user_id = %s AND date >= %s AND date < %s
            GROUP BY category
        """, (user_id, first_prev, first_current))
        prev_category = {r[0]: float(r[1]) for r in cur.fetchall()}

        category_trends = []
        all_cats = set(current_category.keys()) | set(prev_category.keys())
        for cat in sorted(all_cats):
            curr = current_category.get(cat, 0)
            prev = prev_category.get(cat, 0)
            change_pct = round(((curr - prev) / prev * 100), 1) if prev > 0 else (100.0 if curr > 0 else 0)
            category_trends.append({
                "category": cat,
                "current_amount": round(curr, 2),
                "previous_amount": round(prev, 2),
                "change_pct": change_pct,
            })

        cur.execute("""
            SELECT EXTRACT(DOW FROM date) as dow, SUM(amount) as total
            FROM expenses
            WHERE user_id = %s AND date >= %s
            GROUP BY dow ORDER BY dow
        """, (user_id, first_prev))
        weekday_rows = cur.fetchall()

        day_names = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
        weekday_distribution = {}
        weekend_total = 0.0
        weekday_total = 0.0
        weekend_count = 0
        weekday_count = 0
        for dow, total in weekday_rows:
            idx = int(dow)
            label = day_names[idx]
            amt = float(total or 0)
            weekday_distribution[label] = round(amt, 2)
            if idx in (0, 6, 5):
                weekend_total += amt
                weekend_count += 1
            else:
                weekday_total += amt
                weekday_count += 1

        avg_weekend = weekend_total / weekend_count if weekend_count > 0 else 0
        avg_weekday = weekday_total / weekday_count if weekday_count > 0 else 1
        weekend_boost_pct = round(((avg_weekend - avg_weekday) / avg_weekday) * 100, 1) if avg_weekday > 0 else 0

        cur.execute("""
            SELECT category, AVG(amount) as mean_amt, STDDEV(amount) as std_amt
            FROM expenses
            WHERE user_id = %s AND date >= %s
            GROUP BY category
        """, (user_id, first_prev_prev))
        cat_stats = {r[0]: {"mean": float(r[1] or 0), "std": float(r[2] or 0)} for r in cur.fetchall()}

        cur.execute("""
            SELECT id, name, amount, category, date
            FROM expenses
            WHERE user_id = %s AND date >= %s
            ORDER BY date DESC
        """, (user_id, first_prev_prev))
        all_expenses = cur.fetchall()

        anomalies = []
        for exp_id, name, amt, cat, exp_date in all_expenses:
            amt = float(amt)
            stats = cat_stats.get(cat, {"mean": 0, "std": 0})
            if stats["std"] > 0 and stats["mean"] > 0:
                z_score = (amt - stats["mean"]) / stats["std"]
                if abs(z_score) > 2.0:
                    anomalies.append({
                        "id": exp_id,
                        "name": name,
                        "amount": round(amt, 2),
                        "category": cat,
                        "date": exp_date.isoformat(),
                        "z_score": round(z_score, 2),
                        "reason": "Unusually high" if z_score > 0 else "Unusually low",
                    })

        cur.execute("""
            SELECT SUM(amount) FROM expenses
            WHERE user_id = %s AND date >= %s AND date <= %s
        """, (user_id, first_current, today_date))
        current_month_spent = float(cur.fetchone()[0] or 0)

        cur.execute("""
            SELECT AVG(monthly) FROM (
                SELECT SUM(amount) as monthly FROM expenses
                WHERE user_id = %s AND date >= %s
                GROUP BY DATE_TRUNC('month', date)
            ) sub
        """, (user_id, first_prev_prev))
        avg_monthly = float(cur.fetchone()[0] or current_month_spent)

        days_in_month = (first_current.replace(month=first_current.month % 12 + 1, day=1) - timedelta(days=1)).day
        days_elapsed = (today_date - first_current).days + 1
        daily_rate = current_month_spent / days_elapsed if days_elapsed > 0 else 0
        projected_total = round(daily_rate * days_in_month, 2)

        cur.execute("""
            SELECT name, SUM(amount) as total
            FROM expenses
            WHERE user_id = %s AND date >= %s
            GROUP BY name ORDER BY total DESC LIMIT 10
        """, (user_id, first_current))
        top_merchants = [{"name": r[0], "amount": round(float(r[1]), 2)} for r in cur.fetchall()]

    category_breakdown = [{"category": k, "amount": v} for k, v in current_category.items()]

    return {
        "category_breakdown": category_breakdown,
        "category_trends": category_trends,
        "weekday_distribution": weekday_distribution,
        "weekend_boost_pct": weekend_boost_pct,
        "anomalies": anomalies[:20],
        "forecast": {
            "current_spent": round(current_month_spent, 2),
            "projected_total": projected_total,
            "average_monthly": round(avg_monthly, 2),
            "days_elapsed": days_elapsed,
            "days_in_month": days_in_month,
        },
        "top_merchants": top_merchants,
    }


def detect_subscriptions(user_id: int, conn) -> List[Dict[str, Any]]:
    recurring = detect_recurring_transactions(user_id, conn)
    subscription_candidates = [
        s for s in recurring
        if s["is_subscription"] and s["confidence"] >= 0.5
    ]

    with conn.cursor() as cur:
        cur.execute("""
            SELECT id, name FROM recurring_bills
            WHERE user_id = %s AND is_subscription = TRUE AND is_active = TRUE
        """, (user_id,))
        existing_subs = {r[1].strip().lower() for r in cur.fetchall()}

    results = []
    for s in subscription_candidates:
        is_new = s["name"].strip().lower() not in existing_subs
        monthly_cost = s["amount"]
        annual_cost = monthly_cost * 12
        savings_opportunity = annual_cost * 0.3 if is_new else 0

        results.append({
            "name": s["name"],
            "amount": s["amount"],
            "monthly_cost": monthly_cost,
            "annual_cost": round(annual_cost, 2),
            "frequency": s["frequency"],
            "confidence": s["confidence"],
            "is_new": is_new,
            "savings_opportunity": round(savings_opportunity, 2),
            "category": s["category"],
            "bucket": s["bucket"],
            "occurrences": s["occurrences"],
            "last_date": s["last_date"],
        })

    return results


def detect_lapsed_subscriptions(user_id: int, conn) -> List[Dict[str, Any]]:
    with conn.cursor() as cur:
        cur.execute("""
            SELECT rb.id, rb.name, rb.amount, rb.frequency, rb.monthly_equivalent
            FROM recurring_bills rb
            WHERE rb.user_id = %s AND rb.is_subscription = TRUE AND rb.is_active = TRUE
        """, (user_id,))
        bills = cur.fetchall()

    lapsed = []
    for bill_id, name, amount, freq, monthly in bills:
        monthly_amount = float(monthly or amount or 0)
        three_months_ago = date.today() - timedelta(days=90)

        with conn.cursor() as cur:
            cur.execute("""
                SELECT COUNT(*), COALESCE(SUM(amount), 0)
                FROM expenses
                WHERE user_id = %s AND name ILIKE %s AND date >= %s
            """, (user_id, f"%{name}%", three_months_ago))
            row = cur.fetchone()
            charge_count = row[0] if row else 0
            total_charged = float(row[1]) if row and row[1] else 0

        if charge_count == 0:
            annual_cost = monthly_amount * 12
            lapsed.append({
                "bill_id": bill_id,
                "name": name,
                "monthly_cost": monthly_amount,
                "annual_cost": round(annual_cost, 2),
                "status": "lapsed",
                "savings_opportunity": round(annual_cost, 2),
                "days_since_last_charge": (date.today() - three_months_ago).days,
            })
        elif charge_count == 1 and total_charged < monthly_amount * 1.5:
            annual_cost = monthly_amount * 12
            lapsed.append({
                "bill_id": bill_id,
                "name": name,
                "monthly_cost": monthly_amount,
                "annual_cost": round(annual_cost, 2),
                "status": "underused",
                "savings_opportunity": round(annual_cost * 0.5, 2),
                "days_since_last_charge": (date.today() - three_months_ago).days,
            })

    return lapsed


def generate_spending_nudges(user_id: int, conn) -> List[Dict[str, Any]]:
    new_nudges = []
    
    with conn.cursor() as cur:
        cur.execute("SELECT MAX(date) FROM expenses WHERE user_id = %s", (user_id,))
        max_date_row = cur.fetchone()
        today_date = max_date_row[0] if max_date_row and max_date_row[0] else date.today()

    first_current = today_date.replace(day=1)
    first_prev = (first_current - timedelta(days=1)).replace(day=1)

    with conn.cursor() as cur:
        cur.execute("SELECT SUM(amount) FROM incomes WHERE user_id = %s AND frequency = 'monthly'", (user_id,))
        total_income = float(cur.fetchone()[0] or 0)

        cur.execute("""
            SELECT bucket, SUM(amount)
            FROM expenses WHERE user_id = %s AND date >= %s
            GROUP BY bucket
        """, (user_id, first_current))
        bucket_spend = {r[0]: float(r[1]) for r in cur.fetchall()}
        needs_spent = bucket_spend.get("Needs", 0)
        wants_spent = bucket_spend.get("Wants", 0)
        savings_spent = bucket_spend.get("Savings", 0)

        if total_income > 0:
            needs_pct = needs_spent / total_income * 100
            wants_pct = wants_spent / total_income * 100
            savings_pct = savings_spent / total_income * 100

            if needs_pct > 55:
                new_nudges.append({
                    "category": "budget",
                    "severity": "warning",
                    "title": "Needs spending exceeding 50% target",
                    "message": f"Your Needs are at {needs_pct:.0f}% of income (₹{needs_spent:,.0f}). Try to keep it under 50%.",
                    "action_label": "View Budget",
                    "action_link": "budget",
                })

            if wants_pct > 35:
                new_nudges.append({
                    "category": "budget",
                    "severity": "warning",
                    "title": "Wants spending running high",
                    "message": f"You've spent {wants_pct:.0f}% of income on Wants (₹{wants_spent:,.0f}). Target is 30%.",
                    "action_label": "Review Expenses",
                    "action_link": "budget",
                })

            if savings_pct < 10:
                new_nudges.append({
                    "category": "savings",
                    "severity": "warning",
                    "title": "Savings rate below 20% target",
                    "message": f"Only {savings_pct:.0f}% of income saved this month. Try to reach the 20% target.",
                    "action_label": "Set Up SIP",
                    "action_link": "portfolio",
                })

        patterns = analyze_spending_patterns(user_id, conn)
        for trend in patterns["category_trends"]:
            if trend["change_pct"] > 50 and trend["previous_amount"] > 500:
                new_nudges.append({
                    "category": "spending",
                    "severity": "info",
                    "title": f"{trend['category']} spending spike",
                    "message": f"Spending on {trend['category']} is up {trend['change_pct']:.0f}% this month (₹{trend['current_amount']:,.0f} vs ₹{trend['previous_amount']:,.0f} last month).",
                    "action_label": "Review",
                    "action_link": "budget",
                })

        if patterns["weekend_boost_pct"] > 50:
            new_nudges.append({
                "category": "spending",
                "severity": "info",
                "title": "Weekend spending spike detected",
                "message": f"Your weekend spending is {patterns['weekend_boost_pct']:.0f}% higher than weekdays. Consider setting weekend spending limits.",
                "action_label": "View Patterns",
                "action_link": "budget",
            })

        recurring = detect_recurring_transactions(user_id, conn)
        high_conf_recurring = [r for r in recurring if r["confidence"] >= 0.7 and r["occurrences"] >= 3]
        for r in high_conf_recurring[:3]:
            with conn.cursor() as c:
                c.execute(
                    "SELECT 1 FROM recurring_bills WHERE user_id = %s AND name ILIKE %s",
                    (user_id, f"%{r['name']}%")
                )
                already_added = c.fetchone() is not None

            if not already_added:
                new_nudges.append({
                    "category": "subscription",
                    "severity": "info",
                    "title": f"Recurring detected: {r['name']}",
                    "message": f"'{r['name']}' appears every {r['frequency']} (~₹{r['amount']:,.0f}). Would you like to add it as a recurring bill?",
                    "action_label": "Add Recurring Bill",
                    "action_link": "bills",
                })

        for anomaly in patterns["anomalies"][:3]:
            new_nudges.append({
                "category": "spending",
                "severity": "info",
                "title": f"Unusual transaction: {anomaly['name']}",
                "message": f"A {anomaly['reason'].lower()} transaction of ₹{anomaly['amount']:,.0f} in '{anomaly['category']}' (z-score: {anomaly['z_score']}).",
                "action_label": "Review",
                "action_link": "budget",
            })

    return new_nudges


def run_full_analysis(user_id: int, conn) -> Dict[str, Any]:
    recurring = detect_recurring_transactions(user_id, conn)
    patterns = analyze_spending_patterns(user_id, conn)
    subscriptions = detect_subscriptions(user_id, conn)
    lapsed = detect_lapsed_subscriptions(user_id, conn)
    nudges = generate_spending_nudges(user_id, conn)

    with conn.cursor() as cur:
        for nudge in nudges:
            cur.execute(
                """INSERT INTO nudges (user_id, category, severity, title, message, action_label, action_link)
                   SELECT %s, %s, %s, %s, %s, %s, %s
                   WHERE NOT EXISTS (
                       SELECT 1 FROM nudges
                       WHERE user_id = %s AND title = %s AND is_dismissed = FALSE AND created_at > CURRENT_DATE - INTERVAL '7 days'
                   )""",
                (user_id, nudge["category"], nudge["severity"], nudge["title"], nudge["message"],
                 nudge["action_label"], nudge["action_link"],
                 user_id, nudge["title"])
            )
        conn.commit()

    return {
        "recurring_detections": recurring,
        "subscription_candidates": subscriptions,
        "lapsed_subscriptions": lapsed,
        "spending_patterns": patterns,
        "nudges_generated": len(nudges),
    }


def add_confirmed_recurring(user_id: int, confirmation: Dict[str, Any], conn) -> Dict[str, Any]:
    name = confirmation["name"]
    amount = confirmation["amount"]
    frequency = confirmation["frequency"]
    category = confirmation.get("category", "Misc")
    bucket = confirmation.get("bucket", "Wants")
    is_subscription = confirmation.get("is_subscription", frequency == "monthly" and amount <= 2000)

    due_day = confirmation.get("due_day", 1)
    monthly_eq = _calculate_monthly_equivalent(amount, frequency)

    with conn.cursor() as cur:
        cur.execute(
            "SELECT id FROM recurring_bills WHERE user_id = %s AND name ILIKE %s",
            (user_id, name)
        )
        existing = cur.fetchone()
        if existing:
            return {"error": "A recurring bill with this name already exists", "id": existing[0]}

        cur.execute("""
            INSERT INTO recurring_bills
                (user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, TRUE, %s)
            RETURNING id, user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription
        """, (user_id, name, amount, frequency, category, bucket, due_day, monthly_eq, is_subscription))
        row = cur.fetchone()
        conn.commit()

        if row:
            return {
                "id": row[0],
                "name": row[2],
                "amount": float(row[3]),
                "frequency": row[4],
                "category": row[5],
                "bucket": row[6],
                "due_day": row[7],
                "monthly_equivalent": float(row[8]),
                "is_subscription": row[10],
            }
        return {"error": "Could not create recurring bill"}


def _calculate_monthly_equivalent(amount: float, frequency: str) -> float:
    freq_map = {
        "weekly": 4.33,
        "monthly": 1.0,
        "quarterly": 1.0 / 3.0,
        "half_yearly": 1.0 / 6.0,
        "yearly": 1.0 / 12.0,
    }
    multiplier = freq_map.get(frequency, 1.0)
    return round(amount * multiplier, 2)
