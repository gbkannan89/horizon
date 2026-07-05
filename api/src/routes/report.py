import logging
import io
from datetime import date
from fastapi import APIRouter, Depends
from fastapi.responses import StreamingResponse
from ..core.database import get_db
from ..schemas.auth import UserOut
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/report", tags=["Report"])

MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]


def _fmt(v):
    if v >= 10000000:
        return f"{v / 10000000:.1f}Cr"
    if v >= 100000:
        return f"{v / 100000:.1f}L"
    if v >= 1000:
        return f"{v / 1000:.1f}K"
    return f"{v:.0f}"


@router.get("/summary-pdf")
def generate_pdf(
    current_user: UserOut = Depends(get_current_user),
    conn = Depends(get_db)
):
    from reportlab.lib.pagesizes import A4
    from reportlab.lib.units import mm
    from reportlab.lib import colors
    from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle
    from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle

    user_id = current_user.id
    today = date.today()

    with conn.cursor() as cur:
        cur.execute("SELECT name, email FROM users WHERE id = %s", (user_id,))
        u = cur.fetchone()
        user_name = u[0] or "User" if u else "User"

        cur.execute("SELECT SUM(amount) FROM incomes WHERE user_id = %s", (user_id,))
        total_income = float(cur.fetchone()[0] or 0)

        cur.execute("SELECT SUM(amount) FROM assets WHERE user_id = %s", (user_id,))
        total_assets = float(cur.fetchone()[0] or 0)

        cur.execute("SELECT SUM(outstanding) FROM liabilities WHERE user_id = %s", (user_id,))
        total_liabilities = float(cur.fetchone()[0] or 0)

        net_worth = total_assets - total_liabilities

        cur.execute("SELECT SUM(amount) FROM expenses WHERE user_id = %s AND EXTRACT(MONTH FROM date) = %s AND EXTRACT(YEAR FROM date) = %s", (user_id, today.month, today.year))
        monthly_spent = float(cur.fetchone()[0] or 0)

        cur.execute("SELECT name, current_amount, target_amount FROM goals WHERE user_id = %s", (user_id,))
        goals = cur.fetchall()

    buf = io.BytesIO()
    doc = SimpleDocTemplate(buf, pagesize=A4, topMargin=20*mm, bottomMargin=20*mm)
    styles = getSampleStyleSheet()
    title_style = ParagraphStyle("Title2", parent=styles["Title"], fontSize=22, spaceAfter=4, textColor=colors.HexColor("#1E3A8A"))
    heading_style = ParagraphStyle("H2", parent=styles["Heading2"], fontSize=14, spaceBefore=12, spaceAfter=6, textColor=colors.HexColor("#1E3A8A"))
    normal = styles["Normal"]

    elements = []

    # Header
    elements.append(Paragraph(f"Horizon Financial Report", title_style))
    elements.append(Paragraph(f"{user_name} · {today.strftime('%d %B %Y')}", ParagraphStyle("Sub", parent=normal, fontSize=10, textColor=colors.grey)))
    elements.append(Spacer(1, 8*mm))

    # Net Worth Summary
    elements.append(Paragraph("Net Worth Summary", heading_style))
    nw_data = [
        ["Metric", "Amount"],
        ["Total Assets", f"₹{_fmt(total_assets)}"],
        ["Total Liabilities", f"₹{_fmt(total_liabilities)}"],
        ["Net Worth", f"₹{_fmt(net_worth)}"],
        ["Monthly Income", f"₹{_fmt(total_income)}"],
        ["Monthly Spent", f"₹{_fmt(monthly_spent)}"],
    ]
    nw_table = Table(nw_data, colWidths=[120*mm, 60*mm])
    nw_table.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, 0), colors.HexColor("#1E3A8A")),
        ("TEXTCOLOR", (0, 0), (-1, 0), colors.white),
        ("FONTNAME", (0, 0), (-1, 0), "Helvetica-Bold"),
        ("FONTSIZE", (0, 0), (-1, -1), 10),
        ("ALIGN", (1, 0), (-1, -1), "RIGHT"),
        ("GRID", (0, 0), (-1, -1), 0.5, colors.grey),
        ("ROWBACKGROUNDS", (0, 1), (-1, -1), [colors.white, colors.HexColor("#F8FAFC")]),
        ("TOPPADDING", (0, 0), (-1, -1), 6),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 6),
    ]))
    elements.append(nw_table)

    # Goals
    if goals:
        elements.append(Spacer(1, 6*mm))
        elements.append(Paragraph("Financial Goals", heading_style))
        goal_data = [["Goal", "Progress", "Target"]]
        for g in goals:
            pct = (float(g[1]) / float(g[2]) * 100) if float(g[2]) > 0 else 0
            goal_data.append([g[0], f"{pct:.0f}%", f"₹{_fmt(float(g[2]))}"])
        goal_table = Table(goal_data, colWidths=[80*mm, 40*mm, 60*mm])
        goal_table.setStyle(TableStyle([
            ("BACKGROUND", (0, 0), (-1, 0), colors.HexColor("#059669")),
            ("TEXTCOLOR", (0, 0), (-1, 0), colors.white),
            ("FONTNAME", (0, 0), (-1, 0), "Helvetica-Bold"),
            ("FONTSIZE", (0, 0), (-1, -1), 10),
            ("ALIGN", (1, 0), (2, -1), "RIGHT"),
            ("GRID", (0, 0), (-1, -1), 0.5, colors.grey),
            ("ROWBACKGROUNDS", (0, 1), (-1, -1), [colors.white, colors.HexColor("#F0FDF4")]),
            ("TOPPADDING", (0, 0), (-1, -1), 6),
            ("BOTTOMPADDING", (0, 0), (-1, -1), 6),
        ]))
        elements.append(goal_table)

    # Budget Summary
    elements.append(Spacer(1, 6*mm))
    elements.append(Paragraph("Budget Snapshot (50/30/20)", heading_style))
    budget_targets = {"Needs": total_income * 0.50, "Wants": total_income * 0.30, "Savings": total_income * 0.20}
    budget_data = [["Bucket", "Budget", "Spent", "Remaining"]]
    for bucket in ["Needs", "Wants", "Savings"]:
        budget_data.append([bucket, f"₹{_fmt(budget_targets[bucket])}", "—", "—"])
    budget_table = Table(budget_data, colWidths=[50*mm, 50*mm, 50*mm, 50*mm])
    budget_table.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, 0), colors.HexColor("#6B46C1")),
        ("TEXTCOLOR", (0, 0), (-1, 0), colors.white),
        ("FONTNAME", (0, 0), (-1, 0), "Helvetica-Bold"),
        ("FONTSIZE", (0, 0), (-1, -1), 10),
        ("ALIGN", (1, 0), (-1, -1), "RIGHT"),
        ("GRID", (0, 0), (-1, -1), 0.5, colors.grey),
        ("ROWBACKGROUNDS", (0, 1), (-1, -1), [colors.white, colors.HexColor("#F5F3FF")]),
        ("TOPPADDING", (0, 0), (-1, -1), 6),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 6),
    ]))
    elements.append(budget_table)

    # Footer
    elements.append(Spacer(1, 15*mm))
    elements.append(Paragraph("Generated by Horizon — Wealth Intelligence Engine", ParagraphStyle("Footer", parent=normal, fontSize=8, textColor=colors.grey, alignment=1)))

    doc.build(elements)
    buf.seek(0)

    return StreamingResponse(buf, media_type="application/pdf", headers={"Content-Disposition": f"attachment; filename=horizon_report_{today.isoformat()}.pdf"})
