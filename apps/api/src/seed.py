import logging
import bcrypt
import datetime
import random
import json
from .config import settings
from .database import get_db

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def get_password_hash(password: str) -> str:
    pwd_bytes = password.encode("utf-8")
    salt = bcrypt.gensalt()
    return bcrypt.hashpw(pwd_bytes, salt).decode("utf-8")

def seed_data():
    email = "kannan@email.com"
    password = "twelve12"
    name = "Kannan"
    phone = "9578176161"

    logger.info("Connecting to database...")
    db = next(get_db())

    today = datetime.date.today()

    with db.cursor() as cur:
        cur.execute("SELECT id FROM users WHERE email = %s", (email,))
        existing = cur.fetchone()
        if existing:
            logger.info(f"User {email} already exists. Deleting to re-seed...")
            user_id = existing[0]
            cur.execute("DELETE FROM wishlist_items WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM net_worth_snapshots WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM vehicles WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM simulations WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM nudges WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM subscription_insights WHERE bill_id IN (SELECT id FROM recurring_bills WHERE user_id = %s)", (user_id,))
            cur.execute("DELETE FROM recurring_bills WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM expenses WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM goals WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM insurances WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM liabilities WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM assets WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM incomes WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM health_scores WHERE household_id = (SELECT household_id FROM users WHERE id = %s)", (user_id,))
            cur.execute("DELETE FROM contributing_members WHERE household_id = (SELECT household_id FROM users WHERE id = %s)", (user_id,))
            cur.execute("DELETE FROM collection_members WHERE collection_id IN (SELECT id FROM collections WHERE user_id = %s)", (user_id,))
            cur.execute("DELETE FROM collections WHERE user_id = %s", (user_id,))
            cur.execute("DELETE FROM users WHERE id = %s", (user_id,))
            cur.execute("DELETE FROM households WHERE id NOT IN (SELECT household_id FROM users WHERE household_id IS NOT NULL)")

        try:
            # ── HOUSEHOLD ──────────────────────────────────────────────────
            logger.info("Seeding household...")
            cur.execute(
                "INSERT INTO households (name) VALUES (%s) RETURNING id",
                (f"{name}'s Household",)
            )
            household_id = cur.fetchone()[0]

            # ── USER ───────────────────────────────────────────────────────
            logger.info("Seeding user...")
            password_hash = get_password_hash(password)
            cur.execute(
                """
                INSERT INTO users (email, password_hash, name, phone, user_type, risk_profile, household_id)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                RETURNING id
                """,
                (email, password_hash, name, phone, "salaried", "moderate", household_id)
            )
            user_id = cur.fetchone()[0]

            # ── CONTRIBUTING MEMBERS ──────────────────────────────────────
            logger.info("Seeding contributing members...")
            cur.execute(
                """
                INSERT INTO contributing_members (household_id, name, monthly_income, contribution_to_household, relationship)
                VALUES (%s, %s, %s, %s, %s)
                """,
                (household_id, "Priya", 45000.00, 15000.00, "spouse")
            )

            # ── INCOMES ────────────────────────────────────────────────────
            logger.info("Seeding incomes...")
            salary_amount = 116313.23
            cur.execute(
                """
                INSERT INTO incomes (user_id, label, type, amount, frequency, pf_employee, pf_employer, shares_deduction)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                """,
                (user_id, "Monthly Salary", "salary", salary_amount, "monthly", 9884.98, 9884.98, 13889.00)
            )
            cur.execute(
                """
                INSERT INTO incomes (user_id, label, type, amount, frequency, pf_employee, pf_employer, shares_deduction)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                """,
                (user_id, "Rental Income", "passive", 12000.00, "monthly", 0.00, 0.00, 0.00)
            )

            # ── ASSETS ─────────────────────────────────────────────────────
            logger.info("Seeding assets...")
            cur.execute(
                """
                INSERT INTO assets (user_id, type, name, amount, interest_rate, generates_income, income_frequency)
                VALUES
                (%s, 'bank', 'HDFC Savings Account', 185000.00, 3.5, FALSE, NULL),
                (%s, 'bank', 'ICICI Salary Account', 42000.00, 3.0, FALSE, NULL),
                (%s, 'fd', 'SBI Fixed Deposit', 300000.00, 7.1, FALSE, NULL),
                (%s, 'fd', 'HDFC Tax Saver FD', 150000.00, 6.8, FALSE, NULL),
                (%s, 'physical', 'Gold Coins (50g)', 325000.00, 0.0, FALSE, NULL),
                (%s, 'physical', 'Rental Property', 4500000.00, 0.0, TRUE, 'monthly')
                """,
                (user_id, user_id, user_id, user_id, user_id, user_id)
            )

            # ── LIABILITIES ───────────────────────────────────────────────
            logger.info("Seeding liabilities...")
            cur.execute(
                """
                INSERT INTO liabilities (user_id, type, name, outstanding, emi, interest_rate, tenure_months)
                VALUES
                (%s, 'car_loan', 'HDFC Car Loan', 380000.00, 12500.00, 8.5, 48),
                (%s, 'personal_loan', 'Navi Personal Loan', 25000.00, 5000.00, 14.0, 12),
                (%s, 'credit_card', 'HDFC Credit Card', 12400.00, 12400.00, 36.0, 1)
                """,
                (user_id, user_id, user_id)
            )

            # ── GOALS ──────────────────────────────────────────────────────
            logger.info("Seeding goals...")
            cur.execute(
                """
                INSERT INTO goals (user_id, name, target_amount, current_amount, status, color)
                VALUES
                (%s, 'Emergency Fund', 300000.00, 185000.00, 'On Track', '#059669'),
                (%s, 'Europe Trip 2027', 400000.00, 125000.00, 'Behind', '#DC2626'),
                (%s, 'New MacBook Pro', 250000.00, 180000.00, 'Ahead', '#1E3A8A'),
                (%s, 'Down Payment for Home', 2000000.00, 350000.00, 'On Track', '#7C3AED')
                """,
                (user_id, user_id, user_id, user_id)
            )

            # ── RECURRING BILLS ───────────────────────────────────────────
            logger.info("Seeding recurring bills...")
            bills = [
                ('Rent', 18000.00, 'Housing', 'Needs', 1, False, None, None, None, True),
                ('Electricity Bill', 2400.00, 'Utilities', 'Needs', 5, False, None, None, None, True),
                ('Water Bill', 850.00, 'Utilities', 'Needs', 8, False, None, None, None, True),
                ('Internet', 1499.00, 'Utilities', 'Needs', 10, False, None, None, None, True),
                ('Car Loan EMI', 12500.00, 'Debt', 'Needs', 3, True, 48, 14, datetime.date(today.year - 2, 1, 15), False),
                ('Personal Loan EMI', 5000.00, 'Debt', 'Needs', 7, True, 12, 8, datetime.date(today.year - 1, 6, 1), False),
                ('SIP - Index Fund', 15000.00, 'Investment', 'Savings', 12, False, None, None, None, False),
                ('SIP - ELSS', 5000.00, 'Investment', 'Savings', 15, False, None, None, None, False),
                ('PPF Contribution', 5000.00, 'Investment', 'Savings', 28, False, None, None, None, False),
                ('Netflix', 649.00, 'Entertainment', 'Wants', 16, False, None, None, None, True),
                ('Amazon Prime', 299.00, 'Entertainment', 'Wants', 20, False, None, None, None, True),
                ('Spotify Premium', 119.00, 'Entertainment', 'Wants', 22, False, None, None, None, True),
                ('Gym Membership', 1500.00, 'Health', 'Wants', 18, False, None, None, None, True),
                ('YouTube Premium', 139.00, 'Entertainment', 'Wants', 25, False, None, None, None, True),
            ]

            bill_ids = {}
            for name, amt, cat, bucket, due, is_emi, total_months, months_paid, start_date, is_sub in bills:
                cur.execute(
                    """
                    INSERT INTO recurring_bills (user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription, is_emi, emi_total_months, emi_months_paid, start_date)
                    VALUES (%s, %s, %s, 'monthly', %s, %s, %s, %s, TRUE, %s, %s, %s, %s, %s)
                    RETURNING id
                    """,
                    (user_id, name, amt, cat, bucket, due, amt, is_sub, is_emi, total_months, months_paid, start_date)
                )
                bill_ids[name] = cur.fetchone()[0]

            # ── SUBSCRIPTION INSIGHTS ─────────────────────────────────────
            logger.info("Seeding subscription insights...")
            sub_insights = [
                (bill_ids['Netflix'], 'active', None, 649.00, 7788.00, 0.00, 'Regular usage, family plan'),
                (bill_ids['Amazon Prime'], 'active', None, 299.00, 3588.00, 0.00, 'Used for shopping & Prime Video'),
                (bill_ids['Spotify Premium'], 'active', None, 119.00, 1428.00, 0.00, 'Daily usage'),
                (bill_ids['YouTube Premium'], 'active', datetime.date(today.year, today.month, 1) - datetime.timedelta(days=45), 139.00, 1668.00, 0.00, 'Ad-free YouTube on all devices'),
                (bill_ids['Gym Membership'], 'underutilized', datetime.date.today() - datetime.timedelta(days=21), 1500.00, 18000.00, 6000.00, 'Only visited 4 times last month'),
            ]
            for b_id, status, last_used, monthly, annual, savings, notes in sub_insights:
                cur.execute(
                    """
                    INSERT INTO subscription_insights (user_id, bill_id, status, last_used_date, monthly_cost, annual_cost, savings_opportunity, notes)
                    VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                    """,
                    (user_id, b_id, status, last_used, monthly, annual, savings, notes)
                )

            # ── VEHICLES ──────────────────────────────────────────────────
            logger.info("Seeding vehicles...")
            cur.execute(
                """
                INSERT INTO vehicles (user_id, make_model, purchase_cost, insurance_renewal_date)
                VALUES (%s, %s, %s, %s)
                """,
                (user_id, 'Honda City 2023', 1250000.00, datetime.date(today.year, 9, 15))
            )

            # ── INSURANCES ─────────────────────────────────────────────────
            logger.info("Seeding insurances...")
            health_renewal = datetime.date(today.year, 4, 1)
            term_renewal = datetime.date(today.year, 8, 15)
            vehicle_renewal = datetime.date(today.year, 9, 15)
            cur.execute(
                """
                INSERT INTO insurances (user_id, type, provider, policy_name, premium_amount, premium_frequency, coverage_amount, renewal_date)
                VALUES
                (%s, 'health', 'Star Health', 'Star Family Health Optima', 18000.00, 'yearly', 500000.00, %s),
                (%s, 'term', 'HDFC Life', 'HDFC Life Click 2 Protect', 24000.00, 'yearly', 10000000.00, %s),
                (%s, 'vehicle', 'ICICI Lombard', 'ICICI Car Insurance', 14500.00, 'yearly', 1250000.00, %s)
                """,
                (user_id, health_renewal, user_id, term_renewal, user_id, vehicle_renewal)
            )

            # ── HISTORICAL EXPENSES (2 years) ─────────────────────────────
            logger.info("Seeding 2 years of historical expenses...")
            expense_templates = [
                ('Zomato Delivery', 450.00, 'Food', 'Wants', 'fastfood', 0.8, 1.5),
                ('Swiggy Instamart', 650.00, 'Groceries', 'Needs', 'shopping_cart', 0.7, 1.6),
                ('Uber Ride', 320.00, 'Transport', 'Needs', 'directions_car', 0.5, 2.0),
                ('Metro Card Top-up', 500.00, 'Transport', 'Needs', 'subway', 0.9, 1.1),
                ('Blinkit Groceries', 1100.00, 'Groceries', 'Needs', 'shopping_cart', 0.7, 1.4),
                ('Amazon Shopping', 1800.00, 'Shopping', 'Wants', 'shopping_bag', 0.4, 2.5),
                ('BigBasket', 1400.00, 'Groceries', 'Needs', 'shopping_cart', 0.8, 1.3),
                ('PVR Cinemas', 700.00, 'Entertainment', 'Wants', 'movie', 0.7, 1.5),
                ('Pharmacy', 450.00, 'Health', 'Needs', 'local_pharmacy', 0.3, 2.0),
                ('Starbucks Coffee', 350.00, 'Food', 'Wants', 'local_cafe', 0.6, 1.5),
                ('D-Mart Monthly', 3200.00, 'Groceries', 'Needs', 'shopping_cart', 0.8, 1.3),
                ('Petrol', 1800.00, 'Transport', 'Needs', 'local_gas_station', 0.7, 1.4),
                ('Dining Out', 1200.00, 'Food', 'Wants', 'restaurant', 0.5, 2.0),
                ('Railway Ticket', 850.00, 'Travel', 'Wants', 'train', 0.3, 1.0),
                ('Electricity Bill', 2100.00, 'Utilities', 'Needs', 'bolt', 0.8, 1.3),
                ('Mobile Recharge', 499.00, 'Utilities', 'Needs', 'smartphone', 0.9, 1.0),
                ('Myntra Fashion', 2200.00, 'Shopping', 'Wants', 'checkroom', 0.3, 2.0),
                ('Zepto Quick Delivery', 550.00, 'Groceries', 'Needs', 'shopping_cart', 0.6, 1.8),
                ('Domino\'s Pizza', 600.00, 'Food', 'Wants', 'local_pizza', 0.7, 1.5),
                ('UberEats', 380.00, 'Food', 'Wants', 'fastfood', 0.5, 2.0),
                ('BookMyShow', 500.00, 'Entertainment', 'Wants', 'theater_comedy', 0.5, 2.0),
                ('Medibuddy', 300.00, 'Health', 'Needs', 'health_and_safety', 0.4, 1.5),
                ('Rapido Bike', 80.00, 'Transport', 'Needs', 'pedal_bike', 0.5, 3.0),
                ('Urban Company', 600.00, 'Household', 'Needs', 'cleaning_services', 0.4, 1.5),
                ('Cult.fit Pass', 2000.00, 'Health', 'Wants', 'fitness_center', 0.8, 1.2),
            ]

            today = datetime.date.today()
            start_date = today - datetime.timedelta(days=730)

            # Also seed monthly fixed bill payments as expenses on their due days
            bill_expense_map = [
                ('Rent', 18000.00, 'Housing', 'Needs', 1),
                ('Electricity Bill', 2400.00, 'Utilities', 'Needs', 5),
                ('Water Bill', 850.00, 'Utilities', 'Needs', 8),
                ('Internet', 1499.00, 'Utilities', 'Needs', 10),
                ('Car Loan EMI', 12500.00, 'Debt', 'Needs', 3),
                ('Personal Loan EMI', 5000.00, 'Debt', 'Needs', 7),
                ('SIP - Index Fund', 15000.00, 'Investment', 'Savings', 12),
                ('SIP - ELSS', 5000.00, 'Investment', 'Savings', 15),
                ('PPF Contribution', 5000.00, 'Investment', 'Savings', 28),
                ('Netflix', 649.00, 'Entertainment', 'Wants', 16),
                ('Amazon Prime', 299.00, 'Entertainment', 'Wants', 20),
                ('Spotify Premium', 119.00, 'Entertainment', 'Wants', 22),
                ('Gym Membership', 1500.00, 'Health', 'Wants', 18),
                ('YouTube Premium', 139.00, 'Entertainment', 'Wants', 25),
            ]

            last_4_days = set()
            for d_idx in range(730):
                d = start_date + datetime.timedelta(days=d_idx)

                # Skip future dates
                if d > today:
                    break

                # Weekend boost: more spending on Fri-Sun
                is_weekend = d.weekday() >= 4
                is_friday = d.weekday() == 4
                is_saturday = d.weekday() == 5
                is_sunday = d.weekday() == 6

                # Random daily expenses
                if is_saturday:
                    num_expenses = random.randint(2, 5)
                elif is_sunday:
                    num_expenses = random.randint(1, 4)
                elif is_friday:
                    num_expenses = random.randint(1, 3)
                else:
                    num_expenses = random.randint(0, 2)

                for _ in range(num_expenses):
                    if random.random() < 0.25 and d_idx < 700:
                        continue
                    name, base_amt, cat, bucket, icon, low_mult, high_mult = random.choice(expense_templates)
                    final_amt = base_amt * random.uniform(low_mult, high_mult)
                    if is_weekend and bucket == 'Wants':
                        final_amt *= random.uniform(1.0, 1.3)
                    cur.execute(
                        """
                        INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                        VALUES (%s, %s, %s, %s, %s, %s, %s)
                        """,
                        (user_id, name, round(final_amt, 2), cat, bucket, icon, d)
                    )
                    last_4_days.add(d_idx)

                # Monthly bills on their due dates
                for bname, bamt, bcat, bbucket, bdue in bill_expense_map:
                    if d.day == bdue:
                        cur.execute(
                            """
                            INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                            VALUES (%s, %s, %s, %s, %s, 'receipt', %s)
                            """,
                            (user_id, bname, bamt, bcat, bbucket, d)
                        )

            # Ensure at least some expenses in last 4 days (for recent expense display)
            for offset in range(4):
                d = today - datetime.timedelta(days=offset)
                if d not in last_4_days:
                    cur.execute(
                        """
                        INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                        VALUES (%s, %s, %s, %s, %s, %s, %s)
                        """,
                        (user_id, 'D-Mart Groceries', round(random.uniform(800, 2500), 2), 'Groceries', 'Needs', 'shopping_cart', d)
                    )

            # ── WISHLIST ITEMS ────────────────────────────────────────────
            logger.info("Seeding wishlist items...")
            bought_date = today - datetime.timedelta(days=120)
            discarded_date = today - datetime.timedelta(days=90)
            unlocked_date = today - datetime.timedelta(days=45)
            locked_date = today - datetime.timedelta(days=7)
            early_bought = today - datetime.timedelta(days=30)

            wishlist = [
                ('iPhone 15 Pro', 120000.00, bought_date, bought_date + datetime.timedelta(days=60), 'bought', 60),
                ('AirPods Max', 45000.00, discarded_date, discarded_date + datetime.timedelta(days=30), 'discarded', 30),
                ('Dyson V15 Vacuum', 55000.00, unlocked_date, unlocked_date + datetime.timedelta(days=45), 'unlocked', 45),
                ('Herman Miller Aeron Chair', 110000.00, locked_date, locked_date + datetime.timedelta(days=60), 'locked', 60),
                ('PlayStation 5', 45000.00, today - datetime.timedelta(days=60), today + datetime.timedelta(days=30), 'locked', 90),
                ('Samsung 65" OLED TV', 160000.00, early_bought, early_bought + datetime.timedelta(days=90), 'bought_early', 90),
            ]

            for wname, wamt, added, unlock, wstatus, duration in wishlist:
                cur.execute(
                    """
                    INSERT INTO wishlist_items (user_id, name, amount, added_date, unlock_date, status, lock_duration_days)
                    VALUES (%s, %s, %s, %s, %s, %s, %s)
                    """,
                    (user_id, wname, wamt, added, unlock, wstatus, duration)
                )

            # ── NET WORTH SNAPSHOTS (monthly for 2 years) ─────────────────
            logger.info("Seeding net worth snapshots...")
            base_nw = 1800000.00
            for i in range(24):
                snap_date = today.replace(day=1) - datetime.timedelta(days=30 * i)
                month_progress = i / 24.0
                nw = base_nw + (380000 * month_progress) + random.uniform(-30000, 50000)
                cur.execute(
                    """
                    INSERT INTO net_worth_snapshots (user_id, snapshot_date, net_worth)
                    VALUES (%s, %s, %s)
                    ON CONFLICT (user_id, snapshot_date) DO NOTHING
                    """,
                    (user_id, snap_date, round(nw, 2))
                )

            # ── SIMULATIONS ───────────────────────────────────────────────
            logger.info("Seeding simulation history...")
            sims = [
                ("What if I buy a car?", "purchase", {"asset_cost": 800000, "down_payment": 200000, "interest_rate": 9.5, "tenure_months": 60}, {"monthly_impact": -14500, "new_monthly_bills": 48950, "affordable": True}),
                ("Increase SIP to 25k", "sip_increase", {"current_sip": 15000, "new_sip": 25000, "return_rate": 12}, {"monthly_impact": -10000, "projected_5y": 2040000, "affordable": True}),
                ("Pay off personal loan early", "loan_repayment", {"loan_amount": 25000, "lump_sum": 25000}, {"monthly_impact": 5000, "savings_on_interest": 4200, "one_time_cost": 25000, "affordable": True}),
                ("Rental property investment", "purchase", {"asset_cost": 5000000, "down_payment": 1000000, "interest_rate": 8.5, "tenure_months": 240}, {"monthly_impact": -28500, "expected_rental_income": 22000, "net_monthly_impact": -6500, "affordable": False}),
            ]
            for sname, stype, params, results in sims:
                cur.execute(
                    """
                    INSERT INTO simulations (user_id, scenario_name, scenario_type, input_params, results)
                    VALUES (%s, %s, %s, %s, %s)
                    """,
                    (user_id, sname, stype, json.dumps(params), json.dumps(results))
                )

            # ── NUDGES ────────────────────────────────────────────────────
            logger.info("Seeding nudges...")
            nudges = [
                ("spending", "warning", "Entertainment spending is up 40% this month", "You spent ₹4,200 on entertainment this month vs. your usual ₹3,000. Consider cutting back on dining out.", "Review Expenses", "/budget"),
                ("savings", "info", "You haven't invested your monthly surplus yet", "You have ~₹12,000 uninvested this month. Consider adding to your SIP or emergency fund.", "Invest Now", "/advisor"),
                ("subscription", "warning", "Gym membership seems underused", "You've visited the gym only 4 times this month but paid ₹1,500. That's ₹375 per visit!", "Review Subscriptions", "/bills"),
                ("debt", "success", "Personal loan is almost paid off!", "Only ₹25,000 remaining on your personal loan. At ₹5,000/month you'll be debt-free in 5 months.", "Track Progress", "/liabilities"),
                ("emergency", "info", "Emergency fund is at 6.2 months of expenses", "Great job! Your liquid savings of ₹2,27,000 cover 6.2 months of essential expenses.", "View Details", "/discipline"),
            ]
            for cat, severity, title, msg, action_label, action_link in nudges:
                nudge_date = today - datetime.timedelta(days=random.randint(0, 14))
                cur.execute(
                    """
                    INSERT INTO nudges (user_id, category, severity, title, message, action_label, action_link, created_at)
                    VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                    """,
                    (user_id, cat, severity, title, msg, action_label, action_link, nudge_date)
                )

            # ── COLLECTIONS ─────────────────────────────────────────────────────
            logger.info("Seeding collections...")

            collections_data = [
                {
                    "label": "Badminton Court - July 2026",
                    "description": "Monthly badminton court subscription collection from friends",
                    "members": [
                        ("Rajesh", 5000.00),
                        ("Priya", 2500.00),
                        ("Arun", 5000.00),
                        ("Ananya", 5000.00),
                        ("Vikram", 2500.00),
                        ("Sneha", 5000.00),
                        ("Rahul", 5000.00),
                        ("Divya", 5000.00),
                        ("Karthik", 5000.00),
                        ("Meera", 5000.00),
                        ("Ajay", 5000.00),
                        ("Neha", 5000.00),
                        ("Rohit", 2500.00),
                        ("Pooja", 5000.00),
                    ],
                },
                {
                    "label": "Office Birthday Pool - Q3 2026",
                    "description": "Office team collection for birthday celebrations",
                    "members": [
                        ("Amit", 1000.00),
                        ("Sara", 1000.00),
                        ("Ravi", 1000.00),
                        ("Leena", 1000.00),
                        ("Gopal", 1000.00),
                        ("Nisha", 1000.00),
                    ],
                },
                {
                    "label": "Weekend Getaway - Lonavala",
                    "description": "Trip collection for 3-day weekend trip",
                    "members": [
                        ("Varun", 8000.00),
                        ("Kavya", 8000.00),
                        ("Manish", 8000.00),
                        ("Anita", 4000.00),
                        ("Deepak", 4000.00),
                    ],
                },
            ]

            # Track what we insert so we can create matching expense entries
            collection_member_map = {}  # member_name -> (collection_label, amount)

            for c in collections_data:
                total_exp = sum(m[1] for m in c["members"])
                cur.execute("""
                    INSERT INTO collections (user_id, label, description, total_expected, total_collected, status)
                    VALUES (%s, %s, %s, %s, 0, 'active')
                    RETURNING id
                """, (user_id, c["label"], c["description"], total_exp))
                collection_id = cur.fetchone()[0]

                for m_name, m_amount in c["members"]:
                    # Mark some as already paid for realistic look
                    is_paid = random.random() < 0.35  # 35% chance already paid
                    paid_amt = m_amount if is_paid else 0.0
                    paid_date = today - datetime.timedelta(days=random.randint(1, 5)) if is_paid else None
                    m_status = "paid" if is_paid else "pending"

                    cur.execute("""
                        INSERT INTO collection_members (collection_id, name, expected_amount, paid_amount, paid_date, status)
                        VALUES (%s, %s, %s, %s, %s, %s)
                        RETURNING id
                    """, (collection_id, m_name, m_amount, paid_amt, paid_date, m_status))
                    member_id = cur.fetchone()[0]
                    collection_member_map[m_name] = (c["label"], m_amount, member_id)

                    # If paid, also add a matching expense entry (simulates a UPI / bank transfer received)
                    if is_paid:
                        expense_name = random.choice([
                            f"UPI-{m_name.upper()}",
                            f"Transfer from {m_name}",
                            f"Payment from {m_name}",
                            f"{m_name} - Badminton",
                        ])
                        cur.execute("""
                            INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                            VALUES (%s, %s, %s, %s, %s, %s, %s)
                        """, (user_id, expense_name, m_amount, "Transfer Received", "Wants", "payments", paid_date))

            # Update collection totals
            cur.execute("""
                UPDATE collections SET
                    total_collected = (SELECT COALESCE(SUM(paid_amount), 0) FROM collection_members WHERE collection_id = collections.id),
                    total_expected = (SELECT COALESCE(SUM(expected_amount), 0) FROM collection_members WHERE collection_id = collections.id)
            """)

            logger.info(f"  Seeded {sum(len(c['members']) for c in collections_data)} collection members across {len(collections_data)} collections")

            db.commit()
            logger.info(f"Database seeded successfully!")
            logger.info(f"  Email:    {email}")
            logger.info(f"  Password: {password}")
            logger.info(f"  User ID:  {user_id}")
            logger.info(f"  2 years of expenses, all app features populated.")

        except Exception as e:
            db.rollback()
            logger.error(f"Error seeding database: {e}")
            raise e

if __name__ == "__main__":
    seed_data()
