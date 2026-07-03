from decimal import Decimal
from pathlib import Path
from tempfile import NamedTemporaryFile

from openpyxl import Workbook

from app.services.excel_parser import parse_excel


def _make_excel(headers: list[str], rows: list[list]) -> bytes:
    wb = Workbook()
    ws = wb.active
    ws.title = "Sheet1"
    ws.append(headers)
    for row in rows:
        ws.append(row)
    with NamedTemporaryFile(delete=False, suffix=".xlsx") as tmp:
        wb.save(tmp.name)
        tmp_path = tmp.name
    content = Path(tmp_path).read_bytes()
    Path(tmp_path).unlink()
    return content


class TestExcelParser:
    def test_generic_format(self):
        content = _make_excel(
            ["date", "description", "amount"],
            [
                ["2024-01-15", "Salary credit", 50000],
                ["2024-01-16", "Groceries", -2500],
            ],
        )
        result = parse_excel(content)
        assert len(result) == 2
        assert result[0].amount == Decimal("50000")
        assert result[0].event_type == "Income"
        assert result[1].amount == Decimal("2500")
        assert result[1].event_type == "Expense"

    def test_hdfc_format(self):
        headers = ["Date", "Narration", "Chq./Ref.No.", "Value Dat",
                   "Withdrawal Amount (INR)", "Deposit Amount (INR)", "Balance (INR)"]
        content = _make_excel(headers,
            [
                ["01-01-2024", "SALARY", "REF001", "01-01-2024", None, 50000, 100000],
                ["02-01-2024", "AMAZON PAY", "REF002", "02-01-2024", 2500, None, 97500],
            ],
        )
        result = parse_excel(content)
        assert len(result) == 2
        assert result[0].amount == Decimal("50000")
        assert result[0].event_type == "Income"
        assert result[1].amount == Decimal("2500")

    def test_icici_format(self):
        headers = ["Value Date", "Transaction Date", "Transaction Description",
                   "Reference No.", "Withdrawal Amount", "Deposit Amount", "Balance"]
        content = _make_excel(headers,
            [
                ["15-01-2024", "14-01-2024", "NEFT CREDIT", "REF001", None, 50000, 100000],
                ["16-01-2024", "15-01-2024", "SWIGGY", "REF002", 1200, None, 98700],
            ],
        )
        assert len(parse_excel(content)) == 2

    def test_sbi_format(self):
        headers = ["Txn Date", "Value Date", "Description", "Cheque No.", "Debit", "Credit", "Balance"]
        content = _make_excel(headers,
            [
                ["01 Jan 2024", "01 Jan 2024", "SALARY", "", None, 75000, 175000],
                ["02 Jan 2024", "02 Jan 2024", "ZOMATO", "", 1500, None, 173500],
            ],
        )
        result = parse_excel(content)
        assert len(result) == 2

    def test_multiple_sheets(self):
        wb = Workbook()
        ws1 = wb.active
        ws1.title = "Cover"
        ws1.append(["Report"])
        ws2 = wb.create_sheet("Transactions")
        ws2.append(["date", "description", "amount"])
        ws2.append(["2024-01-15", "Income", 10000])

        with NamedTemporaryFile(delete=False, suffix=".xlsx") as tmp:
            wb.save(tmp.name)
            tmp_path = tmp.name
        content = Path(tmp_path).read_bytes()
        Path(tmp_path).unlink()

        result = parse_excel(content)
        assert len(result) == 1
        assert result[0].description == "Income"

    def test_empty_workbook(self):
        wb = Workbook()
        with NamedTemporaryFile(delete=False, suffix=".xlsx") as tmp:
            wb.save(tmp.name)
            tmp_path = tmp.name
        content = Path(tmp_path).read_bytes()
        Path(tmp_path).unlink()

        result = parse_excel(content)
        assert len(result) == 0

    def test_header_detection(self):
        content = _make_excel(
            ["Transaction Date", "Description", "Debit", "Credit"],
            [
                ["2024-01-15", "Sale", None, 5000],
                ["2024-01-16", "Purchase", 1500, None],
            ],
        )
        result = parse_excel(content)
        assert len(result) == 2
        assert result[0].amount == Decimal("5000")
        assert result[1].amount == Decimal("1500")

    def test_bytes_input(self):
        content = _make_excel(
            ["date", "description", "amount"],
            [["2024-01-15", "Test", 100]],
        )
        result = parse_excel(content)
        assert len(result) == 1

    def test_path_input(self):
        wb = Workbook()
        ws = wb.active
        ws.append(["date", "description", "amount"])
        ws.append(["2024-01-15", "Path test", 200])

        with NamedTemporaryFile(delete=False, suffix=".xlsx") as tmp:
            wb.save(tmp.name)
            tmp_path = tmp.name

        result = parse_excel(Path(tmp_path))
        assert len(result) == 1
        assert result[0].amount == Decimal("200")
        Path(tmp_path).unlink()

    def test_decimal_amounts(self):
        content = _make_excel(
            ["date", "description", "amount"],
            [["2024-01-15", "Refund", 1234.56]],
        )
        result = parse_excel(content)
        assert len(result) == 1
        assert result[0].amount == Decimal("1234.56")
