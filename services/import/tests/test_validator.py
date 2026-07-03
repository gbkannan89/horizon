from datetime import date
from decimal import Decimal

from app.models.schemas import ParsedTransaction
from app.services.validator import validate_batch, validate_transaction


class TestValidator:
    def test_valid_transaction(self):
        txn = ParsedTransaction(
            description="Salary",
            amount=Decimal("50000"),
            currency="INR",
            event_date=date(2024, 1, 15),
        )
        errors = validate_transaction(txn)
        assert len(errors) == 0

    def test_missing_description(self):
        txn = ParsedTransaction(
            description="",
            amount=Decimal("100"),
            currency="INR",
            event_date=date(2024, 1, 15),
        )
        errors = validate_transaction(txn)
        assert any("description" in e for e in errors)

    def test_zero_amount(self):
        txn = ParsedTransaction(
            description="Test",
            amount=Decimal("0"),
            currency="INR",
            event_date=date(2024, 1, 15),
        )
        errors = validate_transaction(txn)
        assert any("amount" in e for e in errors)

    def test_negative_amount(self):
        txn = ParsedTransaction(
            description="Test",
            amount=Decimal("-100"),
            currency="INR",
            event_date=date(2024, 1, 15),
        )
        errors = validate_transaction(txn)
        assert any("amount" in e for e in errors)

    def test_invalid_currency(self):
        txn = ParsedTransaction(
            description="Test",
            amount=Decimal("100"),
            currency="INRR",
            event_date=date(2024, 1, 15),
        )
        errors = validate_transaction(txn)
        assert any("currency" in e for e in errors)

    def test_future_date(self):
        txn = ParsedTransaction(
            description="Test",
            amount=Decimal("100"),
            currency="INR",
            event_date=date(2100, 1, 1),
        )
        errors = validate_transaction(txn)
        assert any("future" in e for e in errors)

    def test_past_date(self):
        txn = ParsedTransaction(
            description="Test",
            amount=Decimal("100"),
            currency="INR",
            event_date=date(1969, 12, 31),
        )
        errors = validate_transaction(txn)
        assert any("past" in e for e in errors)

    def test_effective_date_too_far(self):
        txn = ParsedTransaction(
            description="Test",
            amount=Decimal("100"),
            currency="INR",
            event_date=date(2024, 1, 1),
            effective_date=date(2026, 1, 1),
        )
        errors = validate_transaction(txn)
        assert any("effective_date" in e for e in errors)

    def test_batch_valid(self):
        transactions = [
            ParsedTransaction(description="A", amount=Decimal("100"), currency="INR", event_date=date(2024, 1, 1)),
            ParsedTransaction(description="B", amount=Decimal("200"), currency="INR", event_date=date(2024, 1, 2)),
        ]
        valid, errors = validate_batch(transactions)
        assert len(valid) == 2
        assert len(errors) == 0

    def test_batch_with_errors(self):
        transactions = [
            ParsedTransaction(description="", amount=Decimal("100"), currency="INR", event_date=date(2024, 1, 1)),
            ParsedTransaction(description="B", amount=Decimal("-5"), currency="INR", event_date=date(2024, 1, 2)),
        ]
        valid, errors = validate_batch(transactions)
        assert len(valid) == 0
        assert len(errors) >= 2

    def test_batch_partial_errors(self):
        transactions = [
            ParsedTransaction(description="Good", amount=Decimal("100"), currency="INR", event_date=date(2024, 1, 1)),
            ParsedTransaction(description="", amount=Decimal("200"), currency="INR", event_date=date(2024, 1, 2)),
        ]
        valid, errors = validate_batch(transactions)
        assert len(valid) == 1
        assert len(errors) >= 1
