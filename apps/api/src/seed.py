import logging
import bcrypt
import datetime
import random
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
        # Check if user already exists
        cur.execute("SELECT id FROM users WHERE email = %s", (email,))
        if cur.fetchone():
            logger.info(f"User {email} already exists. Deleting to re-seed...")
            cur.execute("DELETE FROM users WHERE email = %s", (email,))
            
        try:
            logger.info("Seeding household...")
            cur.execute(
                "INSERT INTO households (name) VALUES (%s) RETURNING id",
                (f"{name}'s Household",)
            )
            household_id = cur.fetchone()[0]
            
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
            
            logger.info("Seeding income (last 3 months)...")
            salary_amount = 75000.00
            for i in range(3):
                inc_date = today.replace(day=1) - datetime.timedelta(days=30*i)
                cur.execute(
                    """
                    INSERT INTO incomes (user_id, type, amount, frequency, created_at)
                    VALUES (%s, 'salary', %s, 'monthly', %s)
                    """,
                    (user_id, salary_amount, inc_date)
                )
                
            logger.info("Seeding assets...")
            cur.execute(
                """
                INSERT INTO assets (user_id, type, name, amount, interest_rate)
                VALUES 
                (%s, 'bank', 'HDFC Savings Account', 145000.00, 3.5),
                (%s, 'fd', 'SBI Fixed Deposit', 200000.00, 7.1),
                (%s, 'physical', 'Gold Coins', 50000.00, 0.0)
                """,
                (user_id, user_id, user_id)
            )
            
            logger.info("Seeding liabilities...")
            cur.execute(
                """
                INSERT INTO liabilities (user_id, type, name, outstanding, emi, interest_rate)
                VALUES 
                (%s, 'car_loan', 'HDFC Car Loan', 450000.00, 12500.00, 8.5),
                (%s, 'personal_loan', 'Navi Loan', 50000.00, 5000.00, 14.0)
                """,
                (user_id, user_id)
            )
            
            logger.info("Seeding goals...")
            cur.execute(
                """
                INSERT INTO goals (user_id, name, target_amount, current_amount, status, color)
                VALUES 
                (%s, 'Emergency Fund', 100000.00, 42000.00, 'On Track', '#059669'),
                (%s, 'Europe Trip', 300000.00, 54000.00, 'Behind', '#DC2626'),
                (%s, 'New Laptop', 150000.00, 102000.00, 'Ahead', '#1E3A8A')
                """,
                (user_id, user_id, user_id)
            )
            
            logger.info("Seeding recurring bills...")
            bills = [
                ('Rent', 18000.00, 'Housing', 'Needs', 1),
                ('Electricity', 2400.00, 'Utilities', 'Needs', 5),
                ('Car Loan EMI', 12500.00, 'Debt', 'Needs', 10),
                ('SIP - Index Fund', 15000.00, 'Investment', 'Savings', 12),
                ('Netflix', 649.00, 'Entertainment', 'Wants', 15),
                ('Gym Membership', 1500.00, 'Health', 'Wants', 20)
            ]
            
            for name, amt, cat, bucket, due in bills:
                cur.execute(
                    """
                    INSERT INTO recurring_bills (user_id, name, amount, frequency, category, bucket, due_day, monthly_equivalent, is_active, is_subscription)
                    VALUES (%s, %s, %s, 'monthly', %s, %s, %s, %s, TRUE, FALSE)
                    """,
                    (user_id, name, amt, cat, bucket, due, amt)
                )
            
            logger.info("Seeding historical expenses (last 90 days)...")
            expense_templates = [
                ('Zomato Delivery', 450.00, 'Food', 'Wants', 'fastfood'),
                ('Swiggy Instamart', 800.00, 'Groceries', 'Needs', 'shopping_cart'),
                ('Uber Ride', 350.00, 'Transport', 'Needs', 'directions_car'),
                ('Blinkit Groceries', 1200.00, 'Groceries', 'Needs', 'shopping_cart'),
                ('Amazon Shopping', 2500.00, 'Shopping', 'Wants', 'shopping_bag'),
                ('PVR Cinemas', 900.00, 'Entertainment', 'Wants', 'movie'),
                ('Pharmacy', 540.00, 'Health', 'Needs', 'local_pharmacy'),
                ('Starbucks Coffee', 350.00, 'Food', 'Wants', 'local_cafe')
            ]
            
            # Generate past 90 days of random variables expenses
            for i in range(90):
                d = today - datetime.timedelta(days=i)
                # Random 0 to 3 expenses per day
                num_expenses = random.randint(0, 3)
                for _ in range(num_expenses):
                    name, amt, cat, bucket, icon = random.choice(expense_templates)
                    # Add some randomness to amount (+- 20%)
                    final_amt = amt * random.uniform(0.8, 1.2)
                    cur.execute(
                        """
                        INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                        VALUES (%s, %s, %s, %s, %s, %s, %s)
                        """,
                        (user_id, name, final_amt, cat, bucket, icon, d)
                    )
                    
                # Add the fixed monthly bills on their due dates
                for name, amt, cat, bucket, due in bills:
                    if d.day == due:
                        cur.execute(
                            """
                            INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                            VALUES (%s, %s, %s, %s, %s, 'receipt', %s)
                            """,
                            (user_id, name, amt, cat, bucket, d)
                        )
            
            logger.info("Seeding historical wishlist items...")
            
            bought_date = today - datetime.timedelta(days=60)
            discarded_date = today - datetime.timedelta(days=40)
            unlocked_date = today - datetime.timedelta(days=35)
            locked_date = today - datetime.timedelta(days=5)
            
            wishlist = [
                ('iPhone 15 Pro', 120000.00, bought_date, bought_date + datetime.timedelta(days=30), 'bought'),
                ('AirPods Max', 45000.00, discarded_date, discarded_date + datetime.timedelta(days=30), 'discarded'),
                ('Dyson Vacuum', 35000.00, unlocked_date, unlocked_date + datetime.timedelta(days=30), 'unlocked'),
                ('Herman Miller Chair', 110000.00, locked_date, locked_date + datetime.timedelta(days=30), 'locked')
            ]
            
            for name, amt, added, unlock, status in wishlist:
                cur.execute(
                    """
                    INSERT INTO wishlist_items (user_id, name, amount, added_date, unlock_date, status)
                    VALUES (%s, %s, %s, %s, %s, %s)
                    """,
                    (user_id, name, amt, added, unlock, status)
                )

            db.commit()
            logger.info(f"Database seeded successfully with 3 months historical data for user: {email}")
        except Exception as e:
            db.rollback()
            logger.error(f"Error seeding database: {e}")
            raise e

if __name__ == "__main__":
    seed_data()
