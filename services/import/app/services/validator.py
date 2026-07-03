from datetime import date, timedelta

from app.models.schemas import ParsedTransaction


class ValidationError(Exception):
    def __init__(self, message: str, row: int | None = None):
        self.row = row
        super().__init__(f"Row {row}: {message}" if row else message)


def validate_transaction(txn: ParsedTransaction, row: int | None = None) -> list[str]:
    errors: list[str] = []

    if not txn.description or not txn.description.strip():
        errors.append("description is required")

    if txn.amount is None or txn.amount <= 0:
        errors.append(f"amount must be positive, got {txn.amount}")

    if not txn.currency or len(txn.currency) != 3:
        errors.append(f"currency must be 3-letter ISO code, got '{txn.currency}'")

    if txn.event_date is None:
        errors.append("event_date is required")

    if txn.event_date and txn.event_date > date.today() + timedelta(days=1):
        errors.append(f"event_date {txn.event_date} is in the future")

    if txn.event_date and txn.event_date < date(1970, 1, 1):
        errors.append(f"event_date {txn.event_date} is too far in the past")

    if txn.event_date and txn.effective_date:
        diff = abs((txn.effective_date - txn.event_date).days)
        if diff > 365:
            errors.append(
                f"effective_date {txn.effective_date} is more than 1 year from "
                f"event_date {txn.event_date}"
            )

    return errors


def validate_batch(transactions: list[ParsedTransaction]) -> tuple[list[ParsedTransaction], list[tuple[int, str]]]:
    valid: list[ParsedTransaction] = []
    errors: list[tuple[int, str]] = []

    for i, txn in enumerate(transactions):
        txn_errors = validate_transaction(txn, row=i + 1)
        if txn_errors:
            for err in txn_errors:
                errors.append((i + 1, err))
        else:
            valid.append(txn)

    return valid, errors
