from datetime import date
from decimal import Decimal
from enum import StrEnum

from pydantic import BaseModel


class ImportStatus(StrEnum):
    PENDING = "pending"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"


class ParsedTransaction(BaseModel):
    event_type: str = ""
    amount: Decimal
    currency: str = "INR"
    event_date: date
    effective_date: date | None = None
    description: str
    source: str = ""
    destination: str = ""
    reference: str = ""
    notes: str = ""
    category: str = ""


class ImportResult(BaseModel):
    status: ImportStatus
    total: int = 0
    events_created: int = 0
    errors: list[str] = []
    transactions: list[ParsedTransaction] = []


class ColumnMapping(BaseModel):
    date_col: str = "date"
    description_col: str = "description"
    amount_col: str = "amount"
    type_col: str = ""
    reference_col: str = ""
    balance_col: str = ""
    debit_col: str = ""
    credit_col: str = ""
    skip_rows: int = 0
    date_format: str = "%Y-%m-%d"
    delimiter: str = ","
