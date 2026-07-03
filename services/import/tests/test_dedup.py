from datetime import date
from decimal import Decimal

from app.models.schemas import ParsedTransaction
from app.services.dedup import (
    deduplicate_within_batch,
    descriptions_similar,
    find_duplicates,
    is_duplicate,
)


def _txn(desc: str, amt: str, day: int, month: int = 1) -> ParsedTransaction:
    return ParsedTransaction(
        description=desc, amount=Decimal(amt),
        currency="INR", event_date=date(2024, month, day),
    )


class TestDedup:
    def test_exact_duplicate(self):
        a = _txn("Salary", "50000", 1)
        b = _txn("Salary", "50000", 1)
        assert is_duplicate(a, b)

    def test_different_amount(self):
        a = _txn("Salary", "50000", 1)
        b = _txn("Salary", "60000", 1)
        assert not is_duplicate(a, b)

    def test_different_date(self):
        a = _txn("Salary", "50000", 1)
        b = _txn("Salary", "50000", 1, month=5)
        assert not is_duplicate(a, b)

    def test_similar_description(self):
        a = _txn("AMAZON PAY", "2500", 15)
        b = _txn("Amazon Pay", "2500", 15)
        assert is_duplicate(a, b)

    def test_description_similarity(self):
        assert descriptions_similar("salary credit", "SALARY CREDIT") > 0.9
        assert descriptions_similar("amazon pay", "amz prime") > 0.5

    def test_date_tolerance(self):
        a = ParsedTransaction(description="Rent", amount=Decimal("25000"), currency="INR", event_date=date(2024, 1, 1))
        b = ParsedTransaction(description="Rent", amount=Decimal("25000"), currency="INR", event_date=date(2024, 1, 3))
        assert is_duplicate(a, b)

    def test_beyond_date_tolerance(self):
        a = ParsedTransaction(description="Rent", amount=Decimal("25000"), currency="INR", event_date=date(2024, 1, 1))
        b = ParsedTransaction(description="Rent", amount=Decimal("25000"), currency="INR", event_date=date(2024, 1, 10))
        assert not is_duplicate(a, b)

    def test_find_duplicates(self):
        incoming = [_txn("Salary", "50000", 1), _txn("Rent", "25000", 5)]
        existing = [_txn("Salary", "50000", 1)]
        unique, dup_map = find_duplicates(incoming, existing)
        assert len(unique) == 1
        assert len(dup_map) == 1
        assert 0 in dup_map

    def test_deduplicate_within_batch(self):
        txns = [
            ParsedTransaction(description="A", amount=Decimal("100"), currency="INR", event_date=date(2024, 1, 1)),
            ParsedTransaction(description="A", amount=Decimal("100"), currency="INR", event_date=date(2024, 1, 1)),
            ParsedTransaction(description="B", amount=Decimal("200"), currency="INR", event_date=date(2024, 1, 2)),
        ]
        result = deduplicate_within_batch(txns)
        assert len(result) == 2

    def test_no_duplicates(self):
        txns = [
            ParsedTransaction(description="A", amount=Decimal("100"), currency="INR", event_date=date(2024, 1, 1)),
            ParsedTransaction(description="B", amount=Decimal("200"), currency="INR", event_date=date(2024, 1, 2)),
        ]
        result = deduplicate_within_batch(txns)
        assert len(result) == 2

    def test_integration_with_import(self):
        from app.services.import_service import import_csv

        csv_content = "date,description,amount\n2024-01-01,Duplicate,100\n2024-01-01,Duplicate,100\n"
        result = import_csv(csv_content, validate=True, deduplicate=True)
        assert result.total == 1
        assert result.events_created == 1
