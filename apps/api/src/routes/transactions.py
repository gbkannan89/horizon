import csv
import logging
import codecs
from datetime import datetime
from fastapi import APIRouter, Depends, HTTPException, status, UploadFile, File
from ..database import get_db
from ..schemas import ExpenseOut, UserOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/transactions", tags=["transactions"])

def categorize_transaction(description: str, amount: float) -> tuple[str, str, str]:
    """
    Returns (category, bucket, icon) based on the transaction description.
    """
    desc = description.upper()
    
    # Wants
    if any(x in desc for x in ["SWIGGY", "ZOMATO", "NETFLIX", "AMAZON PRIME", "PVR", "BOOKMYSHOW", "STARBUCKS"]):
        return ("Entertainment", "Wants", "movie_outlined")
    if any(x in desc for x in ["MYNTRA", "FLIPKART", "LIFESTYLE", "ZARA", "H&M", "SHOPPERS STOP"]):
        return ("Shopping", "Wants", "shopping_bag_outlined")
        
    # Needs
    if any(x in desc for x in ["RELIANCE FRESH", "DMART", "BIGBASKET", "GROFERS", "SPENCERS", "MORE"]):
        return ("Food", "Needs", "restaurant_outlined")
    if any(x in desc for x in ["UBER", "OLA", "RAPIDO", "IRCTC", "MAKEMYTRIP", "PETROL", "HPCL", "BPCL", "IOCL"]):
        return ("Transport", "Needs", "directions_car_outlined")
    if any(x in desc for x in ["BESCOM", "AIRTEL", "JIO", "VI", "ACT", "ELECTRICITY", "WATER"]):
        return ("Bills", "Needs", "receipt_long_outlined")
        
    # Savings
    if any(x in desc for x in ["ZERODHA", "GROWW", "UPSTOX", "LIC", "EPF", "PPF", "MUTUAL FUND", "SIP"]):
        return ("Savings", "Savings", "savings_outlined")
        
    # Default fallback
    return ("Misc", "Wants", "receipt")

@router.post("/upload", response_model=list[ExpenseOut])
def upload_csv(
    file: UploadFile = File(...),
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    if not file.filename.endswith(".csv"):
        raise HTTPException(status_code=400, detail="Only CSV files are supported.")
        
    try:
        # Read lines and skip preamble
        lines = codecs.iterdecode(file.file, 'utf-8')
        
        # Fast forward until we find a row that looks like a header (contains date/narration)
        header_row_str = ""
        for line in lines:
            lower_line = line.lower()
            if 'date' in lower_line and ('description' in lower_line or 'narration' in lower_line or 'particulars' in lower_line):
                header_row_str = line
                break
                
        if not header_row_str:
            raise HTTPException(status_code=400, detail="Could not find valid CSV headers in the statement.")
            
        # Reconstruct the iterator with the header line and remaining lines
        def line_generator():
            yield header_row_str
            for l in lines:
                yield l
                
        csvReader = csv.DictReader(line_generator())
        
        # Standardize column names
        fieldnames = [f.strip().lower() for f in csvReader.fieldnames or []]
        
        # Find column mappings
        date_col = next((f for f in fieldnames if 'date' in f), None)
        desc_col = next((f for f in fieldnames if 'description' in f or 'narration' in f or 'particulars' in f), None)
        amount_col = next((f for f in fieldnames if 'amount' in f and 'balance' not in f), None)
        debit_col = next((f for f in fieldnames if 'debit' in f or 'withdrawal' in f), None)
        
        if not date_col or not desc_col:
            raise HTTPException(status_code=400, detail=f"Could not identify required columns. Found: {fieldnames}")
            
        inserted_expenses = []
        
        with conn.cursor() as cur:
            for row in csvReader:
                # Normalize row keys
                normalized_row = {k.strip().lower(): v for k, v in row.items() if k}
                
                date_str = normalized_row.get(date_col, "").strip()
                desc = normalized_row.get(desc_col, "").strip()
                
                if not date_str or not desc:
                    continue
                    
                amount = 0.0
                if amount_col and normalized_row.get(amount_col):
                    val = normalized_row[amount_col].replace(',', '').strip()
                    try:
                        amount = float(val)
                    except:
                        pass
                elif debit_col and normalized_row.get(debit_col):
                    val = normalized_row[debit_col].replace(',', '').strip()
                    try:
                        amount = float(val)
                    except:
                        pass
                        
                if amount <= 0:
                    continue # Skip credits/deposits for expenses
                    
                # Parse date - handle basic DD/MM/YYYY or YYYY-MM-DD
                try:
                    if '-' in date_str and len(date_str.split('-')[0]) == 4:
                        parsed_date = datetime.strptime(date_str, "%Y-%m-%d").date()
                    elif '-' in date_str:
                        parsed_date = datetime.strptime(date_str, "%d-%m-%Y").date()
                    elif '/' in date_str:
                        parsed_date = datetime.strptime(date_str, "%d/%m/%Y").date()
                    else:
                        parsed_date = datetime.now().date()
                except:
                    parsed_date = datetime.now().date()
                    
                category, bucket, icon = categorize_transaction(desc, amount)
                
                # Check for duplicates (same date, name, amount)
                cur.execute(
                    "SELECT 1 FROM expenses WHERE user_id = %s AND name = %s AND amount = %s AND date = %s",
                    (current_user.id, desc[:100], amount, parsed_date)
                )
                if cur.fetchone():
                    continue  # Skip duplicate
                
                cur.execute(
                    """
                    INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                    VALUES (%s, %s, %s, %s, %s, %s, %s)
                    RETURNING id, user_id, name, amount, category, bucket, icon, date, created_at
                    """,
                    (current_user.id, desc[:100], amount, category, bucket, icon, parsed_date)
                )
                row = cur.fetchone()
                
                inserted_expenses.append(ExpenseOut(
                    id=row[0], user_id=row[1], name=row[2], amount=float(row[3]),
                    category=row[4], bucket=row[5], icon=row[6], date=row[7], created_at=row[8]
                ))
                
            conn.commit()
            
        return inserted_expenses
            
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        logger.error(f"Error parsing CSV: {e}")
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        file.file.close()
