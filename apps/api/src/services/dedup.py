import logging
from typing import List, Dict, Any, Tuple

logger = logging.getLogger(__name__)


def _normalize_name(name: str) -> str:
    return name.strip().lower()


def _name_words(name: str) -> set:
    words = set()
    for part in name.replace("-", " ").replace("_", " ").replace("/", " ").split():
        cleaned = part.strip(",.!?;:'\"()[]{}")
        if len(cleaned) > 2:
            words.add(cleaned.lower())
    return words


def is_duplicate(user_id: int, name: str, amount: float, txn_date, conn, tolerance: float = 0.01) -> Tuple[bool, str]:
    """
    Check if a transaction already exists using multiple strategies.
    Returns (is_duplicate, match_reason).
    """
    with conn.cursor() as cur:
        # Strategy 1: Exact match (user_id, name, amount, date)
        cur.execute(
            "SELECT id, name FROM expenses WHERE user_id = %s AND name = %s AND amount = %s AND date = %s",
            (user_id, name, amount, txn_date)
        )
        if cur.fetchone():
            return True, "exact match"

        # Strategy 2: Same amount + same date, case-insensitive name match
        cur.execute(
            "SELECT id, name FROM expenses WHERE user_id = %s AND amount = %s AND date = %s",
            (user_id, amount, txn_date)
        )
        existing = cur.fetchall()
        for eid, ename in existing:
            if _normalize_name(ename) == _normalize_name(name):
                return True, f"case-insensitive match: '{ename}'"

        # Strategy 3: Same amount + same date ±1 day, significant name overlap
        cur.execute(
            "SELECT id, name FROM expenses WHERE user_id = %s AND amount BETWEEN %s AND %s AND date BETWEEN %s AND %s",
            (user_id, amount - amount * tolerance, amount + amount * tolerance,
             txn_date, txn_date)
        )
        existing = cur.fetchall()
        new_words = _name_words(name)
        for eid, ename in existing:
            existing_words = _name_words(ename)
            if new_words and existing_words:
                overlap = new_words & existing_words
                # Significant overlap: at least 1 substantial word in common
                if overlap and max(len(w) for w in overlap) >= 3:
                    return True, f"word overlap: '{ename}' ↔ '{name}' (shared: {overlap})"

        # Strategy 4: Same amount + same date, one name contains the other
        cur.execute(
            "SELECT id, name FROM expenses WHERE user_id = %s AND amount = %s AND date = %s",
            (user_id, amount, txn_date)
        )
        existing = cur.fetchall()
        norm_name = _normalize_name(name)
        for eid, ename in existing:
            norm_existing = _normalize_name(ename)
            if len(norm_name) >= 4 and len(norm_existing) >= 4:
                if norm_name in norm_existing or norm_existing in norm_name:
                    return True, f"substring match: '{ename}'"

    return False, ""


def batch_check_duplicates(user_id: int, transactions: List[Dict[str, Any]], conn) -> Tuple[List[Dict], int, int]:
    """
    Check a batch of transactions for duplicates.
    Returns (unique_transactions, duplicate_count, fuzzy_match_count).
    """
    unique = []
    dup_count = 0
    fuzzy_count = 0

    for txn in transactions:
        is_dup, reason = is_duplicate(
            user_id, txn.get("description") or txn.get("name", ""),
            txn["amount"], txn["date"], conn
        )
        if is_dup:
            dup_count += 1
            if "exact" not in reason:
                fuzzy_count += 1
        else:
            unique.append(txn)

    return unique, dup_count, fuzzy_count
