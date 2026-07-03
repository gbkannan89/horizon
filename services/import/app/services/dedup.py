import difflib
from decimal import Decimal

from app.models.schemas import ParsedTransaction

DUPLICATE_AMOUNT_TOLERANCE = Decimal("0.01")
DUPLICATE_DAYS_TOLERANCE = 3
DESCRIPTION_SIMILARITY_THRESHOLD = 0.85


def normalize_description(desc: str) -> str:
    """Normalize description for comparison."""
    return " ".join(desc.strip().lower().split())


def descriptions_similar(a: str, b: str) -> float:
    """Compute similarity ratio between two descriptions."""
    return difflib.SequenceMatcher(None, normalize_description(a), normalize_description(b)).ratio()


def is_duplicate(txn: ParsedTransaction, existing: ParsedTransaction) -> bool:
    """Check if two transactions are likely duplicates."""
    amount_match = abs(txn.amount - existing.amount) <= DUPLICATE_AMOUNT_TOLERANCE
    if not amount_match:
        return False

    date_diff = abs((txn.event_date - existing.event_date).days)
    date_match = date_diff <= DUPLICATE_DAYS_TOLERANCE
    if not date_match:
        return False

    desc_sim = descriptions_similar(txn.description, existing.description)
    return desc_sim >= DESCRIPTION_SIMILARITY_THRESHOLD


def find_duplicates(
    incoming: list[ParsedTransaction],
    existing: list[ParsedTransaction],
) -> tuple[list[ParsedTransaction], dict[int, list[int]]]:
    """
    Find duplicates between incoming and existing transactions.

    Returns:
        unique: Incoming transactions with no duplicate match
        duplicates: Mapping of incoming index to list of matching existing indices
    """
    unique: list[ParsedTransaction] = []
    duplicates: dict[int, list[int]] = {}

    for i, txn in enumerate(incoming):
        matches: list[int] = []
        for j, ext in enumerate(existing):
            if is_duplicate(txn, ext):
                matches.append(j)

        if matches:
            duplicates[i] = matches
        else:
            unique.append(txn)

    return unique, duplicates


def deduplicate_within_batch(
    transactions: list[ParsedTransaction],
) -> list[ParsedTransaction]:
    """Remove duplicates within a single import batch."""
    unique: list[ParsedTransaction] = []
    for txn in transactions:
        if not any(is_duplicate(txn, u) for u in unique):
            unique.append(txn)
    return unique
