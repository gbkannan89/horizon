import csv
import logging
import codecs
import os
import tempfile
from datetime import datetime
from fastapi import APIRouter, Depends, HTTPException, status, UploadFile, File
from ..database import get_db
from ..schemas import ExpenseOut, UserOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/transactions", tags=["transactions"])

ACCEPTED_EXTENSIONS = {'.csv', '.xlsx'}


def categorize_transaction(description: str, amount: float) -> tuple[str, str, str]:
    desc = description.upper()
    if any(x in desc for x in ["SWIGGY", "ZOMATO", "NETFLIX", "AMAZON PRIME", "PVR", "BOOKMYSHOW", "STARBUCKS"]):
        return ("Entertainment", "Wants", "movie_outlined")
    if any(x in desc for x in ["MYNTRA", "FLIPKART", "LIFESTYLE", "ZARA", "H&M", "SHOPPERS STOP"]):
        return ("Shopping", "Wants", "shopping_bag_outlined")
    if any(x in desc for x in ["RELIANCE FRESH", "DMART", "BIGBASKET", "GROFERS", "SPENCERS", "MORE"]):
        return ("Food", "Needs", "restaurant_outlined")
    if any(x in desc for x in ["UBER", "OLA", "RAPIDO", "IRCTC", "MAKEMYTRIP", "PETROL", "HPCL", "BPCL", "IOCL"]):
        return ("Transport", "Needs", "directions_car_outlined")
    if any(x in desc for x in ["BESCOM", "AIRTEL", "JIO", "VI", "ACT", "ELECTRICITY", "WATER"]):
        return ("Bills", "Needs", "receipt_long_outlined")
    if any(x in desc for x in ["ZERODHA", "GROWW", "UPSTOX", "LIC", "EPF", "PPF", "MUTUAL FUND", "SIP"]):
        return ("Savings", "Savings", "savings_outlined")
    return ("Misc", "Wants", "receipt")


def _parse_date(date_str: str) -> datetime.date:
    if not date_str:
        return datetime.now().date()
    date_str = date_str.strip()
    try:
        if '-' in date_str and len(date_str.split('-')[0]) == 4:
            return datetime.strptime(date_str, "%Y-%m-%d").date()
        elif '-' in date_str:
            return datetime.strptime(date_str, "%d-%m-%Y").date()
        elif '/' in date_str:
            return datetime.strptime(date_str, "%d/%m/%Y").date()
    except:
        pass
    return datetime.now().date()


def _parse_csv(file) -> list[dict]:
    lines = codecs.iterdecode(file, 'utf-8-sig')
    header_row_str = ""
    for line in lines:
        lower_line = line.lower()
        if 'date' in lower_line and ('description' in lower_line or 'narration' in lower_line or 'particulars' in lower_line):
            header_row_str = line
            break
    if not header_row_str:
        raise HTTPException(status_code=400, detail="Could not find valid CSV headers in the statement.")

    def line_generator():
        yield header_row_str
        for l in lines:
            yield l

    csvReader = csv.DictReader(line_generator())
    fieldnames = [f.strip().lower() for f in csvReader.fieldnames or []]

    date_col = next((f for f in fieldnames if 'date' in f), None)
    desc_col = next((f for f in fieldnames if 'description' in f or 'narration' in f or 'particulars' in f), None)
    amount_col = next((f for f in fieldnames if 'amount' in f and 'balance' not in f), None)
    debit_col = next((f for f in fieldnames if 'debit' in f or 'withdrawal' in f), None)

    if not date_col or not desc_col:
        raise HTTPException(status_code=400, detail=f"Could not identify required columns. Found: {fieldnames}")

    transactions = []
    for row in csvReader:
        normalized_row = {k.strip().lower(): v for k, v in row.items() if k}
        date_str = (normalized_row.get(date_col) or "").strip()
        desc = (normalized_row.get(desc_col) or "").strip()
        if not date_str or not desc:
            continue

        amount = 0.0
        if amount_col and normalized_row.get(amount_col):
            try: amount = float(normalized_row[amount_col].replace(',', '').strip())
            except: pass
        elif debit_col and normalized_row.get(debit_col):
            try: amount = float(normalized_row[debit_col].replace(',', '').strip())
            except: pass

        if amount <= 0:
            continue

        transactions.append({"date": _parse_date(date_str), "description": desc[:100], "amount": amount})

    return transactions


def _parse_excel(file) -> list[dict]:
    import openpyxl
    suffix = ".xlsx"
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        tmp.write(file.read())
        tmp_path = tmp.name

    try:
        wb = openpyxl.load_workbook(tmp_path, read_only=True, data_only=True)
        ws = wb.active

        all_rows = list(ws.iter_rows(values_only=True))
        header_idx = None
        fieldnames = []

        for idx, row in enumerate(all_rows):
            if row is None or all(c is None for c in row):
                continue
            row_str = ' '.join(str(c).strip().lower() if c is not None else '' for c in row)
            if 'date' in row_str and ('description' in row_str or 'narration' in row_str or 'particulars' in row_str):
                header_idx = idx
                fieldnames = [str(c).strip().lower() if c is not None else '' for c in row]
                break

        if header_idx is None:
            raise HTTPException(status_code=400, detail="Could not find valid headers in the Excel file.")

        date_col = next((i for i, f in enumerate(fieldnames) if 'date' in f), None)
        desc_col = next((i for i, f in enumerate(fieldnames) if any(d in f for d in ['description', 'narration', 'particulars'])), None)
        amount_col = next((i for i, f in enumerate(fieldnames) if 'amount' in f and 'balance' not in f), None)
        debit_col = next((i for i, f in enumerate(fieldnames) if any(d in f for d in ['debit', 'withdrawal'])), None)

        if date_col is None or desc_col is None:
            raise HTTPException(status_code=400, detail=f"Could not identify required columns. Found: {fieldnames}")

        transactions = []
        for row in all_rows[header_idx + 1:]:
            if row is None or all(c is None for c in row):
                continue
            date_val_raw = row[date_col] if date_col is not None else None
            if date_val_raw is not None:
                if isinstance(date_val_raw, datetime):
                    date_val = date_val_raw.date().isoformat()
                else:
                    date_val = str(date_val_raw).strip()
            else:
                date_val = ""
            desc_val = str(row[desc_col]).strip() if desc_col is not None and row[desc_col] is not None else ""
            if not date_val or not desc_val or date_val.lower() == 'none':
                continue

            amount = 0.0
            if amount_col is not None and row[amount_col] is not None:
                try: amount = float(row[amount_col])
                except: pass
            elif debit_col is not None and row[debit_col] is not None:
                try: amount = float(row[debit_col])
                except: pass

            if amount <= 0:
                continue
            transactions.append({"date": _parse_date(date_val), "description": desc_val[:100], "amount": amount})

        wb.close()
        return transactions
    finally:
        os.unlink(tmp_path)


@router.post("/upload", response_model=list[ExpenseOut])
def upload_statement(
    file: UploadFile = File(...),
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    ext = os.path.splitext(file.filename or "")[1].lower()
    if ext not in ACCEPTED_EXTENSIONS:
        raise HTTPException(status_code=400, detail=f"Unsupported file type '{ext}'. Accepted: CSV and XLSX. If you have an .xls file, please save it as .xlsx or .csv.")

    try:
        transactions = _parse_csv(file.file) if ext == '.csv' else _parse_excel(file.file)

        inserted_expenses = []
        with conn.cursor() as cur:
            for txn in transactions:
                category, bucket, icon = categorize_transaction(txn["description"], txn["amount"])
                cur.execute(
                    "SELECT 1 FROM expenses WHERE user_id = %s AND name = %s AND amount = %s AND date = %s",
                    (current_user.id, txn["description"], txn["amount"], txn["date"])
                )
                if cur.fetchone():
                    continue
                cur.execute(
                    """INSERT INTO expenses (user_id, name, amount, category, bucket, icon, date)
                       VALUES (%s, %s, %s, %s, %s, %s, %s)
                       RETURNING id, user_id, name, amount, category, bucket, icon, date, created_at""",
                    (current_user.id, txn["description"], txn["amount"], category, bucket, icon, txn["date"])
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
        logger.error(f"Error parsing file: {e}")
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        file.file.close()
