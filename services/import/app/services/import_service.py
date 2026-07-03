from io import StringIO

from app.models.schemas import ImportResult, ImportStatus
from app.services.csv_parser import parse_csv
from app.services.dedup import deduplicate_within_batch
from app.services.validator import validate_batch


def import_csv(
    content: str | bytes,
    validate: bool = True,
    deduplicate: bool = True,
) -> ImportResult:
    try:
        if isinstance(content, bytes):
            content = content.decode("utf-8", errors="replace")

        file_io = StringIO(content)
        transactions = parse_csv(file_io)

        if not transactions:
            return ImportResult(
                status=ImportStatus.COMPLETED,
                total=0,
                events_created=0,
                transactions=[],
            )

        if validate:
            transactions, validation_errors = validate_batch(transactions)
        else:
            validation_errors = []

        if deduplicate:
            transactions = deduplicate_within_batch(transactions)

        result = ImportResult(
            status=ImportStatus.COMPLETED,
            total=len(transactions),
            events_created=len(transactions),
            transactions=transactions,
            errors=[msg for _, msg in validation_errors],
        )

        return result

    except ValueError as e:
        return ImportResult(
            status=ImportStatus.FAILED,
            total=0,
            events_created=0,
            errors=[str(e)],
        )
    except Exception as e:
        return ImportResult(
            status=ImportStatus.FAILED,
            total=0,
            events_created=0,
            errors=[f"Unexpected error: {e}"],
        )
