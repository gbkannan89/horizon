import csv
import re
from datetime import date, datetime
from decimal import Decimal
from io import StringIO
from typing import TextIO

from app.models.schemas import ColumnMapping, ParsedTransaction

BANK_FORMATS: dict[str, ColumnMapping] = {
    "hdfc": ColumnMapping(
        date_col="Date",
        description_col="Narration",
        amount_col="Withdrawal Amount (INR)",
        debit_col="Withdrawal Amount (INR)",
        credit_col="Deposit Amount (INR)",
        reference_col="Chq./Ref.No.",
        balance_col="Balance (INR)",
        skip_rows=1,
        date_format="%d-%m-%Y",
    ),
    "icici": ColumnMapping(
        date_col="Value Date",
        description_col="Transaction Description",
        reference_col="Reference No.",
        debit_col="Withdrawal Amount",
        credit_col="Deposit Amount",
        balance_col="Balance",
        skip_rows=1,
        date_format="%d-%m-%Y",
    ),
    "sbi": ColumnMapping(
        date_col="Txn Date",
        description_col="Description",
        debit_col="Debit",
        credit_col="Credit",
        balance_col="Balance",
        skip_rows=1,
        date_format="%d %b %Y",
    ),
    "generic": ColumnMapping(),
}


def build_generic_mapping(headers: list[str]) -> ColumnMapping:
    """Auto-detect column mapping from common column names."""
    hl = [h.strip().lower() for h in headers]
    mapping = ColumnMapping()

    for i, h in enumerate(hl):
        if h in ("date", "txn date", "value date", "transaction date", "posting date"):
            mapping.date_col = headers[i]
        elif h in ("description", "narration", "particulars", "transaction description",
                    "transaction remarks", "remarks", "details", "memo"):
            mapping.description_col = headers[i]
        elif h in ("amount", "transaction amount", "txn amount", "value"):
            mapping.amount_col = headers[i]
        elif h in ("debit", "debit amount", "debit_amt", "withdrawal", "withdrawal amount", "dr"):
            mapping.debit_col = headers[i]
        elif h in ("credit", "credit amount", "credit_amt", "deposit", "deposit amount", "cr"):
            mapping.credit_col = headers[i]
        elif h in ("reference", "reference no", "ref no", "ref", "chq/ref.no", "cheque no"):
            mapping.reference_col = headers[i]
        elif h in ("balance", "running balance", "closing balance", "balance (inr)"):
            mapping.balance_col = headers[i]

    return mapping


def detect_format(headers: list[str]) -> str:
    header_lower = [h.strip().lower() for h in headers]
    joined = " ".join(header_lower)

    if "value date" in header_lower and "transaction description" in joined:
        return "icici"
    if "txn date" in header_lower:
        return "sbi"
    if "narration" in header_lower and "withdrawal amount" in joined:
        return "hdfc"
    return "generic"


def parse_date(value: str, fmt: str) -> date:
    value = value.strip()
    for f in [fmt, "%Y-%m-%d", "%d-%m-%Y", "%d/%m/%Y", "%m/%d/%Y", "%d %b %Y", "%d-%b-%Y"]:
        try:
            return datetime.strptime(value, f).date()
        except ValueError:
            continue
    return date.today()


def parse_amount(value: str) -> Decimal:
    cleaned = re.sub(r"[^\d.\-]", "", value.strip().replace(",", ""))
    if not cleaned or cleaned in ("", "-", "."):
        return Decimal("0")
    return Decimal(cleaned)


def detect_event_type(description: str, amount: Decimal, credit_col: str | None = None) -> str:
    desc_lower = description.lower()
    if amount < 0:
        return "Expense"
    income_keywords = ["salary", "credit", "deposit", "refund", "interest", "dividend", "rent received"]
    if any(kw in desc_lower for kw in income_keywords):
        return "Income"
    return "Expense"


def parse_csv(file: TextIO, mapping: ColumnMapping | None = None) -> list[ParsedTransaction]:
    content = file.read()
    if isinstance(content, str) is False:
        content = content.decode("utf-8", errors="replace")
    reader = csv.DictReader(StringIO(content), delimiter=mapping.delimiter if mapping else ",")
    headers = reader.fieldnames or []

    if mapping is None:
        fmt = detect_format(headers)
        mapping = BANK_FORMATS[fmt]
        if fmt == "generic":
            mapping = build_generic_mapping(headers)

    transactions: list[ParsedTransaction] = []
    for row_num, row in enumerate(reader, start=2):
        try:
            txn = _parse_row(row, mapping)
            if txn is not None:
                transactions.append(txn)
        except Exception as e:
            raise ValueError(f"Row {row_num}: {e}")

    return transactions


def _parse_row(row: dict[str, str], mapping: ColumnMapping) -> ParsedTransaction | None:
    date_val = row.get(mapping.date_col, "").strip()
    desc = row.get(mapping.description_col, "").strip()

    if not date_val and not desc:
        return None

    event_date = parse_date(date_val, mapping.date_format)
    effective_date = event_date

    debit_raw = row.get(mapping.debit_col, "").strip() if mapping.debit_col else ""
    credit_raw = row.get(mapping.credit_col, "").strip() if mapping.credit_col else ""
    amount_raw = row.get(mapping.amount_col, "").strip() if mapping.amount_col else ""

    debit = parse_amount(debit_raw) if debit_raw else Decimal("0")
    credit = parse_amount(credit_raw) if credit_raw else Decimal("0")
    amount = parse_amount(amount_raw) if amount_raw else Decimal("0")

    if amount == 0:
        if debit > 0:
            amount = -debit
        elif credit > 0:
            amount = credit

    ref = row.get(mapping.reference_col, "").strip() if mapping.reference_col else ""
    event_type = detect_event_type(desc, amount)

    return ParsedTransaction(
        event_type=event_type,
        amount=abs(amount),
        currency="INR",
        event_date=event_date,
        effective_date=effective_date,
        description=desc,
        reference=ref,
    )
