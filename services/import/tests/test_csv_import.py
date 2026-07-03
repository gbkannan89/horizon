from datetime import date
from decimal import Decimal

from app.services.csv_parser import detect_format, parse_amount, parse_date
from app.services.import_service import import_csv


class TestCSVParser:
    def test_parse_generic_csv(self):
        csv_content = "date,description,amount\n2024-01-15,Salary credit,50000\n2024-01-16,Groceries,-2500\n"
        result = import_csv(csv_content)
        assert result.status.value == "completed"
        assert result.total == 2
        assert len(result.transactions) == 2

    def test_parse_hdfc_format(self):
        csv_content = (
            "Date,Narration,Chq./Ref.No.,Value Dat,Withdrawal Amount (INR),"
            "Deposit Amount (INR),Balance (INR)\n"
            "01-01-2024,SALARY,REF001,01-01-2024,,50000,100000\n"
            "02-01-2024,AMAZON PAY,REF002,02-01-2024,2500,,97500\n"
        )
        result = import_csv(csv_content)
        assert result.status.value == "completed"
        assert result.total == 2
        assert result.transactions[0].amount == Decimal("50000")
        assert result.transactions[0].event_type == "Income"
        assert result.transactions[1].amount == Decimal("2500")
        assert result.transactions[1].event_type == "Expense"

    def test_parse_icici_format(self):
        csv_content = (
            "Value Date,Transaction Date,Transaction Description,Reference No.,"
            "Withdrawal Amount,Deposit Amount,Balance\n"
            "15-01-2024,14-01-2024,NEFT CREDIT,REF001,,50000,100000\n"
            "16-01-2024,15-01-2024,SWIGGY,REF002,1200,,98700\n"
        )
        result = import_csv(csv_content)
        assert result.status.value == "completed"
        assert result.total == 2

    def test_parse_sbi_format(self):
        csv_content = (
            "Txn Date,Value Date,Description,Cheque No.,Debit,Credit,Balance\n"
            "01 Jan 2024,01 Jan 2024,SALARY,,,75000,175000\n"
            "02 Jan 2024,02 Jan 2024,ZOMATO,,1500,,173500\n"
        )
        result = import_csv(csv_content)
        assert result.status.value == "completed"
        assert result.total == 2

    def test_empty_csv(self):
        result = import_csv("")
        assert result.status.value == "completed"
        assert result.total == 0

    def test_invalid_csv(self):
        result = import_csv("not,a,csv\n")
        assert result.status.value == "completed"
        assert result.total == 0

    def test_missing_columns(self):
        csv_content = "name,age\nAlice,30\nBob,25\n"
        result = import_csv(csv_content)
        assert result.status.value == "completed"

    def test_skip_empty_rows(self):
        csv_content = "date,description,amount\n\n2024-01-15,Test,100\n\n"
        result = import_csv(csv_content)
        assert result.status.value == "completed"
        assert result.total == 1

    def test_amount_parsing(self):
        assert parse_amount("1,234.56") == Decimal("1234.56")
        assert parse_amount("₹1,000") == Decimal("1000")
        assert parse_amount("") == Decimal("0")
        assert parse_amount("-500") == Decimal("-500")

    def test_date_parsing(self):
        assert parse_date("2024-01-15", "%Y-%m-%d") == date(2024, 1, 15)
        assert parse_date("15-01-2024", "%Y-%m-%d") == date(2024, 1, 15)
        assert parse_date("15 Jan 2024", "%Y-%m-%d") == date(2024, 1, 15)

    def test_detect_format(self):
        assert detect_format(["Date", "Narration", "Withdrawal Amount (INR)"]) == "hdfc"
        assert detect_format(["Value Date", "Transaction Description"]) == "icici"
        assert detect_format(["Txn Date", "Description"]) == "sbi"
        assert detect_format(["col1", "col2", "col3"]) == "generic"

    def test_credit_debit(self):
        csv_content = "date,description,debit,credit\n2024-01-15,Sale,,5000\n2024-01-16,Purchase,1500,\n"
        result = import_csv(csv_content)
        assert result.status.value == "completed"
        assert result.total == 2
        assert result.transactions[0].amount == Decimal("5000")
        assert result.transactions[1].amount == Decimal("1500")

    def test_bytes_input(self):
        content = b"date,description,amount\n2024-01-15,Test,100\n"
        result = import_csv(content)
        assert result.status.value == "completed"
        assert result.total == 1

    def test_utf8_bom(self):
        csv_content = "date,description,amount\n2024-01-15,Café,100\n"
        result = import_csv(csv_content)
        assert result.status.value == "completed"
        assert result.total == 1
        assert result.transactions[0].description == "Café"
