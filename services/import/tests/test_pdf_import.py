from decimal import Decimal
from pathlib import Path
from tempfile import NamedTemporaryFile

from app.services.pdf_parser import parse_pdf


def _make_simple_pdf(text: str) -> bytes:
    from reportlab.lib.pagesizes import A4
    from reportlab.lib.styles import getSampleStyleSheet
    from reportlab.platypus import Paragraph, SimpleDocTemplate, Spacer, Table, TableStyle

    buf = NamedTemporaryFile(delete=False, suffix=".pdf")
    buf.close()

    doc = SimpleDocTemplate(buf.name, pagesize=A4)
    styles = getSampleStyleSheet()
    elements = []

    for line in text.split("\n"):
        if "\t" in line:
            cells = [c.strip() for c in line.split("\t")]
            t = Table([cells])
            t.setStyle(TableStyle([("FONTSIZE", (0, 0), (-1, -1), 8)]))
            elements.append(t)
        else:
            elements.append(Paragraph(line, styles["Normal"]))
            elements.append(Spacer(1, 4))

    doc.build(elements)
    content = Path(buf.name).read_bytes()
    Path(buf.name).unlink()
    return content


class TestPDFParser:
    def test_hdfc_statement(self):
        pdf_text = (
            "HDFC BANK\n"
            "Savings Account Statement\n"
            "Account No: 12345678\n"
            "Date\tNarration\tChq./Ref.No.\tValue Dat\tWithdrawal Amount (INR)\tDeposit Amount (INR)\tBalance (INR)\n"
            "01-01-2024\tSALARY\t\t01-01-2024\t\t50000\t100000\n"
            "02-01-2024\tAMAZON PAY\t\t02-01-2024\t2500\t\t97500\n"
            "Statement Summary\n"
        )
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) >= 1
        salary_txns = [t for t in result if t.event_type == "Income"]
        assert len(salary_txns) >= 1

    def test_icici_statement(self):
        header = (
            "Value Date\tTransaction Date\tTransaction Description\t"
            "Reference No.\tWithdrawal Amount\tDeposit Amount\tBalance"
        )
        pdf_text = (
            "ICICI Bank\n"
            f"{header}\n"
            "15-01-2024\t14-01-2024\tNEFT CREDIT\tREF001\t\t50000\t100000\n"
            "16-01-2024\t15-01-2024\tSWIGGY\tREF002\t1200\t\t98700\n"
        )
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) == 2

    def test_sbi_statement(self):
        pdf_text = (
            "State Bank of India\n"
            "Txn Date\tValue Date\tDescription\tCheque No.\tDebit\tCredit\tBalance\n"
            "01 Jan 2024\t01 Jan 2024\tSALARY\t\t\t75000\t175000\n"
            "02 Jan 2024\t02 Jan 2024\tZOMATO\t\t1500\t\t173500\n"
        )
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) == 2

    def test_generic_statement(self):
        pdf_text = (
            "Transaction Statement\n"
            "date\tdescription\tamount\n"
            "2024-01-15\tSalary\t50000\n"
            "2024-01-16\tGroceries\t-2500\n"
        )
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) == 2
        assert result[0].event_type == "Income"
        assert result[1].event_type == "Expense"

    def test_empty_pdf(self):
        content = _make_simple_pdf("")
        result = parse_pdf(content)
        assert len(result) == 0

    def test_no_tables_fallback(self):
        pdf_text = (
            "Bank Statement\n"
            "01-01-2024    Salary Credit                  50000\n"
            "02-01-2024    Amazon Purchase                2500\n"
            "03-01-2024    Interest Credited              1200\n"
        )
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) >= 2

    def test_multi_page(self):
        pdf_text = (
            "date\tdescription\tamount\n"
            "2024-01-15\tTransaction 1\t100\n"
            "2024-01-16\tTransaction 2\t200\n"
        )
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) == 2

    def test_debit_credit_columns(self):
        pdf_text = (
            "date\tdescription\tdebit\tcredit\n"
            "2024-01-15\tSale\t\t5000\n"
            "2024-01-16\tPurchase\t1500\t\n"
        )
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) == 2
        assert result[0].amount == Decimal("5000")
        assert result[1].amount == Decimal("1500")

    def test_bytes_input(self):
        pdf_text = "date\tdescription\tamount\n2024-01-15\tTest\t100\n"
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) == 1

    def test_path_input(self):
        pdf_text = "date\tdescription\tamount\n2024-01-15\tPathTest\t200\n"
        content = _make_simple_pdf(pdf_text)
        tmp = NamedTemporaryFile(delete=False, suffix=".pdf")
        tmp.write(content)
        tmp.close()
        result = parse_pdf(Path(tmp.name))
        assert len(result) == 1
        assert result[0].amount == Decimal("200")
        Path(tmp.name).unlink()

    def test_header_footer_removal(self):
        pdf_text = (
            "HDFC BANK\n"
            "Registered Office: Mumbai\n"
            "Page 1 of 1\n"
            "date\tdescription\tamount\n"
            "2024-01-15\tSalary\t50000\n"
            "Continuation Sheet\n"
        )
        content = _make_simple_pdf(pdf_text)
        result = parse_pdf(content)
        assert len(result) == 1
