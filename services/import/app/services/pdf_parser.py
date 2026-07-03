import re
from datetime import date, datetime
from decimal import Decimal
from pathlib import Path
from tempfile import NamedTemporaryFile

import pdfplumber

from app.models.schemas import ParsedTransaction
from app.services.csv_parser import (
    BANK_FORMATS,
    build_generic_mapping,
    detect_event_type,
    detect_format,
    parse_amount,
)

_DATA_HEADER_KEYWORDS = {"date", "narration", "description", "particulars", "transaction"}
_AMOUNT_KEYWORDS = {"amount", "debit", "credit", "withdrawal", "deposit"}


def _is_header_or_footer(text: str) -> bool:
    """Detect common header/footer lines in bank statements."""
    tl = text.strip().lower()
    if not tl:
        return True
    patterns = [
        r"^page\s+\d+",
        r"^(hdfc|icici|sbi|axis|yes\s+bank)",
        r"(registered office|corporate office|toll free|email|website)",
        r"^continuation\s+sheet",
        r"^statement\s+(of\s+)?account",
        r"^(account|customer)\s+(no|number|id)",
        r"^(summary|opening|closing|minimum|average)",
    ]
    if any(re.match(p, tl) for p in patterns):
        return True
    return False


def _is_data_header(text: str) -> bool:
    tl = text.strip().lower()
    words = set(tl.split())
    return bool(words & _DATA_HEADER_KEYWORDS) and bool(words & _AMOUNT_KEYWORDS)


_DATE_LINE_RE = re.compile(r"^\d{1,2}[-/]\d{1,2}[-/]\d{2,4}")


def _clean_cell(value) -> str:
    if value is None:
        return ""
    return str(value).strip()


def parse_pdf(file_content: bytes | str | Path) -> list[ParsedTransaction]:
    if isinstance(file_content, str):
        file_content = file_content.encode("utf-8")
    elif isinstance(file_content, Path):
        file_content = file_content.read_bytes()

    with NamedTemporaryFile(delete=False, suffix=".pdf") as tmp:
        tmp.write(file_content)
        tmp_path = tmp.name

    try:
        with pdfplumber.open(tmp_path) as pdf:
            all_text = ""
            tables = []

            for page in pdf.pages:
                page_text = page.extract_text() or ""
                all_text += page_text + "\n"

                page_tables = page.extract_tables()
                if page_tables:
                    tables.extend(page_tables)

            if not tables:
                return _fallback_text_parse(all_text)

            transactions: list[ParsedTransaction] = []
            for table in tables:
                if not table or len(table) < 2:
                    continue

                headers = [_clean_cell(c) for c in table[0]]
                if not any(h for h in headers):
                    continue

                fmt = detect_format(headers)
                if fmt == "generic":
                    mapping = build_generic_mapping(headers)
                else:
                    mapping = BANK_FORMATS[fmt]

                for row in table[1:]:
                    row_dict = {headers[i]: _clean_cell(v) for i, v in enumerate(row) if i < len(headers)}
                    if all(v == "" for v in row_dict.values()):
                        continue

                    txn = _parse_pdf_row(row_dict, mapping)
                    if txn is not None:
                        transactions.append(txn)

            return transactions

    finally:
        for _ in range(5):
            try:
                Path(tmp_path).unlink()
                break
            except PermissionError:
                import time
                time.sleep(0.1)


def _parse_pdf_row(row: dict[str, str], mapping) -> ParsedTransaction | None:
    date_val = row.get(mapping.date_col, "").strip()
    desc = row.get(mapping.description_col, "").strip()

    if not date_val and not desc:
        return None

    event_date = _parse_pdf_date(date_val, mapping.date_format)
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


DATE_PATTERNS = ["%d-%m-%Y", "%d/%m/%Y", "%Y-%m-%d", "%d %b %Y", "%d-%b-%Y", "%b %d, %Y", "%d %B %Y"]


def _parse_pdf_date(value: str, preferred_fmt: str) -> date:
    value = value.strip()
    for f in [preferred_fmt] + DATE_PATTERNS:
        try:
            return datetime.strptime(value, f).date()
        except ValueError:
            continue
    return date.today()


def _fallback_text_parse(text: str) -> list[ParsedTransaction]:
    """Fallback for PDFs where tables aren't detected — parse line by line."""
    transactions: list[ParsedTransaction] = []
    lines = text.split("\n")

    data_started = False
    for line in lines:
        tl = line.strip()
        if _is_header_or_footer(tl):
            continue
        if not tl:
            continue

        if not data_started:
            if _is_data_header(tl) or _DATE_LINE_RE.match(tl):
                data_started = True
                if _DATE_LINE_RE.match(tl):
                    txn = _parse_fallback_line(tl)
                    if txn is not None:
                        transactions.append(txn)
            continue

        txn = _parse_fallback_line(tl)
        if txn is not None:
            transactions.append(txn)

    return transactions


_FALLBACK_DATE_RE = re.compile(
    r"(\d{1,2})[-/](\d{1,2})[-/](\d{2,4})|"  # 01-01-2024 or 01/01/2024
    r"(\d{1,2})\s+([A-Za-z]{3,9})\s+(\d{2,4})"  # 01 Jan 2024 or 01 January 2024
)


def _parse_fallback_line(line: str) -> ParsedTransaction | None:
    parts = re.split(r"\s{2,}", line)
    if len(parts) < 2:
        parts = line.split()

    date_match = None
    amount_match = None
    desc_parts = []

    # Check if the line as a whole contains a date at the start
    dm = _FALLBACK_DATE_RE.search(line)
    if dm:
        if dm.group(1) and dm.group(2) and dm.group(3):
            date_str = f"{dm.group(1)}-{dm.group(2)}-{dm.group(3)}"
        else:
            date_str = f"{dm.group(4)}-{dm.group(5)}-{dm.group(6)}"
        date_match = _parse_pdf_date(date_str, "%d-%m-%Y")

    # Find amount - last numeric token
    for p in reversed(parts):
        dp = p.strip()
        if not dp:
            continue
        am = re.match(r"^[\-]?\s*[\d,]+\.?\d*$", dp.replace(",", "").strip())
        if am and amount_match is None:
            amount_match = parse_amount(dp)
            break

    # Everything else is description
    for p in parts:
        dp = p.strip()
        if not dp:
            continue
        dm = _FALLBACK_DATE_RE.match(dp)
        if dm and date_match:
            continue
        am = re.match(r"^[\-]?\s*[\d,]+\.?\d*$", dp.replace(",", "").strip())
        if am and amount_match:
            continue
        desc_parts.append(dp)

    if date_match is None or amount_match is None or not desc_parts:
        return None

    desc = " ".join(desc_parts)
    event_type = detect_event_type(desc, amount_match)

    return ParsedTransaction(
        event_type=event_type,
        amount=abs(amount_match) if amount_match else Decimal("0"),
        currency="INR",
        event_date=date_match,
        effective_date=date_match,
        description=desc[:200],
    )
